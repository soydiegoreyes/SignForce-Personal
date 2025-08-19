package objects

import (
	//"errors"
	"fmt"
	"sfmiddle/db"

	"sfmiddle/configs"
	"sfmiddle/models"
	"sfmiddle/utilities"
)

// PRIMER FUNCION PARA REGISTRAR UN NUEVO CLIENTE
// NewUser crea una nueva instancia de User
func RegisterUser(registerReq *models.RegisterRequest, idInst string, passHash string) (string, error) {
	idInstHash, err := utilities.GetHash([]byte(idInst), configs.HashConf)
	if err != nil {
		return "", err
	}
	columns := []string{
		"nameUser",
		"lastNameUser",
		"aliasUser",
		"emailUser",
		"phoneUser",
		"countryPhoneCode",
		"activeUser",
		"roleAppUser_fk",
		"appPassHash",
		"isSignerUser",
		"idTeam_fk",
		"idInstitution_fk",
	}
	values := []interface{}{
		registerReq.LegalSignupName,
		registerReq.LegalSignupLastname,
		fmt.Sprintf("root_%s", idInstHash),
		registerReq.ContactEmailInst,
		registerReq.ContactPhoneInst,
		"52",
		"1",
		"1",
		passHash,
		"0",
		"0",
		idInst,
	}

	return db.DB_con.GenericInsert("users", columns, values)
}
