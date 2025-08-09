package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"

	//"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/objects"

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
		w.Header().Set("Access-Control-Allow-Origin", "*") // En producción cambia por tu dominio específico
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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

	// Registrar rutas
	mux.HandleFunc("/register", registerInst)

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

// Handler HTTP para loguear a un usuario por un ID de usuario y un arreglo de atributos a adquirir
func registerInst(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var registerReq models.RegisterRequest
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
	data, err := db.DB_con.GenericSelect("institutions", "idInstitution", "taxNumInst", []string{registerReq.TaxNumInst}, []string{"statusInst_fk", "activeInst", "typeContractInst", "legalSignupName", "legalSignupLastname"})
	if err != nil {
		json.NewEncoder(respWriter).Encode(models.RegisterResponse{Check: false, InstId: "", Error: "Critical: error al obtener datos de institucion"})
		return
	} else {
		// en caso de no haber ningun registro con el mismo TAXNUMBER  se procede al registro
		if len(data) == 0 || status[data["idInstitution"]["statusInst_fk"]] {
			lastId, err := objects.RegisterInst(&registerReq)
			if err != nil {
				json.NewEncoder(respWriter).Encode(models.RegisterResponse{Check: false, InstId: "", Error: fmt.Sprintf("%s", err)})
				return
			}

			json.NewEncoder(respWriter).Encode(models.RegisterResponse{Check: true, InstId: lastId, Error: ""})
		} else {
			json.NewEncoder(respWriter).Encode(models.RegisterResponse{Check: false, InstId: "", Error: "Ya tiene un registro para su numero de empresa. Revisar estatus de su registro."})
		}
	}
}
