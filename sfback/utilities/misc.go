package utilities

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

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
