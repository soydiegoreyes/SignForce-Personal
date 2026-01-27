package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sfmiddle/auth"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/iapackage"
	"sfmiddle/models"
	"strconv"
	"strings"

	"sfmiddle/utilities"

	"github.com/google/uuid"
)

// =====================================================================================
// Funcion gneérica para subir archivos de cualquier clase
func UploadDocs(respWriter http.ResponseWriter, request *http.Request) {
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
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
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
		var docType string = request.MultipartForm.Value["documentType"][0]

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
					Name:    fileHeader.Filename[:strings.LastIndex(fileHeader.Filename, ".")],
					Ext:     strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:]),
					Size:    fileHeader.Size,
					Hash:    "",
					Path:    "",
					Ok:      false,
				}
				if !configs.AllowedExtensions[strings.ToLower(docData.Ext)] {
					fmt.Println("Extension no aceptada")
					processedFiles = append(processedFiles, docData)
					continue
				}
				file, err = fileHeader.Open()
				if err != nil {
					http.Error(respWriter, fmt.Sprintf("error abriendo archivo %s: %v", fileHeader.Filename, err), http.StatusBadRequest)
					return
				}
				defer file.Close()
				docBytes, err := io.ReadAll(file)
				if err != nil {
					http.Error(respWriter, fmt.Sprintf("error leyendo archivo %s: %v", fileHeader.Filename, err), http.StatusBadRequest)
				}
				// se manda el filepath para comprobar que se guardó el archivo
				docData.Hash, err = utilities.GetHash(docBytes, configs.HashConf)
				if err != nil {
					fmt.Println("Error obteniendo hash de archivo:", err)
					http.Error(respWriter, "Error guardando archivo", http.StatusInternalServerError)
					return
				}

				var tablename string
				var cols []string
				var vals []interface{}
				switch docType {
				case "generic":
					docData.Path = fmt.Sprintf("%s/%s/%s/", os.Getenv("GENERIC_DOC_PATH"), idInst, idUser)
					tablename = "documents"
					cols = []string{"documentHash", "ownerInstDoc_fk", "creatorUserDoc_fk", "documentName", "documentPath", "documentExt", "sizeB", "authUseStatus", "authRoleStatus"}
					vals = []interface{}{docData.Hash, idInst, idUser, docData.Name, docData.Path, docData.Ext, docData.Size, "1", "1"}

				case "template":
					docData.Path = fmt.Sprintf("%s/%s/%s/", os.Getenv("TEMPLATE_DOC_PATH"), idInst, idUser)
					tablename = "doctemplates"
					cols = []string{"ownerInstDoc_fk", "creatorUserDoc_fk", "documentName", "documentPath", "documentExt", "sizeB", "authUseStatus", "authRoleStatus"}
					vals = []interface{}{idInst, idUser, docData.Name, docData.Path, docData.Ext, docData.Size, "1", "1"}

				case "logoinst":
					docData.Path = fmt.Sprintf("%s/%s/logos/", os.Getenv("GENERIC_DOC_PATH"), idInst)
					fileHeader.Filename = uuid.NewString() + "_" + fileHeader.Filename
					tablename = "images"

				case "imgkyc":
					docData.Path = fmt.Sprintf("%s/%s/%s/kyc/", os.Getenv("GENERIC_IMG_PATH"), idInst, idUser)
					tablename = "images"
					cols = []string{"imageHash", "ownerInstImg_fk", "creatorUserImg_fk", "imageName", "imagePath", "imageExt", "sizeB", "authUseStatus", "authRoleStatus"}
					vals = []interface{}{docData.Hash, idInst, idUser, docData.Name, docData.Path, docData.Ext, docData.Size, "1", "1"}

				default:
					fmt.Println("Tipo de documento desconocido")
					return
				}

				// Regreso al inicio del documento
				_, err = file.Seek(0, io.SeekStart)
				if err != nil {
					http.Error(respWriter, "Error reposicionando archivo", http.StatusInternalServerError)
					return
				}
				var filePath string
				filePath, err = utilities.GuardarArchivo(
					file, // archivo completo
					docData.Path,
					fileHeader.Filename,
					idInst,
					idUser,
					false,
				)
				if err != nil || filePath == "" {
					fmt.Println("Error guardando archivo:", err)
					http.Error(respWriter, "Error guardando archivo", http.StatusInternalServerError)
					return
				}

				idDoc, err := db.DB_con.GenericInsert(tablename, cols, vals)
				if err != nil {
					fmt.Printf("Error actualizando datos en DB: %v\n", err)
					http.Error(respWriter, "Error actualizando datos", http.StatusInternalServerError)
					return
				}
				processedFiles = append(processedFiles, docData)
				docIds = append(docIds, idDoc)
				go iapackage.GetAbstractDoc(idDoc)
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

// ====================================================================================================
func DownloadDoc(respWriter http.ResponseWriter, request *http.Request) {
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
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
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

	// variables para datos del usuario y documento
	var docInfo, userInfo map[string]string

	// Validar datos del usuario
	wheres := map[string][]string{
		"idUser":           {idUser},
		"idInstitution_fk": {idInst},
		//"idTeam_fk":        {idTeam},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", []string{"activeUser"}, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener información de usuario", http.StatusInternalServerError)
		return
	}
	// info del usuario
	userInfo = userData[idUser]

	if userInfo["activeUser"] != "1" {
		http.Error(respWriter, "No autorizado. Usuario inactivo", http.StatusUnauthorized)
		return
	}

	// Determinar la tabla y ruta base según el tipo
	var tableName string
	var idColMain string
	var attrs []string
	var docWheres map[string][]string
	switch docRequest.Type {
	case "kyc":
		t := docRequest.PathDoc
		hash := t[:strings.Index(t, "@")]
		pathDoc, _ := strings.CutPrefix(t, hash+"@")
		lif := strings.LastIndex(pathDoc, "/")
		lip := strings.LastIndex(pathDoc, ".")
		basePath := pathDoc[:lif+1]
		ext := pathDoc[lip+1:]
		nameDoc := pathDoc[lif+1 : lip]
		tableName = "kyc"
		attrs = append(attrs, "documentHash", "documentPath", "documentName", "documentExt", "expirationDate")
		idColMain = "documentHash"
		docWheres = map[string][]string{
			"documentHash": {hash}, // Ajusta el nombre de la columna según tu esquema
			"documentPath": {basePath},
			"documentName": {nameDoc},
			"documentExt":  {ext},
		}
	case "uploaded":
		tableName = "documents"
		attrs = append(attrs, "idDocument", "documentPath", "documentName", "documentExt", "documentHash", "sizeB", "abstractDoc")
		idColMain = "idDocument"
		docWheres = map[string][]string{
			"idDocument": {docRequest.IdDoc}, // Ajusta el nombre de la columna según tu esquema
		}
	case "template":
		tableName = "doctemplates"
		attrs = append(attrs, "idTemplate", "documentPath", "documentName", "documentExt", "sizeB", "abstractDoc")
		idColMain = "idTemplate"
		docWheres = map[string][]string{
			"idTemplate": {docRequest.IdDoc}, // Ajusta el nombre de la columna según tu esquema
		}
	default:
		http.Error(respWriter, "Tipo de documento no válido", http.StatusBadRequest)
		return
	}
	// atributos obligatorios
	if docRequest.Type == "uploaded" || docRequest.Type == "template" {
		attrs = append(attrs, "ownerInstDoc_fk", "creatorUserDoc_fk", "authUseStatus", "authRoleStatus", "activeDoc")
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
	for _, doc := range docData {
		docInfo = doc
		break
	}

	switch docRequest.Type {
	case "kyc":
		if docInfo["expirationDate"] != "" {
			http.Error(respWriter, "Documento sin autorización de uso", http.StatusForbidden)
			return
		}
	default:
		// validaciones de permisos y estados de documentos
		if !auth.ValidateUserDocPermissions(idInst, idUser, docInfo, userInfo) {
			http.Error(respWriter, "No autorizado", http.StatusForbidden)
			return
		}
	}

	// Construir la ruta completa del archivo
	var filePath string = fmt.Sprintf("%s/%s%s.%s", os.Getenv("BASE_DIR"), docInfo["documentPath"], docInfo["documentName"], docInfo["documentExt"])

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
	fmt.Println(fileInfo)
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
	fmt.Println(contentType)

	// Resetear el puntero del archivo al inicio
	file.Seek(0, 0)

	// Configurar headers de respuesta
	respWriter.Header().Set("Content-Type", contentType)
	respWriter.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	respWriter.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(filePath)))

	// Copiar el contenido del archivo a la respuesta
	n, err := io.Copy(respWriter, file)
	if err != nil {
		// No podemos usar http.Error aquí porque ya hemos comenzado a escribir la respuesta
		fmt.Printf("Error al enviar archivo: %v", err)
		return
	}
	fmt.Println("numero de bytes ->", n)
}

// ====================================================================================================
// Funcion para devolver datos de uno o varios documentos (no devuelve el documento)
func StatusDocs(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== Autenticación por JWT =====
	cookie, err := request.Cookie("token")
	if err != nil {
		http.Error(respWriter, "No autorizado"+err.Error(), http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(respWriter, "No autorizado"+err.Error(), http.StatusUnauthorized)
		return
	}

	_, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====

	var req models.DocDataRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	var tablename, idCol string
	var attrs = []string{
		"createdAtDoc", "lastModifiedDoc", "deletedAtDoc",
		"documentPath", "documentName", "documentExt", "sizeB",
		"abstractDoc", "authUseStatus", "authRoleStatus", "activeDoc",
	}

	switch req.Type {
	case "docx":
		tablename = "doctemplates"
		idCol = "idTemplate"
	case "pdf":
		tablename = "documents"
		idCol = "idDocument"
		attrs = append([]string{"documentHash", "deletedReasonDoc"}, attrs...)
	default:
		fmt.Println("No debería ocurrir esto.")
		return
	}

	// ===== Parámetros de paginación =====
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "createdAtDoc"
	}
	orderDir := strings.ToUpper(req.OrderDir)
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "DESC"
	}

	// ===== Construcción de filtros =====
	wheres := map[string][]string{
		//"creatorUserDoc_fk": {idUser},
		"ownerInstDoc_fk": {idInst},
	}

	logic := ""

	// Si hay filtros por id, tipo o path -> se priorizan
	switch {
	case len(req.IdDocs) != 0:
		//logic = "creatorUserDoc_fk AND ownerInstDoc_fk AND idDocument"
		logic = "ownerInstDoc_fk AND idDocument"
		wheres[idCol] = req.IdDocs

	case len(req.PathDocs) != 0:
		//logic = "creatorUserDoc_fk AND ownerInstDoc_fk AND documentPath"
		logic = "ownerInstDoc_fk AND documentPath"
		wheres["documentPath"] = req.PathDocs

	case req.Type != "":
		//logic = "creatorUserDoc_fk AND ownerInstDoc_fk AND documentExt"
		logic = "ownerInstDoc_fk AND documentExt"
		wheres["documentExt"] = []string{req.Type}

	case req.DateFrom != "" && req.DateTo != "":
		//logic = "creatorUserDoc_fk AND ownerInstDoc_fk AND createdAtDoc BETWEEN ORDER BY " + orderBy
		logic = "ownerInstDoc_fk AND createdAtDoc BETWEEN ORDER BY " + orderBy
		wheres["createdAtDoc"] = []string{req.DateFrom, req.DateTo}
		wheres[idCol] = []string{orderDir}
		wheres["LOGIC"] = []string{logic, fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)}

	default:
		// Caso general: sin filtros, solo paginación
		//logic = "creatorUserDoc_fk AND ownerInstDoc_fk ORDER BY " + orderBy
		logic = "ownerInstDoc_fk ORDER BY " + orderBy
		wheres[idCol] = []string{orderDir}
		wheres["LOGIC"] = []string{logic, fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)}
	}

	// ===== Ejecutar SELECT =====
	docData, err := db.DB_con.GenericSelect(tablename, idCol, attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener documentos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ===== Obtener total de documentos =====
	//baseWhere := fmt.Sprintf("creatorUserDoc_fk = '%s' AND ownerInstDoc_fk = '%s'", idUser, idInst)
	baseWhere := fmt.Sprintf("ownerInstDoc_fk = '%s'", idInst)

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s;", tablename, baseWhere)
	var total int
	err = db.DB_con.DB.QueryRow(query).Scan(&total)
	if err != nil {
		total = -1
	}

	// ===== Respuesta JSON =====
	resp := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"data":      docData,
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(resp)
}

func EditDoc(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== Autenticación por JWT =====
	cookie, err := request.Cookie("token")
	if err != nil {
		http.Error(respWriter, "No autorizado"+err.Error(), http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil {
		http.Error(respWriter, "No autorizado"+err.Error(), http.StatusUnauthorized)
		return
	}

	idUser, ok1 := claims["uid"].(string)
	idInst, ok2 := claims["iid"].(string)
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====
	var req models.FileMetadata // tiene lo necesario para los datos del documento y validar contra base de datos
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		http.Error(respWriter, "Error al leer la petición", http.StatusBadRequest)
		return
	}

	var tablename, idCol string
	idDoc := strconv.Itoa(req.Id)

	var attrs = []string{
		"createdAtDoc", "lastModifiedDoc", "deletedAtDoc",
		"documentPath", "documentName", "documentExt", "sizeB",
		"abstractDoc", "authUseStatus", "authRoleStatus", "activeDoc",
	}
	wheres := map[string][]string{
		"creatorUserDoc_fk": {idUser},
		"ownerInstDoc_fk":   {idInst},
	}

	switch req.Ext {
	case "docx":
		tablename = "doctemplates"
		idCol = "idTemplate"

		wheres["idTemplate"] = []string{idDoc}

	case "pdf":
		tablename = "documents"
		idCol = "idDocument"
		attrs = append([]string{"documentHash", "deletedReasonDoc"}, attrs...)
		wheres["idDocument"] = []string{idDoc}
	default:
		fmt.Println("No debería ocurrir esto.")
		return
	}

	// ===== Ejecutar SELECT =====
	docData, err := db.DB_con.GenericSelect(tablename, idCol, attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener documentos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if docData[idDoc]["activeDoc"] == "1" {

	}

}

func InteractDoc(respWriter http.ResponseWriter, request *http.Request) {

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
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====
	type docQuery struct {
		IdDocument string `json:"idDocument"`
		Prompt     string `json:"prompt"`
		Type       string `json:"type"`
	}
	var dq docQuery
	err = json.NewDecoder(request.Body).Decode(&dq)
	if err != nil {
		http.Error(respWriter, "Error al decodificar json", http.StatusBadRequest)
	}
	resp := iapackage.InteractDoc(dq.IdDocument, idInst, idUser, dq.Prompt, dq.Type)
	if resp == nil {
		http.Error(respWriter, "Error en iafunctions al obtener respuesta llm", http.StatusInternalServerError)
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(resp)
}

func InteractAgent(respWriter http.ResponseWriter, request *http.Request) {

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
	authInst, ok3 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok3 {
		http.Error(respWriter, "Token inválido", http.StatusUnauthorized)
		return
	}
	if authInst != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	// ===== Estructura de entrada =====
	type agentQuery struct {
		IdDocument string `json:"idDocument,omitempty"`
		Query      string `json:"query"`
		Type       string `json:"type"`
	}
	var aq agentQuery
	err = json.NewDecoder(request.Body).Decode(&aq)
	if err != nil {
		http.Error(respWriter, "Error al decodificar json", http.StatusBadRequest)
	}
	resp := iapackage.InteractDoc(aq.IdDocument, idInst, idUser, aq.Query, aq.Type)
	if resp == nil {
		http.Error(respWriter, "Error en iafunctions al obtener respuesta llm", http.StatusInternalServerError)
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(respWriter).Encode(resp)
}
