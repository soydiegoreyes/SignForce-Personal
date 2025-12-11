import sys
import json
import fitz  # PyMuPDF

def insert_images_to_pdf(pdf_path, json_data):
    try:
        # 1. Parsear los datos recibidos (JSON string a lista de dicts)
        items = json.loads(json_data)
        
        # 2. Abrir el documento PDF
        doc = fitz.open(pdf_path)

        for item in items:
            # Extraer y convertir datos
            # Nota: Asumimos que pageSign viene base 1 (pag 1, 2, 3), 
            # pero fitz usa base 0 (index 0, 1, 2), por eso restamos 1.
            page_num = int(item.get('pageSign', 1)) - 1
            x = float(item.get('xSign', 0))
            y = float(item.get('ySign', 0))
            w = float(item.get('wSign', 50)) # Ancho default por si acaso
            h = float(item.get('hSign', 50)) # Alto default por si acaso
            img_path = item.get('pathImg')

            # Validar que la página existe
            if 0 <= page_num < len(doc):
                page = doc[page_num]
                
                # Definir el rectángulo donde irá la imagen (x0, y0, x1, y1)
                rect = fitz.Rect(x, y, x + w, y + h)
                
                # Insertar la imagen
                page.insert_image(rect, filename=img_path)
            else:
                print(f"Advertencia: Pagina {page_num + 1} fuera de rango.")

        # 3. Guardar cambios
        # save(..., incremental=True) es más rápido y seguro para ediciones pequeñas
        # o puedes guardar en un archivo temporal y renombrar.
        # Aquí sobrescribimos el original usando un archivo temporal implícito de fitz.
        doc.saveIncr() 
        doc.close()
        print("Success")

    except Exception as e:
        print(f"Error: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Uso: python insert_qr.py <pdf_path> <json_data>", file=sys.stderr)
        sys.exit(1)

    path = sys.argv[1]
    data = sys.argv[2]
    
    insert_images_to_pdf(path, data)