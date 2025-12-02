package objects

import (
	//"crypto/rsa"
	//"crypto/x509"
	//"encoding/pem"
	"errors"
	"fmt"
	"log"

	//"os"

	"sfback/db"
	"sfback/utilities"
	//"strings"
	//"time"
)

// Class User
type User struct {
	Uid      string
	Keys     *Keys
	Name     string
	LastName string
	PobID    string
	TaxNum   string
	Email    string
	Active   bool
}

// NewUser crea una nueva instancia de User que ya tiene todos los datos de llaves cargados en DB
func NewUser(idUser string, password string) (*User, error) {
	var attributes = []string{"nameUser", "lastNameUser", "pobUidUser", "taxNumUser", "emailUser", "activeUser", "idKeysUser_fk"}
	var wheres = map[string][]string{
		"idUser": {idUser},
	}
	userData, err := db.DB_con.GenericSelect("users", "idUser", attributes, wheres)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos del usuario %v. error: %v", idUser, err)
	}
	log.Println(userData)

	// Se obtienen las rutas de llave y certificado
	idKeyUser := userData[idUser]["idKeysUser_fk"]
	var keysUserAttr = []string{"keyFilePath", "certFilePath"}
	wheres = map[string][]string{
		"idUserKeys": {idKeyUser},
	}
	keysUser, err := db.DB_con.GenericSelect("userkeys", "idUserKeys", keysUserAttr, wheres)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos del usuario %v. error: %v", idUser, err)
	}
	log.Println(keysUser)

	var usuActivo bool
	if userData[idUser]["activeUser"] == "1" {
		usuActivo = true
	} else {
		usuActivo = false
	}
	user := &User{
		Uid:      idUser,
		Keys:     NewKeys(keysUser[idKeyUser]["keyFilePath"], keysUser[idKeyUser]["certFilePath"]),
		Name:     userData[idUser]["nameUser"],
		LastName: userData[idUser]["lastNameUser"],
		PobID:    userData[idUser]["pobUidUser"],
		TaxNum:   userData[idUser]["taxNumUser"],
		Email:    userData[idUser]["emailUser"],
		Active:   usuActivo,
	}
	if user.Keys == nil {
		return nil, fmt.Errorf("error al obtener hash de las llaves")
	}

	valKeysResp, err := user.Keys.ValidateKeys(password)
	if err != nil {
		return nil, err
	}
	// nunca deberia de suceder ya que un nuevo usuario
	if !valKeysResp.Exists {
		return nil, errors.New("critical: las llaves del usuario no coinciden con el id del usuario en registros de base de datos")
	}
	if !user.Keys.ValidKeys {
		return nil, errors.New("el certificado no está vigente o no coincide con la clave privada")
	}
	fmt.Printf("Propietario: %s, Exp: %s \n", valKeysResp.Owner, valKeysResp.Expiration)
	// comparacion de POBID y TAXNUM que es para validar si el usuario tiene completos esos datos
	/*
		validPobUid := user.PobID == user.Keys.CertMap["SubjectSerialNumber"]
		validTaxUid := user.TaxNum == user.Keys.CertMap["SubjectUniqueId"]

		if !validPobUid && !validTaxUid {
			user.Keys.ValidKeys = false
			return nil, errors.New("el certificado no coincide con el propietario registrado")
		}
	*/
	fmt.Printf("Usuario logueado: %s %s %s\n", user.Uid, user.Name, user.LastName)
	return user, nil
}

func (u *User) GetPublicParams(params []string) map[string]string {
	responseParams := make(map[string]string)
	var map_params = map[string]string{
		"uid":      u.Uid,
		"name":     u.Name,
		"lastName": u.LastName,
		"pobId":    u.PobID,
		"taxNum":   u.TaxNum,
		"email":    u.Email,
	}

	for _, param := range params {
		switch param {
		case "certb64":
			responseParams[param] = utilities.Encode_b64(u.Keys.Certificate.Raw)
		case "subject4514":
			responseParams[param] = u.Keys.CertMap["Subject"].(map[string]string)["subject4514"]
		default:
			responseParams[param] = map_params[param]
		}

	}
	return responseParams
}
