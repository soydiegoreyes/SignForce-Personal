package objects

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"sfback/db"
	"sfback/utilities"
	"strings"
	"time"
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

// NewUser crea una nueva instancia de User
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

	certPathPem, err := utilities.ConvertCertToPem(keysUser[idKeyUser]["certFilePath"])
	if err != nil {
		return nil, fmt.Errorf("error al convertir certificado a PEM %v. error: %v", userData["certFilePath"], err)
	}

	keyPathPem, err := utilities.ConvertKeyToPem(keysUser[idKeyUser]["keyFilePath"], password)
	if err != nil {
		return nil, fmt.Errorf("error al convertir llave a PEM %v. error: %v", userData["keyFilePath"], err)
	}

	var usuActivo bool
	if userData[idUser]["activeUser"] == "1" {
		usuActivo = true
	} else {
		usuActivo = false
	}
	user := &User{
		Uid:      idUser,
		Keys:     NewKeys(keyPathPem, certPathPem),
		Name:     userData[idUser]["nameUser"],
		LastName: userData[idUser]["lastNameUser"],
		PobID:    userData[idUser]["pobUidUser"],
		TaxNum:   userData[idUser]["taxNumUser"],
		Email:    userData[idUser]["emailUser"],
		Active:   usuActivo,
	}

	// Cargar claves y certificado
	err = user.loadPrivateKey(keyPathPem)
	if err != nil {
		return nil, fmt.Errorf("error al cargar la clave privada: %v", err)
	}

	err = user.loadCertificate(certPathPem)
	if err != nil {
		return nil, fmt.Errorf("error al cargar el certificado: %v", err)
	}

	// Validar certificado
	if !user.validateKeys() {
		return nil, errors.New("el certificado no está vigente o no coincide con la clave privada")
	}
	user.Keys.ValidKeys = true
	fmt.Printf("Usuario logueado: %s %s %s\n", user.Uid, user.Name, user.LastName)
	return user, nil
}

// Carga la clave privada desde el archivo de usuario
func (u *User) loadPrivateKey(keyPath string) error {
	if !strings.HasSuffix(keyPath, ".pem") {
		return errors.New("la clave privada no es formato .PEM")
	}
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo de clave privada: %v", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("no se pudo decodificar el bloque PEM de la clave privada")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("error al parsear la clave privada: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return errors.New("la clave privada no es una clave RSA")
	}
	u.Keys.privateKey = rsaPrivateKey

	return nil
}

// Carga el certificado desde el archivo de usuario
func (u *User) loadCertificate(certPath string) error {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("no se pudo decodificar el certificado")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.New("error al parsear el certificado")
	}
	u.Keys.Certificate = cert

	certMap, err := ParseCertificateToMap(cert)
	if err != nil {
		log.Fatal(err)
	}
	u.Keys.CertMap = certMap

	return nil
}

// Valida si la clave y el certificado coinciden y están vigentes
func (u *User) validateKeys() bool {
	if u.Keys.privateKey == nil || u.Keys.Certificate == nil {
		return false
	}

	validBase := u.Keys.privateKey.PublicKey.N.Cmp(u.Keys.Certificate.PublicKey.(*rsa.PublicKey).N) == 0
	validExp := u.Keys.privateKey.PublicKey.E == u.Keys.Certificate.PublicKey.(*rsa.PublicKey).E

	now := time.Now()
	notExpired := now.After(u.Keys.Certificate.NotBefore) && now.Before(u.Keys.Certificate.NotAfter)

	validPobUid := u.PobID == u.Keys.CertMap["Subject"].(map[string]string)[utilities.Coids["x509"]["serialNumber"]]
	validTaxUid := u.TaxNum == u.Keys.CertMap["Subject"].(map[string]string)[utilities.Coids["x509"]["x500UniqueIdentifier"]]

	fmt.Printf("subject--: %v\n", u.Keys.CertMap["Subject"].(map[string]string)["commonName"])

	return validBase && validExp && notExpired && validTaxUid && validPobUid
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
