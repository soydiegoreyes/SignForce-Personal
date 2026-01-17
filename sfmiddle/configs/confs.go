package configs

// HashConfig configuración para el hash
type HashConfig struct {
	Algorithm string // "sha256", "sha512"
	Encoding  string // "hex", "base64"
	ChunkSize int    // Tamaño de chunk para archivos grandes (en bytes)
}

var HashConf = HashConfig{
	Algorithm: "sha256",
	Encoding:  "b64",
	ChunkSize: 1024 * 1024 * 20, // 20MB chunks
}

var AllowedExtensions = map[string]bool{
	"pdf": true, "doc": true, "docx": true, "txt": true, "jpg": true, "jpeg": true,
	"png": true, "gif": true, "zip": true, "rar": false, "xlsx": true, "xls": true,
	"pptx": false, "ppt": false, "pem": true, "key": true, "cer": true, "crt": true,
}
