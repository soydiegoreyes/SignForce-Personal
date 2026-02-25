import hashlib
import json
import os
from datetime import datetime as dt, timedelta

CACHE_DIR = os.path.join(os.path.dirname(__file__), "cache")
CACHE_TTL_HOURS = 24


def _cache_path(doc_path: str, suffix: str = "") -> str:
    """Devuelve la ruta del archivo JSON de caché para un documento.
    El sufijo permite separar caché por método de extracción, ej. '_ocr'.
    """
    key = hashlib.md5(os.path.abspath(doc_path).encode()).hexdigest()
    return os.path.join(CACHE_DIR, f"{key}{suffix}.json")


def get_doc_cache(doc_path: str, suffix: str = "") -> dict | None:
    """
    Devuelve los datos cacheados si existen y son menores a 24h.
    Retorna None si no existe o expiró (borra el archivo expirado en ese caso).
    """
    fpath = _cache_path(doc_path, suffix)
    if not os.path.exists(fpath):
        return None
    try:
        with open(fpath, "r", encoding="utf-8") as f:
            data = json.load(f)
        created = dt.fromisoformat(data["created"])
        if dt.now() - created > timedelta(hours=CACHE_TTL_HOURS):
            os.remove(fpath)
            return None
        return data
    except Exception:
        return None


def save_doc_cache(doc_path: str, text: str, meta: dict, suffix: str = ""):
    """Guarda el texto completo y los metadatos del documento en caché."""
    os.makedirs(CACHE_DIR, exist_ok=True)
    fpath = _cache_path(doc_path, suffix)
    payload = {
        "created": dt.now().isoformat(),
        "path": os.path.abspath(doc_path),
        "text": text,
        "meta": meta,
    }
    with open(fpath, "w", encoding="utf-8") as f:
        json.dump(payload, f, ensure_ascii=False, indent=2)


def cleanup_expired_cache():
    """Borra todos los archivos de caché expirados. Se llama al arrancar el servidor."""
    if not os.path.exists(CACHE_DIR):
        return
    now = dt.now()
    for fname in os.listdir(CACHE_DIR):
        if not fname.endswith(".json"):
            continue
        fpath = os.path.join(CACHE_DIR, fname)
        try:
            with open(fpath, "r", encoding="utf-8") as f:
                data = json.load(f)
            created = dt.fromisoformat(data["created"])
            if now - created > timedelta(hours=CACHE_TTL_HOURS):
                os.remove(fpath)
                print(f"[cache] Eliminado expirado: {fname}")
        except Exception:
            os.remove(fpath)  # archivo corrupto → borrar
