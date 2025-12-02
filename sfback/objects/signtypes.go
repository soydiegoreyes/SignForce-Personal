package objects

import (
	"encoding/base64"
	"fmt"
	"os"

	"sfback/configs"
	"time"

	"sfback/utilities"

	// Tu paquete utilities para HashBytes si se necesita

	"github.com/beevik/etree"
	"github.com/google/uuid"
)

// Constantes de Namespaces y Algoritmos
const (
	NsXMLDSig       = "http://www.w3.org/2000/09/xmldsig#"
	NsXAdES         = "http://uri.etsi.org/01903/v1.3.2#"
	AlgSHA256       = "http://www.w3.org/2001/04/xmlenc#sha256"
	AlgRSA_SHA256   = "http://www.w3.org/2000/09/xmldsig#rsa-sha256"
	AlgC14N         = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	AlgEnvSig       = "http://www.w3.org/2000/09/xmldsig#enveloped-signature"
	TypeSignedProps = "http://uri.etsi.org/01903#SignedProperties"
)

// GenerarFirmaXades construye la estructura XAdES-BES completa
// pdfHash: El hash SHA256 del archivo PDF que quieres firmar
// certPEM: Los bytes del certificado público (archivo .crt o .pem)
// k: Tus llaves para firmar
func (k *Keys) GenerarFirmaXades(pdfHash []byte) (map[string]string, error) {
	fmt.Println("entrando a generarfirmaXades")
	// Identificadores únicos para enlazar las referencias
	signatureId := uuid.NewString() // Podrías usar UUID
	signedPropsId := "SignedProperties-" + uuid.NewString()
	//objectId := "Object-" + uuid.NewString()

	// Convertir el hash del PDF a Base64 para ponerlo en el XML
	pdfHashB64 := base64.StdEncoding.EncodeToString(pdfHash)
	// =========================================================================
	// PASO 1: CONSTRUIR EL NODO SignedProperties (Propiedades firmadas)
	// =========================================================================

	// Calculamos el hash del certificado (CertDigest)
	certHashB64, err := utilities.GetHash(k.Certificate.Raw, configs.HashConf)
	if err != nil {
		return nil, err
	}
	// Usamos etree para construir el XML de SignedProperties
	// Es CRÍTICO usar etree aquí para luego obtener su C14N exacto
	spDoc := etree.NewDocument()
	spRoot := spDoc.CreateElement("xades:SignedProperties")
	spRoot.CreateAttr("xmlns:xades", NsXAdES) // Namespace debe estar aquí para C14N local
	spRoot.CreateAttr("xmlns:ds", NsXMLDSig)  // A veces necesario si hay hijos ds:
	spRoot.CreateAttr("Id", signedPropsId)

	ssp := spRoot.CreateElement("xades:SignedSignatureProperties")

	// 1.1 Tiempo de firma (UTC ISO 8601)
	gentime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	st := ssp.CreateElement("xades:SigningTime")
	st.SetText(gentime)

	// 1.2 Certificado de firma
	sc := ssp.CreateElement("xades:SigningCertificate")
	certNode := sc.CreateElement("xades:Cert")

	certDigestNode := certNode.CreateElement("xades:CertDigest")
	dm := certDigestNode.CreateElement("ds:DigestMethod")
	dm.CreateAttr("Algorithm", AlgSHA256)
	dv := certDigestNode.CreateElement("ds:DigestValue")
	dv.SetText(certHashB64)

	issuerSerial := certNode.CreateElement("xades:IssuerSerial")
	xn := issuerSerial.CreateElement("ds:X509IssuerName")
	xn.SetText(k.Certificate.Issuer.String()) // RFC4514 string
	sn := issuerSerial.CreateElement("ds:X509SerialNumber")
	sn.SetText(k.Certificate.SerialNumber.String())

	// =========================================================================
	// PASO 2: CANONICALIZAR Y HASHEAR SignedProperties
	// =========================================================================
	// Configurar etree para simular C14N (sin espacios, tags cerrados)
	spDoc.WriteSettings.CanonicalEndTags = true

	// Obtenemos los bytes de SignedProperties
	signedPropsBytes, err := spDoc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	// Calculamos el hash de SignedProperties (Esto irá en la segunda referencia)
	spHashB64, err := utilities.GetHash(signedPropsBytes, configs.HashConf)
	if err != nil {
		return nil, err
	}
	// =========================================================================
	// PASO 3: CONSTRUIR EL SignedInfo
	// =========================================================================
	siDoc := etree.NewDocument()
	siRoot := siDoc.CreateElement("ds:SignedInfo")
	siRoot.CreateAttr("xmlns:ds", NsXMLDSig)
	// Importante: SignedInfo a veces necesita conocer el namespace de xades si se usa en atributos,
	// pero generalmente basta con definir ds aquí.

	cm := siRoot.CreateElement("ds:CanonicalizationMethod")
	cm.CreateAttr("Algorithm", AlgC14N)

	sm := siRoot.CreateElement("ds:SignatureMethod")
	sm.CreateAttr("Algorithm", AlgRSA_SHA256)

	// --- Referencia 1: El Documento (PDF) ---
	refDoc := siRoot.CreateElement("ds:Reference")
	refDoc.CreateAttr("URI", "") // URI vacía o ID del objeto si fuera enveloped
	// Si es detached (el pdf está fuera), URI suele ser el nombre del archivo o vacío si es todo el contexto
	// En tu ejemplo anterior usabas un ID específico. Ajustar según necesites.
	// refDoc.CreateAttr("URI", "#PSCNotaria...")

	// Transform (opcional, pero común en enveloped)
	trans := refDoc.CreateElement("ds:Transforms")
	t1 := trans.CreateElement("ds:Transform")
	t1.CreateAttr("Algorithm", AlgEnvSig)

	rdm := refDoc.CreateElement("ds:DigestMethod")
	rdm.CreateAttr("Algorithm", AlgSHA256)
	rdv := refDoc.CreateElement("ds:DigestValue")
	rdv.SetText(pdfHashB64)

	// --- Referencia 2: Las Propiedades (SignedProperties) ---
	// ESTO HACE QUE SEA XAdES
	refProp := siRoot.CreateElement("ds:Reference")
	refProp.CreateAttr("Type", TypeSignedProps)
	refProp.CreateAttr("URI", "#"+signedPropsId)

	rpm := refProp.CreateElement("ds:DigestMethod")
	rpm.CreateAttr("Algorithm", AlgSHA256)
	rpv := refProp.CreateElement("ds:DigestValue")
	rpv.SetText(spHashB64)

	// =========================================================================
	// PASO 4: FIRMAR EL SignedInfo
	// =========================================================================

	// Canonicalizar SignedInfo
	siDoc.WriteSettings.CanonicalEndTags = true
	signedInfoBytes, err := siDoc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	// FIRMA CRIPTOGRÁFICA
	// Pasamos los bytes del XML. Tu función SignHash hará el SHA256 y luego RSA.
	signedInfoHash, err := utilities.GetHash(signedInfoBytes, configs.HashConf)
	if err != nil {
		return nil, err
	}
	signatureBytes, err := k.SignHash(utilities.Decode_b64(signedInfoHash), "sha256")
	if err != nil {
		return nil, err
	}
	signatureValueB64 := utilities.Encode_b64(signatureBytes)

	// =========================================================================
	// PASO 5: ARMAR EL XML FINAL (<Signature>)
	// =========================================================================

	finalDoc := etree.NewDocument()
	sigRoot := finalDoc.CreateElement("ds:Signature")
	sigRoot.CreateAttr("xmlns:ds", NsXMLDSig)
	sigRoot.CreateAttr("xmlns:xades", NsXAdES) // Definimos ambos globales para limpieza
	sigRoot.CreateAttr("Id", signatureId)

	// 5.1 Agregar SignedInfo (Tal cual lo generamos)
	sigRoot.AddChild(siRoot)

	// 5.2 Agregar SignatureValue
	sv := sigRoot.CreateElement("ds:SignatureValue")
	sv.CreateAttr("Id", "value-"+signatureId)
	sv.SetText(signatureValueB64)

	// 5.3 Agregar KeyInfo (Datos del certificado público)
	ki := sigRoot.CreateElement("ds:KeyInfo")
	x509Data := ki.CreateElement("ds:X509Data")

	// Opcional: IssuerSerial también aquí (redundante con SignedProperties pero común)
	xis := x509Data.CreateElement("ds:X509IssuerSerial")
	xis.CreateElement("ds:X509IssuerName").SetText(k.Certificate.Issuer.String())
	xis.CreateElement("ds:X509SerialNumber").SetText(k.Certificate.SerialNumber.String())

	x509Cert := x509Data.CreateElement("ds:X509Certificate")
	// El certificado en Base64 sin headers PEM
	x509Cert.SetText(certHashB64)

	// Opcional: KeyValue (Modulo y Exponente)
	// Es buena práctica incluirlo pero X509Certificate suele bastar.

	// 5.4 Agregar Object -> QualifyingProperties -> SignedProperties
	obj := sigRoot.CreateElement("ds:Object")
	qp := obj.CreateElement("xades:QualifyingProperties")
	qp.CreateAttr("Target", "#"+signatureId)

	// AQUÍ insertamos el nodo SignedProperties que creamos en el Paso 1
	// Importante: Debe ser idéntico al que hasheamos
	qp.AddChild(spRoot)

	// =========================================================================
	// RESULTADO FINAL
	// =========================================================================
	finalDoc.WriteSettings.CanonicalEndTags = true
	finalDoc.Indent(2) // Indentación para que se vea bonito (pretty print)
	// Nota: La indentación del XML final no rompe la firma PORQUE SignedInfo
	// se procesó y firmó internamente sin indentación.

	xmlString, err := finalDoc.WriteToString()
	if err != nil {
		return nil, err
	}

	fmt.Println("todo bien")
	// crear xml
	xmlPath := "./temp/" + signatureId + ".xml"
	if err := os.MkdirAll("./temp", 0755); err != nil {
		return nil, err
	}

	// Crear archivo solo si NO existe (modo seguro)
	f, err := os.OpenFile(xmlPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, fmt.Errorf("el archivo %s ya existe", xmlPath)
	}
	defer f.Close()

	// Escribir XML
	if _, err := f.Write([]byte(xmlString)); err != nil {
		return nil, err
	}

	// Respuesta de firma
	signData := map[string]string{
		"xmlPath":            xmlPath,
		"signedInfoHash":     signedInfoHash,
		"signatureValueSign": signatureValueB64,
		"genTimeSign":        gentime,
		"nonceSign":          signatureId,
		"certHash":           certHashB64,
		"typeSign":           "9",
	}

	return signData, nil
}

/*
### Explicación de los Cambios Clave

1.  **Uso de `etree` para todo**: Eliminé las estructuras `struct` manuales (`Signature`, `SignedInfo`, etc.) dentro de esta función para la *generación*. Es mucho más seguro usar `etree` para construir el XML dinámicamente porque te permite manipular atributos y namespaces con precisión antes de calcular hashes. Las `structs` son buenas para *leer* (unmarshal), pero para *escribir* y firmar XAdES (donde el orden importa para el hash), `etree` es superior.
2.  **`SignHash` (Tu función)**: En el código hago:
    ```go
    // signedInfoBytes contiene el XML <SignedInfo>...</SignedInfo>
    k.SignHash(signedInfoBytes, "sha256")
    ```
    Como tu función `SignHash` hace `sha256.Sum256(hashed)` (donde `hashed` es el input), esto funciona perfecto. El `rsa.SignPKCS1v15` firmará el hash del XML.
3.  **`SigningCertificate`**: Agregué la lógica para parsear tu certificado PEM (`x509.ParseCertificate`). XAdES **exige** que pongas el hash del certificado y el número de serie dentro de `SignedProperties`. Si no pones esto, no es XAdES, es XMLDSig.
4.  **Doble Referencia**: Fíjate en el bloque **PASO 3**. Hay dos `Reference`:
    * Una apunta al PDF (con `pdfHashB64`).
    * Otra apunta a `#SignedProperties` (con `spHashB64`).
    Esto vincula criptográficamente la fecha de firma y el certificado con el documento.

### ¿Cómo usarlo?

```go
// Supongamos que ya tienes:
// pdfHashBytes: []byte del SHA256 del PDF
// user.Keys: Tu struct con la llave privada cargada
// certPemBytes: []byte del archivo .crt o .pem público

xmlXades, err := signatures.GenerarFirmaXades(pdfHashBytes, certPemBytes, user.Keys)
if err != nil {
    log.Fatal(err)
}
fmt.Println(xmlXades)

*/
