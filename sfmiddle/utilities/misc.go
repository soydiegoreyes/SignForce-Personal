package utilities

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Recibe una lista de strings y devuelve un diccionario clave valor [int]string
func Enumerate(list []string) map[int]string {
	result := make(map[int]string)
	for i, l := range list {
		result[i] = l
	}
	return result
}

func Encode_utf8(data []byte) string {
	// Decodificar de ISO-8859-1 a UTF-8
	decoder := charmap.ISO8859_1.NewDecoder()
	utf8Data, _, err := transform.Bytes(decoder, data)
	if err != nil {
		log.Fatalf("Error al decodificar: %v", err)
	}
	return string(utf8Data)
}

func Encode_b64(data []byte) string {
	// Codificar los bytes en Base64
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

func Decode_b64(s string) []byte {
	// Decodificar el string Base64 a bytes
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		fmt.Println("Error al decodificar:", err)

	}

	return data
}

func Latin1ToUTF8(input []byte) string {
	decoder := charmap.ISO8859_1.NewDecoder()
	utf8Str, err := decoder.Bytes(input)
	if err != nil {
		return string(input)
	}
	return string(utf8Str)
}

// Función para guardar archivos en el sistema

func GuardarArchivo(file multipart.File, savepath, filename string, hasUniqueName bool) (string, error) {
	// Crear directorio si no existe
	var uploadDir string
	if savepath == "" {
		return "", fmt.Errorf("%sNo se proporcionó una ruta", "")
	} else {
		savepath, _ = strings.CutSuffix(savepath, "/")
		uploadDir = fmt.Sprintf("%s/", savepath)
	}

	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", err
	}

	var name string
	// Generar nombre único para el archivo
	if hasUniqueName {
		name = fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	} else {
		name = filename

	}

	filePath := path.Join(uploadDir, name)
	_, err = os.Stat(filePath)
	// si el error es nulo es que ya existe el archivo
	if err == nil {
		fmt.Printf("archivo ya existe %s\n", filePath)
		return filePath, nil
	}

	// Crear el archivo
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copiar el contenido
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func GuardarArchivo_cript(file multipart.File, savepath, filename string, hasUniqueName bool, encrypt bool) (string, error) {

	// --- Lógica de directorios y nombres (se mantiene igual) ---
	savepath = strings.TrimSuffix(savepath, "/")
	uploadDir := savepath + "/"

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	name := filename
	if hasUniqueName {
		name = fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	}

	filePath := path.Join(uploadDir, name)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil // El archivo ya existe
	}

	// 2. Crear el archivo destino
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 3. Decidir cómo guardar el contenido
	if encrypt {
		// Regresar al inicio del reader por si acaso fue leído antes
		file.Seek(0, io.SeekStart)

		err = EncryptFile(file, dst)
		if err != nil {
			return "", fmt.Errorf("error al encriptar: %v", err)
		}
	} else {
		file.Seek(0, io.SeekStart)
		_, err = io.Copy(dst, file)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

// AppendQRCodes invoca un script de Python para insertar imágenes en un PDF
func AppendQRCodes(filePath string, SignsPositions map[string]map[string]string) error {

	// 1. Aplanar el mapa a un Slice (Lista) de mapas para el JSON
	// El mapa original tiene un ID como llave principal que no necesitamos enviar al script,
	// solo necesitamos la lista de configuraciones.
	var stampsList []map[string]string

	for _, data := range SignsPositions {
		stampsList = append(stampsList, data)
	}

	// 2. Convertir la lista a un string JSON
	jsonData, err := json.Marshal(stampsList)
	if err != nil {
		fmt.Printf("Error al crear JSON para Python: %v\n", err)
		return err
	}

	// 3. Preparar el comando para ejecutar Python
	// Asegúrate de poner la ruta correcta donde guardaste 'insert_qr.py'
	scriptPath := "./utilities/insert_qr.py"
	fmt.Println(filePath, string(jsonData))
	// Comando: python3 insert_qr.py [RutaPDF] [StringJSON]
	cmd := exec.Command("python", scriptPath, filePath, string(jsonData))

	// Capturar salida estándar y de error para debuggear si Python falla
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println(cmd.Stderr, cmd.Stdout)
	// 4. Ejecutar el script
	fmt.Println("Ejecutando script de Python para insertar QRs...")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error ejecutando script de Python: %v\n", err)
		return err
	}

	fmt.Println("QRs insertados correctamente en:", filePath)
	return nil
}

func CompressZip(pathSource, pathDest string) bool {
	//"7z a pathDest.zip pathSource/*"
	fmt.Println("Comprimiendo datos...")
	cmd := exec.Command("7z", "a", pathDest, pathSource)
	// Capturar salida estándar y de error para debuggear si Python falla
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error ejecutando script de Python: %v\n", err)
		return false
	}
	fmt.Println(cmd.Stderr)
	return true
}

func Bool2Int(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func EuclideanDistance(v1, v2 []float32) float64 {
	if len(v1) != len(v2) {
		return 1.0
	}

	var sum, r float64
	for i := range v1 {
		dist := float64(v1[i]) - float64(v2[i])
		sum += dist * dist
	}
	r = math.Sqrt(sum)
	fmt.Println("distancia: ", r)
	return r
}

func B642ArrFloat(strB64 string) []float32 {
	rawBytes := Decode_b64(strB64)
	var arr []float32
	for i := 0; i < len(rawBytes); i += 4 {
		// Leemos 4 bytes en LittleEndian (estándar de JS)
		bits := binary.LittleEndian.Uint32(rawBytes[i : i+4])
		floatVal := math.Float32frombits(bits)
		arr = append(arr, floatVal)
	}
	return arr
}

func ArrFloat2B64(arr []float32) string {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, arr)
	if err != nil {
		fmt.Println("Error decodificando array")
		return ""
	}

	return Encode_b64(buf.Bytes())
}

func GetClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		fmt.Println("IP de XFF")
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		fmt.Println("IP de XRIP")
		return xrip
	}
	if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" {
		fmt.Println("IP de CFCIP")
		return cfip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		fmt.Println("IP de RA")
		return ip
	}
	return r.RemoteAddr
}
