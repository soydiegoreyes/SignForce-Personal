package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
func GenerateJWT(idUser, roleApp, idInst, authStatusInst string) (string, error) {
	var jwt_exp int
	var err error
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}
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
		"role":     roleApp,
		"jti":      jti,
		"iat":      now.Unix(),
		"exp":      expirationUnix,
	})

	ss, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// Genera un token JWT para el usuario
func GenerateJWTBio(bio, jti string) (string, error) {
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
		"bio": bio,
		"sid": jti,
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
