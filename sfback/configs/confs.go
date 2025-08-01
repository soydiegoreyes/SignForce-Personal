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
