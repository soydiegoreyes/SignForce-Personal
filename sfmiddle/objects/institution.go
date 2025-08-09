package objects

import (
	//"errors"
	"sfmiddle/db"

	//"sfmiddle/utilities"
	"sfmiddle/models"
)

// Class User
type Institution struct {
	Uid                  string
	LegalName            string
	Alias                string
	TaxNum               string
	Email                string
	Phone                string
	Country              int
	City                 int
	Active               bool
	Status               int
	ContractType         int
	RegisteredByName     string
	RegisteredByLastname string
	LegalSignupName      string
	LegalSignupLastname  string
}

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
		1,
		0,
	}

	return db.DB_con.GenericInsert("institutions", columns, values)
}
