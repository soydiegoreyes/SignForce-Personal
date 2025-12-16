package objects

import (
	"fmt"
	"os"
	"sfback/configs"
	"sfback/db"
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
	AlgDetSig       = "http://www.w3.org/2000/09/xmldsig#detached-signature"
	TypeSignedProps = "http://uri.etsi.org/01903#SignedProperties"
)

// GenerarFirmaXades construye la estructura XAdES-BES completa
// pdfHash: El hash SHA256 del archivo PDF que quieres firmar
// certPEM: Los bytes del certificado público (archivo .crt o .pem)
// k: Tus llaves para firmar
func (k *Keys) GenerarFirmaXades(pdfHash []byte, signatureId, docName, idDocument string) (map[string]string, error) {
	// se sacan los datos de la invitación, firma y usuarios involucrados
	signInvData, err := db.DB_con.GenericJoinSelect("signatures", "invites", "signatures.idInvite_fk = invites.idInvite", "idSignature", []string{signatureId}, []string{"idUser_fk"}, []string{"idInvite", "idFolder", "idUserOwnner_fk"})
	//signData, err := db.DB_con.GenericSelect("signatures", "idSignature", attrs, wheres)
	if err != nil {
		return nil, err
	}

	idFolder := signInvData[signatureId]["idFolder"]
	idUserOwnner := signInvData[signatureId]["idUserOwnner_fk"]
	idUserDest := signInvData[signatureId]["idUser_fk"]

	attrs := []string{"nameUser", "lastNameUser", "emailUser", "idInstitution_fk", "activeUser"}
	wheres := map[string][]string{
		"idUser": {idUserOwnner, idUserDest},
	}
	usersData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		return nil, err
	}
	invitePath := fmt.Sprintf("%s/folders/%s/%s/%s/META-INF/", os.Getenv("BASE_DIR"), usersData[idUserDest]["idInstitution_fk"], idFolder, idDocument)
	fmt.Println("invitePath completo: ", invitePath)
	err = CreateManifest(docName, invitePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Identificadores únicos para enlazar las referencias
	signedPropsId := "SignedProperties-" + uuid.NewString()

	// Convertir el hash del PDF a Base64 para ponerlo en el XML
	pdfHashB64 := utilities.Encode_b64(pdfHash)

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
	refDoc.CreateAttr("URI", docName) // URI vacía o ID del objeto si fuera enveloped
	// Si es detached (el pdf está fuera), URI suele ser el nombre del archivo o vacío si es todo el contexto
	// En tu ejemplo anterior usabas un ID específico. Ajustar según necesites.
	// refDoc.CreateAttr("URI", "#PSCNotaria...")

	// Transform (opcional, pero común en enveloped)
	trans := refDoc.CreateElement("ds:Transforms")
	t1 := trans.CreateElement("ds:Transform")
	t1.CreateAttr("Algorithm", AlgDetSig)

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
	sigRoot.CreateAttr("Id", "Id-"+signatureId)

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

	// crear xml
	xmlPath := invitePath + signatureId + ".xml"
	if err := os.MkdirAll(invitePath, 0755); err != nil {
		fmt.Println("error al crear el folder")
		return nil, err
	}
	// Crear archivo solo si NO existe (modo seguro)
	f, err := os.OpenFile(xmlPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Println("el archivo ya existe: ", xmlPath, err)
	} else {
		// Escribir XML
		if _, err := f.Write([]byte(xmlString)); err != nil {
			fmt.Println("error al escribir el archivo")
			return nil, err
		}
	}
	defer f.Close()

	// crear qrcode con link
	qrPath := invitePath + signatureId + ".png"

	if !utilities.GenerateQR(fmt.Sprintf("%sviewSignature?id=%s", os.Getenv("QR_URL_BASE"), signatureId), qrPath) {
		return nil, fmt.Errorf("no se pudo crear qr de firma %s", qrPath)
	}

	if ok := utilities.AddTextToQR(qrPath, fmt.Sprintf("%s %s - %s", usersData[idUserDest]["nameUser"], usersData[idUserDest]["lastNameUser"], gentime)); !ok {
		fmt.Println("No se pudo añadir el nombre al QR")
	}
	updates := map[string]map[string]interface{}{
		signatureId: {
			"pathImg": qrPath,
		},
	}
	err = db.DB_con.GenericBatchUpdate("signstamps", "idSignature_fk", updates)
	if err != nil {
		fmt.Println("error al actualizar qr de firma: ", err)
	}
	// Respuesta de firma
	signData := map[string]string{

		"xmlPath":            xmlPath,
		"hashDoc":            pdfHashB64,
		"signedInfoHash":     signedInfoHash,
		"signatureValueSign": signatureValueB64,
		"genTimeSign":        gentime,
		"nonceSign":          signatureId,
		"certHash":           certHashB64,
		"typeSign":           "9",
	}

	return signData, nil
}

func CreateManifest(docName, folderPath string) error {
	// crear xml
	xmlPath := folderPath + "manifest.xml"

	// Crear archivo solo si NO existe (modo seguro)
	f, err := os.OpenFile(xmlPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Println("el archivo ya existe: ", xmlPath, err)
		return nil
	}
	defer f.Close()

	const manifestVersion string = "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0"
	manifest := etree.NewDocument()
	manRoot := manifest.CreateElement("manifest:manifest")
	manRoot.CreateAttr("xmlns:manifest", manifestVersion)
	fileEntry := manRoot.CreateElement("manifest:file-entry")
	fileEntry.CreateAttr("manifest:full-path", "/")
	fileEntry.CreateAttr("manifest:media-type", "application/vnd.etsi.asic-e+zip")
	docEntry := manRoot.CreateElement("manifest:file-entry")
	docEntry.CreateAttr("manifest:full-path", docName)
	docEntry.CreateAttr("manifest:media-type", "application/pdf")
	// =========================================================================
	// RESULTADO FINAL
	// =========================================================================
	manifest.WriteSettings.CanonicalEndTags = true
	manifest.Indent(2) // Indentación para que se vea bonito (pretty print)
	// Nota: La indentación del XML final no rompe la firma PORQUE SignedInfo
	// se procesó y firmó internamente sin indentación.

	xmlString, err := manifest.WriteToString()
	if err != nil {
		return err
	}

	// Escribir XML
	if _, err := f.Write([]byte(xmlString)); err != nil {
		fmt.Println("error al escribir el archivo")
		return err
	}
	return nil
}
