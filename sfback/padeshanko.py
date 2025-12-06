# sign_pades_b_pyasn1.py
import re
import os
import tempfile
import base64
from datetime import datetime, timezone

from pyasn1.codec.der.decoder import decode as asn1_decode
from pyasn1.codec.der.encoder import encode as asn1_encode
from pyasn1_modules import rfc2459, rfc5652, rfc5035
from pyasn1.type import univ, useful

from cryptography.hazmat.primitives import serialization, hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.backends import default_backend

from PyPDF2 import PdfReader, PdfWriter
from PyPDF2.generic import NameObject, DictionaryObject, ArrayObject, NumberObject, ByteStringObject

# --------------------- Helpers: cert load via pyasn1 ---------------------
def load_cert_pyasn1(cert_path: str):
    """
    Lee certificado PEM/DER y devuelve (cert_asn1_object, der_bytes)
    Usa pyasn1 modules rfc2459.Certificate como asn1Spec para parsing tolerante.
    """
    with open(cert_path, "rb") as f:
        data = f.read()

    if data.strip().startswith(b"-----BEGIN CERTIFICATE-----"):
        # limpiar PEM
        body = data.replace(b"-----BEGIN CERTIFICATE-----", b"").replace(b"-----END CERTIFICATE-----", b"")
        body = b"".join(body.splitlines())
        der = base64.b64decode(body)
    else:
        der = data

    cert_asn1, rest = asn1_decode(der, asn1Spec=rfc2459.Certificate())
    if rest:
        # puede haber datos extra; normalmente se ignoran
        pass
    return cert_asn1, der

# --------------------- Helpers: ESS / SigningCertificateV2 (pyasn1) ---------------------

def build_esscertidv2(cert_der: bytes):
    """
    Construye ESSCertIDv2 (pyasn1) y devuelve el objeto ESSCertIDv2 (no codificado).
    """
    # hashAlgorithm: AlgorithmIdentifier for sha256 -> OID 2.16.840.1.101.3.4.2.1
    alg = rfc2459.AlgorithmIdentifier()
    alg['algorithm'] = univ.ObjectIdentifier('2.16.840.1.101.3.4.2.1')  # sha256

    # certHash: OCTET STRING of SHA256(cert_der)
    digest = hashes.Hash(hashes.SHA256(), backend=default_backend())
    digest.update(cert_der)
    cert_hash = digest.finalize()

    ess = rfc5035.ESSCertIDv2()
    ess.setComponentByName('hashAlgorithm', alg)
    ess.setComponentByName('certHash', univ.OctetString(cert_hash))
    # omit issuerSerial for simplicity
    return ess

def build_signing_cert_v2_struct(cert_der: bytes):
    """
    Devuelve SigningCertificateV2 ASN.1 structure (rfc5035.SigningCertificateV2)
    usando una SEQUENCE que contiene el ESSCertIDv2.
    """
    ess = build_esscertidv2(cert_der)

    # ESSCertIDv2Seq es SEQUENCE OF ESSCertIDv2; usar univ.Sequence y setComponentByPosition es correcto
    seq = univ.Sequence()
    seq.setComponentByPosition(0, ess)

    scv2 = rfc5035.SigningCertificateV2()
    scv2.setComponentByName('certs', seq)
    # policy is optional; omit
    return scv2

# --------------------- Helpers: SignedAttributes builder ---------------------
def build_signed_attributes(content_digest: bytes, cert_der: bytes):
    """
    Construye SET OF Attributes correctamente para pyasn1.
    Hace:
      - content-type (id-data)
      - message-digest
      - signing-time
      - signingCertificateV2
    Retorna un rfc5652.SignedAttributes (SET OF Attribute) listo para encoder.
    """
    attrs_set = rfc5652.SignedAttributes()  # es el tipo SET OF Attribute
    idx = 0

    # 1) content-type = id-data
    attr = rfc5652.Attribute()
    attr['attrType'] = rfc5652.id_contentType  # u otro nombre según tu versión (usa el que tengas)
    # Crear un SET (valores) y poner ObjectIdentifier('1.2.840.113549.1.7.1') dentro
    vals = univ.Set()
    vals.setComponentByPosition(0, univ.ObjectIdentifier('1.2.840.113549.1.7.1'))  # id-data
    attr['attrValues'] = vals
    attrs_set.setComponentByPosition(idx, attr)
    idx += 1

    # 2) message-digest
    attr = rfc5652.Attribute()
    attr['attrType'] = rfc5652.id_messageDigest
    vals = univ.Set()
    vals.setComponentByPosition(0, univ.OctetString(content_digest))
    attr['attrValues'] = vals
    attrs_set.setComponentByPosition(idx, attr)
    idx += 1

    # 3) signing-time (UTCTime)
    attr = rfc5652.Attribute()
    attr['attrType'] = rfc5652.id_signingTime
    vals = univ.Set()
    st = useful.UTCTime(datetime.now(timezone.utc).strftime("%y%m%d%H%M%SZ"))
    vals.setComponentByPosition(0, st)
    attr['attrValues'] = vals
    attrs_set.setComponentByPosition(idx, attr)
    idx += 1

    # 4) signingCertificateV2 OID (rfc5035.id_aa_signingCertificateV2)
    scv2 = build_signing_cert_v2_struct(cert_der)
    attr = rfc5652.Attribute()
    attr['attrType'] = rfc5035.id_aa_signingCertificateV2
    vals = univ.Set()
    vals.setComponentByPosition(0, scv2)
    attr['attrValues'] = vals
    attrs_set.setComponentByPosition(idx, attr)
    idx += 1

    return attrs_set

# --------------------- Helpers: Build CMS SignedData (pyasn1) ---------------------
def build_signed_data_cms_detached(signed_attrs, signer_cert_asn1, signer_der_bytes, signature_bytes):
    """
    Construye ContentInfo -> SignedData con:
      - certificates: el certificado del firmante
      - signerInfos: un SignerInfo con signedAttrs (ya firmado)
    NOTE: signed_attrs must be the pyasn1 object (Attributes) — we will include it (not DER here).
    """

    # Digest alg identifier (sha256)
    digest_alg = rfc2459.AlgorithmIdentifier()
    digest_alg['algorithm'] = univ.ObjectIdentifier('2.16.840.1.101.3.4.2.1')  # sha256

    # signature alg: rsaEncryption OID 1.2.840.113549.1.1.1 (we'll rely on this)
    sig_alg = rfc2459.AlgorithmIdentifier()
    sig_alg['algorithm'] = univ.ObjectIdentifier('1.2.840.113549.1.1.1')  # rsaEncryption

    # Build SignerIdentifier (issuerAndSerialNumber)
    issuer_and_serial = rfc5652.IssuerAndSerialNumber()
    issuer_and_serial['issuer'] = signer_cert_asn1.getComponentByName('tbsCertificate').getComponentByName('issuer')
    issuer_and_serial['serialNumber'] = signer_cert_asn1.getComponentByName('tbsCertificate').getComponentByName('serialNumber')

    sid = rfc5652.SignerIdentifier()
    sid['issuerAndSerialNumber'] = issuer_and_serial

    # SignerInfo
    si = rfc5652.SignerInfo()
    si['version'] = 1
    si['sid'] = sid
    si['digestAlgorithm'] = digest_alg
    si['signedAttrs'] = signed_attrs  # pyasn1 object (SET OF Attributes)
    si['signatureAlgorithm'] = sig_alg
    si['signature'] = univ.OctetString(signature_bytes)
    # unsignedAttrs omitted

    # Put si in a SET OF SignerInfo (actually SignerInfos is a SET OF SignerInfo)
    sis = rfc5652.SignerInfos()
    sis.setComponentByPosition(0, si)

    # Certificates (wrap the raw certificate DER into Certificate structure)
    certs = rfc5652.CertificateSet()
    # CertificateChoices -> we use 'certificate' choice with rfc2459.Certificate
    cc = rfc5652.CertificateChoices()
    # Load signer_cert_asn1 is rfc2459.Certificate instance; put directly
    cc.setComponentByName('certificate', signer_cert_asn1)
    certs.setComponentByPosition(0, cc)

    # Encapsulated content info: data contentType but no eContent (detached)
    eci = rfc5652.EncapsulatedContentInfo()
    eci['eContentType'] = univ.ObjectIdentifier('1.2.840.113549.1.7.1')  # id-data
    # eContent omitted (detached)

    # SignedData
    sd = rfc5652.SignedData()
    sd['version'] = 1
    # digestAlgorithms is a SET OF DigestAlgorithmIdentifier
    da_set = rfc5652.DigestAlgorithmIdentifiers()
    da_set.setComponentByPosition(0, digest_alg)
    sd['digestAlgorithms'] = da_set
    sd['encapContentInfo'] = eci
    sd['certificates'] = certs
    sd['crls'] = univ.Null()  # absent
    sd['signerInfos'] = sis

    ci = rfc5652.ContentInfo()
    ci['contentType'] = univ.ObjectIdentifier('1.2.840.113549.1.7.2')  # signedData
    ci['content'] = sd

    return asn1_encode(ci)  # DER bytes of ContentInfo

# --------------------- PDF helper (create unsigned with placeholders) ---------------------
def _create_unsigned_pdf_with_placeholders(pdf_in, temp_unsigned_path, placeholder_len=131072, reserved_br_ascii=512):
    reader = PdfReader(pdf_in)
    writer = PdfWriter()
    for page in reader.pages:
        writer.add_page(page)

    sig_dict = DictionaryObject()
    sig_dict[NameObject("/Type")] = NameObject("/Sig")
    sig_dict[NameObject("/Filter")] = NameObject("/Adobe.PPKLite")
    sig_dict[NameObject("/SubFilter")] = NameObject("/adbe.pkcs7.detached")
    sig_dict[NameObject("/ByteRange")] = ArrayObject([NumberObject(0), NumberObject(0), NumberObject(0), NumberObject(0)])
    sig_dict[NameObject("/Contents")] = ByteStringObject(b"\x00" * placeholder_len)
    sig_dict[NameObject("/M")] = ByteStringObject(datetime.now(timezone.utc).strftime("D:%Y%m%d%H%M%S+00'00'").encode())
    sig_dict[NameObject("/Reason")] = ByteStringObject("Firma PAdES-B (pyasn1)".encode())

    widget = DictionaryObject()
    widget[NameObject("/Type")] = NameObject("/Annot")
    widget[NameObject("/Subtype")] = NameObject("/Widget")
    widget[NameObject("/FT")] = NameObject("/Sig")
    widget[NameObject("/T")] = ByteStringObject("Signature1".encode())
    widget[NameObject("/Rect")] = ArrayObject([NumberObject(0), NumberObject(0), NumberObject(0), NumberObject(0)])
    widget[NameObject("/V")] = sig_dict

    widget_ref = writer._add_object(widget)
    page0 = writer.pages[0]
    if "/Annots" not in page0:
        page0[NameObject("/Annots")] = ArrayObject()
    page0[NameObject("/Annots")].append(widget_ref)

    root = writer._root_object
    if "/AcroForm" not in root:
        root[NameObject("/AcroForm")] = DictionaryObject()
    if "/Fields" not in root[NameObject("/AcroForm")]:
        root[NameObject("/AcroForm")][NameObject("/Fields")] = ArrayObject()
    root[NameObject("/AcroForm")][NameObject("/Fields")].append(widget_ref)

    # write file (ensure closed handle)
    with open(temp_unsigned_path, "wb") as f:
        writer.write(f)

# --------------------- Main function ---------------------
def sign_pdf_pades(pdf_in, pdf_out, cert_path, key_path):
    """
    Firma PAdES-B usando pyasn1 (para certificados difíciles) + cryptography (firma).
    """
    # 1) load cert via pyasn1 and keep DER
    cert_asn1, cert_der = load_cert_pyasn1(cert_path)

    # 2) load private key (PEM or DER)
    with open(key_path, "rb") as f:
        key_data = f.read()
    try:
        private_key = serialization.load_pem_private_key(key_data, password=None, backend=default_backend())
    except Exception:
        private_key = serialization.load_der_private_key(key_data, password=None, backend=default_backend())

    # 3) create unsigned pdf with placeholders
    PLACEHOLDER_LEN = 131072  # 128 KiB for Contents
    RESERVED_BR_ASCII = 512
    with tempfile.NamedTemporaryFile(delete=False, suffix=".unsigned.pdf") as tmp:
        temp_unsigned = tmp.name
    _create_unsigned_pdf_with_placeholders(pdf_in, temp_unsigned, placeholder_len=PLACEHOLDER_LEN, reserved_br_ascii=RESERVED_BR_ASCII)

    # 4) read bytes, reserve ByteRange ASCII area
    pdf_bytes = open(temp_unsigned, "rb").read()
    br_match = re.search(br"/ByteRange\s*\[(.*?)\]", pdf_bytes)
    if not br_match:
        os.unlink(temp_unsigned)
        raise Exception("No ByteRange found in unsigned PDF")

    # Reserve ASCII space for ByteRange
    pdf_bytes = pdf_bytes[:br_match.start(1)] + (b" " * RESERVED_BR_ASCII) + pdf_bytes[br_match.end(1):]
    with open(temp_unsigned, "wb") as f:
        f.write(pdf_bytes)

    # 5) find Contents placeholder and compute digest excluding it
    pdf_bytes = open(temp_unsigned, "rb").read()
    c_match = re.search(br"/Contents\s*<([0-9A-Fa-f]*)>", pdf_bytes)
    if not c_match:
        os.unlink(temp_unsigned)
        raise Exception("No Contents placeholder found in unsigned PDF")

    hex_start = c_match.start(1)
    hex_end = c_match.end(1)
    contents_start = hex_start - 1
    contents_end = hex_end + 1

    # compute SHA256 digest of the PDF excluding the <...> contents
    digest_ctx = hashes.Hash(hashes.SHA256(), backend=default_backend())
    digest_ctx.update(pdf_bytes[:contents_start])
    digest_ctx.update(pdf_bytes[contents_end:])
    pdf_digest = digest_ctx.finalize()

    # 6) build SignedAttributes (pyasn1) including signingCertificateV2
    signed_attrs = build_signed_attributes(pdf_digest, cert_der)

    # 7) DER-encode signed_attrs (we will sign this DER). According to RFC, the signature is over
    #    the DER encoding of the signedAttrs (SET OF Attributes). pyasn1 encoder returns appropriate bytes.
    signed_attrs_der = asn1_encode(signed_attrs)

    # 8) sign the signed_attrs_der with RSA/SHA256 (PKCS#1 v1.5)
    signature = private_key.sign(
        signed_attrs_der,
        padding.PKCS1v15(),
        hashes.SHA256()
    )

    # 9) build SignedData CMS (pyasn1) including certificate and signerinfo
    cms_der = build_signed_data_cms_detached(signed_attrs, cert_asn1, cert_der, signature)

    # 10) insert CMS DER (HEX) into Contents placeholder
    cms_hex = cms_der.hex().upper().encode()
    available_hex_len = hex_end - hex_start
    if len(cms_hex) > available_hex_len:
        os.unlink(temp_unsigned)
        raise Exception(f"CMS too large for placeholder: {len(cms_hex)} > {available_hex_len}")

    padded_hex = cms_hex + b"0" * (available_hex_len - len(cms_hex))
    pdf_bytes = pdf_bytes[:hex_start] + padded_hex + pdf_bytes[hex_end:]

    # 11) compute final ByteRange and write into reserved ASCII area
    b0 = 0
    b1 = contents_start
    b2 = contents_end
    b3 = len(pdf_bytes) - b2
    final_br = f"{b0} {b1} {b2} {b3}".encode()
    if len(final_br) > RESERVED_BR_ASCII:
        os.unlink(temp_unsigned)
        raise Exception("Final ByteRange ascii length exceeds reserved area. Increase RESERVED_BR_ASCII.")

    final_br_padded = final_br + b" " * (RESERVED_BR_ASCII - len(final_br))
    pdf_bytes = pdf_bytes[:br_match.start(1)] + final_br_padded + pdf_bytes[br_match.start(1) + RESERVED_BR_ASCII:]

    # 12) write final PDF
    with open(pdf_out, "wb") as f:
        f.write(pdf_bytes)

    os.unlink(temp_unsigned)
    print("PDF firmado (PAdES-B, pyasn1) ->", pdf_out)

# ------------------------------------------------------------------
# CLI
# ------------------------------------------------------------------
if __name__ == "__main__":
    pdf_in="C:/Users/USER/Desktop/CV_2025_JJNM.pdf"
    pdf_out="C:/Users/USER/Desktop/firmado_pades.pdf"
    cert_pem="C:/Users/USER/Desktop/certs/numj900112t99.pem"
    key_pem="C:/Users/USER/Desktop/certs/Claveprivada_FIEL_NUMJ900112T99_20220707_112849.pem"
    sign_pdf_pades(pdf_in, pdf_out, cert_pem, key_pem)