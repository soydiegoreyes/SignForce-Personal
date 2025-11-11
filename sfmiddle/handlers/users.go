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
)

// ==========================================================================================================
// funcion que devuelve datos publicos del usuario
func CheckUserStatus(respWriter http.ResponseWriter, request *http.Request) {
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
	idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
		"idTeam_fk":        {idTeam},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", []string{"nameUser", "lastNameUser", "aliasUser", "emailUser", "activeUser", "roleAppUser_fk", "idTeam_fk", "kycUser_fk"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de usuario", http.StatusInternalServerError)
		return
	}
	if userData[idUser]["activeUser"] != "1" || authInst != "1" {
		http.Error(respWriter, "No autorizado. Usuario inactivo", http.StatusUnauthorized)
		return
	}

	userD := &models.UserDataResp{
		Name:     userData[idUser]["nameUser"],
		LastName: userData[idUser]["lastNameUser"],
		Alias:    userData[idUser]["aliasUser"],
		Email:    userData[idUser]["emailUser"],
		Active:   userData[idUser]["activeUser"],
		Role:     userData[idUser]["roleAppUser_fk"],
		Team:     userData[idUser]["idTeam_fk"],
		Kyc:      userData[idUser]["kycUser_fk"],
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(userD)

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
	idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	fmt.Println(idInst, idTeam, idUser, authInst)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
}

func InviteUserTeam(respWriter http.ResponseWriter, request *http.Request) {
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
	idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if authInst != "1" {
		ok4 = false
	}
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	var invite documentflow.InviteUser
	if err := json.NewDecoder(request.Body).Decode(&invite); err != nil {
		return
	}

	invite.IdInst = idInst
	invite.IdTeam = idTeam
	invite.IdUser = idUser
	wheres := map[string][]string{
		"emailUser": {invite.Email},
		"idUser":    {invite.IdGuest},
		"LOGIC":     {"emailUser OR idUser"},
	}
	db.DB_con.GenericSelect("users", "idUser", []string{"nameUser", "emailUser", "activeUser"}, wheres)
	documentflow.InviteNewUser(invite)
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
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
			Team:     usdat["idTeam_fk"],
			Kyc:      usdat["kycUser_fk"],
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
		userCreator := objects.UserInstTeam(teamdat["creatorUser_fk"])
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
