package documentflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/objects"
	"sfmiddle/utilities"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateInvites(invites []models.InviteMail) error {
	whereMap := map[string][]string{
		"nameApp": {"emailServ"},
	}
	data, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, whereMap)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	binDoc, err := os.ReadFile("./templates/sign_invite.html")
	if err != nil {
		fmt.Printf("%s", err)
	}
	for _, invite := range invites {
		body := string(binDoc)
		body = strings.ReplaceAll(body, "{SENDER_NAME}", invite.ReviewerName)
		body = strings.ReplaceAll(body, "{SENDER_EMAIL}", invite.ReviewerEmail)
		body = strings.ReplaceAll(body, "{SENDER_INST}", invite.ReviewerInst)
		body = strings.ReplaceAll(body, "{SENT_DATE}", invite.SentDate)
		body = strings.ReplaceAll(body, "{SENDER_INST}", invite.ReviewerInst)
		body = strings.ReplaceAll(body, "{SENDER_MESSAGE}", invite.SenderMessage)
		body = strings.ReplaceAll(body, "{URL_SIGN_LINK}", invite.UrlSignLink)

		payload := models.EmailRequest{
			IdUser:   invite.IdUser,
			Subject:  fmt.Sprintf("¡Te invitaron a firmar! %s", invite.ReviewerName),
			Body:     body,
			Dest:     []string{invite.ReviewerEmail},
			MimeType: "html",
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return err
		}
		var host, port string
		for _, v := range data {
			host = v["domainApp"]
			port = v["portApp"]
			break
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/mailserv", host, port), bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
		// Ejecutar petición
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s", err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			fmt.Println("TODO OK")
		}

	}
	return nil
}

func LoadInviteInfo(idInvite string) (*models.InviteInfoResp, error) {
	// obtenemos datos de la invitación (solo una por usuario)
	attrs := []string{"idFolder", "idUserOwnner_fk", "idUserDest_fk", "requireAliveProof",
		"expirationDate", "sentAt", "descriptionText"}
	wheres := map[string][]string{
		"idInvite": {idInvite},
	}
	inviteData, err := db.DB_con.GenericSelect("invites", "idInvite", attrs, wheres)

	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos de invitación: %v", err)
	}
	if len(inviteData) == 0 {
		return nil, fmt.Errorf("invitación no encontrada")
	}

	// Obtener datos del folder
	attrs = []string{"secuentialSign", "expirationDate", "creatorUser_fk", "ownerInst_fk"}
	wheres = map[string][]string{
		"idFolder": {inviteData[idInvite]["idFolder"]},
	}
	folderData, err := db.DB_con.GenericSelect("folders", "idFolder", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del folder: %v", err)
	}

	if len(folderData) == 0 {
		return nil, fmt.Errorf("folder no encontrado")
	}

	idFolder := inviteData[idInvite]["idFolder"]

	// Obtener datos del usuario destinatario
	attrs = []string{"activeUser", "idKeysUser_fk", "nameUser", "lastNameUser", "emailUser", "idInstitution_fk"}
	wheres = map[string][]string{
		"idUser": {inviteData[idInvite]["idUserDest_fk"]},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del usuario: %v", err)
	}
	if len(userData) == 0 {
		return nil, fmt.Errorf("usuario destinatario no encontrado")
	}

	// Obtener datos del usuario emisor
	wheres = map[string][]string{
		"idUser": {inviteData[idInvite]["idUserOwnner_fk"]},
	}

	ownerUserData, err := db.DB_con.GenericSelect("users", "idUser",
		[]string{"nameUser", "lastNameUser", "emailUser"}, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del usuario emisor: %v", err)
	}

	// Obtener datos de la institución del emisor
	wheres = map[string][]string{
		"idInstitution": {folderData[idFolder]["ownerInst_fk"]},
	}
	instData, err := db.DB_con.GenericSelect("institutions", "idInstitution",
		[]string{"legalNameInst", "aliasNameInst"}, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos de la institución: %v", err)
	}

	// Obtener datos del equipo del emisor
	/*
		attrs = []string{"nameTeam"}
		wheres = map[string][]string{
			"idTeam": {folderData[idFolder]["ownerTeam_fk"]},
		}
		teamData, err := db.DB_con.GenericSelect("teams", "idTeam", attrs, wheres)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo datos del equipo: %v", err)
		}
	*/

	// Obtener los documentos de la invitación específica
	attrs = []string{"idfolderdocument"}
	wheres = map[string][]string{
		"idInvite": {idInvite},
	}

	inviteDocsData, err := db.DB_con.GenericSelect("invitesdetail", "idfolderdocument", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos de la invitación: %v", err)
	}

	if len(inviteDocsData) == 0 {
		return nil, fmt.Errorf("no se encontraron documentos en la invitación")
	}

	// Obtener los idDocument de los folderdocuments
	var docIds []string
	for _, doc := range inviteDocsData {
		wheres = map[string][]string{
			"idfolderdocument": {doc["idfolderdocument"]},
		}
		folderDocData, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument",
			[]string{"idDocument"}, wheres)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo documentos del folder: %v", err)
		}
		for _, fd := range folderDocData {
			docIds = append(docIds, fd["idDocument"])
		}
	}

	// Obtener datos de los documentos
	attrs = []string{"idDocument", "documentPath", "documentName", "documentExt", "documentHash",
		"authUseStatus", "authRoleStatus", "activeDoc", "abstractDoc", "createdAtDoc", "lastModifiedDoc"}
	wheres = map[string][]string{
		"idDocument": docIds,
	}
	docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del documento: %v", err)
	}

	if len(docData) == 0 {
		return nil, fmt.Errorf("documentos no encontrados")
	}

	// Construir UserInfo del emisor
	userEmisor := models.UserInfo{
		NameInstEmisor: fmt.Sprintf("%s (%s)",
			instData[folderData[idFolder]["ownerInst_fk"]]["legalNameInst"],
			instData[folderData[idFolder]["ownerInst_fk"]]["aliasNameInst"]),
		NameUserEmisor: fmt.Sprintf("%s %s",
			ownerUserData[inviteData[idInvite]["idUserOwnner_fk"]]["nameUser"],
			ownerUserData[inviteData[idInvite]["idUserOwnner_fk"]]["lastNameUser"]),
		//NameTeamEmisor: teamData[folderData[idFolder]["ownerTeam_fk"]]["name"],
	}

	// Construir UserDestInfo del destinatario
	userDest := models.UserDestInfo{
		NameUserDest: fmt.Sprintf("%s %s",
			userData[inviteData[idInvite]["idUserDest_fk"]]["nameUser"],
			userData[inviteData[idInvite]["idUserDest_fk"]]["lastNameUser"]),
		// Nota: Necesitarías obtener nombre de institución y equipo del destinatario si es necesario
	}
	fmt.Println(userDest)
	// Construir la invitación con documentos
	var inviteDocs []models.InviteDoc
	for _, docId := range docIds {
		if doc, exists := docData[docId]; exists {
			if doc["activeDoc"] != "1" {
				return nil, fmt.Errorf("documento %s ya no se encuentra activo", docId)
			}

			documentFullName := fmt.Sprintf("%s.%s", doc["documentName"], doc["documentExt"])
			filePath := fmt.Sprintf("./%s%s", doc["documentPath"], documentFullName)

			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				return nil, fmt.Errorf("archivo no encontrado en el sistema: %s", filePath)
			}

			hash, err := utilities.GetHash(filePath, configs.HashConf)
			if err != nil {
				return nil, fmt.Errorf("error al obtener hash: %v", err)
			}

			if doc["documentHash"] != hash {
				return nil, fmt.Errorf("hash no corresponde al archivo %s", documentFullName)
			}

			inviteDoc := models.InviteDoc{
				Doc: models.Document{
					IdDocument:       docId,
					ActiveDoc:        doc["activeDoc"],
					AuthRoleStatus:   doc["authRoleStatus"],
					AuthUseStatus:    doc["authUseStatus"],
					CreatedAtDoc:     doc["createdAtDoc"],
					DocumentExt:      doc["documentExt"],
					DocumentHash:     doc["documentHash"],
					DocumentName:     doc["documentName"],
					DocumentPath:     doc["documentPath"],
					DocumentFullName: documentFullName,
					Abstract:         doc["abstractDoc"],
					LastModifiedDoc:  doc["lastModifiedDoc"],
				},
				ForSign:   true, // Esto debería venir de la base de datos según el rol del usuario
				ExpiresAt: inviteData[idInvite]["expirationDate"],
				Comment:   inviteData[idInvite]["descriptionText"],
			}
			inviteDocs = append(inviteDocs, inviteDoc)
		}
	}

	// Construir la Invite completa
	invite := models.Invite{
		IdInvite:   idInvite,
		UserDest:   models.Reviewer{User: inviteData[idInvite]["idUserDest_fk"]}, // Solo el ID por ahora
		AliveProof: inviteData[idInvite]["requireAliveProof"] == "1",
		SentAt:     inviteData[idInvite]["sentAt"],
		InviteDocs: inviteDocs,
	}

	// Construir respuesta final
	inviteInfo := &models.InviteInfoResp{
		Folder: models.FolderInvite{
			IdFolder:     idFolder,
			IsSecuential: folderData[idFolder]["secuentialSign"] == "1",
			UserEmisor:   userEmisor,
			Invites:      []models.Invite{invite},
		},
	}

	return inviteInfo, nil
}

// REAHACER FUNCION PARA QUE NO DEPENDA DE LOS teams
func InviteNewUser(idUserDest, emailDest, roleApp, idUser string) {
	whereMap := map[string][]string{
		"nameApp": {"emailServ"},
	}

	dataApp, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, whereMap)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	userData := objects.UserInstJoin(idUser) // info del usuario que invita
	var guestData map[string]string
	binDoc, err := os.ReadFile("./templates/invite_user.html")
	if err != nil {
		fmt.Printf("%s", err)
	}

	if idUserDest == "" {
		idUserDest = idUser
		guestData = map[string]string{
			"nameUser":      emailDest,
			"aliasUser":     emailDest,
			"lastNameUser":  "",
			"legalNameInst": userData["legalNameInst"],
			"aliasNameInst": userData["aliasNameInst"],
			"emailUser":     emailDest,
		}
	} else {
		guestData = objects.UserInstJoin(idUserDest)
	}

	idInviteUser := uuid.NewString()
	hostFullName := fmt.Sprintf("%s %s", userData["nameUser"], userData["lastNameUser"])
	//hostTeamName := userData["nameTeam"]
	body := string(binDoc)
	body = strings.ReplaceAll(body, "{GUEST_ALIAS}", guestData["aliasUser"])
	body = strings.ReplaceAll(body, "{GUEST_EMAIL}", guestData["emailUser"])
	body = strings.ReplaceAll(body, "{HOST_INSTNAME}", userData["legalNameInst"])
	body = strings.ReplaceAll(body, "{HOST_INSTALIAS}", userData["aliasNameInst"])
	body = strings.ReplaceAll(body, "{HOST_FULLNAME}", hostFullName)
	body = strings.ReplaceAll(body, "{HOST_EMAIL}", userData["emailUser"])
	body = strings.ReplaceAll(body, "{EXPIRATION_TIME}", time.Now().Add(5*24*time.Hour).Format("2006-01-02 15:04:05"))
	body = strings.ReplaceAll(body, "{URL_ACCEPT}", fmt.Sprintf("http://%s:%s/viewinviteuser?id=%s", os.Getenv("API_IP"), os.Getenv("API_PORT"), idInviteUser))

	payload := models.EmailRequest{
		IdUser:   idUserDest,
		Subject:  fmt.Sprintf("¡Tienes una invitación! %s quiere que te unas a SIGNFORCE.", guestData["nameUser"]),
		Body:     body,
		Dest:     []string{emailDest},
		MimeType: "html",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}
	var host, port string
	for _, v := range dataApp {
		host = v["domainApp"]
		port = v["portApp"]
		break
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/mailserv", host, port), bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	texp := time.Now().Add(time.Hour * 24 * 7).Format("2006-01-02 15:04:05")
	cols := []string{"idUserInvite", "idUser", "emailDest", "expirationDate", "roleApp"}
	_, err = db.DB_con.GenericInsert("userinvites", cols, []interface{}{idInviteUser, idUser, emailDest, texp, roleApp})
	if err != nil {
		fmt.Printf("%s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
	// Ejecutar petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("TODO OK")
	}
	resp.Body.Close()
}
