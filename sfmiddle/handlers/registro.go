package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/auth"
	"sfmiddle/coms"
	"sfmiddle/configs"
	"sfmiddle/db"
	"sfmiddle/models"
	"sfmiddle/objects"
	"sfmiddle/utilities"
	"strings"
	"time"
)

// ==========================================================================================================
func ProcessPayment(respWriter http.ResponseWriter, request *http.Request) {
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
	//idTeam, ok3 := claims["team"].(string)
	authInst, ok4 := claims["authInst"].(string)
	if !ok1 || !ok2 || !ok4 {
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
	fmt.Println(idUser, idInst, authInst)

	attrs := []string{"statusPayment_fk", "planId_fk", "expirationPlan", "nameOwner"}
	wheres := map[string][]string{"idInstitution": {idInst}}
	data, err := db.DB_con.GenericSelect("payment", "idInstitution", attrs, wheres)
	if err != nil {
		http.Error(respWriter, "Error al obtener informacion de institucion", http.StatusInternalServerError)
		return
	}

	resp := &models.PaymentResp{}
	var idPay string

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
		idPay, err = db.DB_con.GenericInsert("payment", cols, vals)
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
	fmt.Println("Pago procesado con ID: ", idPay)
	// Configurar headers y enviar respuesta
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	respWriter.Write(jsonData)
}

// =======================================================================
// Handler HTTP para loguear a un usuario por un ID de usuario y un arreglo de atributos a adquirir
func RegisterInst(respWriter http.ResponseWriter, request *http.Request) {
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
			userId, err := objects.RegisterRootUser(&registerReq, lastId, passHash)
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

			binDoc, err := os.ReadFile("./templates/welcome_register.html")
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
			}

			body := string(binDoc)
			body = strings.ReplaceAll(body, "{TEMPORAL_USERNAME}", registerReq.ContactEmailInst)
			body = strings.ReplaceAll(body, "{TEMPORAL_PASS}", tempPass)
			body = strings.ReplaceAll(body, "{EXPIRATION_TIME}", time.Now().Add(30*24*time.Hour).Format("2006-01-02 15:04:05"))
			//body = strings.ReplaceAll(body, "{URL_COMPLETAR_REGISTRO}", fmt.Sprintf("%s/login", os.Getenv("API_IP")))
			body = strings.ReplaceAll(body, "{URL_COMPLETAR_REGISTRO}", fmt.Sprintf("http://%s:%s/login", os.Getenv("API_IP"), os.Getenv("API_PORT")))

			payload := models.EmailRequest{
				IdUser:   userId,
				Subject:  fmt.Sprintf("¡Bienvenido a Signforce! Correo de verificación %s", registerReq.TaxNumInst),
				Body:     body,
				Dest:     []string{registerReq.ContactEmailInst},
				MimeType: "html",
			}

			err = coms.EmailCli.SendMail(&payload)
			if err != nil {
				registerResp.Error = fmt.Sprintf("%s", err)
				json.NewEncoder(respWriter).Encode(registerResp)
				return
			}

			registerResp.Check = true

		} else {
			registerResp.Error = "Ya tiene un registro para su numero de empresa. Revisar estatus de su registro."
			json.NewEncoder(respWriter).Encode(registerResp)
			return
		}
	}
	respWriter.Header().Set("Content-Type", "application/json")
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	json.NewEncoder(respWriter).Encode(registerResp)
}

// =======================================================================
// Endpoint para obtener datos de validación
func GetValidationData(respWriter http.ResponseWriter, request *http.Request) {
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
func UpdateValidationData(respWriter http.ResponseWriter, request *http.Request) {
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
	respWriter.Header().Set("X-Content-Type-Options", "nosniff")
	respWriter.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Datos actualizados correctamente"}
	json.NewEncoder(respWriter).Encode(response)
}

// =======================================================================
func CompleteValidation(respWriter http.ResponseWriter, request *http.Request) {
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

// =======================================================================
// Endpoint para obtener datos de validación
func Approvals(respWriter http.ResponseWriter, request *http.Request) {
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
	if idInst != "1" || idUser == "0" {
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
