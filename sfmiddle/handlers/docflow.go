package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/auth"
	"sfmiddle/db"
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

	// validacion deñ hash del documento
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
		docData[d.IdDoc]["idFD"] = idFD
	}
	respWriter.Header().Set("Content-Type", "application/json")
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

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	idTeam, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	fmt.Println(idInst, idUser, idTeam)

	// ===== Leer y parsear el JSON =====
	var requestData models.InviteRequest
	if err := json.NewDecoder(request.Body).Decode(&requestData); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	// ===== Procesar cada documento =====
	for idDocument, doc := range requestData {

		// ===== SELECT en tabla folderdocuments =====
		attrs := []string{"idFolder"}
		wheres := map[string][]string{
			"idfolderdocument": {doc.IdFD},
			"idDocument":       {idDocument},
		}

		folderData, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument", attrs, wheres)
		if err != nil || len(folderData) == 0 {
			http.Error(respWriter, "Error al obtener datos del folder: "+err.Error(), http.StatusInternalServerError)
			return
		}

		idFolder := folderData[doc.IdFD]["idFolder"]

		for _, reviewer := range doc.Reviewers {
			idInvite := uuid.NewString()
			// ===== INSERT en tabla invites =====
			inviteCols := []string{
				"idInvite", "idFolder_fk", "idDocument_fk", "idInstitutionDest_fk",
				"idUserDest_fk", "idTeamDest_fk", "expirationDate", "signerViewer", "descriptionText",
			}

			inviteAttrs := []interface{}{
				idInvite,         // idInvite
				idFolder,         // idFolder_fk
				idDocument,       // idDocument_fk
				idInst,           // idInstitutionDest_fk
				reviewer.User,    // idUserDest_fk
				reviewer.Team,    // idTeamDest_fk
				reviewer.DueDate, // expirationDate
				reviewer.Role,    // signerViewer
				reviewer.Comment, // descriptionText
			}

			_, err := db.DB_con.GenericInsert("invites", inviteCols, inviteAttrs)
			if err != nil {
				http.Error(respWriter, "Error al insertar invite: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// ===== Consultar institución del reviewer =====
			attrs := []string{"idInstitution_fk"}
			wheres := map[string][]string{
				"idUser": {reviewer.User},
			}

			userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
			if err != nil || len(userData) == 0 {
				http.Error(respWriter, "Error al obtener institución del usuario: "+err.Error(), http.StatusInternalServerError)
				return
			}

			reviewerInst := userData[reviewer.User]["idInstitution_fk"]

			// ===== INSERT en tabla signatures =====
			signatureCols := []string{
				"idInvite_fk", "idFolder_fk", "digestValueSign",
				"idInstitution_fk", "idUser_fk", "idTeam_fk",
			}

			signatureAttrs := []interface{}{
				idInvite,         // idInvite_fk
				idFolder,         // idFolder_fk
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
		}
	}

	// ===== Respuesta exitosa =====
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Invitaciones y firmas creadas exitosamente",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
