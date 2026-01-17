package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/auth"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"

	"sfmiddle/utilities"
	"strconv"
	"strings"
)

// Función para subir llaves
func Uploadk(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Validar JWT
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
	if !(ok1 && ok2) {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Validar Content-Type
	contentType := request.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		http.Error(respWriter, "Tipo de contenido inesperado", http.StatusBadRequest)
		return
	}

	// Parsear el form multipart
	var maxmem int
	maxmem, err = strconv.Atoi(os.Getenv("MAX_UPLOAD_MEM"))
	if err != nil {
		maxmem = 10 // las llaves no pesan más de 10 kB
	}
	maxmem = maxmem << 10

	err = request.ParseMultipartForm(int64(maxmem))
	if err != nil {
		http.Error(respWriter, "Error procesando formulario", http.StatusBadRequest)
		return
	}
	// Limpiar recursos del multipart form al finalizar
	defer func() {
		if request.MultipartForm != nil {
			request.MultipartForm.RemoveAll()
		}
	}()

	// Obtener la contraseña
	var passKey string
	if passVal, exists := request.MultipartForm.Value["passKey"]; exists && len(passVal) > 0 {
		passKey = passVal[0]
	} else {
		http.Error(respWriter, "Contraseña requerida", http.StatusBadRequest)
		return
	}

	// Obtener los archivos específicos (keyFile y certFile)
	keyFile, keyHeader, err := request.FormFile("keyFile")
	if err != nil {
		http.Error(respWriter, "Archivo de llave requerido", http.StatusBadRequest)
		return
	}
	defer keyFile.Close()

	certFile, certHeader, err := request.FormFile("certFile")
	if err != nil {
		http.Error(respWriter, "Archivo de certificado requerido", http.StatusBadRequest)
		return
	}
	defer certFile.Close()

	// Crear directorio para guardar los archivos
	basePath := fmt.Sprintf("%s/%s/%s/%s", os.Getenv("KEYS_PATH"), idInst, idUser, strings.Split(certHeader.Filename, ".")[0])
	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		http.Error(respWriter, "Error creando directorio", http.StatusInternalServerError)
		return
	}
	// check para saber si hubo errores en el proceso una vez guardados los archivos para borrar los datos
	var Check bool

	defer func() {
		if !Check {
			os.RemoveAll(basePath)
		}
	}()

	// Guardar archivo de llave
	//keyPath := fmt.Sprintf("%s/%s", basePath, keyHeader.Filename)
	filePath, err := utilities.GuardarArchivo(keyFile, basePath, keyHeader.Filename, idInst, idUser, false)
	if err != nil {
		http.Error(respWriter, "Error creando archivo de llave", http.StatusInternalServerError)
		return
	}
	// Calcular hash de los archivos
	keyHash, err := utilities.GetHash(filePath, configs.HashConf)
	if err != nil {
		http.Error(respWriter, "Error obteniendo hash de la llave", http.StatusInternalServerError)
		return
	}

	// Guardar archivo de certificado
	//certPath := fmt.Sprintf("%s/%s", basePath, certHeader.Filename)
	filePath, err = utilities.GuardarArchivo(certFile, basePath, certHeader.Filename, idInst, idUser, false)
	if err != nil {
		http.Error(respWriter, "Error creando archivo de certificado", http.StatusInternalServerError)
		return
	}
	certHash, err := utilities.GetHash(filePath, configs.HashConf)
	if err != nil {
		http.Error(respWriter, "Error obteniendo hash del certificado", http.StatusInternalServerError)
		return
	}

	// Crear metadata para validación
	keyMetadata := models.KeyMetadata{
		IdInst:  idInst,
		IdUser:  idUser,
		PassKey: passKey,
		NameKey: keyHeader.Filename,
		NameCer: certHeader.Filename,
		HashKey: keyHash,
		HashCer: certHash,
		Path:    basePath,
	}

	// Validar las llaves mediante API externa
	validationData, err := json.Marshal(keyMetadata)
	if err != nil {
		http.Error(respWriter, "Error preparando datos para validación", http.StatusInternalServerError)
		return
	}

	uploadKeysURL := os.Getenv("BACK_URL") + "uploadKeys" // url de sfback
	req, err := http.NewRequest("POST", uploadKeysURL, bytes.NewBuffer(validationData))
	if err != nil {
		http.Error(respWriter, "Error creando solicitud de validación", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(respWriter, "Error conectando al servicio de validación", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		http.Error(respWriter, "LLaves ya existen", http.StatusConflict)
		return
	}
	if resp.StatusCode != http.StatusOK {
		http.Error(respWriter, "Error en la validación de llaves", http.StatusInternalServerError)
		return
	}

	// se busca si el usuario tiene llaves.
	wheres := map[string][]string{
		"idInstitution": {idInst}, // todas las llaves del usuario
	}
	instData, err := db.DB_con.GenericSelect("institutions", "idInstitution", []string{"statusInst_fk", "activeInst"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de llaves", http.StatusInternalServerError)
		return
	}
	// si no tiene contrato aun y esta en el paso de subir llaves entonces es usuario nuevo y debe pasar a firma de contratos
	if instData[idInst]["statusInst_fk"] == "6" || instData[idInst]["activeInst"] == "0" {
		updates := map[string]map[string]interface{}{
			idInst: {
				"statusInst_fk": 7,
			},
		}

		err = db.DB_con.GenericBatchUpdate("institutions", "idInstitution", updates)
		if err != nil {
			http.Error(respWriter, "Error al actualizar valor de llaves", http.StatusInternalServerError)
			return
		}
	}
	Check = true

	var validationResult struct {
		Valid      bool   `json:"valid"`
		Expiration string `json:"expiration"`
		Owner      string `json:"owner"`
	}
	err = json.NewDecoder(resp.Body).Decode(&validationResult)
	if err != nil || !validationResult.Valid {
		http.Error(respWriter, "Las llaves no son válidas", http.StatusBadRequest)
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&validationResult)
}

func Logink(respWriter http.ResponseWriter, request *http.Request) {
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
	//idTeam, ok3 := claims["team"].(string)

	if !ok1 || !ok2 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	attrs := []string{"activeUser", "idKeysUser_fk", "isAliveUser"}
	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
		//"idTeam_fk":        {idTeam},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener datos de usuario", http.StatusInternalServerError)
		return
	}
	if userData[idUser]["activeUser"] != "1" {
		http.Error(respWriter, "Usuario inactivo revise estatus.", http.StatusUnauthorized)
		return
	}

	if userData[idUser]["isAliveUser"] != "1" {
		http.Error(respWriter, "Requiere prueba de vida.", http.StatusUnauthorized)
	}

	type LoginK struct {
		IdUser   string `json:"iduser"`
		Password string `json:"password"`
	}
	var loginReq *LoginK
	err = json.NewDecoder(request.Body).Decode(&loginReq)
	if err != nil {
		http.Error(respWriter, "Error al decodificar peticion.", http.StatusInternalServerError)
		return
	}
	loginReq.IdUser = idUser

	// Convertir a json para enviar como body
	jsonPayload, err := json.Marshal(loginReq)
	if err != nil {
		http.Error(respWriter, "Error al convertir a json.", http.StatusInternalServerError)
		return
	}
	// se obtiene el nombre de la app y si esta activa
	attrs = []string{"domainApp", "portApp"}
	wheres = map[string][]string{
		"nameApp":  {"sfback"},
		"isActive": {"1"},
	}
	appValues, err := db.DB_con.GenericSelect("microapps", "idapp", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener datos de app.", http.StatusInternalServerError)
		return
	}
	if len(appValues) == 0 {
		http.Error(respWriter, "No se encontró la app activa.", http.StatusInternalServerError)
		return
	}
	var host, port string
	for _, app := range appValues {
		host = app["domainApp"]
		port = app["portApp"]
	}
	url := fmt.Sprintf("http://%s:%s/login", host, port)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		http.Error(respWriter, "Error en new request.", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
	// Ejecutar petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(respWriter, "Error en request.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(respWriter, "Error status: "+resp.Status, http.StatusBadRequest)
		return
	}
	type token struct {
		Token string `json:"token"`
	}
	var authk token
	err = json.NewDecoder(resp.Body).Decode(&authk)
	if err != nil {
		http.Error(respWriter, "Error al codificar respuesta.", http.StatusInternalServerError)
		return
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(&authk)
}

func Logoutk(respWriter http.ResponseWriter, request *http.Request) {
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

	if !ok1 || !ok2 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	idkey := make(map[string]string)
	err = json.NewDecoder(request.Body).Decode(&idkey)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	type idUserReq struct {
		IdUser string `json:"iduser"`
		IdKey  string `json:"idkey"`
	}

	// Crear request para logout
	reqData := idUserReq{
		IdUser: idUser,
		IdKey:  idkey["idKey"],
	}

	// Validar las llaves mediante API externa
	userData, err := json.Marshal(reqData)
	if err != nil {
		http.Error(respWriter, "Error preparando datos para validación", http.StatusInternalServerError)
		return
	}

	uploadKeysURL := os.Getenv("BACK_URL") + "logout" // url de sfback
	req, err := http.NewRequest("POST", uploadKeysURL, bytes.NewBuffer(userData))
	if err != nil {
		http.Error(respWriter, "Error creando solicitud de logout", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(respWriter, "Error conectando al servicio de logout", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(respWriter, "Error en logout de firma", http.StatusInternalServerError)
		return
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(map[string]string{"message": "OK"})
}

// =======================================================================
func GetKeysData(respWriter http.ResponseWriter, request *http.Request) {
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

	// Extraer datos del JWT
	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	//idTeam, ok3 := claims["team"].(string)

	if !ok1 || !ok2 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
		//"idTeam_fk":        {idTeam},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", []string{"activeUser", "idKeysUser_fk"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de usuario", http.StatusInternalServerError)
		return
	}
	if userData[idUser]["activeUser"] != "1" {
		http.Error(respWriter, "No autorizado. Usuario inactivo", http.StatusUnauthorized)
		return
	}

	wheres = map[string][]string{
		"idUser_fk": {idUser},
	}
	keys, err := db.DB_con.GenericSelect("userkeys", "idUserKeys", []string{"keyFilePath", "certFilePath", "notValidAfter", "issuerRFC4514", "subjectRFC4514", "subjectUniqueId", "createdAtKey"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de llaves", http.StatusInternalServerError)
		return
	}

	keyStatResp := make([]*models.KeysStatus, 0)
	var selected bool
	for idKey, key := range keys {
		if userData[idUser]["idKeysUser_fk"] == idKey {
			selected = true
		} else {
			selected = false
		}
		ks := &models.KeysStatus{
			IdKey:           idKey,
			NameKey:         key["keyFilePath"][strings.LastIndex(key["keyFilePath"], "/")+1:],
			NameCer:         key["certFilePath"][strings.LastIndex(key["certFilePath"], "/")+1:],
			Expiration:      key["notValidAfter"],
			Owner:           key["subjectRFC4514"][strings.Index(key["subjectRFC4514"], "=")+1 : strings.Index(key["subjectRFC4514"], ",")],
			SubjectUniqueId: key["subjectUniqueId"],
			IssuerRFC4514:   key["issuerRFC4514"],
			UploadedAt:      key["createdAtKey"],
			Selected:        selected,
		}
		keyStatResp = append(keyStatResp, ks)
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(keyStatResp)
}

// =======================================================================
func UpdateKeysData(respWriter http.ResponseWriter, request *http.Request) {
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
	// Extraer datos del JWT

	idUser, ok1 := claims["uid"].(string)

	if !ok1 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	type updateKeyReq struct {
		IdKeyUpdate string `json:"idKeyUpdate"`
	}
	keyReq := updateKeyReq{}
	err = json.NewDecoder(request.Body).Decode(&keyReq)
	if err != nil {
		http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
		return
	}
	updates := map[string]map[string]interface{}{
		idUser: {
			"idKeysUser_fk": strings.TrimLeft(keyReq.IdKeyUpdate, "key-"),
		},
	}

	err = db.DB_con.GenericBatchUpdate("users", "idUser", updates)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de usuario", http.StatusInternalServerError)
		return
	}

	type updateKeys struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
	ks := &updateKeys{
		Status:  true,
		Message: "Valor actualizdo",
	}
	// Convertir a JSON
	jsonData, _ := json.Marshal(ks)

	// Configurar headers y enviar respuesta
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(jsonData)
}
