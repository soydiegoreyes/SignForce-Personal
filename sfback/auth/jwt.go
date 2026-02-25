package auth

import (
	"crypto/rand"
	"encoding/hex"
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

func generateJTI() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Genera un token JWT para el usuario
func GenerateJWT(user *objects.User) (string, error) {
	var jwt_exp int
	var err error
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}
	jwt_exp, err = strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		fmt.Println("Tiempo de expiracion por default")
		jwt_exp = 60
	}
	now := time.Now()
	expirationUnix := now.Add(time.Minute * time.Duration(jwt_exp)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.Uid,
		"jti": jti,
		"iat": now.Unix(),
		"exp": expirationUnix,
	})

	ss, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// Valida un token JWT y devuelve los claims
// Valida un token JWT y devuelve los claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)

	if err != nil {
		return nil, err
	}

	if !parsed.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims inválidos")
	}

	return claims, nil
}
