package utilities

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
func GuardarArchivo(file multipart.File, savepath, filename, idInst, idUser string, hasUniqueName bool) (string, error) {
	// Crear directorio si no existe
	var uploadDir string
	if savepath == "" {
		uploadDir = fmt.Sprintf("%s/%s/%s/", os.Getenv("TEMP_BASE_PATH"), idInst, idUser)
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

func GetClientIP(r *http.Request) string {
	// X-Forwarded-For puede traer varias IPs: client, proxy1, proxy2...
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// La primera IP es la real
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Otro header común
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Si no viene en headers, usamos la IP directa
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
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
