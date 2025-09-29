package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sfback/auth"
	"sfback/configs"
	"sfback/db"
	"sfback/models"
	"sfback/objects"
	"sfback/utilities"

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
	http.HandleFunc("/usersignature", signUser)
	http.HandleFunc("/newinvitation", newInvitation)
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
		user, err = objects.NewUser(loginReq.IdUser, loginReq.Password)
		if err != nil {
			json.NewEncoder(respWriter).Encode(models.LoginResponse{Token: "", Error: "Credenciales inválidas"})
			return
		}

		// Guarda el usuario en activos (opcional, si necesitas tracking)
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

func signUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, err := auth.GetUserFromRequest(request)
	if err != nil {
		http.Error(respWriter, "Usuario no encontrado en usuarios activos", http.StatusNotFound)
		return
	}

	var req models.UserSignRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var signBytes []byte
	signBytes, err = user.Keys.SignData("sha256", utilities.Decode_b64(req.HashedMessage))
	if err != nil {
		http.Error(respWriter, "Usuario no encontrado en usuarios activos", http.StatusNotImplemented)
		return
	}

	signature := utilities.Encode_b64(signBytes)
	fmt.Println("signature: ", signature)
	json.NewEncoder(respWriter).Encode(models.SignatureResponse{Operation: "abcd", Signature: signature, Check: true})
}

func newInvitation(respWriter http.ResponseWriter, request *http.Request) {

}

// logout handler para desloguear
func logoutUser(respWriter http.ResponseWriter, request *http.Request) {
	user, err := auth.GetUserFromRequest(request)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusUnauthorized)
		return
	}

	auth.DeleteUser(user.Uid)
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write([]byte("Sesión cerrada"))
}
