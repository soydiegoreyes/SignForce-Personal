package auth

import (
	"os"
	"sfback/objects"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_KEY = os.Getenv("JWT_KEY")
	EXP, _  = strconv.Atoi(os.Getenv("JWT_EXP"))
	JWT_EXP = time.Minute * time.Duration(EXP)
)

// Genera un token JWT para el usuario
func GenerateJWT(user *objects.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Uid,                       // Incluye datos relevantes
		"exp": time.Now().Add(JWT_EXP).Unix(), // Fecha de expiración
	})

	return token.SignedString([]byte(JWT_KEY))
}

// Valida un token JWT y devuelve los claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_KEY), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
