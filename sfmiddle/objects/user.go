package objects

import (
	//"errors"
	"fmt"
	"sfmiddle/db"

	"sfmiddle/configs"
	"sfmiddle/models"
	"sfmiddle/utilities"
)

type User struct{}

// PRIMER FUNCION PARA REGISTRAR UN NUEVO CLIENTE
func RegisterRootUser(registerReq *models.RegisterRequest, idInst, passHash string) (string, error) {
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
		"isAliveUser",
		"roleAppUser_fk",
		"appPassHash",
		//"idTeam_fk",
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
		"1",
		passHash,
		idInst,
	}
	idUser, err := db.DB_con.GenericInsert("users", columns, values)
	if err != nil {
		fmt.Println("Error: al insertar datos", err)
		return "", err
	}

	dataRole, err := db.DB_con.GenericSelect("userroles", "idUser", []string{"idRole"}, map[string][]string{"idUser": {idUser}, "idRole": {"1"}})
	if err != nil {
		return "", err
	}
	if len(dataRole) != 0 {
		return "", fmt.Errorf("colision de usuario y rol. El usuario ya existe")
	}
	_, err = db.DB_con.GenericInsert("userroles", []string{"idUser", "idRole", "grantedBy"}, []interface{}{idUser, 1, idUser})
	if err != nil {
		return "", fmt.Errorf("no se pudo asignar el rol al usuario")
	}
	return idUser, nil
}

func UserInstJoin(idUser string) map[string]string {
	attrs1 := []string{"idUser", "nameUser", "lastNameUser", "emailUser", "aliasUser", "activeUser"}
	attrs2 := []string{"idInstitution", "legalNameInst", "aliasNameInst", "contactEmailInst", "taxNumInst", "activeInst"}
	userData, err := db.DB_con.GenericJoinSelect("users", "institutions", "users.idInstitution_fk = institutions.idInstitution", "idUser", []string{idUser}, attrs1, attrs2)
	//userData, err := db.DB_con.ExecuteSelect(query)
	if err != nil {
		fmt.Printf("%s", err)
		return nil
	}
	return userData[idUser]
}
