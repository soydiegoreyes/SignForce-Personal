package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/auth"
	"sfmiddle/db"
	"time"

	"github.com/google/uuid"
)

func NewSignProcess(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
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
	idTeam, ok3 := claims["team"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====
	type DocDataRequest struct {
		IdDoc   string `json:"idDoc,omitempty"`
		HashDoc string `json:"hashDoc,omitempty"`
	}

	var req []DocDataRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	// ===== Configuración base =====
	attrs := []string{
		"documentHash", "createdAtDoc", "lastModifiedDoc", "deletedAtDoc",
		"deletedReasonDoc", "documentPath", "documentName", "documentExt",
		"abstractDoc", "authUseStatus", "authRoleStatus", "activeDoc",
	}

	wheres := map[string][]string{
		"idDocument":        {},
		"creatorUserDoc_fk": {idUser},
		"ownerInstDoc_fk":   {idInst},
		"ownerTeamDoc_fk":   {idTeam},
	}

	for _, d := range req {
		wheres["idDocument"] = append(wheres["idDocument"], d.IdDoc)
	}
	// ===== Ejecutar SELECT =====
	docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener documentos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, d := range req {
		if docData[d.IdDoc]["documentHash"] != d.HashDoc {
			docData[d.IdDoc]["activeDoc"] = "-1"
		}
	}
	folderPath := fmt.Sprintf("%s/folders/%s/%s/%s/", os.Getenv("TEMP_BASE_PATH"), idInst, idTeam, idUser)
	expiration := time.Now().AddDate(0, 0, 3).Format("2006-01-02 15:04:05")
	columns := []string{"ownerInst_fk", "ownerTeam_fk", "creatorUser_fk", "pathSerialized", "expirationDate", "numDocs"}
	values := []interface{}{idInst, idTeam, idUser, folderPath, expiration, len(docData)}

	idFolder, err := db.DB_con.GenericInsert("folders", columns, values)
	if err != nil {
		http.Error(respWriter, "Error creando folder", http.StatusInternalServerError)
	}
	columns = []string{"idfolderdocument", "idDocument", "idFolder"}
	for _, d := range req {
		idFD := uuid.NewString()
		values = []interface{}{idFD, d.IdDoc, idFolder}
		_, err := db.DB_con.GenericInsert("folderdocuments", columns, values)
		if err != nil {
			fmt.Println("Error al insertar documento en folderdocuments", err)
		}
		docData[d.IdDoc]["idFD"] = idFD
	}
	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(docData)
}
