package objects

import (
	"fmt"
	"os"
	"sfmiddle/coms"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/utilities"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct{}

func UserInstJoin(idUser string) map[string]string {
	attrs1 := []string{"idUser", "nameUser", "lastNameUser", "emailUser", "aliasUser", "activeUser"}
	attrs2 := []string{"idInstitution", "legalNameInst", "aliasNameInst", "contactEmailInst", "taxNumInst", "activeInst", "rootUser_fk"}
	userData, err := db.DB_con.GenericJoinSelect("users", "institutions", "users.idInstitution_fk = institutions.idInstitution", "idUser", []string{idUser}, attrs1, attrs2)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	return userData[idUser]
}

// PRIMER FUNCION PARA REGISTRAR UN NUEVO CLIENTE
func RegisterRootUser(registerReq *models.RegisterRequest, idInst, passHash string) (string, error) {

	idInstHash, err := utilities.GetHash([]byte(idInst), configs.HashConf)
	if err != nil {
		return "", err
	}

	columns := []string{
		"nameUser",
		"lastNameUser",
		"aliasUser",
		"emailUser",
		"phoneUser",
		"countryPhoneCode",
		"activeUser",
		"isAliveUser",
		"roleAppUser_fk",
		"appPassHash",
		"idInstitution_fk",
	}
	values := []interface{}{
		registerReq.LegalSignupName,
		registerReq.LegalSignupLastname,
		fmt.Sprintf("root_%s", idInstHash),
		registerReq.ContactEmailInst,
		registerReq.ContactPhoneInst,
		"52", // codigo de México
		"1",  // El usuario está activo pero no se considera vivo
		"0",  // si no tiene vector facial no se considera vivo
		"1",  // el role es root
		passHash,
		idInst,
	}
	idUser, err := db.DB_con.GenericInsert("users", columns, values)
	if err != nil {
		fmt.Println("Error: al insertar datos", err)
		return "", err
	}

	cols := []string{"idUser", "embedding", "selected"}
	values = []interface{}{idUser, "", "1"}
	idUserEmb, err := db.DB_con.GenericInsert("faceembeddings", cols, values)
	if err != nil {
		return "", err
	}
	updates := map[string]map[string]interface{}{
		idUser: {"kycUser_fk": idUserEmb},
	}
	err = db.DB_con.GenericBatchUpdate("users", "idUser", updates)
	if err != nil {
		return "", err
	}

	dataRole, err := db.DB_con.GenericSelect("userroles", "idUser", []string{"idRole"}, map[string][]string{"idUser": {idUser}, "idRole": {"1"}})
	if err != nil {
		return "", err
	}
	if len(dataRole) != 0 {
		return "", fmt.Errorf("colision de usuario y rol. El usuario ya existe")
	}
	_, err = db.DB_con.GenericInsert("userroles", []string{"idUser", "idRole", "grantedBy"}, []interface{}{idUser, 1, idUser})
	if err != nil {
		return "", fmt.Errorf("no se pudo asignar el rol al usuario")
	}

	return idUser, nil
}

func CreateUser(reqUser models.UserDataReq) models.RegisterResponse {
	registerResp := models.RegisterResponse{Check: false, InstId: "", Error: ""}

	uinv, err := db.DB_con.GenericJoinSelect("userinvites", "users", "userinvites.idUser=users.idUser", "idUserInvite", []string{reqUser.IdInvite}, []string{"roleApp", "acceptedAt"}, []string{"idUser", "idInstitution_fk"})
	if err != nil {
		registerResp.Error = "error al buscar invitacion"
		return registerResp
	}
	if uinv[reqUser.IdInvite]["acceptedAt"] != "" {
		registerResp.Error = "Usuario ya aceptó la invitación"
		return registerResp
	}

	// obtenemos el facevector
	facevector := utilities.ArrFloat2B64(reqUser.FaceVector)

	// dado que hay vector y la red lo identificó como valido se dice que está vivo
	var isAlive int
	if reqUser.IsAlive {
		isAlive = 1
	}

	ownnerUser := uinv[reqUser.IdInvite]["idUser"]

	tempPass := utilities.PassGenerator(12)
	hashed, err := utilities.GetHash([]byte(tempPass), configs.HashConf)
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}

	cols := []string{"nameUser", "lastNameUser", "taxNumUser", "pobUidUser", "aliasUser", "emailUser", "phoneUser", "appPassHash", "activeUser", "isAliveUser", "roleAppUser_fk", "idInstitution_fk"}
	values := []interface{}{reqUser.Name, reqUser.LastName, reqUser.TaxNum, reqUser.PobUid, reqUser.Alias, reqUser.Email, reqUser.Phone, hashed, 1, isAlive, uinv[reqUser.IdInvite]["roleApp"], uinv[reqUser.IdInvite]["idInstitution_fk"]}
	idNewUser, err := db.DB_con.GenericInsert("users", cols, values)
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}
	fmt.Println("Nuevo usuario: ", idNewUser)

	cols = []string{"idUser", "embedding", "selected"}
	values = []interface{}{idNewUser, facevector, "1"}
	idUserEmb, err := db.DB_con.GenericInsert("faceembeddings", cols, values)
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}
	updates := map[string]map[string]interface{}{
		idNewUser: {"kycUser_fk": idUserEmb},
	}
	err = db.DB_con.GenericBatchUpdate("users", "idUser", updates)
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}
	// invitacion de usuario por email
	binDoc, err := os.ReadFile("./templates/welcome_register.html")
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}

	body := string(binDoc)
	body = strings.ReplaceAll(body, "{TEMPORAL_USERNAME}", reqUser.Email)
	body = strings.ReplaceAll(body, "{TEMPORAL_PASS}", tempPass)
	body = strings.ReplaceAll(body, "{EXPIRATION_TIME}", time.Now().Add(30*24*time.Hour).Format("2006-01-02 15:04:05"))
	body = strings.ReplaceAll(body, "{URL_COMPLETAR_REGISTRO}", fmt.Sprintf("%s/login", os.Getenv("API_IP")))
	//body = strings.ReplaceAll(body, "{URL_COMPLETAR_REGISTRO}", fmt.Sprintf("http://%s:%s/login", os.Getenv("API_IP"), os.Getenv("API_PORT")))

	payload := models.EmailRequest{
		IdUser:   "1",
		Subject:  fmt.Sprintf("¡Bienvenido a Signforce! Correo de verificación %s", reqUser.Name),
		Body:     body,
		Dest:     []string{reqUser.Email},
		MimeType: "html",
	}

	err = coms.EmailCli.SendMail(&payload)
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}

	dataRole, err := db.DB_con.GenericSelect("userroles", "idUser", []string{"idRole"}, map[string][]string{"idUser": {idNewUser}, "idRole": {"1"}})
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}
	if len(dataRole) != 0 {
		registerResp.Error = "colision de usuario y rol. El usuario ya existe"
		return registerResp
	}
	_, err = db.DB_con.GenericInsert("userroles", []string{"idUser", "idRole", "grantedBy"}, []interface{}{idNewUser, 1, ownnerUser})
	if err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}
	// se actualiza que se resolvió la invitación
	if err = db.DB_con.GenericBatchUpdate("userinvites", "idUserInvite", map[string]map[string]interface{}{reqUser.IdInvite: {"acceptedAt": time.Now().Format("2006-01-02 15:04:05")}}); err != nil {
		registerResp.Error = err.Error()
		return registerResp
	}
	registerResp.Check = true
	return registerResp
}

// se actualizan datos del usuario dentro de la plataforma
func UpdateUser(userReq models.UserDataReq, idUser string) (bool, error) {

	var updates = make(map[string]map[string]interface{})
	var values = map[string]interface{}{
		"nameUser":     userReq.Name,
		"lastNameUser": userReq.LastName,
		"aliasUser":    userReq.Alias,
		"emailUser":    userReq.Email,
		"phoneUser":    userReq.Phone,
		"taxNumUser":   userReq.TaxNum,
		"pobUidUser":   userReq.PobUid,
	}
	// los valores a actualizar no deben venir vacíos y deben ser mayores a 3 caracteres
	for k, v := range values {
		val, _ := v.(string)
		if val == "" || len(val) < 3 {
			delete(values, k)
		}
	}
	// solo si el vector facial viene con datos entonces intentamos validar su cambio
	if len(userReq.FaceVector) != 0 {
		attrs := []string{"embedding", "selected"}
		wheres := map[string][]string{
			"idUser": {idUser},
			//"selected": {"1"}, no se usa selected porque solo hay 3 casos: solo hay uno (vacio o lleno) o hay mas de uno (todos llenos)
		}

		// obtenemos los embeddings del usuario
		faceEmbData, err := db.DB_con.GenericSelect("faceembeddings", "idKycUser", attrs, wheres)
		if err != nil {
			return false, err
		}
		// se procede igual que con las personas.
		// el embedding solo se puede actualizar si el usuario no tiene ninguno registrado
		// regularmente el embedding se irá actualizando para tener rasgos mas precisos del usuario
		// pero esa actualización se hara en un proceso separado

		var idKycUser string // es el id que se sustituirá por el nuevo vector (cuando no se ha registrado tiene id pero es un string vacío)

		// si solo hay un vector registrado validamos que esté vacío y esté seleccionado
		// este es el caso de un nuevo registro
		if len(faceEmbData) > 0 {
			// se compara cada uno de los embeddings para promediar si son la misma persona y por tanto actualizar el embedding
			// esta función esta a discusión y puede ser mejorada
			var avg float64 = 1 // completamente alejados
			var selected string // id del embedding seleccionado

			for idkyc, faceData := range faceEmbData {
				if faceData["selected"] == "1" {
					selected = idkyc
					if faceData["embedding"] == "" {
						if len(faceEmbData) == 1 { // caso base donde el usuario es nuevo y se debe insertar su embedding
							avg = 0 // esto hace que el id a actualizar sea el seleccionado
							break
						} else {
							continue
						}
					}
				}
				avg *= utilities.EuclideanDistance(utilities.B642ArrFloat(faceData["embedding"]), userReq.FaceVector)
			}
			// si es la misma persona
			if avg < 0.6 {
				idKycUser = selected
			}
		} else {
			return false, fmt.Errorf("no tiene id para vector facial. Este error nunca debe ocurrir. %s", "")
		}
		// se efectua el update del embedding en caso de que se haya aprobado el cambio
		facevector := utilities.ArrFloat2B64(userReq.FaceVector) // vector tomado del front
		if idKycUser != "" {
			faceupdates := map[string]map[string]interface{}{
				idKycUser: {"embedding": facevector},
			}
			err = db.DB_con.GenericBatchUpdate("faceembeddings", "idKycUser", faceupdates)
			if err != nil {
				return false, err
			}
			values["isAliveUser"] = 1
		}
	}
	// se inserta hasta el final para validar que el usuario está vivo en caso de que facevector sea correcto
	updates[idUser] = values
	err := db.DB_con.GenericBatchUpdate("users", "idUser", updates)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Envía una invitación a un correo para que se una a signforce
func InviteNewUser(idUserDest, emailDest, roleApp, idUser string) (string, error) {

	userData := UserInstJoin(idUser) // info del usuario que invita
	var guestData map[string]string
	binDoc, err := os.ReadFile("./templates/invite_user.html")
	if err != nil {
		fmt.Printf("%s", err)
		return "", err
	}

	if idUserDest == "" {
		idUserDest = idUser // se le asigna el mismo usuario que el anfitrion solo para fines de envio de email
		guestData = map[string]string{
			"nameUser":      emailDest,
			"aliasUser":     emailDest,
			"lastNameUser":  "",
			"legalNameInst": userData["legalNameInst"],
			"aliasNameInst": userData["aliasNameInst"],
			"emailUser":     emailDest,
		}
	} else {
		guestData = UserInstJoin(idUserDest)
	}

	idInviteUser := uuid.NewString()
	hostFullName := fmt.Sprintf("%s %s", userData["nameUser"], userData["lastNameUser"])

	// fecha de expiración de la invitación
	tExp := time.Now().Add(5 * 24 * time.Hour).Format("2006-01-02 15:04:05")

	body := string(binDoc)
	body = strings.ReplaceAll(body, "{GUEST_ALIAS}", guestData["aliasUser"])
	body = strings.ReplaceAll(body, "{GUEST_EMAIL}", guestData["emailUser"])
	body = strings.ReplaceAll(body, "{HOST_INSTNAME}", userData["legalNameInst"])
	body = strings.ReplaceAll(body, "{HOST_INSTALIAS}", userData["aliasNameInst"])
	body = strings.ReplaceAll(body, "{HOST_FULLNAME}", hostFullName)
	body = strings.ReplaceAll(body, "{HOST_EMAIL}", userData["emailUser"])
	body = strings.ReplaceAll(body, "{EXPIRATION_TIME}", tExp)
	body = strings.ReplaceAll(body, "{URL_ACCEPT}", fmt.Sprintf("%s/edituser?id=%s", os.Getenv("API_IP"), idInviteUser))
	//body = strings.ReplaceAll(body, "{URL_ACCEPT}", fmt.Sprintf("http://%s:%s/edituser?id=%s", os.Getenv("API_IP"), os.Getenv("API_PORT"), idInviteUser))

	payload := models.EmailRequest{
		IdUser:   idUserDest,
		Subject:  fmt.Sprintf("¡Tienes una invitación! %s quiere que te unas a SIGNFORCE.", guestData["nameUser"]),
		Body:     body,
		Dest:     []string{emailDest},
		MimeType: "html",
	}

	err = coms.EmailCli.SendMail(&payload)
	if err != nil {
		fmt.Println("No se envió el email de la invitación: ", idInviteUser, err)
		return "", err
	} else {
		cols := []string{"idUserInvite", "idUser", "emailDest", "expirationDate", "roleApp"}
		_, err = db.DB_con.GenericInsert("userinvites", cols, []interface{}{idInviteUser, idUser, emailDest, tExp, roleApp})
		if err != nil {
			fmt.Printf("%s", err)
			return "", err
		}
	}
	return idInviteUser, nil
}
