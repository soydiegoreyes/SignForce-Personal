package auth

import (
	"fmt"
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
	var jwt_exp int
	var err error
	jwt_exp, err = strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		fmt.Println("Tiempo de expiracion por default")
		jwt_exp = 60
	}
	now := time.Now()
	expirationUnix := now.Add(time.Minute * time.Duration(jwt_exp)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Uid,
		"exp": expirationUnix,
	})

	ss, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// Valida un token JWT y devuelve los claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo inesperado: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		fmt.Println("error parseando:", err)
		return nil, err
	}
	fmt.Println("Token válido")
	return parsed.Claims.(jwt.MapClaims), nil
}
