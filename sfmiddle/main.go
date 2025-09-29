package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"time"

	"log"
	"net/http"

	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sfmiddle/auth"
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
	mux.HandleFunc("/uploadDocs", uploadDocs)                 // subir cualquier tipo de documento
	mux.HandleFunc("/getDocs", getDocs)                       // obtener datos de cualquier tipo de documento
	mux.HandleFunc("/uploadk", uploadk)                       // subir llaves
	mux.HandleFunc("/statusk", getKeysData)                   // obtener datos de llaves de usuario
	mux.HandleFunc("/updatek", updateKeysData)                // actualizar datos de llaves de usuario
	mux.HandleFunc("/register", registerInst)                 // registrar nuevo cliente
	mux.HandleFunc("/loginUser", login)                       // loguear usuario
	mux.HandleFunc("/getvaldata", getValidationData)          // validar estatus de usuario en registro
	mux.HandleFunc("/updatevaldata", updateValidationData)    // actualizar estatus de usuario en registro
	mux.HandleFunc("/completevalidation", completeValidation) // completar validacion de usuario en registro
	mux.HandleFunc("/processpayment", processPayment)         // procesar pago de plan
	mux.HandleFunc("/checkUserStatus", checkUserStatus)       // obtener datos de un usuario
	mux.HandleFunc("/approvals", approvals)                   // obtener datos de instituciones que estan en aprovacion

	// Rutas para servir páginas
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/noAuthPage", noauth)
	mux.HandleFunc("/validation", validationPage)
	mux.HandleFunc("/waitapprove", waitApprove)
	mux.HandleFunc("/upload", upload)
	mux.HandleFunc("/uploadKeys", uploadKeys)
	mux.HandleFunc("/payment", payment)
	mux.HandleFunc("/dashboard/approvals", approvalsDash)
	mux.HandleFunc("/mydocs", myDocuments)

	// carpetas publicas
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
	mux.Handle("/administracion/", http.StripPrefix("/administracion/",
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

			http.FileServer(http.Dir("./../sffront/administracion")).ServeHTTP(w, r)
		})))
	// Servir archivos estáticos desde el directorio registro CORREGIDO ("registro")
	mux.Handle("/documentflow/", http.StripPrefix("/documentflow/",
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
	http.ServeFile(respWriter, request, "./../sffront/documentFlow/my_documents.html")
}

// =======================================================================
// Handler HTTP para loguear a un usuario por un ID de usuario y un arreglo de atributos a adquirir
func registerInst(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
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

	// se comprueba que el TAXNUMBER de la empresa no existe en caso de que si, retorna error
	whereMap := map[string][]string{
		"taxNumInst": {registerReq.TaxNumInst},
	}

	data, err := db.DB_con.GenericSelect("institutions", "idInstitution", []string{"statusInst_fk", "activeInst", "typeContractInst", "legalSignupName", "legalSignupLastname"}, whereMap)
	if err != nil {
		registerResp.Error = "Critical: error al obtener datos de institucion"
		json.NewEncoder(respWriter).Encode(registerResp)
		return

	} else {
		whereMap = map[string][]string{
			"idStatusInst": {data["idInstitution"]["statusInst_fk"]},
		}
		status, err := db.DB_con.GenericSelect("statusinstitution", "idStatusInst", []string{"permissionStatusInst_fk"}, whereMap)
		if err != nil {
			registerResp.Error = "Critical: error al obtener estatus de institucion"
			json.NewEncoder(respWriter).Encode(registerResp)
			return
		}
		// en caso de no haber ningun registro con el mismo TAXNUMBER  se procede al registro
		if len(data) == 0 || status["idStatusInst"]["permissionStatusInst_fk"] == "1" {

			// se registra la institucion
			lastId, err := objects.RegisterInst(&registerReq)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

			registerResp.InstId = lastId

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
			fmt.Println("Usuario registrado ", userId, " Inst: ", lastId)

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
			req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
			// Ejecutar petición
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

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

	// Claims de JWT para validar
	/*"uid":      idUser,
	"iid":      idInst,
	"authInst": authStatusInst,
	"team":     idTeam,
	"role":     roleApp,*/

	idUser := claims["uid"].(string)
	idInst := claims["iid"].(string)
	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
	}
	data, err := db.DB_con.GenericSelect("users", "idUser", []string{"activeUser", "roleAppUser_fk", "idTeam_fk", "idInstitution_fk"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener datos del usuario.", http.StatusInternalServerError)
		return
	}
	if len(data) > 0 {
		if data[idUser]["activeUser"] == "1" && data[idUser]["roleAppUser_fk"] == "1" && data[idUser]["idTeam_fk"] == "0" {
			//respWriter.Header().Set("Authorization", "Bearer "+cookie.Value)
			http.ServeFile(respWriter, request, "./../sffront/registro/validation.html")
		}
	} else {
		http.Error(respWriter, "Datos de usuario no encontrados", http.StatusUnauthorized)
		return
	}

}

// =======================================================================
// Endpoint para obtener datos de validación
func getValidationData(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
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

	// Obtener datos de la base de datos usando el user ID
	idUser := claims["uid"].(string)
	idInst := claims["iid"].(string)
	valResp := &models.ValidationResponse{}

	// obtener datos faltantes de la institucion
	var attrs = []string{"legalNameInst", "aliasNameInst", "taxNumInst", "legalSignupName", "legalSignupLastname", "streetAddress", "addressLine", "postalCode", "neighborhood", "locality"}

	// Condiciones WHERE
	wheres := map[string][]string{
		"idInstitution": {idInst},
		"rootUser_fk":   {idUser},
	}
	validationData, err := db.DB_con.GenericSelect("institutions", "idInstitution", attrs, wheres)
	if err != nil {
		http.Error(respWriter, fmt.Sprintf("Error obteniendo datos: %s", err), http.StatusInternalServerError)
		return
	}

	valData := map[string]*string{
		"legalNameInst":       &valResp.LegalName,
		"aliasNameInst":       &valResp.AliasName,
		"taxNumInst":          &valResp.TaxNum,
		"legalSignupName":     &valResp.LegalSignupName,
		"legalSignupLastname": &valResp.LegalSignupLastname,
		"streetAddress":       &valResp.StreetAddress,
		//"addressLine":         &valResp.AddressLine,
		"postalCode":   &valResp.PostalCode,
		"neighborhood": &valResp.Neighborhood,
		"locality":     &valResp.Locality,
	}
	for k, v := range validationData[idInst] {
		// Verificamos si esa clase está en nuestro mapa de punteros
		if ptr, exists := valData[k]; exists && ptr != nil {
			*ptr = v
		} else {
			fmt.Println("Existe: ", exists)
		}
	}

	// Mapeamos el nombre de la clase del documento al puntero del campo correspondiente
	documentClasses := map[string]*string{
		"ActaConstitutiva":   &valResp.ActaConstitutiva,
		"PoderRepresentante": &valResp.PoderRepresentante,
		"IdentidadOficial":   &valResp.IdentidadOficial,
		"PruebaResidencia":   &valResp.PruebaResidencia,
	}

	// Campos que queremos obtener de la tabla
	attrs = []string{"documentPath", "documentName", "documentExt", "documentClass", "expirationDate"}

	// Condiciones WHERE
	wheres = map[string][]string{
		"idInsttitution_fk": {idInst},
		"idUser_fk":         {idUser},
	}

	// Ejecutamos la consulta genérica
	validationDocs, err := db.DB_con.GenericSelect("kyc", "documentHash", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error obteniendo datos", http.StatusInternalServerError)
		return
	}

	// Recorremos cada registro (key = hash, value = fila)
	for _, row := range validationDocs {
		// Obtenemos la clase del documento
		if docClass, ok := row["documentClass"]; !ok {
			// Si no existe la columna, simplemente pasamos al siguiente registro
			continue
		} else {
			// Verificamos si esa clase está en nuestro mapa de punteros
			if ptr, exists := documentClasses[docClass]; exists && ptr != nil {
				// Concatenamos los valores que necesitamos.
				// Nos aseguramos de que cada clave exista antes de usarla.
				path := row["documentPath"]
				name := row["documentName"]
				ext := row["documentExt"]

				concatenated := fmt.Sprintf("%s%s.%s", path, name, ext)

				// Guardamos el resultado en el campo correspondiente de valResp
				*ptr = concatenated
			}
		}
	}

	// Configurar headers de seguridad
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")

	json.NewEncoder(respWriter).Encode(valResp)
}

// =======================================================================
// Endpoint para actualizar datos de validación
func updateValidationData(respWriter http.ResponseWriter, request *http.Request) {
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

	// Obtener user ID e institution ID
	idUser := claims["uid"].(string)
	idInst := claims["iid"].(string)
	valReq := models.ValidationRequest{}

	// Determinar el tipo de contenido
	contentType := request.Header.Get("Content-Type")

	// Si es multipart (contiene documentos)
	if strings.Contains(contentType, "multipart/form-data") {
		// Parsear el form multipart
		err := request.ParseMultipartForm(100 << 20) // 100 MB máximo
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

		// Procesar campos de texto
		if request.MultipartForm.Value != nil {
			if val, exists := request.MultipartForm.Value["streetAddress"]; exists && len(val) > 0 {
				valReq.StreetAddress = val[0]
			}
			if val, exists := request.MultipartForm.Value["postalCode"]; exists && len(val) > 0 {
				valReq.PostalCode = val[0]
			}
			if val, exists := request.MultipartForm.Value["neighborhood"]; exists && len(val) > 0 {
				valReq.Neighborhood = val[0]
			}
			if val, exists := request.MultipartForm.Value["locality"]; exists && len(val) > 0 {
				valReq.Locality = val[0]
			}
		}

		// Procesar archivo (si existe)
		if request.MultipartForm.File != nil {

			// Obtener el tipo de documento
			var docType string
			if docTypeVal, exists := request.MultipartForm.Value["documentType"]; exists && len(docTypeVal) > 0 {
				docType = docTypeVal[0]
			}

			// Buscar el archivo en el campo "document"
			if files, exists := request.MultipartForm.File["document"]; exists && len(files) > 0 {
				file, err := files[0].Open()
				if err != nil {
					http.Error(respWriter, "Error abriendo archivo", http.StatusBadRequest)
					return
				}
				defer file.Close()

				// Guardar el archivo y obtener la ruta
				filePath, err := utilities.GuardarArchivo(file, "", files[0].Filename, idInst, idUser, false)
				if err != nil {
					fmt.Println("Error guardando archivo:", err)
					http.Error(respWriter, "Error guardando archivo", http.StatusInternalServerError)
					return
				}

				// Asignar la ruta al campo correspondiente según el tipo de documento
				switch docType {
				case "docActa":
					valReq.ActaConstitutiva = filePath
				case "docPoder":
					valReq.PoderRepresentante = filePath
				case "docIdentidad":
					valReq.IdentidadOficial = filePath
				case "docResidencia":
					valReq.PruebaResidencia = filePath
				default:
					fmt.Printf("Tipo de documento desconocido: %s\n", docType)
				}
			}
		}
	} else {
		// Es JSON, solo datos de texto
		err = json.NewDecoder(request.Body).Decode(&valReq)
		if err != nil {
			http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
			return
		}
	}

	// Actualizar datos en la base de datos
	err = db.DB_con.UpdateValData(idInst, idUser, &valReq)
	if err != nil {
		fmt.Printf("Error actualizando datos en DB: %v\n", err)
		http.Error(respWriter, "Error actualizando datos", http.StatusInternalServerError)
		return
	}

	// Respuesta exitosa
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Datos actualizados correctamente"}
	json.NewEncoder(respWriter).Encode(response)
}

// =======================================================================
func completeValidation(respWriter http.ResponseWriter, request *http.Request) {
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
	type StatusData struct {
		IdInst string `json:"institution"`
		Status string `json:"status"`
	}
	statusData := StatusData{}
	err = json.NewDecoder(request.Body).Decode(&statusData)
	if err != nil {
		http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
		return
	}
	// Obtener user ID e institution ID
	// si la instutucion signforce hace la actualizacion de alguien entonces debe mandar id
	// si no manda id o la institucion no es signforce entonces se esta haciendo un update de ella misma
	fmt.Println(statusData)
	idInst := claims["iid"].(string)
	if idInst == "1" && statusData.IdInst != "" {
		idInst = statusData.IdInst
	}
	updates := map[string]map[string]interface{}{
		idInst: {
			"statusInst_fk": statusData.Status,
		},
	}
	err = db.DB_con.GenericBatchUpdate("institutions", "idInstitution", updates)
	if err != nil {
		http.Error(respWriter, "No se actualizó el estatus", http.StatusInternalServerError)
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")

	json.NewEncoder(respWriter).Encode(map[string]string{"status": "OK"})

}

func waitApprove(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/registro/waitapprove.html")
}

func uploadKeys(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(respWriter, request, "./../sffront/administracion/uploadkeys.html")
}

// Función para subir llaves
func uploadk(respWriter http.ResponseWriter, request *http.Request) {
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
	idUser, ok := claims["uid"].(string)
	if !ok {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	idInst, ok := claims["iid"].(string)
	if !ok {
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
		if Check {
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

	uploadKeysURL := os.Getenv("BACK_URL") + "uploadKeys"
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
	instData, err := db.DB_con.GenericSelect("institutions", "idInstitution", []string{"statusInst_fk", "typeContractInst"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de llaves", http.StatusInternalServerError)
		return
	}
	// si no tiene contrato aun y esta en el paso de subir llaves entonces es usuario nuevo y debe pasar a firma de contratos
	if instData[idInst]["statusInst_fk"] == "6" && instData[idInst]["typeContractInst"] == "0" {
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

// ========================================================================
// =======================================================================
func getDocs(respWriter http.ResponseWriter, request *http.Request) {
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

	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Estructura para recibir los datos del frontend
	type DocDataRequest struct {
		IdDoc   string `json:"id"`
		Type    string `json:"type"`
		PathDoc string `json:"path"`
	}

	// Leer el body de la petición
	var docRequest DocDataRequest
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&docRequest); err != nil {
		http.Error(respWriter, "Error al leer los datos de la petición", http.StatusBadRequest)
		return
	}

	// Validar datos del usuario
	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
		"idTeam_fk":        {idTeam},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", []string{"activeUser"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener información de usuario", http.StatusInternalServerError)
		return
	}
	if userData[idUser]["activeUser"] != "1" {
		http.Error(respWriter, "No autorizado. Usuario inactivo", http.StatusUnauthorized)
		return
	}

	// Determinar la tabla y ruta base según el tipo
	var tableName string
	var idColMain string
	var basePath string
	var attrs []string
	var docWheres map[string][]string
	switch docRequest.Type {
	case "kyc":
		splitted := strings.Split(docRequest.PathDoc, "@")
		hash := splitted[0]
		pathDoc := splitted[1]
		basePath = pathDoc[:strings.LastIndex(pathDoc, "/")+1]
		fll, _ := strings.CutPrefix(pathDoc, basePath)
		fullname := strings.Split(fll, ".")
		tableName = "kyc"
		attrs = append(attrs, "documentHash", "documentPath", "documentName", "documentExt")
		idColMain = "documentHash"
		docWheres = map[string][]string{
			"documentHash": {hash}, // Ajusta el nombre de la columna según tu esquema
			"documentPath": {basePath},
			"documentName": {fullname[0]},
			"documentExt":  {fullname[1]},
		}
	case "uploaded":
		tableName = "documents"
		attrs = append(attrs, "idDocument", "documentPath", "documentName", "documentExt", "documentHash", "authUseStatus", "authRoleStatus", "activeDoc", "sizeB", "abstractDoc")
		idColMain = "idDocument"
		docWheres = map[string][]string{
			"idDocument": {docRequest.IdDoc}, // Ajusta el nombre de la columna según tu esquema
		}
	default:
		http.Error(respWriter, "Tipo de documento no válido", http.StatusBadRequest)
		return
	}

	// Obtener datos del documento
	docData, err := db.DB_con.GenericSelect(tableName, idColMain, attrs, docWheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener información del documento", http.StatusInternalServerError)
		return
	}

	// Verificar si se encontró el documento
	if len(docData) == 0 {
		http.Error(respWriter, "Documento no encontrado", http.StatusNotFound)
		return
	}

	// Obtener el primer (y único) documento encontrado
	var docInfo map[string]string
	for _, doc := range docData {
		docInfo = doc
		break
	}
	switch docRequest.Type {
	case "uploaded":
		// Verificar permisos y estado del documento
		if docInfo["activeDoc"] != "1" {
			http.Error(respWriter, "Documento inactivo", http.StatusForbidden)
			return
		}

		if docInfo["authRoleStatus"] != "1" {
			http.Error(respWriter, "No autorizado para descargar este documento", http.StatusForbidden)
			return
		}

		if docInfo["authUseStatus"] != "1" {
			http.Error(respWriter, "Documento sin autorización de uso", http.StatusForbidden)
			return
		}
	case "kyc":
		if docInfo["expirationDate"] != "" {
			http.Error(respWriter, "Documento sin autorización de uso", http.StatusForbidden)
			return
		}
	}

	// Verificar permisos adicionales (opcional)
	// Puedes agregar lógica adicional aquí para verificar si el usuario tiene permisos
	// basándose en ownerInstDoc_fk, ownerTeamDoc_fk, creatorUserDoc_fk, etc.

	// Construir la ruta completa del archivo
	var filePath string = fmt.Sprintf("./%s%s.%s", docInfo["documentPath"], docInfo["documentName"], docInfo["documentExt"])
	fmt.Println(filePath)
	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(respWriter, "Archivo no encontrado en el sistema", http.StatusNotFound)
		return
	}

	// Obtener información del archivo
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.Error(respWriter, "Error al obtener información del archivo. Es posible que el recurso no exista", http.StatusNotFound)
		return
	}

	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(respWriter, "Error al acceder al archivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Detectar el tipo MIME del archivo
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		http.Error(respWriter, "Error al leer el archivo", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buffer)

	// Resetear el puntero del archivo al inicio
	file.Seek(0, 0)

	// Configurar headers de respuesta
	respWriter.Header().Set("Content-Type", contentType)
	respWriter.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	respWriter.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(filePath)))

	// Copiar el contenido del archivo a la respuesta
	_, err = io.Copy(respWriter, file)
	if err != nil {
		// No podemos usar http.Error aquí porque ya hemos comenzado a escribir la respuesta
		log.Printf("Error al enviar archivo: %v", err)
		return
	}
}

// =====================================================================================
// Funcion gneérica para subir archivos de cualquier clase
func uploadDocs(respWriter http.ResponseWriter, request *http.Request) {
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
	idUser, ok := claims["uid"].(string)
	if !ok {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	idInst, ok := claims["iid"].(string)
	if !ok {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	idTeam, ok := claims["team"].(string)
	if !ok {
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
		maxmem = 500
	}
	maxmem = maxmem << 20

	err = request.ParseMultipartForm(int64(maxmem)) // 500 MB máxim
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

	var processedFiles []models.FileMetadata
	var docIds []string

	if request.MultipartForm.File != nil {
		var docType string

		if docTypeVal, exists := request.MultipartForm.Value["documentExt"]; exists && len(docTypeVal) > 0 {
			docType = docTypeVal[0]
		}

		// Buscar el archivo en el campo "document"
		if files, exists := request.MultipartForm.File["document"]; exists && len(files) > 0 {
			for i, fileHeader := range files {
				// nuevo archivo
				var file multipart.File

				if fileHeader.Size > int64(maxmem) {
					http.Error(respWriter, "Memoria insuficiente para procesar archivos", http.StatusNotAcceptable)
				}

				docData := models.FileMetadata{
					Id:      i,
					DocType: docType,
					Name:    fileHeader.Filename,
					Ext:     strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:]),
					Size:    fileHeader.Size,
					Hash:    "",
					Path:    "",
					Ok:      false,
				}
				if !configs.AllowedExtensions[docData.Ext] {
					fmt.Println("Extension no aceptada")
					processedFiles = append(processedFiles, docData)
					continue
				}

				switch docType {
				case "generic":
					docData.Path = fmt.Sprintf("%s/%s/%s/%s", os.Getenv("GENERIC_DOC_PATH"), idInst, idTeam, idUser)
				case "template":
					docData.Path = fmt.Sprintf("%s/%s/%s/templates", os.Getenv("GENERIC_DOC_PATH"), idInst, idUser)
				default:
					docData.Path = fmt.Sprintf("%s/%s/%s/%s", os.Getenv("TEMP_BASE_PATH"), idInst, idTeam, idUser)
				}

				file, err = fileHeader.Open()
				if err != nil {
					http.Error(respWriter, fmt.Sprintf("error abriendo archivo %s: %v", fileHeader.Filename, err), http.StatusBadRequest)
					return
				}
				defer file.Close()

				var filePath string
				// Guardar el archivo y obtener la ruta /ruta_general/idInst/idTeam/idUser/midoc.pdf
				filePath, err = utilities.GuardarArchivo(
					file, // archivo completo
					docData.Path,
					docData.Name,
					idInst,
					idUser,
					false,
				)
				if err != nil {
					fmt.Println("Error guardando archivo:", err)
					http.Error(respWriter, "Error guardando archivo", http.StatusInternalServerError)
					return
				}

				docData.Hash, err = utilities.GetHash(filePath, configs.HashConf)
				if err != nil {
					fmt.Println("Error obteniendo hash de archivo:", err)
					http.Error(respWriter, "Error guardando archivo", http.StatusInternalServerError)
					return
				}

				// Actualizar datos en la base de datos
				cols := []string{"documentHash", "ownerInstDoc_fk", "ownerTeamDoc_fk", "creatorUserDoc_fk", "documentName", "documentPath", "documentExt"}
				vals := []interface{}{docData.Hash, idInst, idTeam, idUser, docData.Name, docData.Path, docData.Ext}
				idDoc, err := db.DB_con.GenericInsert("documents", cols, vals)
				if err != nil {
					fmt.Printf("Error actualizando datos en DB: %v\n", err)
					http.Error(respWriter, "Error actualizando datos", http.StatusInternalServerError)
					return
				}
				processedFiles = append(processedFiles, docData)
				docIds = append(docIds, idDoc)
			}
		}
	}

	if len(processedFiles) == 0 {
		http.Error(respWriter, "Error procesando archivos", http.StatusUnprocessableEntity)
		return
	}

	response := models.UploadResponse{
		Success:     true,
		Message:     fmt.Sprintf("Se subieron %d archivo(s) exitosamente", len(processedFiles)),
		DocumentIDs: docIds,
	}
	// Respuesta exitosa
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)

	json.NewEncoder(respWriter).Encode(&response)
}

// =======================================================================
func processPayment(respWriter http.ResponseWriter, request *http.Request) {
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

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	payReq := models.PaymentReq{}
	err = json.NewDecoder(request.Body).Decode(&payReq)
	if err != nil {
		http.Error(respWriter, "Error en los datos", http.StatusBadRequest)
		return
	}
	// borrar cuando ya no se requiera
	fmt.Println(idUser, idInst, idTeam, authInst)

	attrs := []string{"statusPayment_fk", "planId_fk", "expirationPlan", "nameOwner"}
	wheres := map[string][]string{"idInstitution": {idInst}}
	data, err := db.DB_con.GenericSelect("payment", "idInstitution", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de institucion", http.StatusInternalServerError)
		return
	}

	resp := &models.PaymentResp{}

	if len(data) == 0 {
		var plan string
		switch payReq.Plan {
		case "Basico":
			plan = "1"
		case "Profesional":
			plan = "2"
		case "Empresarial":
			plan = "3"
		}

		// aqui se realiza un request a la pasarela de pago que debe regresar un estatus el cual se inserta y se devuelve a la plataforma
		// debe ser una go routine que se encargue de actualizar el estatus en DB independiente del endpoint
		payStatus := "1" // 0 PENDING, 1 COMPLETED, 2 REJECTED, 3 CANCELED, 4 HOLD
		expPlan := time.Now().AddDate(0, 1, 0).Format("2006-01-02 15:04:05")
		cols := []string{"statusPayment_fk", "planId_fk", "expirationPlan", "cardNumber", "expirationDate", "nameOwner"}
		vals := []interface{}{payStatus, plan, expPlan, payReq.CardNum, payReq.Expiration, payReq.NameOwner}
		idPay, err := db.DB_con.GenericInsert("payment", cols, vals)
		if err != nil {
			fmt.Println("Error info pago")
			http.Error(respWriter, "Error al insertar informacion de pago", http.StatusInternalServerError)
		}
		updates := map[string]map[string]interface{}{
			idInst: {
				"statusInst_fk":      "6",
				"paymentDataInst_fk": idPay,
			},
		}

		err = db.DB_con.GenericBatchUpdate("institutions", "idInstitution", updates)
		if err != nil {
			http.Error(respWriter, "Error al actualizar informacion de institucion", http.StatusInternalServerError)
		}
		resp.CurrentSatat = plan
		resp.Success = true
		resp.Status = payStatus
		resp.Plan = plan
		resp.Expiration = expPlan
	} else {
		d := data[idInst]
		resp.CurrentSatat = d["statusPayment_fk"]
		resp.Success = true
		resp.Status = d["statusPayment_fk"]
		resp.Plan = d["planId_fk"]
		resp.Expiration = d["expirationPlan"]
	}
	// Convertir a JSON
	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(respWriter, "Error al generar JSON", http.StatusInternalServerError)
		return
	}

	// Configurar headers y enviar respuesta
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(jsonData)
}

// =======================================================================
func checkUserStatus(respWriter http.ResponseWriter, request *http.Request) {
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
	fmt.Println(claims)
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

// =======================================================================
func getKeysData(respWriter http.ResponseWriter, request *http.Request) {
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
	fmt.Println(claims)
	// Extraer datos del JWT

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	idTeam, ok3 := claims["team"].(string)

	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
		"idTeam_fk":        {idTeam},
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
	keys, err := db.DB_con.GenericSelect("userkeys", "idUserKeys", []string{"keyFilePath", "certFilePath", "notValidAfter", "subjectRFC4514", "subjectUniqueId", "createdAtKey"}, wheres)
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

func updateKeysData(respWriter http.ResponseWriter, request *http.Request) {
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

// =======================================================================
// Endpoint para obtener datos de validación
func approvals(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
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

	// Obtener datos de la base de datos usando el user ID
	idUser := claims["uid"].(string)
	idInst := claims["iid"].(string)
	if idInst != "1" && idUser == "0" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// obtener datos faltantes de la institucion
	var attrs = []string{"legalNameInst", "aliasNameInst", "taxNumInst", "legalSignupName", "legalSignupLastname", "streetAddress", "addressLine", "postalCode", "neighborhood", "locality"}

	// Para que una empresa sea aprovada debe haber subido sus documentos y o estar en estatus de rechazo de documentos y asi mismo debe estar inactivo
	wheres := map[string][]string{
		"statusInst_fk": {"3", "4"},
		"activeInst":    {"0"},
	}
	validationData, err := db.DB_con.GenericSelect("institutions", "idInstitution", attrs, wheres)
	if err != nil {
		http.Error(respWriter, fmt.Sprintf("Error obteniendo datos: %s", err), http.StatusInternalServerError)
		return
	}

	arrayResp := make(map[string]*models.ValidationResponse)
	for instId, instData := range validationData {
		valResp := &models.ValidationResponse{
			LegalName:           instData["legalNameInst"],
			AliasName:           instData["aliasNameInst"],
			TaxNum:              instData["taxNumInst"],
			LegalSignupName:     instData["legalSignupName"],
			LegalSignupLastname: instData["legalSignupLastname"],
			StreetAddress:       instData["streetAddress"],
			//AddressLine,
			PostalCode:   instData["postalCode"],
			Neighborhood: instData["neighborhood"],
			Locality:     instData["locality"],
		}

		// Mapeamos el nombre de la clase del documento al puntero del campo correspondiente
		documentClasses := map[string]*string{
			"ActaConstitutiva":   &valResp.ActaConstitutiva,
			"PoderRepresentante": &valResp.PoderRepresentante,
			"IdentidadOficial":   &valResp.IdentidadOficial,
			"PruebaResidencia":   &valResp.PruebaResidencia,
		}

		// Campos que queremos obtener de la tabla
		attrs = []string{"documentClass", "documentPath", "documentName", "documentExt"}

		// Condiciones WHERE
		wheres = map[string][]string{
			"idInsttitution_fk": {instId},
		}

		// Ejecutamos la consulta genérica
		validationDocs, err := db.DB_con.GenericSelect("kyc", "documentHash", attrs, wheres)
		if err != nil {
			http.Error(respWriter, "Error obteniendo datos", http.StatusInternalServerError)
			return
		}

		// Recorremos cada registro (key = hash, value = fila)
		for docHash, row := range validationDocs {
			// Obtenemos la clase del documento
			if docClass, ok := row["documentClass"]; !ok {
				// Si no existe la columna, pasamos al siguiente registro
				continue
			} else {
				// Verificamos si esa clase está en nuestro mapa de punteros
				if ptr, exists := documentClasses[docClass]; exists && ptr != nil {
					// Concatenamos los valores que necesitamos.
					// Nos aseguramos de que cada clave exista antes de usarla.
					path := row["documentPath"]
					name := row["documentName"]
					ext := row["documentExt"]

					concatenated := fmt.Sprintf("%s@%s%s.%s", docHash, path, name, ext)

					// Guardamos el resultado en el campo correspondiente de valResp
					*ptr = concatenated
				}
			}
		}
		arrayResp[instId] = valResp
	}

	// Configurar headers de seguridad
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")

	json.NewEncoder(respWriter).Encode(arrayResp)
}

// =======================================================================
// unificar el json de respuestas para que mande estatus y lista de documentos
// asegurar que multipart puede recibir uno o muchos archivos subidos de un mismo formulario y sugerir mejoras para subir archivos de distinta ubicacion
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

	// se obtienen los datos de validacion del usuario
	var attrs = []string{"emailUser", "appPassHash", "activeUser", "idTeam_fk", "roleAppUser_fk", "idInstitution_fk"}
	dataUser, err := db.DB_con.GenericSelect("users", "idUser", attrs, map[string][]string{"emailUser": {loginReq.Account}})
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
		nh, err := utilities.GetHash([]byte(loginReq.Password), configs.HashConf)
		if err != nil {
			loginResp.Error = "Error: No se pudo verificar el password"
			json.NewEncoder(respWriter).Encode(loginResp)
			return
		}
		if !(v["appPassHash"] == nh && v["activeUser"] == "1") {
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
	token, err := auth.GenerateJWT(idUser, dataUser[idUser]["idTeam_fk"], dataUser[idUser]["roleAppUser_fk"], idInst, dataInst[idInst]["statusInst_fk"])
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
				location = "/uploadKeys"
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
