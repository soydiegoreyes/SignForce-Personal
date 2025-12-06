package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sfmiddle/auth"
	"sfmiddle/db"
	"sfmiddle/documentflow"
	"sfmiddle/models"
	"sfmiddle/objects"
	//"sfmiddle/objects"
)

// ==========================================================================================================
// funcion que devuelve datos publicos del usuario
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
	//fmt.Println(claims)
	// Extraer datos del JWT

	idUser, ok1 := claims["uid"].(string)
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
		US.ParamType:       US.Params,
	}
	if US.Regex {
		wheres["LOGIC"] = []string{"idInstitution_fk AND REGEXP " + US.ParamType}
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", []string{"idUser", "nameUser", "lastNameUser", "aliasUser", "emailUser", "activeUser", "roleAppUser_fk", "kycUser_fk"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de usuario", http.StatusInternalServerError)
		return
	}

	USData := make(map[string]*models.UserDataResp)
	for idUS, us := range userData {
		if us["activeUser"] != "1" {
			http.Error(respWriter, "No autorizado. Usuario inactivo", http.StatusUnauthorized)
			return
		}

		userD := &models.UserDataResp{
			Id:       userData[idUser]["idUser"],
			Name:     userData[idUser]["nameUser"],
			LastName: userData[idUser]["lastNameUser"],
			Alias:    userData[idUser]["aliasUser"],
			Email:    userData[idUser]["emailUser"],
			Active:   userData[idUser]["activeUser"],
			Role:     userData[idUser]["roleAppUser_fk"],
			//Team:     userData[idUser]["idTeam_fk"],
			Kyc: userData[idUser]["kycUser_fk"],
		}
		USData[idUS] = userD
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

	documentflow.InviteNewUser(invite.IdUserDest, invite.EmailDest, idUser)
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

func AcceptInviteTeam(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	qParams := request.URL.Query()
	idInvite := qParams.Get("idInvite")
	if idInvite == "" {
		http.Error(respWriter, "No tiene id Invitacion", http.StatusBadRequest)
		return
	}

}
