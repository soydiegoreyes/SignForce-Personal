# ═══════════════════════════════════════════════════════════
# SignForce — AES-256-GCM Decryption Utilities
# ═══════════════════════════════════════════════════════════
# Compatible with sfmiddle/utilities/cryptoutils.go
# Key source: DOCS_KEY environment variable (required)
# Algorithm: SHA-256(DOCS_KEY) → 32-byte AES key
# ═══════════════════════════════════════════════════════════

"""
Desencriptación AES-GCM compatible con sfmiddle/utilities/cryptoutils.go.

Algoritmo:
  - Clave: sha256(DOCS_KEY) → 32 bytes
  - Nonce base: primeros 12 bytes del archivo
  - Chunks cifrados de 64*1024 + 16 bytes (datos + tag GCM de 16 bytes)
  - Nonce por chunk: copia del nonce base con los últimos 8 bytes
    reemplazados por el índice del chunk en big-endian uint64
  - Si DOCS_KEY no está configurado o el archivo no es cifrado,
    se devuelven los bytes originales sin modificar.
"""

import hashlib
import os

try:
    from cryptography.hazmat.primitives.ciphers.aead import AESGCM
    _CRYPTO_AVAILABLE = True
except ImportError:
    _CRYPTO_AVAILABLE = False
    print("[crypto] Advertencia: librería 'cryptography' no instalada. "
          "Los archivos cifrados no podrán desencriptarse. "
          "Instalar con: pip install cryptography")

_NONCE_SIZE = 12                     # gcm.NonceSize() en Go
_PLAIN_CHUNK = 64 * 1024             # tamaño del bloque original
_ENC_CHUNK = _PLAIN_CHUNK + 16       # bloque cifrado + 16 bytes de tag GCM


def _build_nonce(base: bytes, index: int) -> bytes:
    """Reconstruye el nonce específico del chunk: base con últimos 8 bytes = índice uint64 BE."""
    n = bytearray(base)
    n[_NONCE_SIZE - 8:] = index.to_bytes(8, "big")
    return bytes(n)


def decrypt_bytes(data: bytes) -> bytes | None:
    """
    Desencripta bytes cifrados con AES-GCM según el esquema de sfmiddle.
    Devuelve los bytes planos, o None si:
      - DOCS_KEY no está disponible
      - la librería cryptography no está instalada
      - el archivo no parece cifrado (tag inválido en el primer chunk)
    """
    if not _CRYPTO_AVAILABLE:
        return None
    key_env = os.getenv("DOCS_KEY", "").encode()
    if not key_env:
        return None
    if len(data) < _NONCE_SIZE + 16:
        return None  # demasiado corto para ser un archivo cifrado válido
    key = hashlib.sha256(key_env).digest()
    aesgcm = AESGCM(key)
    nonce_base = data[:_NONCE_SIZE]
    payload = data[_NONCE_SIZE:]
    plaintext_chunks = []
    offset = 0
    index = 0
    while offset < len(payload):
        chunk = payload[offset:offset + _ENC_CHUNK]
        if not chunk:
            break
        nonce = _build_nonce(nonce_base, index)
        try:
            plain = aesgcm.decrypt(nonce, chunk, None)
        except Exception:
            # Tag inválido → el archivo no está cifrado o la clave es incorrecta
            return None
        plaintext_chunks.append(plain)
        offset += _ENC_CHUNK
        index += 1
    return b"".join(plaintext_chunks)


def read_doc_bytes(path: str) -> bytes:
    """
    Lee un archivo y lo desencripta si es necesario.
    - Si DOCS_KEY está configurado y el archivo está cifrado → devuelve bytes planos.
    - Si el archivo no está cifrado o DOCS_KEY no existe → devuelve bytes tal cual.
    Uso: para pasar a fitz.open(stream=...) o a OCR.
    """
    with open(path, "rb") as f:
        raw = f.read()
    decrypted = decrypt_bytes(raw)
    if decrypted is not None:
        print(f"[crypto] Desencriptado: {path}")
        return decrypted
    print(f"[crypto] Sin cifrado (o clave no disponible): {path}")
    return raw
