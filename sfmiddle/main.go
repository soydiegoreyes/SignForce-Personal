package main

import (
	"encoding/json"
	"fmt"

	"time"

	"log"
	"net/http"

	"os"
	"path/filepath"

	"sfmiddle/auth"
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
	mux.HandleFunc("/register", handlers.RegisterInst)                 // registrar nuevo cliente
	mux.HandleFunc("/loginUser", login)                                // loguear usuario
	mux.HandleFunc("/uploadDocs", handlers.UploadDocs)                 // subir cualquier tipo de documento
	mux.HandleFunc("/downloadDoc", handlers.DownloadDoc)               // obtener datos de cualquier tipo de documento
	mux.HandleFunc("/statusDocs", handlers.StatusDocs)                 // obtener datos de cualquier tipo de documento
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
	mux.HandleFunc("/findUser", handlers.CheckUserStatus)              // obtener datos de un usuario
	mux.HandleFunc("/updateUserStatus", handlers.UpdateUserStatus)     // actualizar datos de un usuario
	mux.HandleFunc("/approvals", handlers.Approvals)                   // obtener datos de instituciones que estan en aprovacion
	mux.HandleFunc("/newSignFolder", handlers.NewSignFolder)           // empezar un proceso de firma desde cero
	mux.HandleFunc("/closeInvite", handlers.CloseAndInvite)            // cierra el folder con todas las invitaciones a firma
	mux.HandleFunc("/getinvite", handlers.GetInvite)                   // obtiene los datos de una invitacion
	mux.HandleFunc("/signDocument", handlers.SignDocument)             // endopoint para firma de documento
	mux.HandleFunc("/getfolder", handlers.GetFolders)                  // obtiene los folders de un usuario
	mux.HandleFunc("/inviteuser", handlers.InviteUser)                 // manda una invitacion a un usuario para formar parte de una institucion
	mux.HandleFunc("/getinviteuser", handlers.GetInviteUser)           // se obtienen datos de la invitacion para unirse a una institucion
	mux.HandleFunc("/createuserate", handlers.CreateUser)              // crea un usuario nuevo dentro de una institucion
	//mux.HandleFunc("/instteams", handlers.InstTeams)                   // obtiene los equipos de una institucion
	//mux.HandleFunc("/teamusers", handlers.TeamUsers)                   // obtinene los usuarios de un equipo
	//mux.HandleFunc("/newteam", handlers.NewTeam)                       // crea un equipo dentro de una institucion por un usuario master o root

	// Rutas para servir páginas
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/noAuthPage", noauth)
	mux.HandleFunc("/validation", validationPage)
	mux.HandleFunc("/waitapprove", waitApprove)
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/mykeys", myKeys)
	mux.HandleFunc("/payment", payment)
	mux.HandleFunc("/dashboard/approvals", approvalsDash)
	mux.HandleFunc("/mydocs", myDocuments)
	mux.HandleFunc("/addSigners", addSigners)
	mux.HandleFunc("/addSignatures", addSignatures)
	mux.HandleFunc("/viewinvite", viewInvite)
	mux.HandleFunc("/viewinviteuser", viewInviteUser)
	mux.HandleFunc("/myfolders", myFolders)
	mux.HandleFunc("/dashboard/users", usersDash)

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

		http.FileServer(http.Dir("./../sffront")).ServeHTTP(w, r)
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
	http.ServeFile(respWriter, request, "./../sffront/index.html")
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
	fmt.Println("data: ", data)
	if len(data) > 0 {
		if data[idUser]["activeUser"] == "1" && data[idUser]["roleAppUser_fk"] == "1" {
			//respWriter.Header().Set("Authorization", "Bearer "+cookie.Value)
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
func viewInviteUser(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/invite_user.html")
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

	nh, err := utilities.GetHash([]byte(loginReq.Password), configs.HashConf)
	if err != nil {
		loginResp.Error = "Error: No se pudo verificar el password"
		json.NewEncoder(respWriter).Encode(loginResp)
		return
	}

	// se obtienen los datos de validacion del usuario
	//var attrs = []string{"emailUser", "appPassHash", "activeUser", "idTeam_fk", "roleAppUser_fk", "idInstitution_fk"}
	var attrs = []string{"emailUser", "appPassHash", "activeUser", "roleAppUser_fk", "idInstitution_fk"}
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
	//token, err := auth.GenerateJWT(idUser, dataUser[idUser]["idTeam_fk"], dataUser[idUser]["roleAppUser_fk"], idInst, dataInst[idInst]["statusInst_fk"])
	token, err := auth.GenerateJWT(idUser, dataUser[idUser]["roleAppUser_fk"], idInst, dataInst[idInst]["statusInst_fk"])
	if err != nil {
		http.Error(respWriter, "Error generando token", http.StatusInternalServerError)
		return
	}

	//var statusVal = map[string]bool{"2": true, "3": true, "4": true, "5": true}
	//var statusContr = map[string]bool{"6": true, "7": true}
	var satusActive = map[string]bool{"8": true}
	var satusInactive = map[string]bool{"9": true, "10": true, "11": true, "12": true}

	var location string
	// Si el usuario esta activo
	if dataUser[idUser]["activeUser"] == "1" {
		// se valida si la institucion no esta activa aun
		if dataInst[idInst]["activeInst"] == "0" {
			// si no esta activa entonces hay que ver que estatus tiene
			if dataInst[idInst]["statusInst_fk"] == "2" {
				// se encuentra en etapa de validacion por lo que se redirige a /validacion
				location = "/validation"
			} else if dataInst[idInst]["statusInst_fk"] == "3" {
				location = "/waitapprove"
			} else if dataInst[idInst]["statusInst_fk"] == "5" {
				location = "/payment"
			} else if dataInst[idInst]["statusInst_fk"] == "6" {
				location = "/mykeys"
			} else if dataInst[idInst]["statusInst_fk"] == "7" {
				location = "/mydocs"
			} else if dataInst[idInst]["statusInst_fk"] == "8" {
				location = "/dashboard"
			} else if satusInactive[dataInst[idInst]["statusInst_fk"]] {
				location = "/noAuthPage"
			} else {
				location = "/noAuthPage"
			}
		} else { // la institucion ya está activa (en un estatus ACTIVO)
			if satusActive[dataInst[idInst]["statusInst_fk"]] {
				location = "/dashboard"
			} else { // La institucion estaba activa pero fue dada de baja, suspendida o revocada
				location = "/noAuthPage"
			}
		}
	} else { // el usuario no esta activo y no tiene pemiso de entrar
		location = "/noAuthPage"
	}
	//fmt.Println(location)

	// Setear cookie con el token
	http.SetCookie(respWriter, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // poner en true en producción con HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(1 * time.Hour),
	})

	// EN LUGAR DE HACER REDIRECT, RETORNAMOS LA INFO AL FRONTEND
	//loginResp.Token = token
	loginResp.RedirectTo = location

	// Configurar headers CORS si es necesario
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("Access-Control-Allow-Credentials", "true")

	json.NewEncoder(respWriter).Encode(loginResp)
}
