package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
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
				docData.Hash, err = utilities.GetHash(docBytes, configs.HashConf, true)
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
					docData.Path = fmt.Sprintf("%s/%s", os.Getenv("GENERIC_IMG_PATH"), idInst)
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
				filePath, err = utilities.GuardarArchivo_cript(
					file, // archivo completo
					docData.Path,
					fileHeader.Filename,
					false,
					true,
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

	// Leer el body de la petición
	var docRequest models.DocDataRequest
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
		if len(docRequest.PathDocs) != 1 {
			http.Error(respWriter, "Error. Demasiados documentos para descarga", http.StatusBadRequest)
			return
		}
		t := docRequest.PathDocs[0]
		hash := t[:strings.Index(t, "@")]
		pathDoc, _ := strings.CutPrefix(t, hash+"@")
		lip := strings.LastIndex(pathDoc, ".")
		ext := pathDoc[lip+1:]
		nameDoc := pathDoc[:lip]
		tableName = "kyc"
		attrs = append(attrs, "documentHash", "documentPath", "documentName", "documentExt", "expirationDate")
		idColMain = "documentHash"
		docWheres = map[string][]string{
			"documentHash": {hash}, // Ajusta el nombre de la columna según tu esquema
			"documentName": {nameDoc},
			"documentExt":  {ext},
		}
	case "uploaded":
		if len(docRequest.IdDocs) != 1 {
			http.Error(respWriter, "Error. Demasiados documentos para descarga", http.StatusBadRequest)
			return
		}
		tableName = "documents"
		attrs = append(attrs, "idDocument", "documentPath", "documentName", "documentExt", "documentHash", "sizeB", "abstractDoc")
		idColMain = "idDocument"
		docWheres = map[string][]string{
			"idDocument": {docRequest.IdDocs[0]}, // Ajusta el nombre de la columna según tu esquema
		}
	case "template":
		if len(docRequest.IdDocs) != 1 {
			http.Error(respWriter, "Error. Demasiados documentos para descarga", http.StatusBadRequest)
			return
		}
		tableName = "doctemplates"
		attrs = append(attrs, "idTemplate", "documentPath", "documentName", "documentExt", "sizeB", "abstractDoc")
		idColMain = "idTemplate"
		docWheres = map[string][]string{
			"idTemplate": {docRequest.IdDocs[0]}, // Ajusta el nombre de la columna según tu esquema
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
	/*
		// Construir la ruta completa del archivo
		var filePath string = fmt.Sprintf("%s%s.%s", docInfo["documentPath"], docInfo["documentName"], docInfo["documentExt"])

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
		fmt.Println("numero de bytes ->", n)*/
	// 1. Construir ruta y verificar (Tu lógica actual)
	filePath := fmt.Sprintf("%s%s.%s", docInfo["documentPath"], docInfo["documentName"], docInfo["documentExt"])

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(respWriter, "Archivo no encontrado", http.StatusNotFound)
		return
	}
	defer file.Close()

	// 2. Determinar el Content-Type (Mejor por extensión si está cifrado)
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 3. Configurar Headers (IMPORTANTE: Quitamos Content-Length)
	// Al ser un flujo dinámico, Go usará "Transfer-Encoding: chunked" automáticamente
	respWriter.Header().Set("Content-Type", contentType)
	respWriter.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(filePath)))

	// 4. EL TRUCO DE MEMORIA: io.Pipe
	// pr: donde leeremos los datos descifrados
	// pw: donde la función DecryptFile escribirá los datos descifrados
	pr, pw := io.Pipe()

	// 5. Lanzamos la desencriptación en una Goroutine
	go func() {
		// Cerramos el pipe al terminar (esto avisa a io.Copy que ya no hay más datos)
		err := utilities.DecryptFile(file, pw)
		if err != nil {
			fmt.Printf("Error en streaming de descifrado: %v\n", err)
			pw.CloseWithError(err) // Notifica el error al lector
			return
		}
		pw.Close()
	}()

	// 6. Copiamos del Pipe directamente al cliente (Navegador)
	// io.Copy usa un buffer interno pequeño (32KB aprox), manteniendo la RAM baja.
	n, err := io.Copy(respWriter, pr)
	if err != nil {
		fmt.Printf("Error enviando flujo al cliente: %v\n", err)
		return
	}

	fmt.Printf("Archivo enviado con éxito. Bytes transferidos: %d\n", n)
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
		"creatorUserDoc_fk": {idUser},
		"ownerInstDoc_fk":   {idInst},
	}

	logic := "(creatorUserDoc_fk AND ownerInstDoc_fk) AND "

	// Si hay filtros por id, tipo o path -> se priorizan
	switch {
	case len(req.IdDocs) != 0:
		wheres[idCol] = req.IdDocs
		logic += " idDocument"

	case len(req.PathDocs) != 0:
		logic += "documentName"
		wheres["documentName"] = req.PathDocs

	case len(req.Tags) > 0:
		// buscamos si hay documentos que contengan los tags
		wtags := map[string][]string{
			"tag":   req.Tags,
			"LOGIC": {"LIKE tag"},
		}
		docTags, err := db.DB_con.GenericSelect("tagsdocuments", "idTag", []string{"idDocument"}, wtags)
		if err != nil {
			fmt.Println("Error al obtener documentos: " + err.Error())
			return
		}
		// si hay resultados se añaden a idDocs
		if len(docTags) > 0 {
			for _, tagDoc := range docTags {
				req.IdDocs = append(req.IdDocs, tagDoc["idDocument"])
			}
			wheres[idCol] = req.IdDocs
			wheres["abstractDoc"] = req.Tags
			wheres["documentName"] = req.Tags
			logic += "(idDocument OR LIKE abstractDoc OR LIKE documentName)"
		} else {
			wheres["abstractDoc"] = req.Tags
			wheres["documentName"] = req.Tags
			logic += "(LIKE abstractDoc OR LIKE documentName)"
		}

	case req.DateFrom != "" && req.DateTo != "":
		logic += "createdAtDoc BETWEEN ORDER BY " + orderBy
		wheres["createdAtDoc"] = []string{req.DateFrom, req.DateTo}
		wheres[idCol] = []string{orderDir}

	default:
		// Caso general: sin filtros, solo paginación
		logic = "creatorUserDoc_fk AND ownerInstDoc_fk ORDER BY " + orderBy
		//wheres[idCol] = []string{orderDir}
	}

	wheres["LOGIC"] = []string{logic, fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)}

	// ===== Ejecutar SELECT =====
	docData, err := db.DB_con.GenericSelect(tablename, idCol, attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener documentos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for _, v := range docData {
		v["documentPath"] = strings.ReplaceAll(v["documentPath"], os.Getenv("GENERIC_DOC_PATH"), "")
	}
	// ===== Obtener total de documentos =====
	baseWhere := fmt.Sprintf("creatorUserDoc_fk = '%s' AND ownerInstDoc_fk = '%s'", idUser, idInst)
	//baseWhere := fmt.Sprintf("ownerInstDoc_fk = '%s'", idInst)

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

// obtiene las plantillas creadas por la IA que aún estan como borrador (no han sido subidas y vivien en archivos aislados para editar)
func IaTemplates(respWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(respWriter, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// ===== 1. Autenticación (Tu lógica se mantiene igual) =====
	cookie, err := request.Cookie("token")
	if err != nil {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}
	claims, err := auth.ValidateJWT(cookie.Value)
	if err != nil || claims["authInst"] != "7" {
		http.Error(respWriter, "No autorizado", http.StatusUnauthorized)
		return
	}

	idUser := claims["uid"].(string)
	idInst := claims["iid"].(string)

	// ===== 2. Parámetros y Rutas =====
	docType := request.URL.Query().Get("type")
	if docType != "templates" && docType != "documents" {
		http.Error(respWriter, "type inválido", http.StatusBadRequest)
		return
	}

	// Ruta base física en el servidor
	basePath := fmt.Sprintf("./../sfia/%s/%s/%s/", idInst, idUser, docType)

	// Limpiamos la ruta
	fileName := filepath.Clean(request.URL.Path)

	// IMPORTANTE: Si fileName es "." o "/" significa que NO pidió un archivo, sino la raíz
	if fileName == "." || fileName == "/" || fileName == "" {
		files, err := os.ReadDir(basePath)
		if err != nil {
			respWriter.Header().Set("Content-Type", "application/json")
			json.NewEncoder(respWriter).Encode([]interface{}{})
			return
		}

		var fileList []map[string]interface{}
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			isDocx := strings.HasSuffix(file.Name(), ".docx")
			isPdf := strings.HasSuffix(file.Name(), ".pdf")

			if (docType == "templates" && isDocx) || (docType == "documents" && isPdf) {
				fileInfo, _ := file.Info()
				fileList = append(fileList, map[string]interface{}{
					"name":    file.Name(),
					"path":    file.Name(), // Pasamos el nombre para que el front lo use en la URL
					"size":    fileInfo.Size(),
					"created": fileInfo.ModTime(),
					"type":    strings.TrimPrefix(filepath.Ext(file.Name()), "."),
				})
			}
		}

		respWriter.Header().Set("Content-Type", "application/json")
		json.NewEncoder(respWriter).Encode(fileList)
		return
	}

	// Caso B: Servir archivo
	// Quitamos cualquier "/" sobrante al inicio para unir rutas correctamente
	fileName = strings.TrimPrefix(fileName, "/")
	fullPath := filepath.Join(basePath, fileName)

	// DEBUG: Descomenta esto para ver en consola qué ruta intenta buscar Go exactamente
	// fmt.Println("Buscando archivo en:", fullPath)

	info, err := os.Stat(fullPath)
	if err != nil || info.IsDir() {
		http.Error(respWriter, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	if request.URL.Query().Get("download") == "1" {
		respWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	}

	http.ServeFile(respWriter, request, fullPath)
}
