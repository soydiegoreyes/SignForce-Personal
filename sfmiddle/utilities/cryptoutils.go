package utilities

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"sfmiddle/configs"
	"strings"
)

// GetHash calcula el hash de datos en memoria o de un archivo
func GetHash(input interface{}, config configs.HashConfig) (string, error) {
	var hasher hash.Hash

	// Seleccionar algoritmo
	switch config.Algorithm {
	case "sha512":
		hasher = sha512.New()
	case "sha256":
		hasher = sha256.New()
	default:
		hasher = sha256.New()
	}

	// Procesar input según el tipo
	switch v := input.(type) {
	case []byte:
		return hashBytes(v, hasher, config.Encoding)
	case string:
		return hashFile(v, hasher, config.Encoding, config.ChunkSize)
	default:
		return "", fmt.Errorf("tipo de entrada no soportado: %T", input)
	}
}

// hashBytes procesa datos en memoria recomendable no mayor a 200 MB
func hashBytes(data []byte, hasher hash.Hash, encoding string) (string, error) {
	_, err := hasher.Write(data)
	if err != nil {
		return "", err
	}

	return Encoder(hasher.Sum(nil), encoding)
}

// hashFile procesa archivos grandes por chunks
func hashFile(filePath string, hasher hash.Hash, encoding string, chunkSize int) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if chunkSize <= 0 {
		chunkSize = 1024 * 1024 * 20 // 20 MB por defecto
	}

	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}

		hasher.Write(buf[:n])
	}

	return Encoder(hasher.Sum(nil), encoding)
}

// encodeHash codifica el hash resultante
func Encoder(hash []byte, encoding string) (string, error) {
	switch encoding {
	case "hex":
		return hex.EncodeToString(hash), nil
	case "b64":
		return base64.StdEncoding.EncodeToString(hash), nil
	default:
		return "", fmt.Errorf("codificación no soportada: %s", encoding)
	}
}

func ConvertCertToPem(rutaCert string) (string, error) {
	var rutaPem string = rutaCert
	var err error = nil

	if !strings.HasSuffix(rutaCert, ".pem") {
		rutaPem = strings.Replace(rutaCert, ".cer", ".pem", 1)
		cmd := exec.Command("openssl", "x509", "-inform", "DER", "-in", rutaCert, "-out", rutaPem)

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error ejecutando OpenSSL:", err)
			rutaPem = rutaCert
		} else {
			fmt.Println("Certificado convertido a formato PEM")
		}
	}
	return rutaPem, err
}

// p=subprocess.run(f'openssl pkcs8 -inform DER -in "{ruta_key}" -out "{ruta_key.replace(".key",".pem")}" -passin pass:{pass_key}')
func ConvertKeyToPem(rutaKey, password string) (string, error) {
	var rutaPem string = rutaKey
	var err error = nil
	if !strings.HasSuffix(rutaKey, ".pem") {
		rutaPem = strings.Replace(rutaKey, ".key", ".pem", 1)
		passin := "pass:" + password
		cmd := exec.Command("openssl", "pkcs8", "-inform", "DER", "-in", rutaKey, "-out", rutaPem, "-passin", passin)

		err = cmd.Run()
		if err != nil {
			fmt.Println("Error ejecutando OpenSSL:", err)
			rutaPem = rutaKey

		} else {
			fmt.Println("Llave convertida a formato PEM")
		}
	}
	return rutaPem, err
}
