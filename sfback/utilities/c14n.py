from lxml import etree
import base64
import hashlib
import sys

def get_xml_node_digest(xml_content: str, node_id: str, hash_alg: str = 'sha256'):
    try:
        parser = etree.XMLParser(remove_blank_text=True)
        # Importante: el XML debe estar bien formado
        root = etree.fromstring(xml_content.encode('utf-8'), parser)
        
        # Buscar el nodo por su atributo Id o id
        referenced_node = root.xpath(f"//*[@Id='{node_id}'] | //*[@id='{node_id}']")
        
        if not referenced_node:
            return None

        node = referenced_node[0]

        # CANONIZACIÓN EXCLUSIVA (Crucial para XAdES y validación ETSI)
        canonicalized_data = etree.tostring(
            node, 
            method="c14n", 
            exclusive=True, 
            with_comments=False
        )

        hash_obj = hashlib.new(hash_alg)
        hash_obj.update(canonicalized_data)
        
        return base64.b64encode(hash_obj.digest()).decode('utf-8')
    except Exception as e:
        # Imprimir error a stderr para no ensuciar la salida del hash
        print(f"Error procesando XML: {e}", file=sys.stderr)
        return None

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("Uso: python get_digest.py <xml_string> <id_node> <algo>", file=sys.stderr)
        sys.exit(1)

    xml_str = sys.argv[1]
    id_node = sys.argv[2]
    algo = sys.argv[3].lower()

    result = get_xml_node_digest(xml_str, id_node, algo)
    
    if result:
        # Imprimimos SOLO el hash para que Go lo lea fácilmente
        print(result, end='')
    else:
        sys.exit(1)