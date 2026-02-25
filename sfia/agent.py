import ollama

from docx import Document
from docx.shared import Inches
from docx.enum.text import WD_PARAGRAPH_ALIGNMENT

import os
import re
import json
import subprocess

path_match = r'(?:[^\s/\\]+[/|\\])*(\w+\.\w{2,8})'
pmatch =re.compile(path_match)

def create_template_entrypoint(**kwargs):
    """
    Normaliza los argumentos para create_template.
    - Caso 1 (ideal): kwargs contiene 'output_path' y 'data' (dict) -> pasa directamente.
    - Caso 2 (común cuando el modelo "aplana" el objeto): kwargs contiene 'output_path' y keys como 'header','body',... -> construye data dict.
    """
    print("FUNC: create_template_entrypoint args:", kwargs)

    # Si el modelo no mandó output_path, generamos uno automáticamente
    if 'output_path' not in kwargs:
        # Intentamos derivar el nombre del header
        raw_header = (
            kwargs.get('header')
            or (kwargs.get('data') or {}).get('header', '')
            or 'documento_generado'
        )
        safe_name = re.sub(r'[^\w\s-]', '', raw_header[:40]).strip().replace(' ', '_').lower()
        safe_name = safe_name or 'documento_generado'
        kwargs['output_path'] = os.path.join(os.path.abspath('.'), f'{safe_name}.docx')
        print(f"[create_template] output_path no proporcionado, usando: {kwargs['output_path']}")

    # Si output_path es solo un nombre (relativo), lo anclamos al cwd
    output_path = kwargs['output_path']
    if not os.path.isabs(output_path):
        output_path = os.path.join(os.path.abspath('.'), output_path)
        kwargs['output_path'] = output_path

    # Si ya viene 'data' como dict -> úsalo
    if 'data' in kwargs and isinstance(kwargs['data'], dict):
        data = kwargs['data']
        return create_template(output_path, data)

    # Si no hay 'data', construimos uno a partir del resto de kwargs (excluyendo output_path)
    data = {}
    for k, v in kwargs.items():
        if k == 'output_path':
            continue
        # 'images' podría ser una cadena JSON que representa una lista; si es str y parece un JSON array, intentar parsear.
        if k == 'images' and isinstance(v, str):
            try:
                parsed = json.loads(v)
                if isinstance(parsed, list):
                    data['images'] = parsed
                    continue
            except Exception:
                # si no se parsea, tratarlo como string simple
                pass
        data[k] = v

    # Aseguramos que al menos 'body' exista (coherente con tools_definition)
    if 'body' not in data:
        # devolvemos error estructurado para que el agente lo muestre
        return {
            "status": "error",
            "message": "create_template requiere al menos 'body' en 'data'.",
            "received": data
        }

    return create_template(output_path, data)

# recibe una ruta donde se guarda el template y el contenido que ira en cada sección del documento.
# ej. output_path = "C:/USERS/USER/Desktop/mi_plantilla.docx"
# ej. data = {"header": "Contrato de compraventa {TITULO_CONTRATO}", "body": "Esto es un contrato entre {COMPRADOR} y {VENDEDOR}.", "footer": "Pie de pagina"}
# devuelve una lista de {PLACEHOLDERS} que se extrae del texto ya que es texto generado con IA y no se pasa directamente
def create_template(output_path: str, data: dict):
    """
    Crea un archivo .docx para ser usado como plantilla pudiendo incorporar si es necesario placeholders {PLACEHOLDER}.
    contenido: {'header': '...', 'body': '...', 'footer': '...', 'images': ['./path.png']}
    """
    print("FUNC: create_template", output_path, data)
    try:
        doc = Document()
        placeholders=list()
        # 1. Configurar Header
        if 'header' in data:
            section = doc.sections[0]
            header = section.header
            header.paragraphs[0].text = data['header']
            placeholders+= [p[1:-1] for p in re.findall(r'{\w+}', data['header'])]
            
        # 2. Configurar Body
        if 'body' in data:
            doc.add_paragraph(data['body'])
            placeholders+= [p[1:-1] for p in re.findall(r'{\w+}', data['body'])]
            
            
        # 3. Añadir Imágenes (si se piden para la plantilla)
        if 'images' in data:
            for img_path in data['images']:
                if os.path.exists(img_path):
                    doc.add_picture(img_path, width=Inches(2.0))
                    placeholders += [img_path]
        
        # 4. Configurar Footer
        if 'footer' in data:
            section = doc.sections[0]
            footer = section.footer
            footer.paragraphs[0].text = data['footer']
            placeholders+= [p[1:-1] for p in re.findall(r'{\w+}', data['footer'])]
            
        doc.save(output_path)
            
        return {
            "status": "ok",
            "ruta": output_path,
            "placeholders": placeholders
        }
    
    except Exception as e:
        return f"Error al crear plantilla: {str(e)}"
    

# recibe una ruta de entrada que es una plantilla y una ruta de salida donde se deposita la plantilla con los {PLACEHOLDERS} sustituidos con el dato real
def format_doc(input_path: str, output_path: str, data: dict):
    print("FUNC: format_doc", input_path, output_path, data)
    if os.path.exists(input_path):
        doc = Document(input_path)
    else:
        return "Error en format_doc"

    def replace_text_with_format(paragraph, data_dict):
        """Reemplaza el texto sin perder el formato."""
        for run in paragraph.runs:
            for match in re.findall(r'{\w+}', run.text):
                key = match[1:-1]
                if key in data_dict:
                    run.text = run.text.replace(match, data_dict[key])

    # Reemplazo en párrafos normales
    for paragraph in doc.paragraphs:
        replace_text_with_format(paragraph, data)
        
    # Reemplazo en tablas
    for table in doc.tables:
        for row in table.rows:
            for cell in row.cells:
                for paragraph in cell.paragraphs:
                    replace_text_with_format(paragraph, data)

    # Añadir una imagen si existe 'ImageBody' en los datos
    if 'ImageBody' in data:
        try:
            image_paragraph = doc.add_paragraph()
            image_paragraph.alignment = WD_PARAGRAPH_ALIGNMENT.LEFT
            image_run = image_paragraph.add_run()
            image_run.add_picture(data['ImageBody'], width=Inches(0.7))
        except Exception as e:
            print(f"Error al incrustar imagen. {e}")

    # Reemplazo en el footer
    for section in doc.sections:
        footer = section.footer
        for paragraph in footer.paragraphs:
            replace_text_with_format(paragraph, data)
    
    try:
        # Guarda el documento modificado
        doc.save(output_path)
        
        return {
            "status": "ok",
            "ruta": output_path
        }
    except Exception as e:
        print(f"Falla al guardar word {output_path} . {e}")
    

# recibe una ruta de un documento en word .docx y lo convierte en pdf (no necesariamente tiene que ser una plantilla)
def word_to_pdf(input_path: str, output_path: str):
    # ruta_docx es la ruta absoluta del documento word a convertir en pdf -> C:/Users/SomeFolder/mi_documento.docx
    # ruta_dest es la ruta absoluta con nombre donde se guardará el pdf -> C:/Users/OtherFolder/mi_documento.pdf
    # el flujo es C:/Users/SomeFolder/mi_documento.docx -> ./mi_documento.pdf -> C:/Users/OtherFolder/mi_documento.pdf
    print("FUNC: word_to_pdf")

    # Resolver rutas relativas
    if not os.path.isabs(input_path):
        input_path = os.path.join(os.path.abspath('.'), input_path)
    if not os.path.isabs(output_path):
        output_path = os.path.join(os.path.abspath('.'), output_path)

    if not os.path.exists(input_path):
        return f"Falla al convertir a pdf: no existe el archivo de entrada '{input_path}'"

    try:
        cmd = ['soffice', '--headless', '--convert-to', 'pdf', input_path]
        p = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        p.communicate()
        print(f"Return code: {p.returncode}")
        if p.returncode not in [0, None]:
            raise Exception(f"soffice retornó {p.returncode}")

        name, *_ = re.match(r'(?:.+[/\\])?(.+)(?:\.\w{2,4})$', input_path).groups()
    except Exception as e:
        return f"Falla al convertir a pdf {input_path}: {e}"

    try:
        local_pdf = os.path.join(os.path.abspath("."), name + ".pdf")
        if not os.path.exists(local_pdf):
            return f"Falla al convertir: LibreOffice no generó el PDF en '{local_pdf}'"

        with open(local_pdf, 'rb') as fr, open(output_path, 'wb') as fw:
            fw.write(fr.read())
        if local_pdf != output_path:
            os.remove(local_pdf)

        return {"status": "ok", "ruta": output_path}

    except Exception as e:
        return f"Falla al guardar archivo {output_path}: {e}"
    

# sfia cuenta con un directorio local con los artefactos generados para el usuario donde otros programas obtendrán esa info
def copiar_a_destino(input_path, idInst, idUser: str):
    r = {"status": "", "ruta": ""}

    if not input_path or not idInst or not idUser:
        r["status"] = "error: falta input_path, idInst o idUser"
        r["ruta"] = input_path or ""
        return r

    # Resolver a ruta absoluta
    if not os.path.isabs(input_path):
        input_path = os.path.join(os.path.abspath('.'), input_path)

    if not os.path.exists(input_path):
        r["status"] = f"error: Archivo no encontrado en '{input_path}'"
        r["ruta"] = input_path
        return r

    destino = "temp"
    base_path = os.path.abspath(".")
    if input_path.endswith(".docx"):
        destino = "templates"
    elif input_path.endswith(".pdf"):
        destino = "documents"

    ruta_dir = os.path.join(base_path, idInst, idUser, destino)
    os.makedirs(ruta_dir, 0o755, exist_ok=True)

    # Si ya está en el directorio correcto no hay nada que mover
    if os.path.abspath(os.path.dirname(input_path)) == os.path.abspath(ruta_dir):
        r["status"] = "ok"
        r["ruta"] = input_path
        return r

    filename = os.path.basename(input_path)
    ruta_destino = os.path.join(ruta_dir, filename)
    with open(input_path, "rb") as fr, open(ruta_destino, "wb") as fw:
        fw.write(fr.read())
    r["status"] = "ok"
    r["ruta"] = ruta_destino
    os.remove(input_path)
    return r
        

def notificar_error(mensaje: str):
    print("FUNC: notificar_error")
    """Útil para que el agente informe al usuario si algo falta o falla."""
    return f"Atención: {mensaje}"

# 2. Mapeo para ejecución
available_functions = {
    'notificar_error': notificar_error,
    'create_template': create_template_entrypoint,
    'word_to_pdf': word_to_pdf,
    'format_doc': format_doc,
}

# 3. DEFINICIÓN EN FORMATO JSON
# Esto es lo que frameworks como LangChain hacen por detrás
tools_definition = [
    {
        'type': 'function',
        'function': {
            'name': 'create_template',
            'description': 'Crea un archivo .docx nuevo que servirá como plantilla, incluyendo texto con placeholders.',
            'parameters': {
                'type': 'object',
                'properties': {
                    'output_path': {'type': 'string', 'description': 'Ruta absoluta donde se guardará la plantilla .docx'},
                    'data': {
                        'type': 'object',
                        'properties': {
                            'header': {'type': 'string', 'description': 'Útil para títulos o encabezados del documento'},
                            'body': {'type': 'string', 'description': 'Contenido completo del documento'},
                            'footer': {'type': 'string', 'description': 'Útil para avisos o pies de página o derechos de autor.'},
                            'images': {'type': 'array', 'items': {'type': 'string', 'description':'Cualquier imagen que se tenga que poner sobre el documento'}}
                        },
                        'required': ['body']
                    }
                },
                'required': ['output_path', 'data']
            }
        }
    },
    {
        'type': 'function',
        'function': {
            'name': 'format_doc',
            'description': 'Rellena placeholders en formato {KEY} de una plantilla Word (.docx) usando un diccionario.',
            'parameters': {
                'type': 'object',
                'properties': {
                    'input_path': {
                        'type': 'string', 
                        'description': 'Ruta del archivo .docx original (ej: plantilla.docx)'
                    },
                    'output_path': {
                        'type': 'string', 
                        'description': 'Ruta donde se guardará el archivo generado (ej: contrato_final.docx)'
                    },
                    'data': {
                        'type': 'object',
                        'description': 'Diccionario donde las llaves son los placeholders y los valores el texto a insertar',
                        'additionalProperties': {'type': 'string'}
                    },
                },
                'required': ['input_path', 'output_path', 'data'],
            }
        }
    },
    {
        'type': 'function',
        'function': {
            'name': 'word_to_pdf',
            'description': 'Convierte un documento word (.docx) a un documento pdf (.pdf)',
            'parameters': {
                'type': 'object',
                'properties': {
                    'input_path': {'type': 'string', 'description': 'Ruta absoluta del documento word a convertir en pdf'},
                    'output_path': {'type': 'string', 'description': 'Ruta absoluta del documento pdf convertido'},
                },
                'required': ['input_path', 'output_path'],
            },
        },
    },
]


def run_agent(prompt, idInst, idUser, model, path=None):
    base_dir = os.path.abspath('.')
    templates_dir = os.path.join(base_dir, idInst, idUser, 'templates')
    documents_dir = os.path.join(base_dir, idInst, idUser, 'documents')
    os.makedirs(templates_dir, exist_ok=True)
    os.makedirs(documents_dir, exist_ok=True)

    system_lines = [
        "Eres un redactor profesional de documentos legales.",
        "OBJETIVO: Genera documentos COMPLETOS, coherentes y útiles para un humano.",
        "PLACEHOLDERS: Usa {EN_MAYUSCULAS} SOLO para datos variables concretos (nombres, montos, fechas).",
        "  NO rellenes todo con placeholders; el documento debe tener contenido real.",
        "",
        "DIRECTORIOS DE TRABAJO — usa SIEMPRE rutas absolutas de estos directorios:",
        f"  Plantillas .docx → {templates_dir}",
        f"  Documentos PDF  → {documents_dir}",
        "",
        "FLUJO PARA NUEVA PLANTILLA:",
        "  1. Redacta header, body y footer del documento.",
        f" 2. Llama a 'create_template' con output_path='{templates_dir}/<nombre>.docx'",
        "  3. Si se solicita PDF, llama a 'word_to_pdf' usando el .docx recién creado.",
        "  Ejecuta los pasos uno por uno esperando el resultado anterior.",
    ]

    if path and os.path.exists(path):
        stem = os.path.splitext(os.path.basename(path))[0]
        system_lines += [
            "",
            f"ARCHIVO EXISTENTE: {path}",
            "FLUJO PARA USAR ARCHIVO EXISTENTE:",
            f"  · Para LLENAR placeholders: 'format_doc' con input_path='{path}'"
            f" y output_path='{documents_dir}/{stem}_relleno.docx'",
            f"  · Para CONVERTIR a PDF: 'word_to_pdf' con input_path='{path}'"
            f" y output_path='{documents_dir}/{stem}.pdf'",
            "  · Si se pide llenar Y convertir: primero 'format_doc', luego 'word_to_pdf' con el docx resultante.",
        ]

    results = dict()
    messages = [
        {'role': 'system', 'content': '\n'.join(system_lines)},
        {'role': 'user', 'content': prompt},
    ]

    # Bucle de ejecución (máximo 6 iteraciones para cubrir create+format+pdf)
    for i in range(6):
        response = ollama.chat(model=model, messages=messages, tools=tools_definition)
        if not response["message"].get("tool_calls"):
            return response["message"]["content"]

        messages.append(response["message"])
        for call in response["message"]["tool_calls"]:
            function_name = call["function"]["name"]
            args = call["function"]["arguments"]
            results.setdefault(function_name, {})["args"] = args

            print(f"--- Paso {i+1}: Ejecutando {function_name} ---")

            if function_name not in available_functions:
                print(f"--- Función desconocida: {function_name} ---")
                result = f"Error: función '{function_name}' no existe. Usa solo: {list(available_functions.keys())}"
                messages.append({'role': 'tool', 'content': result, 'name': function_name})
                continue

            result = available_functions[function_name](**args)

            # Mover a directorio del usuario SOLO si el archivo quedó fuera de él
            if isinstance(result, dict) and result.get("status") == "ok" and "ruta" in result:
                ruta_abs = os.path.abspath(result["ruta"])
                ruta_dir = os.path.abspath(os.path.dirname(ruta_abs))
                user_dirs = {os.path.abspath(templates_dir), os.path.abspath(documents_dir)}
                if ruta_dir not in user_dirs:
                    r = copiar_a_destino(result["ruta"], idInst=idInst, idUser=idUser)
                    result["ruta"] = r["ruta"]

            results["result"] = result
            messages.append({
                'role': 'tool',
                'content': json.dumps(result),
                'name': function_name,
            })

    return "El agente completó las iteraciones sin producir una respuesta de texto final."

if __name__=="__main__":
    # Prueba tu agentez
    prompt= '''Realiza un contrato de venta de terreno entre las partes Juan Jesús como el comprador y Guadalupe Maria como la vendedora de terreno ubicado en Toluca que estipule el monto por $10.00 MXN (diez pesos mexicanos ) con clausulas comunes en un contrato de compra y venta justo entre ambas partes y guardalo enn la ruta C:/Users/USER/Desktop/contrato_prueba.docx'''
    print(run_agent(prompt, "qwen3:4b"))

