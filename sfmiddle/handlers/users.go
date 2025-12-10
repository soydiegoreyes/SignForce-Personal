package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sfmiddle/auth"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/documentflow"
	"sfmiddle/models"
	"sfmiddle/objects"
	"sfmiddle/utilities"
	"strconv"
	"time"
	//"sfmiddle/objects"
)

// ==========================================================================================================
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
			wheres["LOGIC"] = []string{"idInstitution_fk AND REGEXP " + US.ParamType}
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
	//fmt.Println(claims)
	// Extraer datos del JWT

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	//idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	fmt.Println(idInst, idUser, authInst)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
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
	//_, ok3 := claims["team"].(string)
	//authInst, ok4 := claims["authInst"].(string)

	if !ok1 || !ok2 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
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
		var userData map[string]string
		for id, data := range invUser {
			if invite.IdUserDest == "" {
				invite.IdUserDest = id
			} else {
				invite.EmailDest = data["emailUser"]
			}
			userData = data
			break
		}

		if userData["activeUser"] != "1" {
			http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
			return
		}
	}

	documentflow.InviteNewUser(invite.IdUserDest, invite.EmailDest, invite.RoleApp, idUser)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(map[string]string{"message": "OK"})
}

func TeamUsers(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := request.Cookie("token")
	if err != nil {
		http.Error(respWriter, "No autorizado (no token)", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(respWriter, "No autorizado (claims no validas)", http.StatusUnauthorized)
		return
	}

	// Extraer datos del JWT
	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)

	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	fmt.Println(idUser, idTeam, idInst, authInst)

	var req models.GetTeamUsersRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return
	}

	var attrs []string
	if len(req.Fields) > 0 {
		attrs = append(attrs, req.Fields...)
	} else {
		attrs = []string{"nameUser", "lastNameUser", "aliasUser", "emailUser", "activeUser", "roleAppUser_fk", "idTeam_fk", "kycUser_fk"}
	}

	wheres := map[string][]string{
		"idTeam_fk": {req.IdTeam},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener datos de usuario", http.StatusBadRequest)
		return
	}
	var usersStat []models.UserDataResp
	for idUsr, usdat := range userData {
		usersStat = append(usersStat, models.UserDataResp{
			Id:       idUsr,
			Name:     usdat["nameUser"],
			LastName: usdat["lastNameUser"],
			Alias:    usdat["aliasUser"],
			Email:    usdat["emailUser"],
			Active:   usdat["activeUser"],
			Role:     usdat["roleAppUser_fk"],
			//Team:     usdat["idTeam_fk"],
			Kyc: usdat["kycUser_fk"],
		})
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&usersStat)
}

func InstTeams(respWriter http.ResponseWriter, request *http.Request) {
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
	idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	fmt.Println("JWT claims:", idUser, idTeam, idInst, authInst)

	attrs := []string{"creatorUser_fk", "nameTeam", "limitSigners", "limitUsers", "deletedAt", "description"}
	wheres := map[string][]string{"idInstitution_fk": {idInst}}

	teamData, err := db.DB_con.GenericSelect("teams", "idTeam", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener datos de usuario", http.StatusBadRequest)
		return
	}
	fmt.Printf("teamData -> %+v\n", teamData)

	var teamsStat []models.TeamDataResp
	for idTeamInst, teamdat := range teamData {
		fmt.Printf("Fila -> idTeam: %s, data: %+v\n", idTeamInst, teamdat)
		userCreator := objects.UserInstJoin(teamdat["creatorUser_fk"])
		if len(userCreator) == 0 {
			userCreator = make(map[string]string)
			userCreator["nameUser"] = "Sin dato"
		}
		fmt.Printf("Creador: %+v\n", userCreator)

		teamsStat = append(teamsStat, models.TeamDataResp{
			Id:           idTeamInst,
			CreatorUser:  userCreator["nameUser"],
			Name:         teamdat["nameTeam"],
			LimitSigners: teamdat["limitSigners"],
			LimitUsers:   teamdat["limitUsers"],
			DeletedAt:    teamdat["deletedAt"],
			Description:  teamdat["description"],
		})
	}

	fmt.Println("Total equipos:", len(teamsStat))

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(respWriter).Encode(teamsStat); err != nil {
		fmt.Println("Error al codificar JSON:", err)
	}
}

func NewTeam(respWriter http.ResponseWriter, request *http.Request) {
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
	_, ok3 := claims["team"].(string)
	//_, ok4 := claims["authInst"].(string)

	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	type newTeamReq struct {
		NameTeam    string `json:"nameTeam"`
		Description string `json:"description"`
	}
	var req newTeamReq
	err = json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(respWriter, "Error decodificando json", http.StatusBadRequest)
		return
	}

	attrs1 := []string{"emailUser", "activeUser", "roleAppUser_fk", "idTeam_fk", "idInstitution_fk"}
	attrs2 := []string{"statusInst_fk", "activeInst", "rootUser_fk"}

	join := "users.idInstitution_fk = institutions.idInstitution"

	attrs := []string{"creatorUser_fk", "nameTeam", "deletedAt"}
	wheres := map[string][]string{
		"idInstitution_fk": {idInst},
	}
	userInstData, err := db.DB_con.GenericJoinSelect("users", "institutions", join, "idUser", []string{idUser}, attrs1, attrs2)
	if err != nil {
		http.Error(respWriter, "Error obteniendo institucion y usuario", http.StatusInternalServerError)
		return
	}
	teamsData, err := db.DB_con.GenericSelect("teams", "idTeam", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error obteniendo datos de equipos", http.StatusInternalServerError)
		return
	}
	fmt.Println(userInstData, teamsData)
	cols := []string{"idInstitution_fk", "creatorUser_fk", "nameTeam", "description"}
	values := []interface{}{idInst, idUser, req.NameTeam, req.Description}
	idNewTeam, err := db.DB_con.GenericInsert("teams", cols, values)
	if err != nil {
		http.Error(respWriter, "Error insertando equipo", http.StatusInternalServerError)
		return
	}
	fmt.Println(idNewTeam)

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(respWriter).Encode(map[string]string{"message": "OK"}); err != nil {
		fmt.Println("Error al codificar JSON:", err)
	}
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

	attrs := []string{"emailDest", "idUser", "createdAt", "expirationDate", "roleApp"}
	wheres := map[string][]string{
		"idUserInvite": {idInvite},
	}
	invData, err := db.DB_con.GenericSelect("userinvites", "idUserInvite", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener invitación", http.StatusInternalServerError)
		return
	}
	if exptime, _ := time.Parse(invData[idInvite]["expirationDate"], "2006-01-02T15:04:05Z"); exptime.After(time.Now()) {
		http.Error(respWriter, "Error invitación expirada", http.StatusForbidden)
		return
	}

	ownnerData := objects.UserInstJoin(invData[idInvite]["idUser"])
	invData["host"] = ownnerData
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(respWriter).Encode(invData); err != nil {
		fmt.Println("Error al codificar JSON:", err)
	}
}

func CreateUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var req models.UserDataReq
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
	}
	uinv, err := db.DB_con.GenericJoinSelect("userinvites", "users", "userinvites.idUser=users.idUser", "idUserInvite", []string{req.IdInvite}, []string{"acceptedAt"}, []string{"idInstitution_fk"})
	if err != nil {
		fmt.Println("error al buscar invitacion")
		return
	}
	if uinv[req.IdInvite]["acceptedAt"] != "" {
		fmt.Println("Usuario ya acepto la invitación")
		return
	}
	tempPass := utilities.PassGenerator(12)
	fmt.Println("pasa: ", tempPass)
	hashed, err := utilities.GetHash([]byte(tempPass), configs.HashConf)
	if err != nil {
		fmt.Println(err)
		return
	}
	cols := []string{"nameUser", "lastNameUser", "taxNumUser", "pobUidUser", "aliasUser", "emailUser", "phoneUser", "appPassHash", "activeUser", "roleAppUser_fk", "idInstitution_fk"}
	values := []interface{}{req.Name, req.LastName, req.TaxNum, req.PobUid, req.Alias, req.Email, req.Phone, hashed, 1, req.Role, uinv[req.IdInvite]["idInstitution_fk"]}
	idNewUser, err := db.DB_con.GenericInsert("users", cols, values)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(idNewUser)

	if err = db.DB_con.GenericBatchUpdate("userinvites", "idUserInvite", map[string]map[string]interface{}{req.IdInvite: {"acceptedAt": time.Now().Format("2006-01-02 15:04:05")}}); err != nil {
		fmt.Println(err)
		return
	}
}
