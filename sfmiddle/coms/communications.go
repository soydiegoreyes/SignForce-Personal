package coms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sfmiddle/db"
	"sfmiddle/models"
	"time"
)

type EmailClient struct {
	url       string
	nextCheck time.Time
	retrys    int
}

func (EC *EmailClient) SendMail(payload *models.EmailRequest) error {
	defer func() {
		EC.retrys = 5
	}()
	if EC.nextCheck.After(time.Now()) {
		EmailCli = ConfEmail()
	}
	var check bool
	var returnError error

	for !check && EC.retrys > 0 {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error en payload al convertir a JSON:", err)
			return err
		}

		req, err := http.NewRequest("POST", EC.url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			EC.retrys -= 1
			fmt.Printf("%s", err)
			returnError = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("token", "3") // cambiar por bearer************************** importante!!
		// Ejecutar petición
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			EC.retrys -= 1
			fmt.Printf("%s", err)
			returnError = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			EC.retrys -= 1
			returnError = fmt.Errorf("Error en servicio email: %s", resp.Status)
		} else {
			fmt.Println("Email Enviado")
			break
		}
	}

	return returnError
}

func ConfEmail() EmailClient {
	var EC EmailClient
	whereMap := map[string][]string{
		"nameApp": {"emailServ"},
	}
	data, err := db.DB_con.GenericSelect("microapps", "idapp", []string{"domainApp", "portApp"}, whereMap)
	if err != nil {
		fmt.Printf("%s", err)
		return EC
	}

	for _, v := range data {
		EC.url = fmt.Sprintf("http://%s:%s/mailserv", v["domainApp"], v["portApp"])
		EC.nextCheck = time.Now().Add(time.Hour)
		EC.retrys = 5
		break
	}
	return EC
}

var EmailCli EmailClient
