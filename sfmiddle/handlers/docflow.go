package handlers

import (
	"bytes"
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

// Funcion para devolver datos de folders
func GetFolders(respWriter http.ResponseWriter, request *http.Request) {
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
	_, ok2 := claims["iid"].(string)
	idTeam, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====
	var req models.FolderRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	// ===== Validaciones básicas =====
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// ===== Llamar a la función principal =====
	folderData, err := documentflow.LoadFolderInfo(
		req.IdFolder,
		idTeam,
		idUser,
		req.OnlyShared,
		req.OnlyTeam,
		req.OnlyUser,
		req.Page,
		req.PageSize,
		req.OrderBy,
		req.OrderDir,
	)

	if err != nil {
		http.Error(respWriter, "Error al obtener folders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ===== Respuesta JSON =====
	respWriter.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(respWriter).Encode(folderData); err != nil {
		http.Error(respWriter, "Error al generar respuesta", http.StatusInternalServerError)
		return
	}
}

// ===============================================================================================
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
	_, ok2 := claims["iid"].(string)
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
	// se toman los reviewers unicos de todos los documentos y se le añaden los documentos que debe firmar
	// se obtienen los datos para la invitacion por email del reviewer
	// una invitacion por usuario y cada invitacion tiene un link, al entrar al link debe tener los documentos del folder que el usuario va a firmar solamente
	// doc tiene datos del documento como hash, ruta, id

	// lista de invitaciones que serán identificadas cada una por el id del usuario
	invites := make(map[string]*models.Invite)
	invitesMails := make(map[string]models.InviteMail) // cada usuario tiene su invite mail y se deben de enviar como lista
	var numSigners, numReviewers int
	var idFolder string
	for _, docrev := range requestData { // -> por cada documento dentro del request
		idFolder = docrev.IdFolder
		for i, reviewer := range docrev.Reviewers { // -> por cada revisor dentro del documento
			// si dentro de la lista de invitaciones NO esta el idUser entonces se crea una nueva invitación
			if _, ok := invites[reviewer.User]; !ok {
				invites[reviewer.User] = &models.Invite{ //-> dado que es una invitacion por usuario entonces el id es el del usuario como director de la invitacion
					IdInvite:   uuid.NewString(), // -> el uuid es para que la invitacion se identifique en base de datos
					UserDest:   reviewer,         // -> dado que el userDest es el reviewer y son la misma estructura se asigna directamente
					SentAt:     time.Now().Format("2006-01-02 15:04:05"),
					InviteDocs: []models.InviteDoc{}, // -> el usuario solo podrá ver los documentos que esten dentro de esta lista aunque sea el mismo folder
				}

				// ===== INSERT en tabla invites =====
				inviteCols := []string{"idInvite", "idFolder", "idUserOwnner_fk", "idUserDest_fk", "expirationDate", "sentAt", "descriptionText"}

				inviteAttrs := []interface{}{
					invites[reviewer.User].IdInvite, // idInvite
					idFolder,                        //idFolder
					idUser,                          // idUserOwnner
					reviewer.User,                   // idUserDest_fk
					reviewer.DueDate,                // expirationDate
					invites[reviewer.User].SentAt,   // sentAt
					reviewer.Comment,                // descriptionText
				}

				_, err = db.DB_con.GenericInsert("invites", inviteCols, inviteAttrs)
				if err != nil {
					http.Error(respWriter, "Error al insertar invite: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}
			// se inserta en InviteDocs los datos del documento que hay que firmar (uno por documento)
			var invdoc = models.InviteDoc{
				ForSign:   reviewer.Role == 1,
				ExpiresAt: reviewer.DueDate,
				Comment:   reviewer.Comment,
				Order:     i,
				Doc:       docrev.Document,
			}

			// se inserta la invitacion con el id del usuario
			invites[reviewer.User].InviteDocs = append(invites[reviewer.User].InviteDocs, invdoc)

			// ===== Obtener datos del reviewer para personalizar la invitacion de email =====
			// como solo se cuenta con el id del usuario al que se dirije la inv se debe sacar su institucion y correo y nombre
			attrs := []string{"idInstitution_fk", "nameUser", "lastNameUser", "emailUser"}
			wheres := map[string][]string{
				"idUser": {reviewer.User},
			}
			userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
			if err != nil || len(userData) == 0 {
				http.Error(respWriter, "Error al obtener institución del usuario: "+err.Error(), http.StatusInternalServerError)
				return
			}
			reviewerName := fmt.Sprintf("%s %s", userData[reviewer.User]["nameUser"], userData[reviewer.User]["lastNameUser"])
			reviewerEmail := userData[reviewer.User]["emailUser"]

			// ===== Consultar institución del reviewer para personalizar email =====
			attrs = []string{"legalNameInst", "aliasNameInst", "activeInst"}
			wheres = map[string][]string{
				"idInstitution": {userData[reviewer.User]["idInstitution_fk"]},
			}
			InstData, err := db.DB_con.GenericSelect("institutions", "idInstitution", attrs, wheres)
			if err != nil || len(InstData) == 0 {
				http.Error(respWriter, "Error al obtener institución del usuario: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// ===== INSERT en tabla signatures =====
			signatureCols := []string{
				"idInvite_fk", "digestValueSign", "idUser_fk",
			}

			signatureAttrs := []interface{}{
				invites[reviewer.User].IdInvite, // idInvite_fk
				docrev.Document.DocumentHash,    // digestValueSign
				reviewer.User,                   // idUser_fk
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

			wheres = map[string][]string{
				"idFolder": {docrev.IdFolder},
			}

			folderdocData, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument", []string{"idDocument"}, wheres)
			if err != nil {
				http.Error(respWriter, "Error al obtener documentos del folder: "+err.Error(), http.StatusInternalServerError)
				return
			}

			inviteDetCols := []string{"idInvite", "idfolderdocument"}

			for idfolderdoc := range folderdocData {
				inviteDetAttrs := []interface{}{
					invites[reviewer.User].IdInvite,
					idfolderdoc,
				}

				_, err = db.DB_con.GenericInsert("invitesdetail", inviteDetCols, inviteDetAttrs)
				if err != nil {
					http.Error(respWriter, "Error al insertar invite: "+err.Error(), http.StatusInternalServerError)
					return
				}
			}

			if _, ok := invitesMails[reviewer.User]; !ok {
				reviewerInst := userData[reviewer.User]["idInstitution_fk"]

				inviteEmail := models.InviteMail{
					IdUser:        reviewer.User,
					ReviewerName:  reviewerName,
					ReviewerEmail: reviewerEmail,
					ReviewerInst:  fmt.Sprintf("%s (%s)", InstData[reviewerInst]["legalNameInst"], InstData[reviewerInst]["aliasNameInst"]),
					SentDate:      invites[reviewer.User].SentAt,
					SenderMessage: reviewer.Comment,
					UrlSignLink:   fmt.Sprintf("http://%s:%s/viewinvite?idInvite=%s", os.Getenv("API_IP"), os.Getenv("API_PORT"), invites[reviewer.User].IdInvite),
				}
				invitesMails[reviewer.User] = inviteEmail
			}
			if reviewer.Role == 1 {
				numSigners += 1
			} else {
				numReviewers += 1
			}
		}

		updates := map[string]map[string]interface{}{
			idFolder: {
				"numSigners":   numSigners,
				"numReceivers": numReviewers,
			},
		}
		err := db.DB_con.GenericBatchUpdate("folders", "idFolder", updates)
		if err != nil {
			http.Error(respWriter, "Error actualizando folder", http.StatusInternalServerError)
		}

		imails := []models.InviteMail{}
		for _, invs := range invitesMails {
			imails = append(imails, invs)
		}
		err = documentflow.GenerateInvites(imails)
		if err != nil {
			http.Error(respWriter, "Error al generar invite: "+err.Error(), http.StatusInternalServerError)
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
	_, ok2 := claims["iid"].(string)
	_, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Decodificar request
	var req models.GetByInviteRequest
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
	if idUser != inviteInfo.Folder.Invites[0].UserDest.User {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(inviteInfo)
}

// firma de documento mediante invitacion
func SignDocument(respWriter http.ResponseWriter, request *http.Request) {
	// Validar método
	if request.Method != http.MethodPost {
		respWriter.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Método no permitido",
		})
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
	_, ok2 := claims["iid"].(string)
	_, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	// Decodificar request
	var reqdoc models.SignDoc
	if err := json.NewDecoder(request.Body).Decode(&reqdoc); err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Request inválido",
		})
		return
	}

	// Validar que idInvite no esté vacío
	if reqdoc.IdInvite == "" || reqdoc.IdFolder == "" || reqdoc.IdKey == "" || len(reqdoc.Documents) == 0 {
		respWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(respWriter).Encode(map[string]string{
			"error": "Datos faltantes para completar la operación",
		})
		return
	}

	// ============= Validar que el usuario coincide con la invitación =============
	// Cargar información de la invitación
	inviteInfo, err := documentflow.LoadInviteInfo(reqdoc.IdInvite)
	if err != nil {
		http.Error(respWriter, "Invitación no encontrada", http.StatusBadRequest)
		return
	}

	if idUser != inviteInfo.Folder.Invites[0].UserDest.User {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// ============= Validar que la llave pertenezca al usuario =======================
	attrs := []string{"activeUser", "idKeysUser_fk", "isAliveUser"}
	wheres := map[string][]string{
		"idUser": {idUser},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil || len(userData) == 0 {
		http.Error(respWriter, "Error al obtener institución del usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if userData[idUser]["activeUser"] != "1" {
		http.Error(respWriter, "Usuario no se encuentra activo", http.StatusUnauthorized)
		return
	}
	if userData[idUser]["isAliveUser"] != "1" {
		http.Error(respWriter, "Se necesita prueba de vida", http.StatusUnauthorized)
		return
	}
	if userData[idUser]["idKeysUser_fk"] != reqdoc.IdKey {
		http.Error(respWriter, "La llave no pertenece al usuario autenticado", http.StatusUnauthorized)
		return
	}
	//=====================================================================
	whereMap := map[string][]string{
		"nameApp": {"sfback"},
	}
	appdata, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, whereMap)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	var host, port string
	for _, v := range appdata {
		host = v["domainApp"]
		port = v["portApp"]
		break
	}

	jsonPayload, err := json.Marshal(reqdoc) // se envia tal cual llegó la peticion
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	reqsign, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/signdocument", host, port), bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	reqsign.Header.Set("Content-Type", "application/json")
	reqsign.Header.Set("Authorization", "Bearer "+reqdoc.Aut)
	// Ejecutar petición
	client := &http.Client{}
	resp, err := client.Do(reqsign)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer resp.Body.Close()

	var signResult models.SignResponse

	if resp.StatusCode == http.StatusOK {
		fmt.Println("TODO OK")
		err = json.NewDecoder(resp.Body).Decode(&signResult)
		if err != nil {
			fmt.Println("Error decodificando:", err)
		}
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&signResult)
}
