import uvicorn
from typing import List, Optional

from fastapi import FastAPI, HTTPException ,Body, Depends, Header
from fastapi.responses import HTMLResponse
from fastapi.middleware.cors import CORSMiddleware

import fitz

import re
import os
from datetime import datetime as dt
import dotenv
import ollama


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
    "Resume brevemente: ",
    "Resume detalladamente: ",
    "Obtén las ideas clave: ",
    "Genera un itinerario con esta información: ",
    "Obtén los nombres de todas las personas: ",
    "Realiza una tabla con los datos de todas las personas que encuentres: ",
    "Realiza una tabla con los siguientes datos: ",
    'Realiza un <table class="genTable" id="genTable"></table> sin estilos css, sin tailwind, sin saltos de linea, sin caracteres de escape, para insertar los siguientes datos: '
    ]))
print(prompts)

#API-------------------------------------------------------------------------------------------
#=================================================================================
app = FastAPI()

#CORS-------------------------------------------------------------------------------------------
#=================================================================================
app.add_middleware(CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=['*'],
    allow_headers=['*'])


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
    if os.path.exists(data["path"]):
        if data["action"] != "":
            msj = prompts[data["action"]]
        msj += prompts[0]
        doc = fitz.open(data["path"])
        for page in doc.pages():
            msj += page.get_text()
        
        resp = LLMsession.ask(msj)
        
        if resp:
            return {"status":resp["done"], "message":clean_text(resp["response"]), "date": resp["created_at"]}
        else:
            return {"status":"error","message": "Error al generar respuesta", "response": None}
    else:
        print("No se encontró el documento")
        return {"status":"error","message": "No existe el documento", "date": dt.now().isoformat()}

@app.post("/interact")
async def interact(data: dict = Body(...)):
    # esta funcion se puede optimizar con memoria de conversacion con redis
    if os.path.exists(data["path"]):
        if data["query"] != "":
            msj = data["query"]
        msj += prompts[0]
        doc = fitz.open(data["path"])
        for page in doc.pages():
            msj += page.get_text()
        
        resp = LLMsession.ask(msj)
        
        if resp:
            return {"status":resp["done"], "message":clean_text(resp["response"]), "date": resp["created_at"]}
        else:
            return {"status":"error","message": "Error al generar respuesta", "response": None}
    else:
        print("No se encontró el documento")
        return {"status":"error","message": "No existe el documento", "date": dt.now().isoformat()}

if __name__ == "__main__":
    uvicorn.run(app, host = "0.0.0.0", port = 4999)