package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"time"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	//"sfmiddle/configs"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/objects"
	"sfmiddle/utilities"

	//"sfmiddle/utilities"

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
	MIDDLEWARE
	-- Este modulo es para comunicar las operaciones que tengan que ver con la webapp
	-- Recibe datos de usuario, datos de sesion, documentos y llaves de usuarios
	   estos datos depende de la etapa o la operacion son procesados o redirigidos
	   a las apis correspondientes por ejemplo: un usuario quiere realizar una firma,
	   entonces usando la sesion de su cuenta se mandan los identificadores a la api de firma
	-- Se comunica por un token que debe contener el Id de la aplicacion, el id del usuario
	   y datos relacionados con el tipo de operacion que se desea realizar

	-- LOS ENDPOINTS ESTAN DISEÑADOS PARA SEGUIR LOS PASOS:
	   - Registro de institucion
	   - Validacion de documentos de registro
	   - Aprobacion de registro
	   - Declinacion de registro
	   - Cancelacion de registro
	   - Acceso a paneles de gestion
	   - Acceso a flujo de firma
	   - Acceso a visualizador de documentos
	   - Acceso a crear, editar y guardar documentos
	   - Panel de estatus y pagos
	   - Cancelacion de contratos

*/

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Configurar headers CORS
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:8000" || origin == "http://127.0.0.1:8000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")

		// Manejar preflight request (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continuar con el handler siguiente
		next.ServeHTTP(w, r)
	})
}

// ============================================
func main() {
	defer db.DB_con.Desconectar()

	// Crear un nuevo ServeMux
	mux := http.NewServeMux()

	// Registrar rutas API
	mux.HandleFunc("/", home)
	mux.HandleFunc("/register", registerInst)
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/loginUser", login)
	mux.HandleFunc("/validation", validationInst)

	mux.Handle("/home/", http.StripPrefix("/home/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Determinar el Content-Type basado en la extensión del archivo
			switch filepath.Ext(r.URL.Path) {
			case ".css":
				w.Header().Set("Content-Type", "text/css")
			case ".js":
				w.Header().Set("Content-Type", "application/javascript")
			case ".html":
				w.Header().Set("Content-Type", "text/html")
			case ".png":
				w.Header().Set("Content-Type", "image/png")
			case ".jpg", ".jpeg":
				w.Header().Set("Content-Type", "image/jpeg")
			case ".ico":
				w.Header().Set("Content-Type", "image/x-icon")
			default:
				w.Header().Set("Content-Type", "text/plain")
			}

			http.FileServer(http.Dir("./../sffront")).ServeHTTP(w, r)
		})))

	// Servir archivos estáticos desde el directorio registro CORREGIDO ("registro")
	mux.Handle("/registro/", http.StripPrefix("/registro/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Determinar el Content-Type basado en la extensión del archivo
			switch filepath.Ext(r.URL.Path) {
			case ".css":
				w.Header().Set("Content-Type", "text/css")
			case ".js":
				w.Header().Set("Content-Type", "application/javascript")
			case ".html":
				w.Header().Set("Content-Type", "text/html")
			case ".png":
				w.Header().Set("Content-Type", "image/png")
			case ".jpg", ".jpeg":
				w.Header().Set("Content-Type", "image/jpeg")
			case ".ico":
				w.Header().Set("Content-Type", "image/x-icon")
			default:
				w.Header().Set("Content-Type", "text/plain")
			}

			http.FileServer(http.Dir("./../sffront/registro")).ServeHTTP(w, r)
		})))

	// Aplicar middleware CORS
	handler := corsMiddleware(mux)

	// Inicia el servidor
	log.Println("Servidor corriendo en el puerto", os.Getenv("API_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("API_PORT"), handler)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

// ====================================== Handlers ======================================== //
// landing page
func home(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/index.html")
}

// Handler HTTP para loguear a un usuario por un ID de usuario y un arreglo de atributos a adquirir
func registerInst(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var registerReq models.RegisterRequest
	registerResp := models.RegisterResponse{Check: false, InstId: "", Error: ""}
	var err error

	err = json.NewDecoder(request.Body).Decode(&registerReq)
	if err != nil {
		http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
		return
	}

	var status = map[string]bool{
		"2": true, "3": true, "4": true, "5": false, "6": true,
		"7": true, "8": true, "9": true, "10": false, "11": false,
		"12": true, "13": false,
	}

	// se comprueba que el TAXNUMBER de la empresa no existe en caso de que si, retorna error
	var whereMap = map[string][]string{
		"taxNumInst": {registerReq.TaxNumInst},
	}

	data, err := db.DB_con.GenericSelect("institutions", "idInstitution", []string{"statusInst_fk", "activeInst", "typeContractInst", "legalSignupName", "legalSignupLastname"}, whereMap)
	if err != nil {
		registerResp.Error = "Critical: error al obtener datos de institucion"
		json.NewEncoder(respWriter).Encode(registerResp)
		return

	} else {
		// en caso de no haber ningun registro con el mismo TAXNUMBER  se procede al registro
		if len(data) == 0 || status[data["idInstitution"]["statusInst_fk"]] {

			// se registra la institucion
			lastId, err := objects.RegisterInst(&registerReq)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

			registerResp.InstId = lastId
			fmt.Println("inst: ", lastId)

			// se genera un password temporal y se hashea
			tempPass := utilities.PassGenerator(12)
			passHash, err := utilities.GetHash([]byte(tempPass), configs.HashConf)
			if err != nil {
				registerResp.Error = "Error: No se pudo generar el password temporal"
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

			// se registra el usuario root
			userId, err := objects.RegisterUser(&registerReq, lastId, passHash)
			if err != nil {
				fmt.Println(err)
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}
			updates := map[string]map[string]interface{}{
				lastId: {
					"rootUser_fk": userId,
				},
			}
			err = db.DB_con.GenericBatchUpdate("institutions", "idInstitution", updates)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}
			fmt.Println("Usuario registrado ", userId)

			whereMap = map[string][]string{
				"nameApp": {"emailServ"},
			}

			data, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, whereMap)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

			binDoc, err := os.ReadFile("./templates/welcome_register.html")
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
			}

			body := string(binDoc)
			body = strings.ReplaceAll(body, "{TEMPORAL_USERNAME}", registerReq.ContactEmailInst)
			body = strings.ReplaceAll(body, "{TEMPORAL_PASS}", tempPass)
			body = strings.ReplaceAll(body, "{EXPIRATION_TIME}", time.Now().Add(30*24*time.Hour).Format("2006-01-02 15:04:05"))
			body = strings.ReplaceAll(body, "{URL_COMPLETAR_REGISTRO}", fmt.Sprintf("http://%s:%s/login", os.Getenv("API_IP"), os.Getenv("API_PORT")))

			payload := models.EmailRequest{
				IdUser:   userId,
				Subject:  fmt.Sprintf("¡Bienvenido a Signforce! Correo de verificación %s", registerReq.TaxNumInst),
				Body:     body,
				Dest:     []string{registerReq.ContactEmailInst},
				MimeType: "html",
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			var host, port string
			for _, v := range data {
				host = v["domainApp"]
				port = v["portApp"]
				break
			}

			req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/mailserv", host, port), bytes.NewBuffer(jsonPayload))
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("token", "3") // cambiar por bearer
			// Ejecutar petición
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

			fmt.Println("Código de respuesta:", resp.Status)
			if strings.Contains(resp.Status, "200 OK") {
				registerResp.Check = true
			}
			resp.Body.Close()

		} else {
			registerResp.Error = "Ya tiene un registro para su numero de empresa. Revisar estatus de su registro."
			json.NewEncoder(respWriter).Encode(registerResp)
			return
		}
	}

	json.NewEncoder(respWriter).Encode(registerResp)
}

// pagina de login
func loginPage(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/login.html")
}

// subir documentos para validacion
func validationInst(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/validation.html")

}

func login(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var loginReq models.LoginRequest
	loginResp := models.LoginResponse{Token: "", Error: ""}

	err = json.NewDecoder(request.Body).Decode(&loginReq)
	if err != nil {
		http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
		return
	}
	// se obtienen los datos de validacion del usuario
	var attrs = []string{"emailUser", "appPassHash", "activeUser", "authStatusUser", "idTeam_fk", "idInstitution_fk"}
	dataUser, err := db.DB_con.GenericSelect("users", "idUser", attrs, map[string][]string{"emailUser": {loginReq.Account}})
	if err != nil {
		loginResp.Error = fmt.Sprintf("%s", err)
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}
	if len(dataUser) == 0 {
		loginResp.Error = fmt.Sprintf("Error: Usuario %s no existe", loginReq.Account)
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}
	if len(dataUser) > 1 {
		loginResp.Error = fmt.Sprintf("Error: Usuario %s tiene mas de un registro", loginReq.Account)
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}

	// se valida el hash del password y se saca el id del usuario
	var idUser string
	for idU, v := range dataUser {
		idUser = idU
		nh, err := utilities.GetHash([]byte(loginReq.Password), configs.HashConf)
		if err != nil {
			loginResp.Error = "Error: No se pudo verificar el password"
			json.NewEncoder(respWriter).Encode(loginResp)
			return
		}
		if !(v["appPassHash"] == nh && v["activeUsere"] == "1") {
			loginResp.Error = "Error: Usuario no autorizado. Verificar Password o verifique el estado de su cuenta."
			json.NewEncoder(respWriter).Encode(loginResp)
			return
		}
		break
	}

	// se obtienen los datos de la institucion a la que pertenece el usuario
	attrs = []string{"statusInst_fk", "activeInst", "contactEmailInst"}
	dataInst, err := db.DB_con.GenericSelect("institutions", "idInstitution", attrs, map[string][]string{"idInstitution": {dataUser[idUser]["idInstitution_fk"]}})
	if err != nil {
		loginResp.Error = "Critical: error al obtener datos de institucion"
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}

	var redirectUrl = ""
	if dataInst[dataUser[idUser]["idInstitution_fk"]]["activeInst"] == "1" {
		switch dataInst[dataUser[idUser]["idInstitution_fk"]]["statusInst_fk"] {
		case "0":
			redirectUrl = "/validation"
		default:
			redirectUrl = "/loginUser"
		}
	}

	loginResp.Token = redirectUrl
	json.NewEncoder(respWriter).Encode(loginResp)
	//b := []string{"streetAddress", "addressLine", "postalCode", "neighborhood", "locality", "stateCodeInst_fk", "countryCodeInst_fk", "formattedAddress", "typeContractInst", "paymentDataInst_fk", "logoUrlInst", "rootUser_fk"}

}
