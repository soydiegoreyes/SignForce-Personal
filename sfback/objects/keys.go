package objects

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"sfback/utilities"
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
	privateKey  *rsa.PrivateKey
	Certificate *x509.Certificate
	CertMap     CertificateMap
	ValidKeys   bool
}

// Constructor para Keys
func NewKeys(keypath string, certpath string) *Keys {
	return &Keys{
		keyfile:  keypath,
		Certfile: certpath,
	}
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
	result["NotBefore"] = cert.NotBefore
	result["NotAfter"] = cert.NotAfter
	result["SerialNumber"] = cert.SerialNumber.String()
	result["Signature"] = fmt.Sprintf("%x", cert.Signature)
	result["SignatureAlgorithm"] = cert.SignatureAlgorithm.String()
	result["PublicKey"] = cert.PublicKey.(*rsa.PublicKey).N
	result["Exponent"] = cert.PublicKey.(*rsa.PublicKey).E
	result["PublicKeyAlgorithm"] = cert.PublicKeyAlgorithm.String()
	result["BasicConstraintsValid"] = cert.BasicConstraintsValid
	result["ExtKeyUsage"] = cert.ExtKeyUsage //[{2.5.29.19 true [48 0]} {2.5.29.15 false [3 2 3 216]} {2.16.840.1.113730.1.1 false [3 2 5 160]} {2.5.29.37 false [48 20 6 8 43 6 1 5 5 7 3 4 6 8 43 6 1 5 5 7 3 2]}]
	result["Extensions"] = cert.Extensions
	result["ExtraExtensions"] = cert.ExtraExtensions
	result["IsCA"] = cert.IsCA
	result["KeyUsage"] = cert.KeyUsage
	result["UnhandledCriticalExtensions"] = cert.UnhandledCriticalExtensions
	result["UnknownExtKeyUsage"] = cert.UnknownExtKeyUsage

	// Procesar Subject
	subjectMap, err := parseRDNsToMap(cert.RawSubject)
	if err != nil {
		return nil, fmt.Errorf("error parsing subject: %v", err)
	}
	result["Subject"] = subjectMap

	// Procesar Issuer
	issuerMap, err := parseRDNsToMap(cert.RawIssuer)
	if err != nil {
		return nil, fmt.Errorf("error parsing issuer: %v", err)
	}
	result["Issuer"] = issuerMap

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

	var subject4514 string = ""
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
			subject4514 = subject4514 + utilities.Oids[name]["rfc4514"] + "=" + valDecoded + ","
		}
	}
	attrMap["subject4514"] = subject4514

	return attrMap, nil
}
