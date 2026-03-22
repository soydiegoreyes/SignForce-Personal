// ═══════════════════════════════════════════════════════════
// SignForce — Middleware Gateway (sfmiddle)
// ═══════════════════════════════════════════════════════════
// Service:     sfmiddle
// Port:        ${API_PORT} (default 5002)
// Description: Authentication gateway, session management,
//              request routing, and JWT token handling.
//              Proxies requests to sfback API.
// ═══════════════════════════════════════════════════════════

package main

import (
	"encoding/json"
	"fmt"

	"time"

	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"sfmiddle/auth"
	"sfmiddle/coms"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/utilities"

	"sfmiddle/handlers"

	"github.com/joho/godotenv"
)

// Se ejecuta antes de main()
func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error cargando .env: %v", err)
	}
	log.Println(".env cargado correctamente")

	db.DB_con = db.NewConn()
	coms.EmailCli = coms.ConfEmail()
	go func() {
		var lastUpdate = time.Now()
		var nextUpdt = lastUpdate
		for true {
			if time.Now().After(nextUpdt) {
				db.Update_algos()
				lastUpdate = time.Now()
				nextUpdt = lastUpdate.Add(24 * time.Hour)
				time.Sleep(24 * time.Hour)
			}
		}
	}()
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
		origin := r.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://localhost:8000": true,
			"http://127.0.0.1:8000": true,

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

// ============================================
func main() {
	defer db.DB_con.Desconectar()

	// Crear un nuevo ServeMux
	mux := http.NewServeMux()

	// Registrar rutas API
	mux.HandleFunc("/", home)
	mux.HandleFunc("/register", handlers.RegisterInst) // registrar nuevo cliente
	mux.HandleFunc("/loginUser", login)                // loguear usuario
	mux.HandleFunc("/logoutUser", logout)
	mux.HandleFunc("/uploadDocs", handlers.UploadDocs)                 // subir cualquier tipo de documento
	mux.HandleFunc("/downloadDoc", handlers.DownloadDoc)               // obtener datos de cualquier tipo de documento
	mux.HandleFunc("/statusDocs", handlers.StatusDocs)                 // obtener datos de cualquier tipo de documento
	mux.HandleFunc("/editDoc", handlers.EditDoc)                       // editar atributos de algún documento
	mux.HandleFunc("/interactDoc", handlers.InteractDoc)               // obtener datos de cualquier tipo de documento
	mux.HandleFunc("/uploadk", handlers.Uploadk)                       // subir llaves
	mux.HandleFunc("/statusk", handlers.GetKeysData)                   // obtener datos de llaves de usuario
	mux.HandleFunc("/updatek", handlers.UpdateKeysData)                // actualizar datos de llaves de usuario
	mux.HandleFunc("/logink", handlers.Logink)                         // login para firma
	mux.HandleFunc("/logoutk", handlers.Logoutk)                       // login para firma
	mux.HandleFunc("/getvaldata", handlers.GetValidationData)          // validar estatus de usuario en registro
	mux.HandleFunc("/updatevaldata", handlers.UpdateValidationData)    // actualizar estatus de usuario en registro
	mux.HandleFunc("/completevalidation", handlers.CompleteValidation) // completar validacion de usuario en registro
	mux.HandleFunc("/processpayment", handlers.ProcessPayment)         // procesar pago de plan
	mux.HandleFunc("/approvalsDash", handlers.Approvals)               // obtener datos de instituciones que estan en aprovacion
	mux.HandleFunc("/newSignFolder", handlers.NewSignFolder)           // empezar un proceso de firma desde cero
	mux.HandleFunc("/closeInvite", handlers.CloseAndInvite)            // cierra el folder con todas las invitaciones a firma
	mux.HandleFunc("/getinvite", handlers.GetInvite)                   // obtiene los datos de una invitación a firma
	mux.HandleFunc("/signDocument", handlers.SignDocument)             // endopoint para firma de documento
	mux.HandleFunc("/viewSign", handlers.ViewSign)                     // obtiene los equipos de una institucion
	mux.HandleFunc("/getAsice", handlers.BuildAsice)                   // obtiene los equipos de una institucion
	mux.HandleFunc("/getfolders", handlers.GetFolders)                 // obtiene los folders de un usuario
	mux.HandleFunc("/statusFolder", handlers.StatusFolder)             // obtener datos relacionados con las firmas del folder
	mux.HandleFunc("/inviteuser", handlers.InviteUser)                 // manda una invitacion a un usuario para formar parte de una institucion
	mux.HandleFunc("/getinviteuser", handlers.GetInviteUser)           // se obtienen datos de la invitacion para unirse a una institucion
	mux.HandleFunc("/createuser", handlers.CreateUser)                 // crea un usuario nuevo dentro de una institucion
	mux.HandleFunc("/updateUserStatus", handlers.UpdateUserStatus)     // actualizar datos de un usuario
	mux.HandleFunc("/findUser", handlers.CheckUserStatus)              // obtener datos de un usuario
	mux.HandleFunc("/validateFace", handlers.ValidateFace)
	mux.HandleFunc("/dashStats", handlers.GetStatsDash)

	// Rutas para servir páginas
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/noAuthPage", noauth)
	mux.HandleFunc("/validation", validationPage)
	mux.HandleFunc("/waitapprove", waitApprove)
	mux.HandleFunc("/check-email", checkEmail)
	mux.HandleFunc("/users", usersDash)
	mux.HandleFunc("/edituser", editUser)
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/mykeys", myKeys)
	mux.HandleFunc("/payment", payment)
	mux.HandleFunc("/approvals", approvalsDash)
	mux.HandleFunc("/mydocs", myDocuments)
	mux.HandleFunc("/myfolders", myFolders)
	mux.HandleFunc("/mytemplates", myTemplates)
	mux.HandleFunc("/addSigners", addSigners)
	mux.HandleFunc("/addSignatures", addSignatures)
	mux.HandleFunc("/viewSignature", viewSignature) // obtiene los equipos de una institucion
	mux.HandleFunc("/viewinvite", viewInvite)

	// carpetas publicas
	mux.Handle("/home/", http.StripPrefix("/home/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		http.FileServer(http.Dir("./../sffront/index")).ServeHTTP(w, r)
	})))
	mux.Handle("/images/", http.StripPrefix("/images/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determinar el Content-Type basado en la extensión del archivo
		switch filepath.Ext(r.URL.Path) {
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case ".ico":
			w.Header().Set("Content-Type", "image/x-icon")
		default:
			w.Header().Set("Content-Type", "image/svg")
		}

		http.FileServer(http.Dir("./../sffront/images")).ServeHTTP(w, r)
	})))
	// Servir archivos estáticos desde el directorio registro CORREGIDO ("registro")
	mux.Handle("/registro/", http.StripPrefix("/registro/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	mux.Handle("/administracion/", http.StripPrefix("/administracion/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		http.FileServer(http.Dir("./../sffront/administracion")).ServeHTTP(w, r)
	})))
	// Servir archivos estáticos desde el directorio registro CORREGIDO ("registro")
	mux.Handle("/documentflow/", http.StripPrefix("/documentflow/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		http.FileServer(http.Dir("./../sffront/documentFlow")).ServeHTTP(w, r)
	})))

	// sirve el modelo de reconocimiento facial
	mux.Handle("/facevector/", http.StripPrefix("/facevector/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Determinar el Content-Type basado en la extensión del archivo
		switch filepath.Ext(r.URL.Path) {
		case ".bin":
			w.Header().Set("Content-Type", "application/octet-stream")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		default:
			w.Header().Set("Content-Type", "text/plain")
		}

		http.FileServer(http.Dir("./../sffront/iamodels/facevector/")).ServeHTTP(w, r)
	})))

	// sirve los archivos temporales de templates
	mux.Handle("/iatemplates/", http.StripPrefix("/iatemplates/", http.HandlerFunc(handlers.IaTemplates)))

	// Aplicar middleware CORS
	handler := corsMiddleware(mux)

	// Inicia el servidor
	log.Println("Servidor corriendo en el puerto", os.Getenv("API_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("API_PORT"), handler)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

// ====================================== Handlers ========================================
// landing page
func home(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/index/index.html")
}

// pagina de no autorizacion
func noauth(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/noauth.html")
}

// pagina de login
func upload(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/upload_zone.html")
}
func approvalsDash(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/dashboards/dashboard_approvals.html")
}
func payment(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/plan_pay.html")
}
func myDocuments(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/documents.html")
}
func myTemplates(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/template_man.html")
}
func myFolders(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/folders.html")
}
func usersDash(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/dashboards/dashboard_admin.html")
}

// =======================================================================
// pagina de login
func loginPage(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/login.html")
}

// =======================================================================
// subir documentos para validacion
func validationPage(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := request.Cookie("token")
	if err != nil {
		fmt.Println("validation page: No cookie ", err)
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	idUser := claims["uid"].(string)
	idInst := claims["iid"].(string)
	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
	}
	data, err := db.DB_con.GenericSelect("users", "idUser", []string{"activeUser", "roleAppUser_fk", "idInstitution_fk"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener datos del usuario.", http.StatusInternalServerError)
		return
	}

	if len(data) > 0 {
		if data[idUser]["activeUser"] == "1" && data[idUser]["roleAppUser_fk"] == "1" {
			http.ServeFile(respWriter, request, "./../sffront/registro/validation.html")
		}
	} else {
		http.Error(respWriter, "Datos de usuario no encontrados", http.StatusUnauthorized)
		return
	}

}

// ==========================================================================================================
// sirve la pagina de espera de aprovacion
func waitApprove(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/waitapprove.html")
}

// ==========================================================================================================
// sirve la pagina de gestion de llaves
func myKeys(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/administracion/mykeys.html")
}

// ==========================================================================================================
// sirve la pagina para añadir firmantes
func addSigners(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/add_signers.html")
}

// ==========================================================================================================
// sirve la pagina para añadir firmas
func addSignatures(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/add_signs.html")
}

func viewSignature(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/view_signature.html")
}

// ==========================================================================================================
// sirve la pagina para ver una invitacion
func viewInvite(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/view_invite.html")
}

// ==========================================================================================================
// sirve la pagina para ver una invitacion de usuario nuevo
func editUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/edit_user.html")
}

// =======================================================================
// unificar el json de respuestas para que mande estatus y lista de documentos
// asegurar que multipart puede recibir uno o muchos archivos subidos de un mismo formulario y sugerir mejoras para subir archivos de distinta ubicacion
// =======================================================================
// =======================================================================
func login(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var loginReq models.LoginRequest
	loginResp := models.LoginResponse{RedirectTo: "noAuthPage", Error: ""} //se inicializa vacio para manejo mas sencillo

	err = json.NewDecoder(request.Body).Decode(&loginReq)
	if err != nil {
		http.Error(respWriter, "Error: JSON no valido.", http.StatusBadRequest)
		return
	}

	nh, err := utilities.GetHash([]byte(loginReq.Password), configs.HashConf, false)
	if err != nil {
		loginResp.Error = "Error: No se pudo verificar el password"
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}

	// se obtienen los datos de validacion del usuario
	//var attrs = []string{"emailUser", "appPassHash", "activeUser", "idTeam_fk", "roleAppUser_fk", "idInstitution_fk"}
	var attrs = []string{"emailUser", "appPassHash", "activeUser", "roleAppUser_fk", "idInstitution_fk", "nameUser"}
	var wheres = map[string][]string{
		"emailUser":   {loginReq.Account},
		"appPassHash": {nh},
	}
	dataUser, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		loginResp.Error = fmt.Sprintf("%s", err)
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}

	// no existe el usuario por lo que no puede hacer login
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
		if v["activeUser"] != "1" {
			loginResp.Error = "Error: Usuario no autorizado. Verificar Password o verifique el estado de su cuenta."
			json.NewEncoder(respWriter).Encode(loginResp)
			return
		}
		fmt.Println("hash valido y usuario activo")
		loginResp.UserName = v["nameUser"]
		break
	}
	// se asigna el id de la institucion
	idInst := dataUser[idUser]["idInstitution_fk"]

	// se obtienen los datos de la institucion a la que pertenece el usuario
	attrs = []string{"statusInst_fk", "activeInst", "contactEmailInst"}
	dataInst, err := db.DB_con.GenericSelect("institutions", "idInstitution", attrs, map[string][]string{"idInstitution": {idInst}})
	if err != nil {
		loginResp.Error = fmt.Sprintf("Critical: error al obtener datos de institucion %s", err)
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}

	// Generar JWT
	token, err := auth.GenerateJWT(idUser, dataUser[idUser]["roleAppUser_fk"], idInst, dataInst[idInst]["statusInst_fk"])
	if err != nil {
		http.Error(respWriter, "Error generando token", http.StatusInternalServerError)
		return
	}

	var satusActive = map[string]string{
		"2": "/validation", "3": "/waitapprove", "4": "/noAuthPage",
		"5": "/payment", "6": "/mykeys", "7": "/users", "8": "/noAuthPage",
		"9": "/noAuthPage", "10": "/noAuthPage", "11": "/noAuthPage",
	}
	//var satusInactive = map[string]bool{"9": true, "10": true, "11": true, "12": true}

	var location string
	// Si el usuario esta activo
	if dataUser[idUser]["activeUser"] == "1" {
		// se valida si la institucion no esta activa aun
		if dataInst[idInst]["activeInst"] == "0" {
			switch dataInst[idInst]["statusInst_fk"] {
			case "2":
				location = "/validation"
			case "3":
				location = "/waitapprove"
			case "5":
				location = "/payment"
			case "6":
				location = "/mykeys"
			default:
				location = "/noAuthPage"
				token = ""
			}
		} else {
			location = satusActive[dataInst[idInst]["statusInst_fk"]]
		}
	} else { // el usuario no esta activo y no tiene pemiso de entrar
		fmt.Println("no es usuario activo.")
		location = "/noAuthPage"
		token = ""
	}

	// Setear cookie con el token
	http.SetCookie(respWriter, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // poner en true en producción con HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(8 * time.Hour),
	})

	// EN LUGAR DE HACER REDIRECT, RETORNAMOS LA INFO AL FRONTEND
	loginResp.RedirectTo = location

	// Configurar headers CORS si es necesario
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("Access-Control-Allow-Credentials", "true")

	json.NewEncoder(respWriter).Encode(loginResp)
}

func logout(respWriter http.ResponseWriter, request *http.Request) {
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
	if !ok1 || !ok2 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	attrs := []string{"activeUser", "idKeysUser_fk"}
	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
	}

	http.SetCookie(respWriter, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // true en producción con https
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // fecha en el pasado
		MaxAge:   -1,              // MUY IMPORTANTE
	})

	userData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Datos de usuario no encontrado", http.StatusNotFound)
		return
	}
	if userData[idUser]["activeUser"] != "1" {
		http.Error(respWriter, "No es un usuario activo", http.StatusUnauthorized)
		return
	}
	idKey := userData[idUser]["idKeysUser_fk"]

	type idUserReq struct {
		IdUser string `json:"iduser"`
		IdKey  string `json:"idkey"`
	}

	// Crear request para logout
	reqData := idUserReq{
		IdUser: idUser,
		IdKey:  idKey,
	}

	// Validar las llaves mediante API externa
	userRequest, err := json.Marshal(reqData)
	if err != nil {
		http.Error(respWriter, "Error preparando datos para validación", http.StatusInternalServerError)
		return
	}

	uploadKeysURL := os.Getenv("BACK_URL") + "logout" // url de sfback
	req, err := http.NewRequest("POST", uploadKeysURL, bytes.NewBuffer(userRequest))
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

func checkEmail(respWriter http.ResponseWriter, request *http.Request) {
	http.ServeFile(respWriter, request, "./../sffront/check-email.html")
}
