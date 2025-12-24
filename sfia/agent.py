import ollama

from docx import Document
from docx.shared import Inches
from docx.enum.text import WD_PARAGRAPH_ALIGNMENT

import os
import re
import json
import subprocess

# recibe una ruta donde se guarda el template y el contenido que ira en cada sección del documento.
# ej. output_path = "C:/USERS/USER/Desktop/mi_plantilla.docx"
# ej. data = {"header": "Contrato de compraventa {TITULO_CONTRATO}", "body": "Esto es un contrato entre {COMPRADOR} y {VENDEDOR}.", "footer": "Pie de pagina"}
# devuelve una lista de {PLACEHOLDERS} que se extrae del texto ya que es texto generado con IA y no se pasa directamente
def create_template(output_path: str, data: dict) -> list:
    """
    Crea un archivo .docx base con placeholders para ser usado como plantilla.
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

def format_doc(input_path: str, output_path: str, data: dict) ->str:
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
        
        for m in re.findall('{\w+}', paragraph.text):
            try:
                paragraph.text=paragraph.text.replace(m, data[m[1:-1]])
                
            except Exception as e:
                print(f"Atributo {paragraph.text} no encontrado en parrafo para sustituir. {e}")

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
    try:
        #p = subprocess.run(['soffice', '--headless', '--convert-to', 'pdf', f'{ruta_dest}.docx'], capture_output=True, text=True)
        #cmd = ['sudo', 'soffice', '--headless', '--convert-to', 'pdf', ruta_docx]
        cmd = ['soffice', '--headless', '--convert-to', 'pdf', input_path]
        p = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        p.communicate()
        # Esto lanzará un error si el proceso falla
        print(f"Return code: {p.returncode}")
        if p.returncode not in [0, None]:
            raise Exception  
        
        name, *_ = re.match('(?:.+[/\\\])?(.+)(?:\.\w{2,4})$', input_path).groups()
    except Exception as e:
        return f"Falla al convertir a pdf {input_path} . {e}"
    
    try:
        # se valida que el pdf se haya creado con éxito en la ruta local    
        if os.path.exists(f'{os.path.abspath(".")}\\{name}.pdf'):
            with open(name+'.pdf', 'rb') as fr, open(output_path, 'wb') as fw:
                fw.write(fr.read())
            # se elimina el documento creado en la carpeta local
            os.remove(f'{os.path.abspath(".")}\\{name}.pdf')
            #os.remove(ruta_dest+'.docx')
        if p.returncode:
            return f"Hubo un error al guardar el documento word, revise rutas del archivo {output_path}"
        else:
            return {
                "status": "ok",
                "ruta": output_path
            }
        
    except Exception as e:
        return f"Falla al guardar archivo {output_path} -> {e}"
    

def notificar_error(mensaje: str):
    print("FUNC: notificar_error")
    """Útil para que el agente informe al usuario si algo falta o falla."""
    return f"Atención: {mensaje}"

# 2. Mapeo para ejecución
available_functions = {
    'notificar_error': notificar_error,
    'create_template': create_template,
    'word_to_pdf': word_to_pdf,
    'format_doc': format_doc
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
                    'input_path': {'type': 'string', 'description': 'Ruta absoluta del documento word a convertir en pdf. ej: C:/Users/SomeFolder/mi_documento.docx'},
                    'output_path': {'type': 'string', 'description': 'Ruta absoluta del documento pdf convertido -> C:/Users/OtherFolder/mi_documento.pdf'},
                },
                'required': ['input_path', 'output_path'],
            },
        },
    }
]


def run_agent(prompt, model):
    messages = [
        {
            'role': 'system', 
            'content': (
                '''Eres un redactor profesional de documentos legales.
                OBJETIVO PRINCIPAL:
                - Genera un documento COMPLETO, coherente y útil para un humano.
                - El texto debe poder leerse y entenderse incluso sin rellenar datos.
                - El documento debe poder convertirse de word (.docx) a (.pdf)
                USO DE PLACEHOLDERS:
                - SOLO usa placeholders {EN_MAYUSCULAS} cuando un dato específico sea variable y no sea una ruta, path o nombre de un archivo.
                - NO reemplaces todo el texto por placeholders.
                - El documento debe contener frases, cláusulas, datos útiles y contexto real.
                FLUJO:
                0. Si se pide explicitamente una tarea que se resuelve con el llamado a una función concreta, omitir el flujo.
                1. Redacta el texto completo del documento (header, body y footer).
                2. Inserta placeholders SOLO para datos variables (nombres, montos, fechas, ubicaciones, o datos específicados explícitamente).
                3. Llama a 'create_template' con el texto generado y por defecto si no se especifica el nombre de la plantilla inventa un nombre sin ningun path.
                4. Si hay una lista de placeholders devuelta, usalos para llamar a 'format_doc'.
                5. Si EXPLICITAMENTE se indica que un documento debe convertirse a pdf, usa 'word_to_pdf' y por default con output_path y nombre de documento igual a la de entrada pero con extensión pdf.
                6. Ejecuta los pasos uno por uno, esperando el resultado anterior.'''
            )
        },
        {'role': 'user', 'content': prompt}
    ]

    # Bucle de ejecución (máximo 5 iteraciones para evitar bucles infinitos)
    for i in range(5):
        response = ollama.chat(model=model, messages=messages, tools=tools_definition)
        
        # Si el modelo ya no quiere llamar a más funciones, terminamos
        if not response["message"].get("tool_calls"):
            return response["message"]["content"]

        # Si hay llamadas a funciones
        messages.append(response["message"])
        
        for call in response["message"]["tool_calls"]:
            function_name = call["function"]["name"]
            args = call["function"]["arguments"]
            
            print(f"--- Paso {i+1}: Ejecutando {function_name} ---")
            
            # Buscamos la función en nuestro diccionario
            # OJO: Asegúrate de que 'create_template' esté en available_functions
            result = available_functions[function_name](**args)
            print(result)
            messages.append({
                'role': 'tool',
                'content': json.dumps(result),
                'name': function_name
            })
            
        # El bucle continúa: enviamos los resultados de vuelta a Ollama 
        # para que decida qué sigue (ej. ya creó la plantilla, ahora le toca llenarla).
    print(messages)
    return messages

if __name__=="__main__":
    # Prueba tu agentez
    prompt= '''Realiza un contrato de venta de terreno entre las partes Juan Jesús como el comprador y Guadalupe Maria como la vendedora de terreno ubicado en Toluca que estipule el monto por $10.00 MXN (diez pesos mexicanos ) con clausulas comunes en un contrato de compra y venta justo entre ambas partes y guardalo enn la ruta C:/Users/USER/Desktop/contrato_prueba.docx'''
    print(run_agent(prompt, "qwen3:4b"))

