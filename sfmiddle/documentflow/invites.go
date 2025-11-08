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
	"sfmiddle/utilities"
	"strings"
)

type InviteMail struct {
	IdUser        string
	DocumentName  string
	ReviewerName  string
	ReviewerEmail string
	ReviewerInst  string
	SentDate      string
	DeadlineDate  string
	DocumentClass string
	Description   string
	SenderMessage string
	UrlSignLink   string
}

func GenerateInvites(invites []InviteMail) error {
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
		body = strings.ReplaceAll(body, "{DOCUMENT_NAME}", invite.DocumentName)
		body = strings.ReplaceAll(body, "{SENDER_NAME}", invite.ReviewerName)
		body = strings.ReplaceAll(body, "{SENDER_EMAIL}", invite.ReviewerEmail)
		body = strings.ReplaceAll(body, "{SENDER_INST}", invite.ReviewerInst)
		body = strings.ReplaceAll(body, "{SENT_DATE}", invite.SentDate)
		body = strings.ReplaceAll(body, "{SENDER_INST}", invite.ReviewerInst)
		body = strings.ReplaceAll(body, "{DEADLINE_DATE}", invite.DeadlineDate)
		body = strings.ReplaceAll(body, "{DOCUMENT_CLASS}", invite.DocumentClass)
		body = strings.ReplaceAll(body, "{DOCUMENT_DESCRIPTION}", invite.Description)
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

type InviteInfoResp struct {
	Folder Folder `json:"folder"`
}

type Folder struct {
	IdFolder       string       `json:"idFolder"`
	IsSecuential   bool         `json:"isSecuential"`
	ExpirationDate string       `json:"expirationDate"`
	UserEmisor     UserInfo     `json:"userEmisor"`
	UserDest       UserDestInfo `json:"userDest"`
	Invites        []Invite     `json:"invites"`
}

type UserInfo struct {
	NameInstEmisor string `json:"nameInstEmisor"`
	NameUserEmisor string `json:"nameUserEmisor"`
	NameTeamEmisor string `json:"nameTeamEmisor"`
}

type UserDestInfo struct {
	NameInstDest string               `json:"nameInstDest"`
	NameUserDest string               `json:"nameUserDest"`
	NameTeamDest string               `json:"nameTeamDest"`
	Keys         []*models.KeysStatus `json:"keys"`
}

type Invite struct {
	IdInvite      string   `json:"idInvite"`
	AliveProof    bool     `json:"aliveProof"`
	IsSigner      bool     `json:"isSigner"`
	SentAt        string   `json:"sentAt"`
	MessageEmisor string   `json:"messageEmisor"`
	Document      Document `json:"document"`
}

type Document struct {
	IdDocument       string `json:"idDocument"`
	DocumentFullName string `json:"documentFullName"`
	AbstractDoc      string `json:"abstractDoc"`
}

// Request structure
type SignByInviteRequest struct {
	IdInvite string `json:"idInvite"`
}

func LoadInviteInfo(idInvite string) (*InviteInfoResp, error) {
	attrs := []string{"idInstitutionDest_fk", "idUserDest_fk", "idTeamDest_fk", "requireAliveProof",
		"expirationDate", "signerViewer", "sentAt", "descriptionText"}

	inviteData, err := db.DB_con.GenericJoinSelect("invites", "folderdocuments",
		"folderdocuments.idfolderdocument = invites.idFolderDocument", "idInvite",
		[]string{idInvite}, attrs, []string{"idDocument", "idFolder"},
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos de invitación: %v", err)
	}

	if len(inviteData) == 0 {
		return nil, fmt.Errorf("invitación no encontrada")
	}

	// Obtener datos del folder
	attrs = []string{"secuentialSign", "expirationDate"}
	wheres := map[string][]string{
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

	// Obtener datos del documento
	attrs = []string{"idDocument", "documentPath", "documentName", "documentExt", "documentHash",
		"authUseStatus", "authRoleStatus", "activeDoc", "abstractDoc"}
	wheres = map[string][]string{
		"idDocument": {inviteData[idInvite]["idDocument"]},
	}

	docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del documento: %v", err)
	}

	if len(docData) == 0 {
		return nil, fmt.Errorf("documento no encontrado")
	}

	// Validar documento
	for _, doc := range docData {
		if doc["activeDoc"] != "1" {
			return nil, fmt.Errorf("documento ya no se encuentra activo")
		}

		if doc["authRoleStatus"] != "1" {
			return nil, fmt.Errorf("no autorizado para acceder a este documento")
		}

		if doc["authUseStatus"] != "1" {
			return nil, fmt.Errorf("documento sin autorización de uso")
		}

		filePath := fmt.Sprintf("./%s%s.%s", doc["documentPath"], doc["documentName"], doc["documentExt"])

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
	}

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
	userEmisor := UserInfo{
		NameInstEmisor: "Nombre Institución Emisor",
		NameUserEmisor: "Nombre Usuario Emisor",
		NameTeamEmisor: "Nombre Team Emisor",
	}

	// TODO: Obtener datos reales de institución y equipo del destinatario
	userDest := UserDestInfo{
		NameInstDest: "Nombre Institución Destino",
		NameUserDest: userData[inviteData[idInvite]["idUserDest_fk"]]["nameUser"],
		NameTeamDest: "Nombre Team Destino",
		Keys:         keyStatResp,
	}

	// Crear invite con documento
	idDoc := inviteData[idInvite]["idDocument"]
	doc := docData[idDoc]
	documentFullName := fmt.Sprintf("%s.%s", doc["documentName"], doc["documentExt"])

	invite := Invite{
		IdInvite:      idInvite,
		AliveProof:    inviteData[idInvite]["requireAliveProof"] == "1",
		IsSigner:      inviteData[idInvite]["signerViewer"] == "1",
		SentAt:        inviteData[idInvite]["sentAt"],
		MessageEmisor: inviteData[idInvite]["descriptionText"],
		Document: Document{
			IdDocument:       idDoc,
			DocumentFullName: documentFullName,
			AbstractDoc:      doc["abstractDoc"],
		},
	}

	// Construir respuesta completa
	inviteInfo := &InviteInfoResp{
		Folder: Folder{
			IdFolder:       idFolder,
			IsSecuential:   folderData[idFolder]["secuentialSign"] == "1",
			ExpirationDate: folderData[idFolder]["expirationDate"],
			UserEmisor:     userEmisor,
			UserDest:       userDest,
			Invites:        []Invite{invite},
		},
	}

	return inviteInfo, nil
}

// Request structure
type SignByFolderRequest struct {
	IdFolder string `json:"idFolder"`
}

func LoadInviteInfoByFolder(idFolder string) (*InviteInfoResp, error) {
	// Obtener datos del folder
	attrs := []string{"secuentialSign", "expirationDate", "idUserEmisor_fk", "idInstitutionEmisor_fk", "idTeamEmisor_fk"}
	wheres := map[string][]string{
		"idFolder": {idFolder},
	}

	folderData, err := db.DB_con.GenericSelect("folders", "idFolder", attrs, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo datos del folder: %v", err)
	}

	if len(folderData) == 0 {
		return nil, fmt.Errorf("folder no encontrado")
	}

	folderInfo := folderData[idFolder]

	// Obtener todas las invitaciones del folder
	attrs = []string{"idInvite", "idInstitutionDest_fk", "idUserDest_fk", "idTeamDest_fk", "requireAliveProof",
		"signerViewer", "sentAt", "descriptionText", "idDocument"}
	
	invitesWheres := map[string][]string{
		"idFolder": {idFolder},
	}

	invitesData, err := db.DB_con.GenericSelect("invites", "idInvite", attrs, invitesWheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo invitaciones del folder: %v", err)
	}

	if len(invitesData) == 0 {
		return nil, fmt.Errorf("no se encontraron invitaciones para este folder")
	}

	// Obtener IDs de documentos únicos
	documentIDs := make([]string, 0)
	for _, invite := range invitesData {
		if docID, exists := invite["idDocument"]; exists && docID != "" {
			documentIDs = append(documentIDs, docID)
		}
	}

	// Obtener datos de todos los documentos
	var allDocData map[string]map[string]string
	if len(documentIDs) > 0 {
		attrs = []string{"idDocument", "documentPath", "documentName", "documentExt", "documentHash",
			"authUseStatus", "authRoleStatus", "activeDoc", "abstractDoc"}
		wheres = map[string][]string{
			"idDocument": documentIDs,
		}

		allDocData, err = db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo datos de los documentos: %v", err)
		}
	}

	// Validar documentos y construir mapa de documentos válidos
	validDocuments := make(map[string]map[string]string)
	for idDoc, doc := range allDocData {
		if doc["activeDoc"] != "1" {
			continue // Saltar documentos no activos
		}

		if doc["authRoleStatus"] != "1" {
			continue // Saltar documentos sin autorización
		}

		if doc["authUseStatus"] != "1" {
			continue // Saltar documentos sin autorización de uso
		}

		filePath := fmt.Sprintf("./%s%s.%s", doc["documentPath"], doc["documentName"], doc["documentExt"])

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // Saltar archivos no encontrados
		}

		hash, err := utilities.GetHash(filePath, configs.HashConf)
		if err != nil {
			continue // Saltar documentos con error de hash
		}

		if doc["documentHash"] != hash {
			continue // Saltar documentos con hash no correspondiente
		}

		validDocuments[idDoc] = doc
	}

	// Obtener datos del usuario destinatario (tomamos el primero como referencia)
	// En un escenario real, podrías necesitar manejar múltiples destinatarios
	for k := range invitesData {
		break
	}
	firstInvite := invitesData[k]
	wheres = map[string][]string{
		"idUser":           {firstInvite["idUserDest_fk"]},
		"idInstitution_fk": {firstInvite["idInstitutionDest_fk"]},
		"idTeam_fk":        {firstInvite["idTeamDest_fk"]},
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
		"idUser_fk": {firstInvite["idUserDest_fk"]},
	}
	keys, err := db.DB_con.GenericSelect("userkeys", "idUserKeys",
		[]string{"keyFilePath", "certFilePath", "notValidAfter", "subjectRFC4514",
			"subjectUniqueId", "createdAtKey"}, wheres)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo llaves del usuario: %v", err)
	}

	keyStatResp := make([]*models.KeysStatus, 0)
	for idKey, key := range keys {
		selected := userData[firstInvite["idUserDest_fk"]]["idKeysUser_fk"] == idKey

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

	// TODO: Obtener datos reales del emisor desde la BD usando folderInfo["idUserEmisor_fk"]
	userEmisor := UserInfo{
		NameInstEmisor: "Nombre Institución Emisor",
		NameUserEmisor: "Nombre Usuario Emisor", 
		NameTeamEmisor: "Nombre Team Emisor",
	}

	// TODO: Obtener datos reales de institución y equipo del destinatario
	userDest := UserDestInfo{
		NameInstDest: "Nombre Institución Destino",
		NameUserDest: userData[firstInvite["idUserDest_fk"]]["nameUser"],
		NameTeamDest: "Nombre Team Destino",
		Keys:         keyStatResp,
	}

	// Construir array de invites
	invites := make([]Invite, 0)
	for inviteID, invite := range invitesData {
		docID := invite["idDocument"]
		
		// Solo incluir invites con documentos válidos
		if doc, exists := validDocuments[docID]; exists {
			documentFullName := fmt.Sprintf("%s.%s", doc["documentName"], doc["documentExt"])

			invite := Invite{
				IdInvite:      inviteID,
				AliveProof:    invite["requireAliveProof"] == "1",
				IsSigner:      invite["signerViewer"] == "1",
				SentAt:        invite["sentAt"],
				MessageEmisor: invite["descriptionText"],
				Document: Document{
					IdDocument:       docID,
					DocumentFullName: documentFullName,
					AbstractDoc:      doc["abstractDoc"],
				},
			}
			invites = append(invites, invite)
		}
	}

	if len(invites) == 0 {
		return nil, fmt.Errorf("no se encontraron documentos válidos para firmar en este folder")
	}

	// Construir respuesta completa
	inviteInfo := &InviteInfoResp{
		Folder: Folder{
			IdFolder:       idFolder,
			IsSecuential:   folderInfo["secuentialSign"] == "1",
			ExpirationDate: folderInfo["expirationDate"],
			UserEmisor:     userEmisor,
			UserDest:       userDest,
			Invites:        invites,
		},
	}

	return inviteInfo, nil
}
