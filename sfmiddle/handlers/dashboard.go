package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sfmiddle/auth"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/objects"
	"strconv"
)

func GetStatsDash(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== Autenticación por JWT =====
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
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Preaparamos estructura de respuesta

	type StatsResp struct {
		UserData       models.UserDataResp
		InstData       models.InstResponse
		NumUsers       int
		NumUsersMonth  int
		TotalInvites   int
		PendingInvites int
		InvitesMonth   int
		FolderStats    map[string]int
	}

	userInstData := objects.UserInstJoin(idUser)
	fmt.Println(userInstData)
	var resp = StatsResp{
		UserData: models.UserDataResp{
			Name:       userInstData["nameUser"],
			LastName:   userInstData["lastNameUser"],
			Alias:      userInstData["aliasUser"],
			Email:      userInstData["emailUser"],
			Active:     userInstData["activeUser"],
			TotalSigns: 0,
			PendSigns:  0,
			SignedDocs: 0,
		},
		InstData: models.InstResponse{
			LegalName: userInstData["legalNameInst"],
			AliasName: userInstData["aliasNameInst"],
			Email:     userInstData["contactEmailInst"],
			IsActive:  userInstData["activeInst"] == "1",
		},
		NumUsers:       0,
		NumUsersMonth:  0,
		TotalInvites:   0,
		PendingInvites: 0,
		InvitesMonth:   0,
	}

	// users
	query := "SELECT COUNT(*) AS totalUsers, SUM(CASE WHEN createdAtUser >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS newUsersMonth FROM users WHERE idInstitution_fk = " + idInst + ";"
	numUsers, err := db.DB_con.ExecuteSelect(query)
	for _, num := range numUsers {
		if resp.NumUsers, err = strconv.Atoi(num["totalUsers"]); err != nil {
			resp.NumUsers = 0
		}
		if resp.NumUsersMonth, _ = strconv.Atoi(num["newUsersMonth"]); err != nil {
			resp.NumUsersMonth = 0
		}
		break
	}
	// firmas
	query = "SELECT COUNT(*) AS totalSigns, SUM(CASE WHEN genTimeSign IS NULL THEN 1 ELSE 0 END) AS pendingSigns, SUM(CASE WHEN genTimeSign >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS signsMonth FROM signatures WHERE idUser_fk = " + idUser + ";"
	numSigns, err := db.DB_con.ExecuteSelect(query)

	for _, numS := range numSigns {
		if resp.UserData.TotalSigns, err = strconv.Atoi(numS["totalSigns"]); err != nil {
			resp.UserData.TotalSigns = 0
		}
		if resp.UserData.SignedDocs, err = strconv.Atoi(numS["signsMonth"]); err != nil {
			resp.UserData.SignedDocs = 0
		}
		if resp.UserData.PendSigns, err = strconv.Atoi(numS["pendingSigns"]); err != nil {
			resp.UserData.PendSigns = 0
		}
		break
	}
	// invites
	query = "SELECT COUNT(*) AS totalInvites, SUM(CASE WHEN closedAt IS NULL THEN 1 ELSE 0 END) AS pendingInvites, SUM(CASE WHEN closedAt >= DATE_FORMAT(CURDATE(), '%Y-%m-01') THEN 1 ELSE 0 END) AS InvitesMonth FROM invites WHERE idUserOwnner_fk = " + idUser + ";"
	NumInvites, err := db.DB_con.ExecuteSelect(query)

	for _, numI := range NumInvites {
		if resp.TotalInvites, err = strconv.Atoi(numI["totalInvites"]); err != nil {
			resp.TotalInvites = 0
		}
		if resp.InvitesMonth, err = strconv.Atoi(numI["InvitesMonth"]); err != nil {
			resp.InvitesMonth = 0
		}
		if resp.PendingInvites, err = strconv.Atoi(numI["pendingInvites"]); err != nil {
			resp.PendingInvites = 0
		}
		break
	}

	// folders stats
	var foldersHour = make(map[string]int)
	for i := range 24 {
		var m string = ""
		m = strconv.Itoa(i)
		if len(m) == 1 {
			m = "0" + m
		}
		foldersHour[m+":00"] = 0
	}
	folderData, err := db.DB_con.ExecuteSelect("SELECT CONCAT(LPAD(FLOOR(HOUR(creationAt)), 2, '0'),':00') AS intervalo, COUNT(*) AS totalfolders FROM folders WHERE creationAt >= CURDATE() AND creationAt < CURDATE() + INTERVAL 1 DAY GROUP BY CONCAT(LPAD(FLOOR(HOUR(creationAt)), 2, '0'),':00') ORDER BY intervalo;")
	if err != nil {
		fmt.Println("Error al obtrener datos del folder")
	}

	for k, v := range folderData {
		count, err := strconv.Atoi(v["totalfolders"])
		if err != nil {
			count = -1
		}
		foldersHour[k] = count
	}
	resp.FolderStats = foldersHour
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&resp)
}
