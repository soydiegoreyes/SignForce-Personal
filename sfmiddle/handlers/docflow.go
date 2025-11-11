package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/auth"
	"sfmiddle/db"
	"sfmiddle/documentflow"
	"sfmiddle/models"
	"time"

	"github.com/google/uuid"
)

// Handler que crea un nuevo folder para firma con los documentos que estarán dentro
func NewSignFolder(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== Autenticación por JWT =====
	cookie, err := request.Cookie("token")
	if err != nil {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	idTeam, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====
	type DocDataRequest struct {
		IdDoc   string `json:"idDoc,omitempty"`
		HashDoc string `json:"hashDoc,omitempty"`
	}

	var req []DocDataRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	attrs := []string{
		"documentHash", "createdAtDoc", "lastModifiedDoc",
		"documentPath", "documentName", "documentExt",
		"authUseStatus", "authRoleStatus", "activeDoc",
	}

	wheres := map[string][]string{
		"idDocument":        {},
		"creatorUserDoc_fk": {idUser},
		"ownerInstDoc_fk":   {idInst},
		"ownerTeamDoc_fk":   {idTeam},
	}

	for _, d := range req {
		wheres["idDocument"] = append(wheres["idDocument"], d.IdDoc)
	}
	// ===== Ejecutar SELECT =====
	docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener documentos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// validacion del hash del documento
	for _, d := range req {
		if docData[d.IdDoc]["documentHash"] != d.HashDoc {
			docData[d.IdDoc]["activeDoc"] = "-1"
		}
	}

	// se crea el nuevo folder
	folderPath := fmt.Sprintf("%s/folders/%s/%s/%s/", os.Getenv("TEMP_BASE_PATH"), idInst, idTeam, idUser)
	expiration := time.Now().AddDate(0, 0, 3).Format("2006-01-02 15:04:05")
	columns := []string{"ownerInst_fk", "ownerTeam_fk", "creatorUser_fk", "pathSerialized", "expirationDate", "numDocs"}
	values := []interface{}{idInst, idTeam, idUser, folderPath, expiration, len(docData)}

	idFolder, err := db.DB_con.GenericInsert("folders", columns, values)
	if err != nil {
		http.Error(respWriter, "Error creando folder", http.StatusInternalServerError)
	}

	// se asocian los documentos enviados al folder creado
	columns = []string{"idfolderdocument", "idDocument", "idFolder"}
	for _, d := range req {
		idFD := uuid.NewString()
		values = []interface{}{idFD, d.IdDoc, idFolder}
		_, err := db.DB_con.GenericInsert("folderdocuments", columns, values)
		if err != nil {
			fmt.Println("Error al insertar documento en folderdocuments", err)
		}
		docData[d.IdDoc]["idFolder"] = idFolder
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	json.NewEncoder(respWriter).Encode(docData)
}

// Handler para cerrar folder y enviar invitaciones a firma
func CloseAndInvite(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== Autenticación por JWT =====
	cookie, err := request.Cookie("token")
	if err != nil {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	_, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	_, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	// ===== Leer y parsear el JSON =====
	var requestData models.InviteRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	// ===== Procesar cada documento =====
	for _, doc := range requestData {
		// invitaciones para el documento
		invites := []models.InviteMail{}

		idInvite := uuid.NewString()
		sentAt := time.Now().Format("2006-01-02 15:04:05")
		for i, reviewer := range doc.Reviewers {

			// ===== Consultar datos del reviewer =====
			attrs := []string{"idInstitution_fk", "nameUser", "lastNameUser", "emailUser"}
			wheres := map[string][]string{
				"idUser": {reviewer.User},
			}

			userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
			if err != nil || len(userData) == 0 {
				http.Error(respWriter, "Error al obtener institución del usuario: "+err.Error(), http.StatusInternalServerError)
				return
			}

			reviewerInst := userData[reviewer.User]["idInstitution_fk"]
			reviewerName := fmt.Sprintf("%s %s", userData[reviewer.User]["nameUser"], userData[reviewer.User]["lastNameUser"])
			reviewerEmail := userData[reviewer.User]["emailUser"]

			// ===== Consultar institución del reviewer =====
			attrs = []string{"legalNameInst", "aliasNameInst", "activeInst"}
			wheres = map[string][]string{
				"idInstitution": {reviewerInst},
			}

			InstData, err := db.DB_con.GenericSelect("institutions", "idInstitution", attrs, wheres)
			if err != nil || len(InstData) == 0 {
				http.Error(respWriter, "Error al obtener institución del usuario: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// ===== INSERT en tabla signatures =====
			signatureCols := []string{
				"idInvite_fk", "idFolder", "digestValueSign",
				"idInstitution_fk", "idUser_fk", "idTeam_fk",
			}

			signatureAttrs := []interface{}{
				idInvite,         // idInvite_fk
				doc.IdFolder,     // idFolder
				doc.DocumentHash, // digestValueSign
				reviewerInst,     // idInstitution_fk
				reviewer.User,    // idUser_fk
				reviewer.Team,    // idTeam_fk
			}

			idSignature, err := db.DB_con.GenericInsert("signatures", signatureCols, signatureAttrs)
			if err != nil {
				http.Error(respWriter, "Error al insertar signature: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// ===== INSERT en tabla signstamps (para cada posición del reviewer) =====
			for _, pos := range reviewer.SignPositions {
				signstampCols := []string{
					"idSignature_fk", "xSign", "ySign", "wSign", "hSign", "pageSign",
				}

				signstampAttrs := []interface{}{
					idSignature, // idSignature_fk
					pos.X,       // xSign
					pos.Y,       // ySign
					pos.Width,   // wSign
					pos.Height,  // hSign
					pos.Page,    // pageSign
				}

				_, err := db.DB_con.GenericInsert("signstamps", signstampCols, signstampAttrs)
				if err != nil {
					http.Error(respWriter, "Error al insertar signstamp: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
			if i == len(doc.Reviewers)-1 {
				// ===== INSERT en tabla invites =====
				inviteCols := []string{
					"idInvite", "idFolder", "idInstitutionDest_fk",
					"idUserDest_fk", "idTeamDest_fk", "expirationDate", "sentAt", "descriptionText",
				}

				inviteAttrs := []interface{}{
					idInvite,         // idInvite
					doc.IdFolder,     // idFolder
					idInst,           // idInstitutionDest_fk
					reviewer.User,    // idUserDest_fk
					reviewer.Team,    // idTeamDest_fk
					reviewer.DueDate, // expirationDate
					sentAt,           // sentAt
					reviewer.Comment, // descriptionText
				}

				_, err := db.DB_con.GenericInsert("invites", inviteCols, inviteAttrs)
				if err != nil {
					http.Error(respWriter, "Error al insertar invite: "+err.Error(), http.StatusInternalServerError)
					return
				}

				invite := models.InviteMail{
					IdUser:        reviewer.User,
					ReviewerName:  reviewerName,
					ReviewerEmail: reviewerEmail,
					ReviewerInst:  fmt.Sprintf("%s (%s)", InstData[reviewerInst]["legalNameInst"], InstData[reviewerInst]["aliasNameInst"]),
					SentDate:      sentAt,
					SenderMessage: reviewer.Comment,
					UrlSignLink:   fmt.Sprintf("http://%s:%s/viewinvite?idInvite=%s", os.Getenv("API_IP"), os.Getenv("API_PORT"), idInvite),
				}
				invites = append(invites, invite)
			}
		}

		err = documentflow.GenerateInvites(invites)
		if err != nil {
			fmt.Println(err)
		}
	}

	// ===== Respuesta exitosa =====
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Invitaciones y firmas creadas exitosamente",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func GetInvite(respWriter http.ResponseWriter, request *http.Request) {
	// Validar método
	if request.Method != http.MethodPost {
		respWriter.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	// Decodificar request
	var req models.SignByInviteRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Request inválido",
		})
		return
	}

	// Validar que idInvite no esté vacío
	if req.IdInvite == "" {
		respWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "idInvite es requerido",
		})
		return
	}

	// Cargar información de la invitación
	inviteInfo, err := documentflow.LoadInviteInfo(req.IdInvite)
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(inviteInfo)
}

func GetFolder(respWriter http.ResponseWriter, request *http.Request) {
	// Validar método
	if request.Method != http.MethodPost {
		respWriter.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Método no permitido",
		})
		return
	}

	// Decodificar request

	var req models.SignByInviteRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Request inválido",
		})
		return
	}

	// Validar que idInvite no esté vacío
	if req.IdInvite == "" {
		respWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "idInvite es requerido",
		})
		return
	}

	// Cargar información de la invitación
	folderInfo, err := documentflow.LoadInviteInfo(req.IdInvite)
	if err != nil {
		respWriter.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(folderInfo)
}
