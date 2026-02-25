package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"os/exec"
	"sfmiddle/configs"
	"strings"
)

var symbols = [70]string{
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "Ñ", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "ñ", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"!", "#", "$", "%", "&", ".", "/", "(", ")", "-", "_", ";", ",", "*", "+", ":",
}

// GetHash calcula el hash de datos en memoria o de un archivo
func GetHash(input interface{}, config configs.HashConfig, isEncrypted bool) (string, error) {
	var hasher hash.Hash

	switch config.Algorithm {
	case "sha512":
		hasher = sha512.New()
	case "sha256":
		hasher = sha256.New()
	default:
		hasher = sha256.New()
	}

	switch v := input.(type) {
	case []byte:
		return hashBytes(v, hasher, config.Encoding)
	case string:
		// Decidimos qué función usar según el flag
		if isEncrypted {
			return hashFileEncrypted(v, hasher, config.Encoding)
		}
		return hashFilePlain(v, hasher, config.Encoding)
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
// hashFilePlain para archivos normales (Grandes o pequeños)
func hashFilePlain(filePath string, hasher hash.Hash, encoding string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// io.Copy lee el archivo por partes y lo manda al hasher eficientemente
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return Encoder(hasher.Sum(nil), encoding)
}

// hashFileEncrypted usa el Pipe que ya tenías (es la mejor forma de hacerlo)
func hashFileEncrypted(filePath string, hasher hash.Hash, encoding string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	pr, pw := io.Pipe()
	errChan := make(chan error, 1)

	go func() {
		// Importante: CloseWithError para que io.Copy se entere si algo falló
		err := DecryptFile(file, pw)
		if err != nil {
			pw.CloseWithError(err)
			errChan <- err
			return
		}
		pw.Close()
		errChan <- nil
	}()

	// io.Copy absorberá los datos descifrados que vienen del pipe
	if _, err := io.Copy(hasher, pr); err != nil {
		return "", fmt.Errorf("error procesando hash cifrado: %v", err)
	}

	if decryptErr := <-errChan; decryptErr != nil {
		return "", decryptErr
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

// Usa GCM y lee de un reader y escribe el cifrado en un writer
func EncryptFile(src io.Reader, dst io.Writer) error {
	key := []byte(os.Getenv("DOCS_KEY"))
	key_hash := sha256.Sum256(key)
	block, _ := aes.NewCipher(key_hash[:])
	gcm, _ := cipher.NewGCM(block)

	// 1. Generar un Nonce base aleatorio
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(crand.Reader, nonce)

	// Escribir el nonce al inicio del archivo destino
	dst.Write(nonce)

	buf := make([]byte, 64*1024)
	var i uint64 = 0

	for {
		n, err := src.Read(buf)
		if n > 0 {
			// 2. Modificar el nonce para cada bloque (evita ataques de repetición)
			// Usamos el índice del bloque para que cada uno sea único
			currentNonce := make([]byte, len(nonce))
			copy(currentNonce, nonce)
			binary.BigEndian.PutUint64(currentNonce[len(nonce)-8:], i)

			// 3. Cifrar el bloque y escribirlo
			ciphertext := gcm.Seal(nil, currentNonce, buf[:n], nil)
			dst.Write(ciphertext)
			i++
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptFile lee el archivo cifrado en GCM y escribe el original en dst
func DecryptFile(src io.Reader, dst io.Writer) error {
	key := []byte(os.Getenv("DOCS_KEY"))
	key_hash := sha256.Sum256(key)
	block, err := aes.NewCipher(key_hash[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 1. Leer el nonce base (los primeros 12 bytes del archivo)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(src, nonce); err != nil {
		return err
	}

	// El tamaño de lo que vamos a leer es el bloque original + el tag de seguridad
	encryptedChunkSize := (64 * 1024) + 16
	buf := make([]byte, encryptedChunkSize)
	var i uint64 = 0

	for {
		// Intentamos leer el bloque completo
		n, err := io.ReadFull(src, buf)

		// Si hay un error y no es EOF ni UnexpectedEOF (un bloque final parcial), salimos
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return err
		}

		if n > 0 {
			// Reconstruir el nonce (tu lógica actual está bien)
			currentNonce := make([]byte, len(nonce))
			copy(currentNonce, nonce)
			binary.BigEndian.PutUint64(currentNonce[len(nonce)-8:], i)

			// Descifrar solo los n bytes leídos
			plaindata, decryptErr := gcm.Open(nil, currentNonce, buf[:n], nil)
			if decryptErr != nil {
				return fmt.Errorf("fallo de integridad en bloque %d: %v", i, decryptErr)
			}

			if _, writeErr := dst.Write(plaindata); writeErr != nil {
				return writeErr
			}
			i++
		}

		// Si el error fue EOF o UnexpectedEOF, significa que ya no hay más datos
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
	}
	return nil
}
