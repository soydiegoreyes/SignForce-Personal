package objects

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"sfback/configs"
	"sfback/db"
	"sfback/models"
	"sfback/utilities"
	"strings"
	"time"
)

/*
  KEYS:
  -- Este programa modela las llaves publica y privada de un usuario, por lo que puede ser invocada desde cualquier parte que el usuario necesite obtener sus llaves
  -- Recibe
*/

// UserKeys estructura que representa las llaves y el certificado de un usuario
type Keys struct {
	keyfile     string
	Certfile    string
	KeyHash     string
	CertHash    string
	privateKey  *rsa.PrivateKey
	Certificate *x509.Certificate
	CertMap     CertificateMap
	ValidKeys   bool
}

// Constructor para Keys
func NewKeys(keypath string, certpath string) *Keys {
	keyHash, err := utilities.GetHash(keypath, configs.HashConf)
	if err != nil {
		fmt.Println("error al obtener hash de llave privada")
		return nil
	}
	cerHash, err := utilities.GetHash(certpath, configs.HashConf)
	if err != nil {
		fmt.Println("error al obtener hash de certificado")
		return nil
	}
	return &Keys{
		keyfile:  keypath,
		Certfile: certpath,
		KeyHash:  keyHash,
		CertHash: cerHash,
	}
}

// Carga el certificado desde el archivo de usuario
func (k *Keys) loadCertificate(certPath string) error {
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
	k.Certificate = cert

	certMap, err := ParseCertificateToMap(cert)
	if err != nil {
		return err
	}
	k.CertMap = certMap
	fmt.Println(certMap)
	return nil
}

// Carga la clave privada desde el archivo de usuario
func (k *Keys) loadPrivateKey(keyPath string) error {
	if !strings.HasSuffix(keyPath, ".pem") {
		return errors.New("la llave privada no es formato .PEM")
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
	k.privateKey = rsaPrivateKey

	return nil
}

// Valida si la clave y el certificado coinciden y están vigentes
func (k *Keys) TestKeys() bool {
	if k.privateKey == nil || k.Certificate == nil {
		return false
	}

	validBase := k.privateKey.PublicKey.N.Cmp(k.Certificate.PublicKey.(*rsa.PublicKey).N) == 0
	validExp := k.privateKey.PublicKey.E == k.Certificate.PublicKey.(*rsa.PublicKey).E

	now := time.Now()
	notExpired := now.After(k.Certificate.NotBefore) && now.Before(k.Certificate.NotAfter)

	return validBase && validExp && notExpired
}

func (k *Keys) ValidateKeys(password string) (*models.UploadKeysResponse, error) {
	var valid bool = false

	certPathPem, err := utilities.ConvertCertToPem(k.Certfile)
	if err != nil {
		return nil, fmt.Errorf("error al convertir certificado a PEM %v. error: %v", k.Certfile, err)
	}

	err = k.loadCertificate(certPathPem)
	if err != nil {
		return nil, fmt.Errorf("error al cargar el certificado: %v", err)
	}

	keyPathPem, err := utilities.ConvertKeyToPem(k.keyfile, password)
	if err != nil {
		return nil, fmt.Errorf("error al convertir llave a PEM %v. error: %v", k.keyfile, err)
	}
	// Cargar claves y certificado
	err = k.loadPrivateKey(keyPathPem)
	if err != nil {
		return nil, fmt.Errorf("error al cargar la clave privada: %v", err)
	}
	err = os.Remove(keyPathPem)
	if err != nil {
		fmt.Println(err)
	}

	resp := &models.UploadKeysResponse{}
	if valid = k.TestKeys(); valid {
		fmt.Println("validas: ", valid)
		wheres := map[string][]string{
			"hashKey": {k.KeyHash},
			"hashCer": {k.CertHash},
		}
		keys, err := db.DB_con.GenericSelect("userkeys", "idUserKeys", []string{"serialNumber", "signature", "hashKey", "hashCer", "subjectUniqueId", "subjectSerialNumber"}, wheres)
		if err != nil {
			fmt.Println(err)
			resp.Exists = true
		}
		if len(keys) == 0 {
			resp.Exists = false
		} else {
			for _, key := range keys {
				if k.CertMap["SerialNumber"] == key["serialNumber"] && k.CertMap["Signature"] == key["signature"] && k.CertMap["SubjectUniqueId"] == key["subjectUniqueId"] && k.CertMap["SubjectSerialNumber"] == key["subjectSerialNumber"] {
					resp.Exists = true
					break
				}
			}
		}

		resp.Owner = k.CertMap["Subject"].(map[string]string)[utilities.Coids["x509"]["commonName"]]
		resp.Expiration = k.Certificate.NotAfter.Format("2006-01-02 15:04:05")
		k.ValidKeys = true
		resp.Valid = true
	}

	return resp, nil
}

// SignData genera una firma con la clave privada
func (k *Keys) SignData(digestMethod string, data []byte) ([]byte, error) {
	if k.privateKey == nil {
		return nil, errors.New("clave privada no disponible")
	}

	var hashed []byte
	var digestAlgo crypto.Hash

	switch digestMethod {
	case "sha256":
		hashed_ := sha256.Sum256(data)
		hashed = hashed_[:]
		digestAlgo = crypto.SHA256
	case "sha512":
		hashed_ := sha512.Sum512(data)
		hashed = hashed_[:]
		digestAlgo = crypto.SHA512
	default:
		hashed_ := sha256.Sum256(data)
		hashed = hashed_[:]
		digestAlgo = crypto.SHA256
	}

	signature, err := rsa.SignPKCS1v15(nil, k.privateKey, digestAlgo, hashed)
	if err != nil {
		return nil, fmt.Errorf("error al firmar los datos: %v", err)
	}

	return signature, nil
}

// VerifySignature verifica una firma usando la clave pública
func (k *Keys) VerifySignature(digestMethod string, data, signature []byte) error {
	if k.Certificate == nil {
		return errors.New("certificado no disponible")
	}

	publicKey, ok := k.Certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("error al obtener la clave pública del certificado")
	}

	var hashed []byte
	var digestAlgo crypto.Hash

	switch digestMethod {
	case "sha256":
		hashed_ := sha256.Sum256(data)
		hashed = hashed_[:]
		digestAlgo = crypto.SHA256
	case "sha512":
		hashed_ := sha512.Sum512(data)
		hashed = hashed_[:]
		digestAlgo = crypto.SHA512
	default:
		hashed_ := sha256.Sum256(data)
		hashed = hashed_[:]
		digestAlgo = crypto.SHA256
	}

	return rsa.VerifyPKCS1v15(publicKey, digestAlgo, hashed, signature)
}

// =============================================================================================================
// Diccionario para representar todos los campos de un certificado y parsearlos a una estructura mas manejable
type CertificateMap map[string]interface{}

// ParseCertificateToNestedMap convierte un certificado x509 a nuestro formato de mapas anidados
func ParseCertificateToMap(cert *x509.Certificate) (CertificateMap, error) {
	result := make(CertificateMap)

	// Campos estándar
	result["Version"] = cert.Version
	result["NotBefore"] = cert.NotBefore.Format("2006-01-02 15:04:05")
	result["NotAfter"] = cert.NotAfter.Format("2006-01-02 15:04:05")
	result["SerialNumber"] = cert.SerialNumber.String()
	result["Signature"] = utilities.Encode_b64(cert.Signature) // FIRMA B64
	result["SignatureAlgorithm"] = strings.ToLower(cert.SignatureAlgorithm.String())
	result["PublicKey"] = cert.PublicKey.(*rsa.PublicKey).N
	result["Exponent"] = cert.PublicKey.(*rsa.PublicKey).E
	result["PublicKeyAlgorithm"] = cert.PublicKeyAlgorithm.String()
	result["ExtKeyUsage"] = cert.ExtKeyUsage //[{2.5.29.19 true [48 0]} {2.5.29.15 false [3 2 3 216]} {2.16.840.1.113730.1.1 false [3 2 5 160]} {2.5.29.37 false [48 20 6 8 43 6 1 5 5 7 3 4 6 8 43 6 1 5 5 7 3 2]}]
	result["Extensions"] = cert.Extensions
	result["ExtraExtensions"] = cert.ExtraExtensions
	result["IsCA"] = cert.IsCA
	result["KeyUsage"] = cert.KeyUsage
	result["UnhandledCriticalExtensions"] = cert.UnhandledCriticalExtensions
	result["UnknownExtKeyUsage"] = cert.UnknownExtKeyUsage

	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		result["CurveName"] = ""
		result["KeySize"] = pub.Size() * 8
	case *ecdsa.PublicKey:
		result["CurveName"] = pub.Curve.Params().Name
		result["KeySize"] = pub.Curve.Params().BitSize
	case ed25519.PublicKey:
		result["CurveName"] = "Clave Ed25519"
		result["KeySize"] = 0
	default:
		result["CurveName"] = ""
		result["KeySize"] = 0
		fmt.Printf("Tipo de clave desconocido: %T\n", pub)
	}

	// Procesar Subject
	subjectMap, err := parseRDNsToMap(cert.RawSubject)
	if err != nil {
		return nil, fmt.Errorf("error parsing subject: %v", err)
	}
	// subject important attrs
	result["Subject"] = subjectMap
	result["SubjectUniqueId"] = subjectMap[utilities.Coids["x509"]["x500UniqueIdentifier"]]
	result["SubjectSerialNumber"] = cert.Subject.SerialNumber
	//result["SubjectEmailAddress"] = subjectMap[utilities.Coids["x509"]["emailAddress"]]

	// Procesar Issuer
	issuerMap, err := parseRDNsToMap(cert.RawIssuer)
	if err != nil {
		return nil, fmt.Errorf("error parsing issuer: %v", err)
	}
	result["Issuer"] = issuerMap
	result["IssuerUniqueId"] = issuerMap[utilities.Coids["x509"]["x500UniqueIdentifier"]]

	if len(cert.OCSPServer) > 0 {
		result["OCSP"] = cert.OCSPServer[0]
	} else {
		result["OCSP"] = ""
	}
	if len(cert.CRLDistributionPoints) > 0 {
		result["CRLS"] = cert.CRLDistributionPoints[0]
	} else {
		result["CRLS"] = ""
	}

	return result, nil
}

// parseRDNsToMap procesa los RDNSequence y devuelve un mapa con los atributos
func parseRDNsToMap(rawBytes []byte) (map[string]string, error) {
	attrMap := make(map[string]string)

	// los bytes de RawIssuer y RawSubject se codifican a asn1 y se guardan en rdnSeq
	var rdnSeq pkix.RDNSequence
	_, err := asn1.Unmarshal(rawBytes, &rdnSeq)
	if err != nil {
		return nil, err
	}

	var RFC4514 string = ""
	// rdnSeq es de tipo map[string]string
	for _, rdnSet := range rdnSeq {
		for _, atv := range rdnSet {
			name := atv.Type.String()
			valDecoded := utilities.Latin1ToUTF8([]byte(fmt.Sprint(atv.Value)))

			// name siempre será un OID
			// Si el atributo ya existe, concatenamos los valores
			if existing, exists := attrMap[name]; exists {
				attrMap[name] = existing + ", " + valDecoded
			} else {
				attrMap[name] = valDecoded
			}
			RFC4514 = RFC4514 + utilities.Oids[name]["rfc4514"] + "=" + valDecoded + ","
		}
	}
	attrMap["RFC4514"] = RFC4514

	return attrMap, nil
}
