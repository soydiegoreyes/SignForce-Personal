package documentflow

import (
	"fmt"
	"os"
	"sfmiddle/coms"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"

	"sfmiddle/utilities"
	"strings"
)

func GenerateInvites(invites []models.InviteMail) error {
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
		err = coms.EmailCli.SendMail(&payload)
		if err != nil {
			fmt.Println("No se envió el email de la invitación: ", invite.UrlSignLink)
		}
	}
	return nil
}

// funcion de invitación a firma de documentos
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

	// Obtener los documentos de la invitación específica
	attrs = []string{"digestValueSign", "signatureValueSign"}
	wheres = map[string][]string{
		"idInvite_fk": {idInvite},
		"idUser_fk":   {inviteData[idInvite]["idUserDest_fk"]},
	}

	signsData, err := db.DB_con.GenericSelect("signatures", "idSignature", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo documentos de la invitación: %v", err)
	}

	if len(signsData) == 0 {
		return nil, fmt.Errorf("no se encontraron documentos en la invitación")
	}

	var docHashes []string
	idSigns := make(map[string]string)

	for idSgin, dD := range signsData {
		docHashes = append(docHashes, "'"+dD["digestValueSign"]+"'")
		if dD["signatureValueSign"] != "" {
			idSigns[dD["digestValueSign"]] = idSgin
		}
	}
	attrs = []string{"D.idDocument", "documentPath", "documentName", "documentExt", "documentHash", "authUseStatus", "authRoleStatus", "activeDoc", "abstractDoc", "createdAtDoc", "lastModifiedDoc"}
	var q = fmt.Sprintf("SELECT %s FROM documents as D inner join folderdocuments as FD ON D.idDocument=FD.idDocument WHERE D.documentHash IN (%s) AND FD.idFolder IN (%s);", strings.Join(attrs, ","), strings.Join(docHashes, ","), idFolder)
	docData, err := db.DB_con.ExecuteSelect(q)

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
	/*
		// Construir UserDestInfo del destinatario
		userDest := models.UserDestInfo{
			NameUserDest: fmt.Sprintf("%s %s",
				userData[inviteData[idInvite]["idUserDest_fk"]]["nameUser"],
				userData[inviteData[idInvite]["idUserDest_fk"]]["lastNameUser"]),
			// Nota: Necesitarías obtener nombre de institución y equipo del destinatario si es necesario
		}
	*/

	// Construir la invitación con documentos
	var inviteDocs []models.InviteDoc
	for docId, doc := range docData {
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
			ForSign:     true, // Esto debería venir de la base de datos según el rol del usuario
			ExpiresAt:   inviteData[idInvite]["expirationDate"],
			Comment:     inviteData[idInvite]["descriptionText"],
			IdSignature: idSigns[doc["documentHash"]],
		}
		inviteDocs = append(inviteDocs, inviteDoc)

	}

	// Construir la Invite completa
	invite := models.Invite{
		IdInvite: idInvite,
		UserDest: models.Reviewer{
			User:       inviteData[idInvite]["idUserDest_fk"],
			AliveProof: inviteData[idInvite]["requireAliveProof"] == "1",
		}, // Solo el ID por ahora
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
