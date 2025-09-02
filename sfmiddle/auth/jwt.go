package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Genera un token JWT para el usuario
func GenerateJWT(idUser, idTeam, roleApp, idInst, authStatusInst string) (string, error) {
	var jwt_exp int
	var err error
	jwt_exp, err = strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		fmt.Println("Tiempo de expiracion por default")
		jwt_exp = 300
	}
	now := time.Now()
	expirationUnix := now.Add(time.Minute * time.Duration(jwt_exp)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":      idUser,
		"iid":      idInst,
		"authInst": authStatusInst,
		"team":     idTeam,
		"role":     roleApp,
		"exp":      expirationUnix,
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
