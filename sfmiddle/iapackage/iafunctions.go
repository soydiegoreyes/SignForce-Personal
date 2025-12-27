package iapackage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sfmiddle/db"
	"sfmiddle/models"
)

func GetAbstractDoc(idDoc string) {

	attrs := []string{"documentPath", "documentName", "documentExt", "activeDoc"}
	wheres := map[string][]string{
		"idDocument": {idDoc},
	}

	docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
	if err == nil && len(docData) > 0 {
		fileName := fmt.Sprintf("%s/%s%s.%s", os.Getenv("BASE_DIR"), docData[idDoc]["documentPath"], docData[idDoc]["documentName"], docData[idDoc]["documentExt"])
		_, err = os.Stat(fileName)
		if err != nil {
			fmt.Println("No se encontró el archivo: ", fileName)
		} else {

			wheres = map[string][]string{
				"nameApp": {"iaServ"},
			}

			appData, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, wheres)
			if err != nil {
				return
			}

			// Modelo para request a api de LLM
			payload := models.LLMrequest{
				Path:   fileName,
				Action: 1,
			}
			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("Error al convertir a JSON:", err)
				return
			}
			var host, port string
			for _, v := range appData {
				host = v["domainApp"]
				port = v["portApp"]
				break
			}

			req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/actions", host, port), bytes.NewBuffer(jsonPayload))
			if err != nil {
				fmt.Println("Error en request: ", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
			// Ejecutar petición
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error en respuesta: ", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				var llamaresp = models.LLMresp{}
				err := json.NewDecoder(resp.Body).Decode(&llamaresp)
				if err != nil {
					fmt.Println("Error al convertir respuesta a JSON:", err)
					return
				}

				updates := map[string]map[string]interface{}{
					idDoc: {
						"abstractDoc": llamaresp.Message,
					},
				}

				err = db.DB_con.GenericBatchUpdate("documents", "idDocument", updates)
				if err != nil {
					fmt.Println("No se insertó el resumen en documento: ", idDoc)
					return
				}
			}
		}
	} else {
		fmt.Println("Error en obtener datos de documentos o no hay información")
	}
}

func InteractDoc(idDoc, idInst, idUser, prompt, typeIA string) *models.LLMresp {
	var fileName string
	if idDoc != "" {
		attrs := []string{"documentPath", "documentName", "documentExt", "activeDoc"}
		wheres := map[string][]string{
			"idDocument": {idDoc},
		}

		docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
		if err == nil && len(docData) > 0 {
			fileName = fmt.Sprintf("%s/%s%s.%s", os.Getenv("BASE_DIR"), docData[idDoc]["documentPath"], docData[idDoc]["documentName"], docData[idDoc]["documentExt"])
			_, err = os.Stat(fileName)
			if err != nil {
				fmt.Println("No se encontró el archivo: ", fileName)
			}
		} else {
			fmt.Println("Error en obtener datos de documentos o no hay información")
			return nil
		}
	}

	wheres := map[string][]string{
		"nameApp": {"iaServ"},
	}

	appData, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, wheres)
	if err != nil {
		return nil
	}

	// Modelo para request a api de LLM
	payload := models.LLMrequest{
		IdInst: idInst,
		IdUser: idUser,
		Path:   fileName,
		Prompt: prompt,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return nil
	}
	var host, port string
	for _, v := range appData {
		host = v["domainApp"]
		port = v["portApp"]
		break
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/%s", host, port, typeIA), bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error en request: ", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
	// Ejecutar petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error en respuesta: ", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var llamaresp = models.LLMresp{}
		err := json.NewDecoder(resp.Body).Decode(&llamaresp)
		if err != nil {
			fmt.Println("Error al convertir respuesta a JSON:", err)
			return nil
		}
		fmt.Println("Message: ", llamaresp.Message)
		return &llamaresp
	} else {
		fmt.Println("Error en servicio llm: ", err)
		return nil
	}
}

func interactDoc(idDoc, query string) *models.LLMresp {
	if idDoc != "" {
		attrs := []string{"documentPath", "documentName", "documentExt", "activeDoc"}
		wheres := map[string][]string{
			"idDocument": {idDoc},
		}

		docData, err := db.DB_con.GenericSelect("documents", "idDocument", attrs, wheres)
		if err != nil || len(docData) == 0 {
			fmt.Println("Error al obtener datos del documento: ", err)
			return nil
		}
		fileName := fmt.Sprintf("%s/%s%s.%s", os.Getenv("BASE_DIR"), docData[idDoc]["documentPath"], docData[idDoc]["documentName"], docData[idDoc]["documentExt"])
		_, err = os.Stat(fileName)
		if err != nil {
			fmt.Println("No se encontró el archivo: ", fileName)
		}
	}
	return nil
}
