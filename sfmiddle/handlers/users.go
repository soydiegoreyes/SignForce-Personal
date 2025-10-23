package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sfmiddle/auth"
	"sfmiddle/db"
	"sfmiddle/models"
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
	jsonData, err := json.Marshal(userD)
	if err != nil {
		http.Error(respWriter, "Error preparando datos para validación", http.StatusInternalServerError)
		return
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(jsonData)
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
}
