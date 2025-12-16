package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sfback/auth"
	"sfback/configs"
	"sfback/db"
	"sfback/models"
	"sfback/objects"
	"sfback/utilities"
	"strings"

	"github.com/joho/godotenv"
)

// Se ejecuta antes de main()
func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}
	log.Println(".env cargado correctamente")

	// Debug temporal:
	//log.Println("DB_USER:", os.Getenv("DB_USER"))

	db.DB_con = db.NewConn()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:8000": true,
			"http://127.0.0.1:8000": true,
			"http://localhost:8001": true,
			"http://127.0.0.1:8001": true,
			"http://localhost:8002": true,
			"http://127.0.0.1:8002": true,
			"http://localhost:4999": true,
			"http://127.0.0.1:4999": true,

			// IMPORTANTE: agrega TU dominio de túnel
			//"https://0153mh84-8000.usw3.devtunnels.ms/": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Vary", "Origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
	USER AND SIGNATURE CENTER
	-- Este modulo es para que un usuario pueda firmar y realizar operaciones criptográficas
	-- Recibe un login de un usuario y si este etá activo y es válido le permite realizar
	   operaciones criptográficas como firma, hash, obtener datos de su certificado publico
	   asi como conocer datos de su usuario
	-- Se comunica por un token que debe contener el Id de la aplicacion, el id del usuario
	   y datos relacionados con el tipo de operacion que se desea realizar
*/

// ============================================
func main() {
	defer db.DB_con.Desconectar()
	// Rutas
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/getuser", getUserData)
	http.HandleFunc("/signdocument", signFolderUser)
	http.HandleFunc("/hashdatab64", hashDataB64)
	http.HandleFunc("/uploadKeys", uploadKeys)
	http.HandleFunc("/logout", logoutUser)

	// Inicia el servidor
	log.Println("Servidor corriendo en el puerto", os.Getenv("API_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("API_PORT"), nil)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}

// ====================================== Handlers ======================================== //
/*
body:{
    id: "1",
	password: "mipass"
}
*/
//====================================================================================================
// Handler HTTP para loguear a un usuario por un ID de usuario y un arreglo de atributos a adquirir
func loginUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var user *objects.User
	var loginReq models.LoginUserReq
	var err error
	var token string

	err = json.NewDecoder(request.Body).Decode(&loginReq)
	if err != nil {
		http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
		return
	}

	_, isActive := auth.GetUser(loginReq.IdUser)
	if !isActive {
		fmt.Println("Usuario no está activo")
		user, err = objects.NewUser(loginReq.IdUser, loginReq.Password)
		if err != nil {
			json.NewEncoder(respWriter).Encode(models.LoginResponse{Token: "", Error: "Credenciales inválidas"})
			return
		}

		// Guarda el usuario en activos (opcional, si se necesita tracking)
		auth.AddUser(user.Uid, user)

		// Genera el token JWT y lo evuelve
		token, err = auth.GenerateJWT(user)
		if err != nil {
			http.Error(respWriter, "Error generando token", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(respWriter).Encode(models.LoginResponse{Token: token, Error: ""})
	} else {
		json.NewEncoder(respWriter).Encode(models.LoginResponse{Token: "", Error: "Ya tiene un usuario logueado. Desloguear para obtener token nuevo"})
	}
}

// ====================================================================================================
// Handler HTTP para obtener datos de un usuario por un ID de usuario y un arreglo de atributos a adquirir
func getUserData(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.GetUserRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if req.IdUser == "" {
		http.Error(respWriter, "ID requerido", http.StatusBadRequest)
		return
	}

	user, ok := auth.GetUser(req.IdUser)
	if !ok {
		http.Error(respWriter, "Usuario no encontrado o expirado", http.StatusUnauthorized)
		return
	}

	response := user.GetPublicParams(req.Fields)

	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(response)
}

// ====================================================================================================
// hash data in b64 by chunks
func hashDataB64(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
	}
	var req models.HashDataB64
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Datos inválidos", http.StatusBadRequest)
		return
	}
	configs.HashConf.Algorithm = req.DigestAlg
	hash_b64, err := utilities.GetHash(utilities.Decode_b64(req.DataB64), configs.HashConf)
	if err != nil {
		http.Error(respWriter, "Error al obtener hash de los datos", http.StatusNotImplemented)
		return
	}
	fmt.Println("hashed message: ", hash_b64)
	json.NewEncoder(respWriter).Encode(models.HashDataResponse{Operation: "abcd", HashedMessage: hash_b64, Check: true})
}

// ====================================================================================================
// recibe solicitud desde cualquier endpoint para validar llaves que ya estan almacenadas
func uploadKeys(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.UploadKeysReq
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if req.IdUser == "" || req.IdInst == "" {
		http.Error(respWriter, "IDs requeridos", http.StatusBadRequest)
		return
	}

	kn := fmt.Sprintf("%s/%s", req.Path, req.NameKey)
	cn := fmt.Sprintf("%s/%s", req.Path, req.NameCer)

	keys := objects.NewKeys(kn, cn)
	if keys == nil {
		http.Error(respWriter, "Error al obtener hash de las llaves", http.StatusInternalServerError)
		return
	}
	if keys.CertHash != req.HashCer && keys.KeyHash != req.HashKey {
		fmt.Println("Error Hashes no coinciden")
		http.Error(respWriter, "Hashes no coinciden", http.StatusExpectationFailed)
		return
	}
	valResp, err := keys.ValidateKeys(req.PassKey)
	if err != nil {
		fmt.Println("Error Llaves no encontrado o expirado ", err)
		http.Error(respWriter, "Llaves no encontrado o expirado", http.StatusInternalServerError)
		return
	}
	if !valResp.Valid {
		fmt.Println("Error Llaves no válidas o expiradas")
		http.Error(respWriter, "Llaves no válidas o expiradas", http.StatusNotFound)
		return
	}
	if valResp.Exists {
		fmt.Println("Error usuario ya ha subido la llave previamente")
		http.Error(respWriter, "Llaves no válidas o expiradas", http.StatusConflict)
		return
	}
	var validKeys int
	if keys.ValidKeys {
		validKeys = 1
	}

	cols := []string{"idUser_fk", "keyFilePath", "certFilePath", "serialNumber", "certVersion", "issuerRFC4514",
		"notValidAfter", "notValidBefore", "subjectRFC4514", "ocspUrl", "crlsUrl", "signature",
		"signAlgo", "validKeys", "keyLenKey", "hashKey", "hashCer", "subjectUniqueId", "subjectSerialNumber",
	}

	values := []interface{}{
		req.IdUser, kn, cn, keys.CertMap["SerialNumber"], keys.CertMap["Version"], keys.CertMap["Issuer"].(map[string]string)["RFC4514"],
		keys.CertMap["NotAfter"], keys.CertMap["NotBefore"], keys.CertMap["Subject"].(map[string]string)["RFC4514"], keys.CertMap["OCSP"], keys.CertMap["CRLS"], keys.CertMap["Signature"],
		keys.CertMap["SignatureAlgorithm"], validKeys, keys.CertMap["KeySize"], keys.KeyHash, keys.CertHash, keys.CertMap["SubjectUniqueId"], keys.CertMap["SubjectSerialNumber"],
	}

	keysId, err := db.DB_con.GenericInsert("userkeys", cols, values)
	if err != nil {
		http.Error(respWriter, "Error insertando nuevo registro de llaves", http.StatusInternalServerError)
		return
	}

	updates := map[string]map[string]interface{}{
		req.IdUser: {
			"idKeysUser_fk": keysId,
		},
	}
	err = db.DB_con.GenericBatchUpdate("users", "idUser", updates)
	if err != nil {
		http.Error(respWriter, "Error al actualizar valor de llaves", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Propietario: %s, Exp: %s \n", valResp.Owner, valResp.Expiration)
	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(valResp)
}

//====================================================================================================

// =================================== FIRMA DE DOCUMENTO ==============================================
// Recibe el id del folder y los documentos que solo el usuario puede firmar y los hace en secuencia
func signFolderUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== Autenticación =====
	user, err := auth.GetUserFromRequest(request)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	// ===== Parseo del body =====
	var req models.SignDocRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if req.IdInvite == "" || req.IdKey == "" || req.IdFolder == "" || len(req.Documents) == 0 {
		http.Error(respWriter, "Datos incompletos", http.StatusBadRequest)
		return
	}
	//  Validar que no exista firma previa ==========================
	attrs := []string{"idInvite_fk", "idUserKeys_fk", "digestValueSign", "signatureValueSign", "genTimeSign", "pathSign", "typeSign_fk", "nonceSign", "ipSignerSign", "notifyCreatorSign"}
	wheres := map[string][]string{
		"idInvite_fk": {req.IdInvite},
	}
	signatureData, err := db.DB_con.GenericSelect("signatures", "idSignature", attrs, wheres) // devuelve map[string]map[string]string
	if err != nil {
		http.Error(respWriter, "No se pudo acceder a datos de firma", http.StatusInternalServerError)
		return
	}

	// REVISAR ESTA SECCION YA QUE PUEDE TENER PROBLEMAS AL FIRMAR SELECTIVAMENTE
	// Recorre todas las firmas de la invitacion (del usuario)
	for _, s := range signatureData {
		if s["idUserKeys_fk"] != "" || s["signatureValueSign"] != "" {
			http.Error(respWriter, "Ya existe una firma registrada", http.StatusBadRequest)
			return
		}
	}

	// Cargar documentos desde DB ==============================
	docIds := []string{}
	for _, d := range req.Documents {
		docIds = append(docIds, d.IdDocument)
	}

	docAttrs := []string{"documentHash", "documentPath", "documentName", "documentExt", "activeDoc"}
	docWhere := map[string][]string{"idDocument": docIds}

	docData, err := db.DB_con.GenericSelect("documents", "idDocument", docAttrs, docWhere)
	if err != nil {
		http.Error(respWriter, "Error cargando documentos", http.StatusInternalServerError)
		return
	}

	// Procesar cada documento ==============================
	signaturesXML := make(map[string]map[string]string)

	for _, docReq := range req.Documents {
		docInfo, ok := docData[docReq.IdDocument]
		if !ok {
			fmt.Println("Documento no encontrado: ", docReq.IdDocument)
			continue
		}

		// Validar hash enviado vs BD
		if docInfo["documentHash"] != docReq.DocumentHash {
			fmt.Println("El hash del documento no coincide con la base de datos", docReq.IdDocument)
			continue
		}

		// Cargar archivo subido
		fileName := fmt.Sprintf("%s/%s%s.%s", os.Getenv("BASE_DIR"), docInfo["documentPath"], docInfo["documentName"], docInfo["documentExt"])
		fileBytes, err := os.ReadFile(fileName)
		if err != nil {
			http.Error(respWriter, "No se pudo leer el archivo real", http.StatusInternalServerError)
			return
		}

		// Obtener hash de original para comparación
		realHash, err := utilities.GetHash(fileBytes, configs.HashConf)
		if err != nil {
			http.Error(respWriter, "No se pudo obtener hash del archivo real", http.StatusInternalServerError)
			return
		}
		if realHash != docInfo["documentHash"] {
			http.Error(respWriter, "El archivo ha sido alterado (hash mismatch)", http.StatusBadRequest)
			return
		}

		// Actualizar datos de firmas ========================================
		// "idUser_fk", "idInvite_fk", "idUserKeys_fk", "digestValueSign",  "signatureValueSign", "genTimeSign", "pathSign", "typeSign_fk", "nonceSign", "ipSignerSign"

		var idSign string
		for idS, s := range signatureData {
			if s["digestValueSign"] == docInfo["documentHash"] {
				idSign = idS
				s["digestValueSign"] = s["digestValueSign"] + "."
				break
			}
		}

		// Generar firma XAdES
		docName := docInfo["documentName"] + "." + docInfo["documentExt"]

		xmlData, err := user.Keys.GenerarFirmaXades(utilities.Decode_b64(realHash), idSign, docName, docReq.IdDocument)
		if err != nil {
			http.Error(respWriter, "Error generando firma XAdES", http.StatusInternalServerError)
			return
		}
		xmlData["idDocument"] = docReq.IdDocument
		xmlData["nameDocument"] = docName
		// crear archivo para firmas qr (entregable)
		var folderPath string
		reg := regexp.MustCompile(`.*/folders/\d+/\d+/\d+/`)
		if baseFolderPath := reg.FindAllString(xmlData["xmlPath"], 1); len(baseFolderPath) > 0 {
			folderPath = baseFolderPath[0] + docReq.IdDocument + "_" + docName
			if _, err = os.Stat(folderPath); err != nil { // si hay error creamos el documento ya que no existe
				err = os.WriteFile(folderPath, fileBytes, 0644)
				if err != nil {
					fmt.Println("Error generando archivo entregable: ", err)
				}
			} else {
				fmt.Println("El documento ya existe: ", err)
			}

		} else {
			http.Error(respWriter, "Error critico error en path para archivo entregable", http.StatusInternalServerError)
			return
		}

		// se añade a la respuesta
		signaturesXML[idSign] = xmlData

		fmt.Println("resultado signatures: ", signaturesXML)
		updates := map[string]map[string]interface{}{
			idSign: {
				"idUserKeys_fk":      req.IdKey,
				"signatureValueSign": xmlData["signatureValueSign"],
				"genTimeSign":        strings.ReplaceAll(xmlData["genTimeSign"], "Z", ""),
				"pathSign":           strings.ReplaceAll(xmlData["xmlPath"], os.Getenv("BASE_DIR")+"/", ""),
				"typeSign_fk":        xmlData["typeSign"],
				"nonceSign":          xmlData["nonceSign"],
			},
		}

		err = db.DB_con.GenericBatchUpdate("signatures", "idSignature", updates)
		if err != nil {
			fmt.Println("Error en update firma:", idSign)
			continue
		}

	}

	//=============================================================================================
	respWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(respWriter).Encode(map[string]interface{}{
		"message": "OK",
		"signed":  signaturesXML,
	})
}

// ====================================================================================================
// =================================== LOGOUT DE USUARIO ==============================================
// logout handler para desloguear
func logoutUser(respWriter http.ResponseWriter, request *http.Request) {
	//user, err := auth.GetUserFromRequest(request)
	type idUserReq struct {
		IdUser string `json:"iduser"`
		IdKey  string `json:"idkey"`
	}
	var logoutReq *idUserReq
	err := json.NewDecoder(request.Body).Decode(&logoutReq)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusUnauthorized)
		return
	}
	userData, _ := db.DB_con.GenericSelect("users", "idUser", []string{"idKeysUser_fk"}, map[string][]string{"idUser": {logoutReq.IdUser}})
	if userData[logoutReq.IdUser]["idKeysUser_fk"] != logoutReq.IdKey {
		fmt.Println("Advertencia: las llaves en uso no corresponden a las llaves recibidas")
	}
	auth.DeleteUser(logoutReq.IdUser)
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("Sesión cerrada"))
}
