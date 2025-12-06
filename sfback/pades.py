'''openssl pkcs12 -export -in numj900112t99.pem -inkey Claveprivada_FIEL_NUMJ900112T99_20220707_112849.pem -out numj900112t99.p12 -passout pass:Xpressiceover1'''
import subprocess
import tempfile
import re
import os
from PyPDF2 import PdfReader, PdfWriter
from PyPDF2.generic import NameObject, DictionaryObject, ArrayObject, NumberObject, ByteStringObject


def sign_pdf_pades_manual(pdf_in, pdf_out, cert_pem, key_pem):
    # ===============================================================
    # 1) Crear PDF UNSIGNED con placeholder de firma
    # ===============================================================

    PLACEHOLDER_LEN = 131072   # 128 KB de hueco — suficiente para firmas grandes
    PLACEHOLDER_BYTES = b"\x00" * PLACEHOLDER_LEN

    reader = PdfReader(pdf_in)
    writer = PdfWriter()

    for page in reader.pages:
        writer.add_page(page)

    # --- Construcción del diccionario de firma ---
    sig_dict = DictionaryObject()
    sig_dict[NameObject("/Type")] = NameObject("/Sig")
    sig_dict[NameObject("/Filter")] = NameObject("/Adobe.PPKLite")
    sig_dict[NameObject("/SubFilter")] = NameObject("/adbe.pkcs7.detached")
    sig_dict[NameObject("/ByteRange")] = ArrayObject([NumberObject(0), NumberObject(0),
                                                      NumberObject(0), NumberObject(0)])
    sig_dict[NameObject("/Contents")] = ByteStringObject(PLACEHOLDER_BYTES)
    sig_dict[NameObject("/M")] = ByteStringObject("D:20250101000000+00'00'".encode())
    sig_dict[NameObject("/Reason")] = ByteStringObject("Firma PAdES".encode())

    # Widget
    widget = DictionaryObject()
    widget[NameObject("/Type")] = NameObject("/Annot")
    widget[NameObject("/Subtype")] = NameObject("/Widget")
    widget[NameObject("/FT")] = NameObject("/Sig")
    widget[NameObject("/T")] = ByteStringObject("Signature1".encode())
    widget[NameObject("/Rect")] = ArrayObject([NumberObject(0), NumberObject(0),
                                               NumberObject(0), NumberObject(0)])
    widget[NameObject("/V")] = sig_dict

    # Añadir al writer
    widget_ref = writer._add_object(widget)

    # Página
    page = writer.pages[0]
    if "/Annots" not in page:
        page[NameObject("/Annots")] = ArrayObject()
    page[NameObject("/Annots")].append(widget_ref)

    # AcroForm
    root = writer._root_object
    if "/AcroForm" not in root:
        root[NameObject("/AcroForm")] = DictionaryObject()

    if "/Fields" not in root[NameObject("/AcroForm")]:
        root[NameObject("/AcroForm")][NameObject("/Fields")] = ArrayObject()

    root[NameObject("/AcroForm")][NameObject("/Fields")].append(widget_ref)

    # Guardar PDF unsigned temporal
    with tempfile.NamedTemporaryFile(delete=False) as tmp:
        temp_unsigned_path = tmp.name
        # IMPORTANTÍSIMO: el archivo queda CERRADO tras salir del "with"
        pass
    with open(temp_unsigned_path, "wb") as f:
        writer.write(f)

    # ===============================================================
    # 2) ABRIR EL UNSIGNED Y RESERVAR ESPACIO PARA EL BYTERANGE
    # ===============================================================
    pdf_bytes = open(temp_unsigned_path, "rb").read()

    # Encontrar el ByteRange
    br = re.search(br"/ByteRange\s*\[(.*?)\]", pdf_bytes)
    if not br:
        os.unlink(temp_unsigned_path)
        raise Exception("No se encontró ByteRange después de crear el PDF.")

    # Reservar 256 bytes ASCII dentro del Array del ByteRange
    RESERVED_BR_ASCII = 256
    reserved_ascii = b" " * RESERVED_BR_ASCII

    pdf_bytes = (
        pdf_bytes[:br.start(1)]
        + reserved_ascii
        + pdf_bytes[br.end(1):]
    )

    # Reescribir el unsigned (ahora con espacio real para ByteRange)
    with open(temp_unsigned_path, "wb") as f:
        f.write(pdf_bytes)

    # ===============================================================
    # 3) Localizar el hueco /Contents y calcular hash excluyéndolo
    # ===============================================================
    pdf_bytes = open(temp_unsigned_path, "rb").read()

    contents = re.search(br"/Contents\s*<([0-9A-Fa-f]+)>", pdf_bytes)
    if not contents:
        os.unlink(temp_unsigned_path)
        raise Exception("No se encontró campo /Contents en el PDF.")

    hex_start = contents.start(1)
    hex_end   = contents.end(1)

    contents_start = contents.start(0) + pdf_bytes[contents.start(0):].find(b"<")
    contents_end   = contents_start + (hex_end - hex_start) + 2  # incluye "< >"

    # Hash
    from cryptography.hazmat.primitives import hashes
    h = hashes.Hash(hashes.SHA256())
    h.update(pdf_bytes[:contents_start])
    h.update(pdf_bytes[contents_end:])
    pdf_hash = h.finalize()

    # ===============================================================
    # 4) GENERAR PKCS#7/CMS detached CON OPENSSL
    # ===============================================================
    tmp_hash = tempfile.NamedTemporaryFile(delete=False)
    tmp_hash.write(pdf_hash)
    tmp_hash.close()

    cms_path = tempfile.NamedTemporaryFile(delete=False).name

    cmd = [
        "openssl", "smime", "-sign",
        "-binary",
        "-nocerts",
        "-noattr",
        "-in", tmp_hash.name,
        "-signer", cert_pem,
        "-inkey", key_pem,
        "-outform", "DER",
        "-out", cms_path,
    ]

    subprocess.run(cmd, check=True)

    cms_bytes = open(cms_path, "rb").read()

    os.unlink(tmp_hash.name)
    os.unlink(cms_path)

    # ===============================================================
    # 5) INSERTAR CMS en el /Contents (en HEX)
    # ===============================================================
    cms_hex = cms_bytes.hex().upper().encode()

    # tamaño disponible dentro del <...>
    hueco_hex_len = hex_end - hex_start

    if len(cms_hex) > hueco_hex_len:
        os.unlink(temp_unsigned_path)
        raise Exception(
            f"El CMS ({len(cms_hex)} chars) no cabe en el hueco ({hueco_hex_len} chars). "
            f"Aumenta PLACEHOLDER_LEN."
        )

    padded_hex = cms_hex + b"0" * (hueco_hex_len - len(cms_hex))

    pdf_bytes = pdf_bytes[:hex_start] + padded_hex + pdf_bytes[hex_end:]

    # ===============================================================
    # 6) ESCRIBIR EL BYTERANGE REAL dentro del espacio reservado
    # ===============================================================
    # ByteRange estructura: 0 b1 b2 b3
    b0 = 0
    b1 = contents_start
    b2 = contents_end
    b3 = len(pdf_bytes) - b2

    final_br = f"{b0} {b1} {b2} {b3}".encode()

    if len(final_br) > RESERVED_BR_ASCII:
        os.unlink(temp_unsigned_path)
        raise Exception("El ByteRange final NO cabe en el espacio reservado!! Sube RESERVED_BR_ASCII.")

    # Rellenar con espacios
    final_br_padded = final_br + b" " * (RESERVED_BR_ASCII - len(final_br))

    # Reemplazar dentro de los bytes
    pdf_bytes = (
        pdf_bytes[:br.start(1)]
        + final_br_padded
        + pdf_bytes[br.start(1) + RESERVED_BR_ASCII:]
    )

    # ===============================================================
    # 7) GUARDAR EL PDF FINAL
    # ===============================================================
    with open(pdf_out, "wb") as f:
        f.write(pdf_bytes)

    os.unlink(temp_unsigned_path)

    print("ÉXITO: PDF firmado manualmente (PAdES-BASIC).")

if __name__ == "__main__":
    sign_pdf_pades_manual(
        "C:/Users/USER/Desktop/CV_2025_JJNM.pdf", 
        "C:/Users/USER/Desktop/firmado.pdf", 
        "C:/Users/USER/Desktop/certs/numj900112t99.pem", 
        "C:/Users/USER/Desktop/certs/Claveprivada_FIEL_NUMJ900112T99_20220707_112849.pem",
    )
    print("PDF firmado correctamente!")




