import uvicorn
from typing import List, Optional

from fastapi import FastAPI, HTTPException ,Body, Depends, Header
from fastapi.responses import HTMLResponse
from fastapi.middleware.cors import CORSMiddleware

import fitz

import re
import json
import os
from datetime import datetime as dt
import dotenv
import ollama
from agent import run_agent
from doc_cache import get_doc_cache, save_doc_cache, cleanup_expired_cache
from crypto_utils import read_doc_bytes

dotenv.load_dotenv()
#________________________________________________________________________________________________________ 
 # mapa de sustitucion para acentos y diéresis
 # usado para quitar emojis
allowedch={225: 97,
  233: 101, 237: 105, 243: 111, 250: 117, 193: 97,
  201: 101, 205: 105, 211: 111, 218: 117, 228: 97,
  235: 101, 239: 105, 246: 111, 252: 117, 196: 97,
  203: 101, 207: 105, 214: 111, 220: 117, 224: 97, 
  232: 101, 236: 105, 242: 111, 249: 117}

month_parse={'01':'ENERO','02':'FEBRERO','03':'MARZO','04':'ABRIL','05':'MAYO','06':'JUNIO','07':'JULIO','08':'AGOSTO','09':'SEPTIEMBRE','10':'OCTUBRE','11':'NOVIEMBRE','12':'DICIEMBRE'}

def clean_text(text: str) -> str:
    # quita saltos de linea y saltos con puntos comunes en IG y espacios seguidos
    text=re.sub('\n', ' ', text)
    # Bloque para quitar emojis
    # pasamos los caracteres con puntuación a sin puntuación (á->a) depende del diccionaio que se tenga
    #text=text.translate(allowedch)
    # en el mismo paso se convierte a utf-8 y se reemplazan las ñ por ñ (truco para que al quitar caracteres raros no reemplace las ñ)
    #text=re.sub('(\\\\xc3\\\\xb1)|(\\\\xc3\\\\x91)', 'ñ',str(text.encode('utf-8')))
    # se quitan los caracteres raros pero como la ñ no queda con este formato no se reemplaza :P
    #text=re.sub('\\\\x[a-z\d]{2}', '',text)
    # quita los caracteres extra generados en el proceso anterior
    # originalmente son comillas pero varía si la comilla es " o '
    #text=re.findall("^b\W(.*)\W$", text)[0]
    # quita 2 o mas espacios y los deja como un solo espacio
    text=re.sub('[\s\*\+=\|]{2,}', ' ', text)
    return text


#LLM-------------------------------------------------------------------------------------------
#=================================================================================
class OllamaSession:
    def __init__(self, model):
        self.model = model
        self.client = ollama.Client()
        self.client.generate(
            model=model,
            prompt="Hola :)",
            keep_alive=-1,
            options={"gpu": True}
        )

    def ask(self, prompt):
        return self.client.generate(model=self.model, prompt=prompt)


LLMsession = OllamaSession("deepseek-r1:8b")
#PROMPTS--------------------------------------------------------------------------------------
#=================================================================================
prompts = dict(enumerate([
    " Evita introduccione, comentarios, consejos u opiniones, proporciona solo hechos y da una respuesta seca.",
    " Resume brevemente: ",
    " Resume detalladamente: ",
    " Obtén las ideas clave: ",
    " Genera un itinerario con esta información: ",
    " Obtén los nombres de todas las personas: ",
    " Realiza una tabla con los datos de todas las personas que encuentres: ",
    " Realiza una tabla con los siguientes datos: ",
    ' Realiza un <table class="genTable" id="genTable"></table> sin estilos css, sin tailwind, sin saltos de linea, sin caracteres de escape, para insertar los siguientes datos: ',
    " Categoriza mediante un máximo de 15 etiquetas o 'labels' en español de una sola palabra en minúsculas (no palabras compuestas como 'palabracompuesta' o 'palabra-compuesta' o 'palabra compuesta') separadas por comas y sin ningún tipo de símbolo o información extra mas que las etiquetas la siguiente información: "
    ]))
print(prompts)

PROMPT_META = (
    "Extrae SOLO en JSON válido (sin texto extra, sin bloques markdown) los siguientes campos del documento:\n"
    '"nombres": lista de nombres de personas,\n'
    '"fechas": lista de fechas mencionadas,\n'
    '"direcciones": lista de direcciones o ubicaciones,\n'
    '"proposito": frase corta con el propósito del documento,\n'
    '"tags": lista de hasta 10 palabras clave en minúsculas.\n'
    "Si un campo no aplica usa lista vacía o cadena vacía.\n"
    "Documento:\n"
)

#OCR-------------------------------------------------------------------------------------------
#=================================================================================
import base64

OCR_MODEL = "deepseek-ocr:3b"
IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".bmp", ".tiff", ".tif", ".webp"}
OCR_PROMPT = "What text is written in this image? Write out all text exactly as it appears, word by word."
_OCR_TAG_RE = re.compile(r"<\|[^|]+\|>", re.DOTALL)


def _clean_ocr_output(text: str) -> str:
    """Elimina artefactos de detección de regiones del output del modelo de visión."""
    # Formato <|ref|>label<|/ref|><|det|>[[coords]]<|/det|>
    text = _OCR_TAG_RE.sub("", text)
    # Formato "Etiqueta: [[x, y, x2, y2]]" — líneas que solo tienen label + coordenadas
    text = re.sub(r"^[^\n\[]*:\s*\[\[[\d,\s]+\]\]\s*$", "", text, flags=re.MULTILINE)
    # Coordenadas sueltas [[x, y, x2, y2]] que queden
    text = re.sub(r"\[\[[\d,\s]+\]\]", "", text)
    # Colapsar líneas vacías múltiples
    text = re.sub(r"\n{2,}", "\n", text)
    return text.strip()

# Umbral mínimo de caracteres por página para considerar que el PDF tiene texto real.
# PDFs escaneados devuelven 0-10 chars/página con fitz; documentos de texto tienen miles.
MIN_CHARS_PER_PAGE = 50

#API-------------------------------------------------------------------------------------------
#=================================================================================
app = FastAPI()

#CORS-------------------------------------------------------------------------------------------
#=================================================================================
app.add_middleware(CORSMiddleware,
    allow_origins=["http://localhost:8000", "http://127.0.0.1:8000"],
    allow_credentials=True,
    allow_methods=['POST', 'GET'],
    allow_headers=['*'])

@app.on_event("startup")
async def startup_event():
    cleanup_expired_cache()
    print("[cache] Limpieza de caché completada al arrancar.")

def get_doc_text(path: str) -> tuple[str, dict]:
    """
    Devuelve (texto, metadatos) del documento.
    - Cache HIT: carga texto y metadatos del archivo JSON sin leer el PDF.
    - Cache MISS: parsea el PDF con fitz. Si el texto extraído es insuficiente
      (PDF escaneado / imagen), cae automáticamente al modo OCR.
    """
    cached = get_doc_cache(path)
    if cached:
        print(f"[cache] HIT: {path}")
        return cached["text"], cached.get("meta", {})

    print(f"[cache] MISS: {path} — leyendo PDF")
    doc_bytes = read_doc_bytes(path)
    doc = fitz.open(stream=doc_bytes, filetype="pdf")
    num_pages = max(doc.page_count, 1)
    text = "".join(page.get_text() for page in doc.pages())

    # Detección automática de PDFs escaneados o basados en imágenes.
    avg_chars = len(text.strip()) / num_pages
    if avg_chars < MIN_CHARS_PER_PAGE:
        print(
            f"[cache] Texto insuficiente ({avg_chars:.0f} chars/pág, mín {MIN_CHARS_PER_PAGE}) "
            f"— activando modo OCR automáticamente"
        )
        return get_doc_text_ocr(path)

    # Extracción de metadatos con LLM (ocurre solo en el primer acceso)
    print(f"[cache] Texto OK ({avg_chars:.0f} chars/pág) — extrayendo metadatos")
    meta = {}
    try:
        meta_resp = LLMsession.ask(PROMPT_META + text)
        raw = clean_text(meta_resp["response"])
        try:
            meta = json.loads(raw)
        except Exception:
            m = re.search(r'\{.*\}', raw, re.DOTALL)
            if m:
                meta = json.loads(m.group())
    except Exception as e:
        print(f"[cache] Error extrayendo metadatos: {e}")

    save_doc_cache(path, text, meta)
    return text, meta


def extract_text_ocr(path: str) -> str:
    """
    Extrae texto de un archivo de imagen o PDF escaneado usando deepseek-ocr.
    - Imágenes (.jpg, .png, etc.): lee bytes directamente.
    - PDFs: convierte cada página a imagen PNG con fitz.get_pixmap() y las envía al modelo.
    Devuelve el texto concatenado de todas las páginas/imágenes.
    """
    ext = os.path.splitext(path)[1].lower()
    pages_bytes = []
    file_bytes = read_doc_bytes(path)  # desencripta si es necesario

    if ext in IMAGE_EXTENSIONS:
        pages_bytes.append(file_bytes)
    else:
        doc = fitz.open(stream=file_bytes, filetype="pdf")
        for page in doc.pages():
            pix = page.get_pixmap(dpi=300)
            pages_bytes.append(pix.tobytes("png"))

    texts = []
    for img_bytes in pages_bytes:
        img_b64 = base64.b64encode(img_bytes).decode()
        resp = ollama.chat(
            model=OCR_MODEL,
            messages=[{
                "role": "user",
                "content": OCR_PROMPT,
                "images": [img_b64],
            }]
        )
        page_text = _clean_ocr_output(resp["message"]["content"])
        texts.append(page_text)

    return "\n".join(texts)


def get_doc_text_ocr(path: str) -> tuple[str, dict]:
    """
    Devuelve (texto, metadatos) del documento usando OCR.
    - Cache HIT: carga de {md5}_ocr.json sin llamar a deepseek-ocr.
    - Cache MISS: extrae con deepseek-ocr, luego extrae metadatos con LLM y guarda en caché.
    """
    cached = get_doc_cache(path, suffix="_ocr")
    if cached:
        print(f"[cache-ocr] HIT: {path}")
        return cached["text"], cached.get("meta", {})

    print(f"[cache-ocr] MISS: {path} — extrayendo con OCR")
    text = extract_text_ocr(path)

    meta = {}
    try:
        meta_resp = LLMsession.ask(PROMPT_META + text)
        raw = clean_text(meta_resp["response"])
        try:
            meta = json.loads(raw)
        except Exception:
            m = re.search(r'\{.*\}', raw, re.DOTALL)
            if m:
                meta = json.loads(m.group())
    except Exception as e:
        print(f"[cache-ocr] Error extrayendo metadatos: {e}")

    save_doc_cache(path, text, meta, suffix="_ocr")
    return text, meta

#ENDPOINTS--------------------------------------------------------------------------------------
#=================================================================================
@app.get("/")
async def get():
    return "hola :9"

@app.post("/chat")
async def chat(data: dict = Body(...)):
    if data["action"] != "":
        msj = prompts[data["action"]] + data["payload"]
    else:
        msj = data["payload"]
    msj += prompts[0]

    resp = LLMsession.ask(msj)
    if resp:
        return {"status":resp["done"], "message":clean_text(resp["response"]), "date": resp["created_at"]}
    else:
        return {"status":False, "message": "Error al generar respuesta", "date": dt.now().isoformat()}

@app.post("/actions")
async def actions(data: dict = Body(...)):
    print(f"ACTIONS data-> : {data}")
    if "path" not in data and "prompt" not in data:
        return {"status": "error", "message": "Se requiere 'path' o 'prompt'", "date": dt.now().isoformat()}
    if "action" not in data:
        return {"status": "error", "message": "Se requiere 'action'", "date": dt.now().isoformat()}
    try:
        if os.path.exists(data["path"]):
            action = data.get("action", 0)
            msj = prompts[action] if action else ""
            msj += prompts[0]
            doc_text, doc_meta = get_doc_text(data["path"])
            if doc_meta:
                msj += (
                    f"[Contexto: Propósito: {doc_meta.get('proposito', '')}. "
                    f"Personas: {', '.join(doc_meta.get('nombres', []))}. "
                    f"Fechas: {', '.join(doc_meta.get('fechas', []))}]\n"
                )
            msj += doc_text
            
            resp = LLMsession.ask(msj)
            
            if resp:
                return {"status":resp["done"], "message":clean_text(resp["response"]), "date": resp["created_at"]}
            else:
                return {"status":"error","message": "Error al generar respuesta", "response": None}
        else:
            if "prompt" in data:
                if data["action"] != "" and data["prompt"] != "":
                    msj = prompts[data["action"]] + data["prompt"]
                    resp = LLMsession.ask(msj)
                    if resp:
                        return {"status":resp["done"], "message":resp["response"], "date": resp["created_at"]}
                    else:
                        return {"status":"error","message": "Error al generar respuesta", "response": None}
                
            print("No se encontró el documento")
            return {"status":"error","message": "No existe el documento", "date": dt.now().isoformat()}
    except Exception as e:
        print(f"ACTIONS exception: {e}")
        return {"status": "error", "message": f"Error interno: {str(e)}", "date": dt.now().isoformat()}

@app.post("/interact")
async def interact(data: dict = Body(...)):
    # esta funcion se puede optimizar con memoria de conversacion con redis
    if "path" not in data:
        return {"status": "error", "message": "Se requiere 'path'", "date": dt.now().isoformat()}
    if "prompt" not in data:
        return {"status": "error", "message": "Se requiere 'prompt'", "date": dt.now().isoformat()}
    if os.path.exists(data["path"]):
        msj = data.get("prompt", "")
        msj += prompts[0]
        doc_text, doc_meta = get_doc_text(data["path"])
        if doc_meta:
            msj += (
                f"[Contexto: Propósito: {doc_meta.get('proposito', '')}. "
                f"Personas: {', '.join(doc_meta.get('nombres', []))}. "
                f"Fechas: {', '.join(doc_meta.get('fechas', []))}]\n"
            )
        msj += doc_text
        
        resp = LLMsession.ask(msj)
        
        if resp:
            return {"status":resp["done"], "message":clean_text(resp["response"]), "date": resp["created_at"]}
        else:
            return {"status":"error","message": "Error al generar respuesta", "response": None}
    else:
        print("No se encontró el documento")
        return {"status":"error","message": "No existe el documento", "date": dt.now().isoformat()}

@app.post("/agent")
async def agent(data: dict = Body(...)):
    for field in ("prompt", "idInst", "idUser"):
        if field not in data:
            return {"status": False, "message": f"Se requiere '{field}'", "date": dt.now().isoformat()}
    path = data.get("path") or None
    resp = run_agent(data["prompt"], data["idInst"], data["idUser"], "qwen3:4b", path=path)
    return {"status": True, "message": resp}

@app.post("/ocr")
async def ocr(data: dict = Body(...)):
    print(f"OCR data-> : {data}")
    if "path" not in data:
        return {"status": "error", "message": "Se requiere 'path'", "date": dt.now().isoformat()}
    try:
        path = data["path"]
        if not os.path.exists(path):
            return {"status": "error", "message": "No existe el documento", "date": dt.now().isoformat()}

        doc_text, doc_meta = get_doc_text_ocr(path)

        # Construir prompt: action tiene prioridad, luego prompt libre, luego solo instrucción base
        action = data.get("action", 0)
        prompt = data.get("prompt", "")

        if action:
            msj = prompts[action]
        elif prompt:
            msj = prompt
        else:
            msj = ""
        msj += prompts[0]

        if doc_meta:
            msj += (
                f"[Contexto: Propósito: {doc_meta.get('proposito', '')}. "
                f"Personas: {', '.join(doc_meta.get('nombres', []))}. "
                f"Fechas: {', '.join(doc_meta.get('fechas', []))}]\n"
            )
        msj += doc_text

        resp = LLMsession.ask(msj)
        if resp:
            return {"status": resp["done"], "message": clean_text(resp["response"]), "date": resp["created_at"]}
        else:
            return {"status": "error", "message": "Error al generar respuesta", "date": dt.now().isoformat()}

    except Exception as e:
        print(f"OCR exception: {e}")
        return {"status": "error", "message": f"Error interno: {str(e)}", "date": dt.now().isoformat()}

if __name__ == "__main__":
    uvicorn.run(app, host = "0.0.0.0", port = 4999)