package models

// ==================== REQUESTS ======================//
type GetUserRequest struct {
	IdUser string   `json:"iduser"`
	Fields []string `json:"fields"`
}

type UserSignRequest struct {
	IdUser        string `json:"iduser"`
	Password      string `json:"password"`
	HashedMessage string `json:"hash"`
	SignType      string `json:"signtype"`
}

// Decodifica cuerpo de hashDataB64
type HashDataB64 struct {
	DataB64   string `json:"data_b64"`
	DigestAlg string `json:"digest_alg"`
}

type UploadKeysReq struct {
	IdInst  string `json:"IdInst"`
	IdUser  string `json:"IdUser"`
	PassKey string `json:"PassKey"`
	NameKey string `json:"NameKey"`
	NameCer string `json:"NameCer"`
	HashKey string `json:"HashKey"`
	HashCer string `json:"HashCer"`
	Path    string `json:"Path"`
}

// Decodificar el cuerpo (ej: {"id": "123", "password": "secret"})
type LoginUserReq struct {
	IdUser   string `json:"iduser"`
	Password string `json:"password"`
}

// ==================== RESPONSES ====================//
// LoginResponse estructura para la respuesta del login
type LoginResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// UserResponse devuelve datos del usuario (sin sensibles)
type UserResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

// respuesta al endpoint SignUser
type SignatureResponse struct {
	Operation string `json:"Operation"`
	Signature string `json:"signature"`
	Check     bool   `json:"check"`
}

// respuesta al endpoint Hash
type HashDataResponse struct {
	Operation     string `json:"operation"`
	HashedMessage string `json:"hashed_message"`
	Check         bool   `json:"check"`
}

type UploadKeysResponse struct {
	Valid      bool   `json:"valid"`
	Expiration string `json:"expiration"`
	Owner      string `json:"owner"`
	KeysId     string `json:"keysid"`
}
