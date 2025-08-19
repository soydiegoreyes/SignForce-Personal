package main

import (
	"crypto/tls"
	"emailServ/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Cargar las variables de entorno desde .env
func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}
}

// EmailData representa la estructura de datos para el envío de emails
type EmailData struct {
	IdUser   string   `json:"iduser"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	Dest     []string `json:"dest"`
	MimeType string   `json:"mimetype"`
}

// Estructura para cliente
type EmailClient struct {
	SMTPsender   string
	SMTPuser     string
	SMTPpassword string
	SMTPServer   string
	tlsConfig    *tls.Config
	conn         *tls.Conn
	client       *smtp.Client
}

func (cl *EmailClient) Connect() bool {
	var check bool

	if cl.SMTPsender == "" || cl.SMTPuser == "" || cl.SMTPpassword == "" || cl.SMTPServer == "" {
		log.Println("Faltan variables de entorno necesarias")
	} else {
		log.Printf("Conectando con %s", cl.SMTPServer)

		// Configuración TLS
		cl.tlsConfig = &tls.Config{
			InsecureSkipVerify: true, // En producción deberías validar el certificado
			ServerName:         cl.SMTPServer,
		}

		// Conectar al servidor SMTP
		conn, err := tls.Dial("tcp", cl.SMTPServer+":465", cl.tlsConfig)
		if err != nil {
			log.Println("Error al conectar al servidor SMTP:", err)
		} else {
			cl.conn = conn
		}

		// Crear cliente SMTP
		client, err := smtp.NewClient(conn, cl.SMTPServer)
		if err != nil {
			log.Println("Error al crear cliente SMTP:", err)
		} else {
			cl.client = client
		}

		// Autenticación
		auth := smtp.PlainAuth("", cl.SMTPuser, cl.SMTPpassword, cl.SMTPServer)
		if err := cl.client.Auth(auth); err != nil {
			log.Println("Error de autenticación:", err)
		}

		if err == nil {
			log.Println("Conectado al servidor SMTP")
			check = true
		}
	}

	return check
}

func (cl *EmailClient) CloseAll() {

	if cl.client != nil {
		// Cerrar la conexión correctamente
		cl.client.Quit()
		//cl.client.Close()
	} else {
		log.Println("cliente de email no existe.")
	}
	if cl.conn != nil {
		cl.conn.Close()
	} else {
		log.Println("conexion de email no existe.")
	}
}

func (cl *EmailClient) EmailSender(idApp string, data EmailData) bool {
	log.Println("Intento de envío de email")
	if len(data.Dest) == 0 || data.Body == "" || data.IdUser == "" {
		log.Println("Datos insuficientes para envio de email")
		return false
	}

	// Verificar si la conexión SMTP sigue viva
	if cl.client == nil {
		log.Println("Cliente SMTP no inicializado, conectando...")
		if !cl.Connect() {
			return false
		}
	} else {
		if err := cl.client.Noop(); err != nil {
			log.Println("Conexión SMTP caducada, reconectando...")
			cl.CloseAll()
			if !cl.Connect() {
				return false
			}
		}
	}

	// se obtiene el nombre y correo del usuario emisor y se valida que sea un usuario activo
	var userAttributes = []string{"nameUser", "lastNameUser", "emailUser", "activeUser"}
	var wheres = map[string][]string{
		"idUser": {data.IdUser},
	}
	userValues, err := db.DB_con.GenericSelect("users", "idUser", userAttributes, wheres)
	if err != nil {
		fmt.Println("error: no se pudo obtener datos del usuario ", data.IdUser)
		return false
	}
	if userValues[data.IdUser]["activeUser"] != "1" {
		fmt.Println("error: Usuario inactivo: ", data.IdUser)
		return false
	}

	// se construye el complemento del cuerpo indicando quien envía ese correo
	data.Body = fmt.Sprintf("Correo enviado de %s %s\n\n%s\n",
		userValues[data.IdUser]["nameUser"], userValues[data.IdUser]["lastNameUser"], data.Body)

	recipients := data.Dest

	// Configurar los headers del email
	headers := make(map[string]string)
	headers["From"] = cl.SMTPsender
	headers["To"] = strings.Join(recipients, ",")
	headers["Subject"] = data.Subject

	// Determinar el Content-Type según el mime_type
	contentType := "text/plain"
	if data.MimeType == "html" {
		contentType = "text/html"
	}
	headers["Content-Type"] = contentType + "; charset=\"UTF-8\""

	// Construir el mensaje
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n" + data.Body)

	// Establecer remitente
	if err := cl.client.Mail(cl.SMTPsender); err != nil {
		log.Println("Error al establecer remitente:", err)
		return false
	}
	// Establecer destinatario
	for _, recipient := range recipients {
		if err := cl.client.Rcpt(recipient); err != nil {
			log.Println("Error al establecer destinatario:", recipient, err)
			return false
		}
	}

	// Enviar el cuerpo del email
	w, err := cl.client.Data()
	if err != nil {
		log.Println("Error al preparar el cuerpo del email:", err)
		return false
	}
	_, err = w.Write([]byte(message.String()))
	if err != nil {
		log.Println("Error al escribir el cuerpo del email:", err)
		return false
	}
	err = w.Close()
	if err != nil {
		log.Println("Error al cerrar el escritor del cuerpo del email:", err)
		return false
	}

	// insertar registro en base de datos
	attributes := []string{"idUserSender_fk", "reason", "subjectEmail", "bodyEmail", "idAppSource_fk"}
	values := []interface{}{data.IdUser, data.Subject, cl.SMTPsender, message.String(), idApp}

	row, err := db.DB_con.InsertEmailAttributes(attributes, values)
	if err != nil {
		log.Println("Error al insertar en base de datos:", err)
		return false
	}
	log.Println("Email enviado exitosamente: ", row)
	return true
}

func NewEmailClient() *EmailClient {
	return &EmailClient{
		SMTPsender:   os.Getenv("EMAIL_DIR"),
		SMTPuser:     os.Getenv("EMAIL_USER"),
		SMTPpassword: os.Getenv("EMAIL_PASS"),
		SMTPServer:   os.Getenv("SMTP_SERVER"),
	}
}

// ============================================
var Eclient *EmailClient = NewEmailClient()

// Ejemplo de uso
func main() {
	loadEnv()
	Eclient.Connect()
	defer Eclient.CloseAll()
	http.HandleFunc("/mailserv", mailServ)

	// Inicia el servidor
	if err := http.ListenAndServe(":"+os.Getenv("API_PORT"), nil); err != nil {
		log.Fatal("Error en el servidor:", err)
	}
}

// Manejador para el envío de emails
func mailServ(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var data EmailData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error en los datos: "+err.Error(), http.StatusBadRequest)
		return
	}
	// se obtiene el nombre de la app y si esta activa
	appAttributes := []string{"isActive", "domainApp", "portApp", "currPathApp", "nameApp"}
	var wheres = map[string][]string{
		"idapp": []string{token},
	}
	appValues, err := db.DB_con.GenericSelect("microapps", "idapp", appAttributes, wheres)
	if err != nil {
		fmt.Println("error: no se pudo autenticar la app", token)
		return
	}
	//fmt.Println(userValues)
	if appValues[token]["isActive"] != "1" {
		fmt.Println("error: app inactiva: ", token)
		return
	}

	success := Eclient.EmailSender(token, data)
	if !success {
		http.Error(w, "Error al enviar el email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email enviado con éxito")
}
