package objects

import (
	"fmt"
	"os"
	"sfback/configs"
	"sfback/db"

	"strings"
	"time"

	"sfback/utilities"

	"github.com/beevik/etree"
	//"github.com/google/uuid"
)

const (
	NsXMLDSig       = "http://www.w3.org/2000/09/xmldsig#"
	NsXAdES         = "http://uri.etsi.org/01903/v1.3.2#"
	NsASiC          = "http://uri.etsi.org/02918/v1.2.1#" // Namespace ASiC
	AlgSHA256       = "http://www.w3.org/2001/04/xmlenc#sha256"
	AlgRSA_SHA256   = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
	AlgExcC14N      = "http://www.w3.org/2001/10/xml-exc-c14n#" // Canonicalización Exclusiva
	TypeSignedProps = "http://uri.etsi.org/01903#SignedProperties"
)

func (k *Keys) GenerarFirmaXades(pdfHash []byte, signatureId, docName, idDocument string) (map[string]string, error) {
	fmt.Println("Generando Firma XADES")
	// 1. Obtención de datos previos (Igual a tu código original)
	signInvData, err := db.DB_con.GenericJoinSelect("signatures", "invites", "signatures.idInvite_fk = invites.idInvite", "idSignature", []string{signatureId}, []string{"idUser_fk"}, []string{"idInvite", "idFolder", "idUserOwnner_fk"})
	if err != nil {
		return nil, err
	}

	idFolder := signInvData[signatureId]["idFolder"]
	idUserDest := signInvData[signatureId]["idUser_fk"]

	attrs := []string{"nameUser", "lastNameUser", "emailUser", "idInstitution_fk", "activeUser"}
	wheres := map[string][]string{"idUser": {signInvData[signatureId]["idUserOwnner_fk"], idUserDest}}
	usersData, err := db.DB_con.GenericSelect("users", "idUser", attrs, wheres)
	if err != nil {
		return nil, err
	}

	invitePath := fmt.Sprintf("%s/%s/%s/%s/META-INF/", os.Getenv("GENERIC_FOLDER_PATH"), usersData[idUserDest]["idInstitution_fk"], idFolder, idDocument)
	fmt.Println("Invite path -> ", invitePath)
	// Crear Manifiesto con versión 1.2
	err = CreateManifest(docName, signatureId, invitePath)
	if err != nil {
		return nil, err
	}

	// Identificadores consistentes
	sigId := "id-" + signatureId
	refPdfId := "r-" + sigId + "-1"
	signedPropsId := "xades-" + sigId
	pdfHashB64 := utilities.Encode_b64(pdfHash)

	// =========================================================================
	// PASO 1: CONSTRUIR SignedProperties (XAdES-BES)
	// =========================================================================
	certHashB64, _ := utilities.GetHash(k.Certificate.Raw, configs.HashConf, false)
	certB64 := utilities.Encode_b64(k.Certificate.Raw)

	spDoc := etree.NewDocument()
	spRoot := spDoc.CreateElement("xades:SignedProperties")
	spRoot.CreateAttr("xmlns:xades", NsXAdES)
	spRoot.CreateAttr("xmlns:ds", NsXMLDSig)
	spRoot.CreateAttr("Id", signedPropsId)

	// 1.1 SignedSignatureProperties
	ssp := spRoot.CreateElement("xades:SignedSignatureProperties")
	gentime := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	ssp.CreateElement("xades:SigningTime").SetText(gentime)

	// Cambiado a SigningCertificateV2
	scv2 := ssp.CreateElement("xades:SigningCertificateV2")
	cert := scv2.CreateElement("xades:Cert")
	digest := cert.CreateElement("xades:CertDigest")
	digest.CreateElement("ds:DigestMethod").CreateAttr("Algorithm", AlgSHA256)
	digest.CreateElement("ds:DigestValue").SetText(certHashB64)

	// IssuerSerialV2 (Opcional simplificado, la estructura IssuerSerial está obsoleta)
	isv2 := cert.CreateElement("xades:IssuerSerialV2")
	issuerName := k.Certificate.Issuer.String()
	isv2.CreateElement("ds:X509IssuerName").SetText(issuerName)
	// El Serial Number debe ser el número decimal
	serialNumber := k.Certificate.SerialNumber.String()
	isv2.CreateElement("ds:X509SerialNumber").SetText(serialNumber)

	// 1.2 SignedDataObjectProperties (REQUERIDO para ASiC-E)
	sdop := spRoot.CreateElement("xades:SignedDataObjectProperties")
	dof := sdop.CreateElement("xades:DataObjectFormat")
	dof.CreateAttr("ObjectReference", "#"+refPdfId)
	dof.CreateElement("xades:MimeType").SetText("application/pdf")

	// Canonicalizar SignedProperties para el Hash
	spDoc.WriteSettings.CanonicalEndTags = true
	spDoc.WriteSettings.CanonicalText = true
	signedPropsString, _ := spDoc.WriteToString()
	spHashB64, err := utilities.GetHashC14N(signedPropsString, signedPropsId, "sha256")
	if err != nil {
		fmt.Println("error en canonización", err)
	}

	// =========================================================================
	// PASO 2: CONSTRUIR SignedInfo
	// =========================================================================
	siDoc := etree.NewDocument()
	siRoot := siDoc.CreateElement("ds:SignedInfo")
	siRoot.CreateAttr("xmlns:ds", NsXMLDSig)

	// Canonicalización Exclusiva
	siRoot.CreateElement("ds:CanonicalizationMethod").CreateAttr("Algorithm", AlgExcC14N)
	siRoot.CreateElement("ds:SignatureMethod").CreateAttr("Algorithm", AlgRSA_SHA256)

	// Referencia al PDF (Sin transformaciones según recomendación)
	refDoc := siRoot.CreateElement("ds:Reference")
	refDoc.CreateAttr("Id", refPdfId)
	refDoc.CreateAttr("URI", docName)
	refDoc.CreateElement("ds:DigestMethod").CreateAttr("Algorithm", AlgSHA256)
	refDoc.CreateElement("ds:DigestValue").SetText(pdfHashB64)

	// Referencia a SignedProperties
	refProp := siRoot.CreateElement("ds:Reference")
	refProp.CreateAttr("Type", TypeSignedProps)
	refProp.CreateAttr("URI", "#"+signedPropsId)

	// Se recomienda añadir transform de C14N exclusiva a la referencia interna
	tps := refProp.CreateElement("ds:Transforms")
	tps.CreateElement("ds:Transform").CreateAttr("Algorithm", AlgExcC14N)
	refProp.CreateElement("ds:DigestMethod").CreateAttr("Algorithm", AlgSHA256)
	refProp.CreateElement("ds:DigestValue").SetText(spHashB64)

	// Firmar SignedInfo
	siDoc.WriteSettings.CanonicalEndTags = true
	siBytes, _ := siDoc.WriteToBytes()
	siHash, _ := utilities.GetHash(siBytes, configs.HashConf, false)
	sigBytes, _ := k.SignHash(utilities.Decode_b64(siHash), "sha256")
	signatureValueB64 := utilities.Encode_b64(sigBytes)

	// =========================================================================
	// PASO 3: ARMAR XML FINAL (<asic:XAdESSignatures>)
	// =========================================================================
	finalDoc := etree.NewDocument()
	finalDoc.WriteSettings = etree.WriteSettings{
		CanonicalEndTags: true,
	}
	finalDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8" standalone="no"`)
	asicRoot := finalDoc.CreateElement("asic:XAdESSignatures")
	asicRoot.CreateAttr("xmlns:asic", NsASiC)
	asicRoot.CreateAttr("xmlns:ds", NsXMLDSig)
	asicRoot.CreateAttr("xmlns:xades", NsXAdES)

	sig := asicRoot.CreateElement("ds:Signature")
	sig.CreateAttr("Id", sigId)
	sig.AddChild(siRoot)

	sv := sig.CreateElement("ds:SignatureValue")
	sv.CreateAttr("Id", "value-"+sigId)
	sv.SetText(signatureValueB64)

	ki := sig.CreateElement("ds:KeyInfo")
	xd := ki.CreateElement("ds:X509Data")
	xd.CreateElement("ds:X509Certificate").SetText(certB64)

	obj := sig.CreateElement("ds:Object")
	qp := obj.CreateElement("xades:QualifyingProperties")
	qp.CreateAttr("Target", "#"+sigId)
	qp.AddChild(spRoot)

	// Guardar y Finalizar (Igual que tu flujo original)
	xmlFinal, _ := finalDoc.WriteToBytes()

	// crear xml
	xmlPath := invitePath + signatureId + ".xml"
	if err := os.MkdirAll(invitePath, 0755); err != nil {
		fmt.Println("error al crear el folder")
		return nil, err
	}
	fmt.Println("xml path -> ", xmlPath)
	// Crear archivo solo si NO existe (modo seguro)
	f, err := os.OpenFile(xmlPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fmt.Println("el archivo ya existe: ", xmlPath, err)
	} else {
		// Escribir XML
		if _, err := f.Write(xmlFinal); err != nil {
			fmt.Println("error al escribir el archivo")
			return nil, err
		}
	}
	defer f.Close()

	// crear qrcode con link
	qrPath := strings.ReplaceAll(invitePath, "META-INF/", "") + signatureId + ".png"

	if !utilities.GenerateQR(fmt.Sprintf("%sviewSignature?id=%s", os.Getenv("QR_URL_BASE"), signatureId), qrPath) {
		return nil, fmt.Errorf("no se pudo crear qr de firma %s", qrPath)
	}

	if ok := utilities.AddTextToQR(qrPath, fmt.Sprintf("%s %s \n %s", usersData[idUserDest]["nameUser"], usersData[idUserDest]["lastNameUser"], gentime)); !ok {
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
		"signatureValueSign": signatureValueB64,
		"genTimeSign":        gentime,
		"nonceSign":          signatureId,
		"certHash":           certHashB64,
		"typeSign":           "9", // rsa Signature con hash 256 -> será mapeada para front
	}

	return signData, nil
}

func CreateManifest(docName, idSignature, folderPath string) error {
	xmlPath := folderPath + "manifest.xml"
	manifest := etree.NewDocument()

	var manRoot *etree.Element

	// Intentar abrir el archivo existente
	if err := manifest.ReadFromFile(xmlPath); err == nil {
		// Si existe, buscamos el nodo raíz
		manRoot = manifest.SelectElement("manifest:manifest")
	}

	// Si el archivo no existe o no tiene el nodo raíz, creamos la base
	if manRoot == nil {
		manifest = etree.NewDocument() // Reiniciar por si ReadFromFile dejó algo
		manRoot = manifest.CreateElement("manifest:manifest")
		manRoot.CreateAttr("xmlns:manifest", "urn:oasis:names:tc:opendocument:xmlns:manifest:1.0")
		manRoot.CreateAttr("manifest:version", "1.2")

		// 1. Entrada raíz (/)
		eRoot := manRoot.CreateElement("manifest:file-entry")
		eRoot.CreateAttr("manifest:full-path", "/")
		eRoot.CreateAttr("manifest:media-type", "application/vnd.etsi.asic-e+zip")

		// 2. Entrada del Documento PDF original
		eDoc := manRoot.CreateElement("manifest:file-entry")
		eDoc.CreateAttr("manifest:full-path", docName)
		eDoc.CreateAttr("manifest:media-type", "application/pdf")
	}

	// 3. Añadir la entrada de la nueva firma específica
	// Verificamos si ya existe para evitar duplicados (opcional pero recomendado)
	signaturePath := fmt.Sprintf("META-INF/%s.xml", idSignature)
	exists := false
	for _, entry := range manRoot.SelectElements("manifest:file-entry") {
		if entry.SelectAttrValue("manifest:full-path", "") == signaturePath {
			exists = true
			break
		}
	}

	if !exists {
		eSig := manRoot.CreateElement("manifest:file-entry")
		eSig.CreateAttr("manifest:full-path", signaturePath)
		eSig.CreateAttr("manifest:media-type", "application/x-erts-sig+xml")
	}

	// Asegurarse de que el directorio existe
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	return manifest.WriteToFile(xmlPath)
}
