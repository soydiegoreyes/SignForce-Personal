package handlers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sfmiddle/auth"
	"sfmiddle/db"

	"sfmiddle/models"
	"sfmiddle/objects"
	"sfmiddle/utilities"
	"strconv"
	"time"
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
	//idTeam, ok3 := claims["team"].(string)
	//authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 {
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

	json.NewEncoder(respWriter).Encode(USData)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
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
	fmt.Println(claims)

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)

	if !(ok1 && ok2) {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
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
	reqUpdate.Role = ""
	updateResp, err := objects.UpdateUser(reqUpdate, idUser)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	if !updateResp {
		http.Error(respWriter, "Error al actualizar datos del usuario", http.StatusInternalServerError)
		return
	}

	// se actualiza que se resolvió la invitación
	if err = db.DB_con.GenericBatchUpdate("userinvites", "idUserInvite", map[string]map[string]interface{}{reqUpdate.IdInvite: {"acceptedAt": time.Now().Format("2006-01-02 15:04:05")}}); err != nil {
		http.Error(respWriter, "Error al aceptar la invitación: "+err.Error(), http.StatusInternalServerError)
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
	_, ok2 := claims["iid"].(string)
	//authInst, ok3 := claims["authInst"].(string)

	if !ok1 || !ok2 {
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

	// buscamos por email o id de usuario destinatario
	wheres := map[string][]string{
		"emailUser": {invite.EmailDest},
		"idUser":    {invite.IdUserDest},
		"LOGIC":     {"emailUser OR idUser"},
	}
	invUser, err := db.DB_con.GenericSelect("users", "idUser", []string{"nameUser", "emailUser", "activeUser"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error obteniendo datos de usuario para invitacion", http.StatusInternalServerError)
		return
	}

	if len(invUser) > 0 {
		// se rellenan los campos del destinatario
		for id, userData := range invUser {
			if invite.IdUserDest == "" {
				invite.IdUserDest = id
			} else {
				invite.EmailDest = userData["emailUser"]
			}
			if userData["activeUser"] != "1" {
				http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
				return
			}
			break
		}

	}

	idInvite, err := objects.InviteNewUser(invite.IdUserDest, invite.EmailDest, invite.RoleApp, idUser)
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

	// obtenemos datos del usuario que envía la invitación
	ownnerData := objects.UserInstJoin(invData[idInvite]["idUser"])

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

	// si es el mismo usuario root de la institucion el que envía la invitación entonces está en el proceso de registro
	if ownnerData["rootUser_fk"] == invData[idInvite]["idUser"] && ownnerData["activeInst"] == "0" {
		ownnerData["exists"] = invData[idInvite]["idUser"]
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

// Función para obtener estadísticas de uso de la plataforma para dashboard
func StatsDash(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
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
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	fmt.Println("JWT claims:", idUser, idInst, authInst)

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
	if !ok1 || !ok2 {
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
	embRaw := utilities.Decode_b64(b64emb)

	var storedVector []float64
	for i := 0; i < len(embRaw); i += 4 {
		// Leemos 4 bytes en LittleEndian (estándar de JS)
		bits := binary.LittleEndian.Uint32(embRaw[i : i+4])
		floatVal := math.Float32frombits(bits)
		storedVector = append(storedVector, float64(floatVal))
	}

	// 2. Decodificar el vector que viene del Frontend
	var inputData struct {
		FaceVector []float64 `json:"facevector"`
	}
	if err := json.NewDecoder(request.Body).Decode(&inputData); err != nil {
		http.Error(respWriter, "Request inválido", http.StatusBadRequest)
		return
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, inputData.FaceVector)
	if err != nil {
		fmt.Println(err)

	}
	// 3. Comparar
	fmt.Println(inputData.FaceVector)
	distancia := utilities.EuclideanDistance(storedVector, inputData.FaceVector)

	// 4. Umbral (Threshold)
	// En face-api.js / tensorflow, un umbral de 0.6 suele ser el estándar
	samePerson := distancia < 0.6

	fmt.Printf("Distancia calculada: %f - Match: %v\n", distancia, samePerson)

	// Enviar respuesta
	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(map[string]interface{}{
		"match":    samePerson,
		"distance": distancia,
	})

	//=====================================================================

}
