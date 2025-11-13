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

		if strings.Contains(resp.Status, "200 OK") {
			fmt.Println("TODO OK")
		}
		resp.Body.Close()

		//fmt.Println(inviteResp, dataInvites)
	}
	return nil
}

func LoadInviteInfo(idInvite string) (*models.InviteInfoResp, error) {
	// obtenemos datos de la invitación (solo una por usuario)
	attrs := []string{"idFolder", "idInstitutionDest_fk", "idUserDest_fk", "idTeamDest_fk", "requireAliveProof",
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
	attrs = []string{"secuentialSign", "expirationDate"}
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
	wheres = map[string][]string{
		"idUser":           {inviteData[idInvite]["idUserDest_fk"]},
		"idInstitution_fk": {inviteData[idInvite]["idInstitutionDest_fk"]},
		"idTeam_fk":        {inviteData[idInvite]["idTeamDest_fk"]},
	}

	userData, err := db.DB_con.GenericSelect("users", "idUser",
		[]string{"activeUser", "idKeysUser_fk", "nameUser"}, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del usuario: %v", err)
	}

	if len(userData) == 0 {
		return nil, fmt.Errorf("usuario destinatario no encontrado")
	}

	// Obtener llaves RSA del usuario
	wheres = map[string][]string{
		"idUser_fk": {inviteData[idInvite]["idUserDest_fk"]},
	}
	keys, err := db.DB_con.GenericSelect("userkeys", "idUserKeys",
		[]string{"keyFilePath", "certFilePath", "notValidAfter", "subjectRFC4514",
			"subjectUniqueId", "createdAtKey"}, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo llaves del usuario: %v", err)
	}

	keyStatResp := make([]*models.KeysStatus, 0)
	for idKey, key := range keys {
		selected := userData[inviteData[idInvite]["idUserDest_fk"]]["idKeysUser_fk"] == idKey

		ks := &models.KeysStatus{
			IdKey:           idKey,
			NameKey:         key["keyFilePath"][strings.LastIndex(key["keyFilePath"], "/")+1:],
			NameCer:         key["certFilePath"][strings.LastIndex(key["certFilePath"], "/")+1:],
			Expiration:      key["notValidAfter"],
			Owner:           key["subjectRFC4514"][strings.Index(key["subjectRFC4514"], "=")+1 : strings.Index(key["subjectRFC4514"], ",")],
			SubjectUniqueId: key["subjectUniqueId"],
			UploadedAt:      key["createdAtKey"],
			Selected:        selected,
		}
		keyStatResp = append(keyStatResp, ks)
	}

	// Aquí deberías obtener los datos del emisor, institución y equipo
	// (esto depende de tu esquema de base de datos)
	// Por ahora dejo valores de ejemplo - debes ajustar según tu BD

	// TODO: Obtener datos reales del emisor desde la BD
	userEmisor := models.UserInfo{
		NameInstEmisor: "Nombre Institución Emisor",
		NameUserEmisor: "Nombre Usuario Emisor",
		NameTeamEmisor: "Nombre Team Emisor",
	}

	// TODO: Obtener datos reales de institución y equipo del destinatario
	userDest := models.UserDestInfo{
		NameInstDest: "Nombre Institución Destino",
		NameUserDest: userData[inviteData[idInvite]["idUserDest_fk"]]["nameUser"],
		NameTeamDest: "Nombre Team Destino",
		Keys:         keyStatResp,
	}

	// Obtener los id de los documentos del folder
	attrs = []string{"idDocument"}
	wheres = map[string][]string{
		"idFolder": {inviteData[idInvite]["idFolder"]},
	}

	folderDocs, err := db.DB_con.GenericSelect("folderdocuments", "idfolderdocument", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del folder: %v", err)
	}

	if len(folderDocs) == 0 {
		return nil, fmt.Errorf("documentos no encontrados")
	}

	// Obtener datos de los documentos del folder
	attrs = []string{"idDocument", "documentPath", "documentName", "documentExt", "documentHash",
		"authUseStatus", "authRoleStatus", "activeDoc", "abstractDoc"}
	wheres = map[string][]string{
		"idDocument": {},
	}
	for _, d := range folderDocs {
		wheres["idDocument"] = append(wheres["idDocument"], d["idDocument"])
	}
	docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del documento: %v", err)
	}

	if len(docData) == 0 {
		return nil, fmt.Errorf("documento no encontrado")
	}

	// Construir respuesta completa
	inviteInfo := &models.InviteInfoResp{
		Folder: models.FolderInvite{
			IdFolder:       idFolder,
			IsSecuential:   folderData[idFolder]["secuentialSign"] == "1",
			ExpirationDate: folderData[idFolder]["expirationDate"],
			UserEmisor:     userEmisor,
			UserDest:       userDest,
			Invites:        []models.Invite{},
		},
	}

	// Validar documento y añadir datos de la invitacion
	for idDoc, doc := range docData {
		if doc["activeDoc"] != "1" {
			return nil, fmt.Errorf("documento ya no se encuentra activo")
		}

		if doc["authRoleStatus"] != "1" {
			return nil, fmt.Errorf("no autorizado para acceder a este documento")
		}

		if doc["authUseStatus"] != "1" {
			return nil, fmt.Errorf("documento sin autorización de uso")
		}

		documentFullName := fmt.Sprintf("%s.%s", doc["documentName"], doc["documentExt"])
		filePath := fmt.Sprintf("./%s%s", doc["documentPath"], documentFullName)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return nil, fmt.Errorf("archivo no encontrado en el sistema")
		}

		hash, err := utilities.GetHash(filePath, configs.HashConf)
		if err != nil {
			return nil, fmt.Errorf("error al obtener hash: %v", err)
		}

		// Corregido: debería ser != en lugar de ==
		if doc["documentHash"] != hash {
			return nil, fmt.Errorf("hash no corresponde al archivo")
		}

		// Crear invite con documento
		invite := models.Invite{
			IdInvite:      idInvite,
			AliveProof:    inviteData[idInvite]["requireAliveProof"] == "1",
			IsSigner:      inviteData[idInvite]["signerViewer"] == "1",
			SentAt:        inviteData[idInvite]["sentAt"],
			MessageEmisor: inviteData[idInvite]["descriptionText"],
			Document: models.Document{
				IdDocument:       idDoc,
				DocumentFullName: documentFullName,
				Abstract:         doc["abstractDoc"],
			},
		}
		inviteInfo.Folder.Invites = append(inviteInfo.Folder.Invites, invite)
	}

	return inviteInfo, nil
}

type InviteUser struct {
	IdGuest string
	Email   string
	Alias   string
	IdInst  string
	IdUser  string
	IdTeam  string
}

func InviteNewUser(newInvite InviteUser) {
	userData := objects.UserInstTeam(newInvite.IdUser)
	whereMap := map[string][]string{
		"nameApp": {"emailServ"},
	}

	dataApp, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, whereMap)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	binDoc, err := os.ReadFile("./templates/invite_team.html")
	if err != nil {
		fmt.Printf("%s", err)
	}
	hostFullName := fmt.Sprintf("%s %s", userData["nameUser"], userData["lastNameUser"])
	hostTeamName := userData["nameTeam"]
	body := string(binDoc)
	body = strings.ReplaceAll(body, "{GUEST_ALIAS}", newInvite.Alias)
	body = strings.ReplaceAll(body, "{GUEST_EMAIL}", newInvite.Email)
	body = strings.ReplaceAll(body, "{HOST_INSTNAME}", userData["legalNameInst"])
	body = strings.ReplaceAll(body, "{HOST_INSTALIAS}", userData["aliasNameInst"])
	body = strings.ReplaceAll(body, "{HOST_FULLNAME}", hostFullName)
	body = strings.ReplaceAll(body, "{GUEST_TEAMNAME}", userData["nameTeam"])
	body = strings.ReplaceAll(body, "{HOST_EMAIL}", userData["emailUser"])
	body = strings.ReplaceAll(body, "{EXPIRATION_TIME}", time.Now().Add(5*24*time.Hour).Format("2006-01-02 15:04:05"))
	body = strings.ReplaceAll(body, "{URL_ACCEPT}", fmt.Sprintf("http://%s:%s/login", os.Getenv("API_IP"), os.Getenv("API_PORT")))

	payload := models.EmailRequest{
		IdUser:   newInvite.IdUser,
		Subject:  fmt.Sprintf("¡Tienes una invitación! El equipo %s quiere que te unas a ellos.", hostTeamName),
		Body:     body,
		Dest:     []string{newInvite.Email},
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
	// Ejecutar petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	if strings.Contains(resp.Status, "200 OK") {
		fmt.Println("TODO OK")
	}
	resp.Body.Close()
}
