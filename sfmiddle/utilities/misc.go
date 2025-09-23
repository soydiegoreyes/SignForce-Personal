package utilities

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
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
		uploadDir = fmt.Sprintf("./%s/%s/%s/", os.Getenv("TEMP_BASE_PATH"), idInst, idUser)
	} else {
		savepath, _ = strings.CutSuffix(savepath, "/")
		uploadDir = fmt.Sprintf("%s/", savepath)
	}

	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		return "", err
	}

	var uniqueName string
	// Generar nombre único para el archivo
	if hasUniqueName {
		uniqueName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	} else {
		uniqueName = filename

	}

	filePath := path.Join(uploadDir, uniqueName)
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
