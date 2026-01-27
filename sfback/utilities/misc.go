package utilities

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font/gofont/goregular" // Fuente de Go estándar
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

func GenerateQR(url, fullpath string) bool {
	qrCode, _ := qrcode.New(url, qrcode.Medium)
	err := qrCode.WriteFile(256, fullpath)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if _, err := os.Stat(fullpath); err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Printf("QR code generated and saved as %s\n", fullpath)
	return true
}

// AddTextToQR compone el código QR existente en 'qrPath' con el texto proporcionado debajo.
// Sobrescribe el archivo QR original con la nueva imagen compuesta.
func AddTextToQR(qrPath string, text string) bool {
	// 1. Cargar la imagen del QR existente
	qrFile, err := os.Open(qrPath)
	if err != nil {
		fmt.Printf("Error al abrir el archivo QR: %v\n", err)
		return false
	}
	defer qrFile.Close()

	qrImage, err := png.Decode(qrFile)
	if err != nil {
		fmt.Printf("Error al decodificar la imagen QR: %v\n", err)
		return false
	}

	// 2. Definir las dimensiones y el tamaño de la fuente
	qrBounds := qrImage.Bounds()
	qrWidth := qrBounds.Dx()
	qrHeight := qrBounds.Dy()
	textHeightSpace := 50 // Espacio extra para el texto
	newHeight := qrHeight + textHeightSpace

	// 3. Crear el nuevo lienzo (imagen)
	// Creamos un lienzo RGBA del tamaño del QR más el espacio para el texto.
	newBounds := image.Rect(0, 0, qrWidth, newHeight)
	newImage := image.NewRGBA(newBounds)

	// Llenar el fondo del lienzo de blanco
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(newImage, newBounds, image.NewUniform(white), image.Point{}, draw.Src)

	// 4. Dibujar la imagen del QR en la parte superior
	// El punto de inicio es (0, 0)
	draw.Draw(newImage, qrBounds, qrImage, image.Point{}, draw.Src)

	// 5. Preparar la fuente y el contexto Freetype
	ft, err := truetype.Parse(goregular.TTF)
	if err != nil {
		fmt.Printf("Error al cargar la fuente: %v\n", err)
		return false
	}

	// Crear un contexto de Freetype
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(ft)
	c.SetFontSize(16.0) // Tamaño de la fuente para el texto
	c.SetClip(newImage.Bounds())
	c.SetDst(newImage)
	c.SetSrc(image.NewUniform(color.Black)) // Color del texto

	// 6. Calcular la posición del texto (centrado horizontalmente)
	pt := freetype.Pt(0, 0) // Punto de dibujo (se recalcula abajo)

	// Usamos un valor fijo para simplificar, idealmente se usaría la función MeasureString.
	// Asumimos que el texto es corto y centramos la posición Y debajo del QR.
	// La coordenada Y es el borde inferior del QR + un pequeño margen
	textY := qrHeight + 20

	// Calcular la posición X para centrar el texto (aproximadamente)
	// Esta es una aproximación, ya que la medición precisa es más compleja.
	// Se divide el ancho total del QR y se resta un estimado del ancho del texto
	textWidthEstimate := len(text) * 8 // Estimación basada en 8px por caracter
	textX := (qrWidth / 2) - (textWidthEstimate / 2)

	// Configurar el punto de inicio para el dibujo del texto
	pt = freetype.Pt(textX, textY)

	// 7. Dibujar el texto
	_, err = c.DrawString(text, pt)
	if err != nil {
		fmt.Printf("Error al dibujar el texto: %v\n", err)
		return false
	}

	// 8. Sobrescribir el archivo original
	outFile, err := os.Create(qrPath)
	if err != nil {
		fmt.Printf("Error al crear el archivo de salida: %v\n", err)
		return false
	}
	defer outFile.Close()

	if err := png.Encode(outFile, newImage); err != nil {
		fmt.Printf("Error al codificar la nueva imagen PNG: %v\n", err)
		return false
	}

	return true
}

func GetHashC14N(xmlFull string, idNode string, algo string) (string, error) {
	scriptPath := "./utilities/c14n.py"

	// Pasamos los argumentos. Ojo: si el XML es muy grande,
	// algunos SO pueden tener límites en el tamaño de los argumentos.
	cmd := exec.Command("python", scriptPath, xmlFull, idNode, algo)

	// Capturamos la salida (stdout)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("python error: %v, output: %s", err, string(out))
	}

	// El resultado es el string en Base64
	hashBase64 := strings.TrimSpace(string(out))
	return hashBase64, nil
}

func Bool2Int(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func GetClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" {
		return cfip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}
	return r.RemoteAddr
}
