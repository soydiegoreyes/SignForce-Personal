package objects

import (
	//"errors"
	"sfmiddle/db"

	//"sfmiddle/utilities"
	"sfmiddle/models"
)

// PRIMER FUNCION PARA REGISTRAR UN NUEVO CLIENTE
// NewUser crea una nueva instancia de User
func RegisterInst(registerReq *models.RegisterRequest) (string, error) {
	columns := []string{
		"legalNameInst",
		"aliasNameInst",
		"taxNumInst",
		"contactEmailInst",
		"contactPhoneInst",
		"legalSignupName",
		"legalSignupLastname",
		"registerSignupName",
		"registerSignupLastname",
		"statusInst_fk",
		"activeInst",
	}
	values := []interface{}{
		registerReq.NameInst,
		registerReq.AliasNameInst,
		registerReq.TaxNumInst,
		registerReq.ContactEmailInst,
		registerReq.ContactPhoneInst,
		registerReq.LegalSignupName,
		registerReq.LegalSignupLastname,
		registerReq.RegisterSignupName,
		registerReq.RegisterSignupLastname,
		2, // PENDIENTE_REGISTRO
		0, // A1 esta activa
	}

	return db.DB_con.GenericInsert("institutions", columns, values)
}
