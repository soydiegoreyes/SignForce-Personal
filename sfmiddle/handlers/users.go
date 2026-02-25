package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sfmiddle/auth"
	"sfmiddle/configs"
	"sfmiddle/db"
	"strings"

	"sfmiddle/models"
	"sfmiddle/objects"
	"sfmiddle/utilities"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// ==========================================================================================================
// handler para crear un usuario nuevo
func CreateUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var reqUser models.UserDataReq
	err := json.NewDecoder(request.Body).Decode(&reqUser)
	if err != nil {
		fmt.Println(err)
	}
	var registerResp models.RegisterResponse
	registerResp = objects.CreateUser(reqUser)

	// se obtiene el estatus de la institución
	instData, err := db.DB_con.GenericSelect("institutions", "idInstitution", []string{"statusInst_fk", "activeInst"}, map[string][]string{"idInstitution": {registerResp.InstId}})

	// Generar JWT
	if instData[registerResp.InstId]["activeInst"] == "1" {
		token, err := auth.GenerateJWT(registerResp.IdUser, registerResp.Role, registerResp.InstId, instData[registerResp.InstId]["statusInst_fk"]) // role root y statusinst
		if err != nil {
			http.Error(respWriter, "Error generando token", http.StatusInternalServerError)
			return
		}

		// Setear cookie con el token
		http.SetCookie(respWriter, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true, // poner en true en producción con HTTPS
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})
	}

	registerResp.IdUser = ""
	registerResp.Role = ""
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(respWriter).Encode(registerResp); err != nil {
		http.Error(respWriter, "Error al codificar JSON"+err.Error(), http.StatusInternalServerError)
		return
	}
}

// funcion que devuelve datos publicos de usuarios
func CheckUserStatus(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
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
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// se busca una lista de parametros dentro de una columna por ejemplo ["Juan", "Pedro", "Luis"], type: "nameUser" o "aliasUser"
	// otro ejemplo ["2", "43", "31"], type: "idUser" o ["correo@dominio.com", "mi_email@yahoo.com"] type: "emailUser" o ["fulanito"] type
	type usersSearch struct {
		Params    []string `json:"params,omitempty"`
		ParamType string   `json:"type,omitempty"`
		Regex     bool     `json:"likeop,omitempty"`
	}
	var US usersSearch
	err = json.NewDecoder(request.Body).Decode(&US)
	if err != nil {
		fmt.Println("Advertencia: no se pudo obtener el body")
		return
	}
	// parametros de busqueda
	wheres := map[string][]string{
		"idInstitution_fk": {idInst},
	}
	if len(US.Params) > 0 {
		wheres[US.ParamType] = US.Params
		if US.Regex {
			wheres["LOGIC"] = []string{"idInstitution_fk AND LIKE " + US.ParamType}
		}
	}

	attrs := []string{"idUser", "nameUser", "lastNameUser", "aliasUser", "emailUser", "phoneUser", "countryPhoneCode", "activeUser", "roleAppUser_fk", "createdAtUser", "lastModifiedUser", "deletedAtUser", "kycUser_fk"}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de usuario", http.StatusInternalServerError)
		return
	}

	var ids string
	USData := make(map[string]*models.UserDataResp)
	for idUS, us := range userData {
		ids = ids + idUS + ","
		userD := &models.UserDataResp{
			Id:           us["idUser"],
			Name:         us["nameUser"],
			LastName:     us["lastNameUser"],
			Alias:        us["aliasUser"],
			Email:        us["emailUser"],
			Phone:        us["countryPhoneCode"] + us["phoneUser"],
			Active:       us["activeUser"],
			Role:         us["roleAppUser_fk"],
			Kyc:          us["kycUser_fk"],
			CreatedAt:    us["createdAtUser"],
			DeletedAt:    us["deletedAtUser"],
			LastModified: us["lastModifiedUser"],
		}
		USData[idUS] = userD
	}
	query := "SELECT idUser_fk, COUNT(*) AS total_signs, SUM(CASE WHEN signatureValueSign IS NOT NULL THEN 1 ELSE 0 END) AS finished FROM signatures WHERE idUser_fk IN (%s) GROUP BY idUser_fk;"
	if len(ids) > 0 {
		ids = ids[:len(ids)-1]
	}
	signsData, _ := db.DB_con.ExecuteSelect(fmt.Sprintf(query, ids))
	for idus, sData := range signsData {
		USData[idus].TotalSigns, _ = strconv.Atoi(sData["total_signs"])
		USData[idus].SignedDocs, _ = strconv.Atoi(sData["finished"])
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(USData)
}

func UpdateUserStatus(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
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
	authInst, ok3 := claims["authInst"].(string)
	roleapp, ok4 := claims["role"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	if !(authInst == "2" || authInst == "6" || authInst == "7") {
		http.Error(respWriter, "Actualice su pago para acceder", http.StatusPaymentRequired)
		return
	}

	userInstData := objects.UserInstJoin(idUser)
	fmt.Println(userInstData)
	if userInstData["idInstitution"] != idInst {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// se reutiliza el UserDataReq para aprovechar los mismos campos que en create User
	var reqUpdate models.UserDataReq
	json.NewDecoder(request.Body).Decode(&reqUpdate)
	reqUpdate.Role = roleapp
	if authInst == "2" {
		fmt.Println("Usuario está en validacion: ", idUser)
		reqUpdate.OldPass = strings.ReplaceAll(reqUpdate.IdInvite, "-", "")
	}
	updateResp, err := objects.UpdateUser(reqUpdate, idUser)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if !updateResp {
		http.Error(respWriter, "Error al actualizar datos del usuario", http.StatusInternalServerError)
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
}

// rehacer logica para que no dependa de teams
func InviteUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
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

	// Extraer datos del JWT
	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	// obtenemos datos del usuario que está invitando
	userInst := objects.UserInstJoin(idUser)
	if userInst["activeUser"] != "1" || userInst["activeInst"] != "1" {
		http.Error(respWriter, "Usuario o institucion no estan activas.", http.StatusInternalServerError)
		return
	}
	var invite models.InviteUser
	if err := json.NewDecoder(request.Body).Decode(&invite); err != nil {
		http.Error(respWriter, "Error decodificando json", http.StatusInternalServerError)
		return
	}

	// buscamos por email o id de usuario destinatario (solo un email o id por institución)
	// alguno de los 2 campos (email o id) debe venir vacío
	wheres := map[string][]string{
		"emailUser":        {invite.EmailDest},
		"idUser":           {invite.IdUserDest},
		"idInstitution_fk": {idInst},
		"LOGIC":            {"idInstitution_fk AND emailUser OR idUser"},
	}
	invUser, err := db.DB_con.GenericSelect("users", "idUser", []string{"nameUser", "emailUser", "activeUser"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error obteniendo datos de usuario para invitacion", http.StatusInternalServerError)
		return
	}
	// si existe el usuario entonces se le invita para que actualice sus datos
	if len(invUser) == 1 {
		// se rellenan los campos del destinatario
		for id, userData := range invUser {
			if userData["activeUser"] != "1" {
				http.Error(respWriter, "Invitado no autorizado", http.StatusUnauthorized)
				return
			}
			if invite.IdUserDest == "" {
				invite.IdUserDest = id
			} else {
				invite.EmailDest = userData["emailUser"]
			}
			break
		}
		// si hay varias veces un correo en la misma institucion O se busco por email y id al mismo tiempo puede dar resultados multiples
	} else if len(invUser) > 1 {
		http.Error(respWriter, "Error resultado de usuario multiple", http.StatusNotAcceptable)
		return
	} else {
		fmt.Println("Email sin registro previo bajo mismo cliente")
	}
	idInvite := uuid.NewString()
	err = objects.InviteNewUser(invite.IdUserDest, invite.EmailDest, invite.RoleApp, idUser, idInvite)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Invitación creada: ", idInvite)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
}

func GetInviteUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	qParams := request.URL.Query()
	idInvite := qParams.Get("id")
	if idInvite == "" {
		http.Error(respWriter, "No tiene id Invitacion", http.StatusBadRequest)
		return
	}

	attrs := []string{"emailDest", "idUser", "createdAt", "acceptedAt", "expirationDate", "roleApp"}
	wheres := map[string][]string{
		"idUserInvite": {idInvite},
	}
	invData, err := db.DB_con.GenericSelect("userinvites", "idUserInvite", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener invitación", http.StatusInternalServerError)
		return
	}

	// obtenemos datos del usuario y de la institucion que envía la invitación
	ownnerData := objects.UserInstJoin(invData[idInvite]["idUser"])

	// buscamos por email de usuario destinatario (solo un email por institución)
	wheres = map[string][]string{
		"emailUser":        {invData[idInvite]["emailDest"]},
		"idInstitution_fk": {ownnerData["idInstitution"]},
	}
	invUser, err := db.DB_con.GenericSelect("users", "idUser", []string{"nameUser", "emailUser", "activeUser", "roleAppUser_fk"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error obteniendo datos de usuario para invitacion", http.StatusInternalServerError)
		return
	}
	// el usuario existe
	if len(invUser) == 1 {
		for idUserInv := range invUser {
			// si es el mismo usuario root de la institucion el que envía la invitación entonces está en el proceso de registro
			if ownnerData["rootUser_fk"] == invData[idInvite]["idUser"] && idUserInv == ownnerData["rootUser_fk"] {
				ownnerData["exists"] = "1" // es primera vez que edita sus datos
			} else {
				ownnerData["exists"] = "-1" // no es primera vez que edita sus datos
			}
			break
		}
		// si el usuario existe y es único se le da un token para que pueda actualizar sus datos en updateUserStatus
		token, err := auth.GenerateJWT(invData[idInvite]["idUser"], invData[idInvite]["roleAppUser_fk"], ownnerData["idInstitution"], ownnerData["statusInst_fk"]) // role root y statusinst
		if err != nil {
			http.Error(respWriter, "Error generando token", http.StatusInternalServerError)
			return
		}

		// Setear cookie con el token
		http.SetCookie(respWriter, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true, // poner en true en producción con HTTPS
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(1 * time.Hour),
		})

	} else if len(invUser) > 1 {
		http.Error(respWriter, "Error Usuario cuenta con mas de un correo registrado bajo el mismo cliente", http.StatusNotAcceptable)
		return
	} else {
		ownnerData["exists"] = "0" // nunca ha introducido sus datos
		fmt.Println("Email sin registro previo bajo mismo cliente")
	}

	// si hoy es despues del tiempo de expiración no se deja entrar
	if exptime, _ := time.Parse("2006-01-02T15:04:05Z", invData[idInvite]["expirationDate"]); time.Now().After(exptime) {
		http.Error(respWriter, "Error invitación expirada", http.StatusForbidden)
		return
	}
	if invData[idInvite]["acceptedAt"] != "" {
		// si ahora es antes que el tiempo de aceptado algo anda mal
		if acceptedAt, _ := time.Parse("2006-01-02T15:04:05Z", invData[idInvite]["acceptedAt"]); time.Now().Before(acceptedAt) {
			http.Error(respWriter, "Error invitación no coincide con fecha", http.StatusForbidden)
			return
		}
		http.Error(respWriter, "La invitación ya fue aceptada", http.StatusForbidden)
		return
	}

	// datos del usuario que envió la envió la invitación
	invData["host"] = ownnerData

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(respWriter).Encode(invData); err != nil {
		http.Error(respWriter, "Error al codificar JSON"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func ValidateFace(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
	}
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
	_, ok3 := claims["authInst"].(string)
	_, ok4 := claims["role"].(string)
	jti, ok5 := claims["jti"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	var attrs = []string{"kycUser_fk", "isAliveUser", "activeUser"}
	var wheres = map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	//=====================================================================
	whereMap := map[string][]string{
		"idUser":   {idUser},
		"selected": {"1"},
	}
	fvData, err := db.DB_con.GenericSelect("faceembeddings", "idKycUser", []string{"embedding"}, whereMap)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	if len(fvData) == 0 {
		fmt.Printf("Usuario no tiene biometría registrada")
		return
	}

	b64emb := fvData[userData[idUser]["kycUser_fk"]]["embedding"]
	var storedVector = utilities.B642ArrFloat(b64emb)
	embhash, err := utilities.GetHash([]byte(b64emb), configs.HashConf, false)
	if err != nil {
		fmt.Println("Error decodificando embedding")
	}
	// 2. Decodificar el vector que viene del Frontend
	var inputData struct {
		FaceVector []float32 `json:"facevector"`
	}
	if err := json.NewDecoder(request.Body).Decode(&inputData); err != nil {
		http.Error(respWriter, "Request inválido", http.StatusBadRequest)
		return
	}

	// 3. Comparar
	distancia := utilities.EuclideanDistance(storedVector, inputData.FaceVector)

	// 4. Umbral (Threshold)
	// En face-api.js / tensorflow, un umbral de 0.6 suele ser el estándar
	samePerson := distancia < 0.6
	fmt.Printf("Distancia calculada: %f - Match: %v\n", distancia, samePerson)
	if !samePerson {
		http.Error(respWriter, "Biometría inválida", http.StatusUnauthorized)
		return
	}
	biotoken, err := auth.GenerateJWTBio(embhash, jti)
	if err != nil {
		http.Error(respWriter, "Error generando biotoken", http.StatusInternalServerError)
		return
	}
	// Enviar respuesta
	// Setear cookie con el token
	http.SetCookie(respWriter, &http.Cookie{
		Name:     "biotoken",
		Value:    biotoken,
		Path:     "/signDocument",
		HttpOnly: true,
		Secure:   true, // poner en true en producción con HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(5 * time.Minute),
	})
	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(map[string]interface{}{
		"match":    samePerson,
		"distance": distancia,
	})
}
