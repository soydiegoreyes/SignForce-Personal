package auth

import (
	//"sfback/db"
	"fmt"
	"net/http"
	"sfmiddle/objects"
	"strings"
	"sync"
	"time"
)

// tipo que guarda un usuario y su tiempo de expiracion de sesion
type ActiveUser struct {
	User      *objects.User
	ExpiresAt time.Time
}

// diccionario que guarda los usuarios activos mediante su indice y tiene un mutex para comunicarse
var (
	activeUsers = make(map[string]ActiveUser)
	mu          sync.Mutex
)

// AddUser añade un usuario activo
func AddUser(id string, user *objects.User) {
	mu.Lock()
	defer mu.Unlock()
	activeUsers[id] = ActiveUser{
		User:      user,
		ExpiresAt: time.Now().Add(time.Hour), // Expira en 1 hora
	}
}

// GetUser verifica si un usuario está activo dentro del arreglo de usuarios activos
func GetUser(id string) (*objects.User, bool) {
	mu.Lock()
	defer mu.Unlock()

	activeUser, exists := activeUsers[id]
	if !exists || time.Now().After(activeUser.ExpiresAt) {
		return nil, false
	}
	return activeUser.User, true
}

// DeleteUser elimina un usuario activo (para logout)
func DeleteUser(id string) {
	mu.Lock()
	defer mu.Unlock()
	delete(activeUsers, id)
}

func GetUserFromRequest(request *http.Request) (*objects.User, error) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header faltante")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("formato inválido del token")
	}

	tokenStr := parts[1]
	claims, err := ValidateJWT(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("token inválido")
	}
	// jwt-> {"uid": "id_user", "expires": "time"}
	uid, ok := claims["uid"].(string)
	if !ok {
		return nil, fmt.Errorf("UID no presente en el token")
	}

	user, ok := GetUser(uid)
	if !ok {
		return nil, fmt.Errorf("usuario no activo o no encontrado")
	}

	return user, nil
}

func ValidateUserDocPermissions(idInst, idUser string, docInfo, userInfo map[string]string) bool {
	var check bool = true
	if docInfo["activeDoc"] != "1" {
		check = false
		return check
	}

	switch docInfo["authUseStatus"] {
	// 0: publico, 1: misma institución, 2: institución y usuario, 3: depende del rol del usuario
	case "1":
		if docInfo["ownerInstDoc_fk"] != idInst {
			check = false
			return check
		}

	case "2":
		if docInfo["ownerInstDoc_fk"] != idInst || docInfo["creatorUserDoc_fk"] != idUser {
			check = false
			return check
		}
	case "3":
		// 1: misma institución cualquier rol
		// 2: misma institución root, 3: misma institución admin, 4: misma institución regular user,
		if docInfo["authRoleStatus"] == "1" || docInfo["authRoleStatus"] == "2" || docInfo["authRoleStatus"] == "3" || docInfo["authRoleStatus"] == "4" {
			// lo ve la institución pero no el usuario
			if docInfo["ownerInstDoc_fk"] != idInst {
				check = false
				return check
			}
			switch docInfo["authRoleStatus"] {
			case "2":
				if userInfo["roleAppUser_fk"] != "1" {
					check = false
					return check
				}
			case "3":
				if !(userInfo["roleAppUser_fk"] == "1" || userInfo["roleAppUser_fk"] == "2") {
					check = false
					return check
				}
			case "4":
				if !(userInfo["roleAppUser_fk"] == "1" || userInfo["roleAppUser_fk"] == "2" || userInfo["roleAppUser_fk"] == "3") {
					check = false
					return check
				}
			}
		}
		// 5: institución y usuario root, 6: institución y usuario admin, 7: institución y usuario regular user
		if docInfo["authRoleStatus"] == "5" || docInfo["authRoleStatus"] == "6" || docInfo["authRoleStatus"] == "7" {
			if docInfo["ownerInstDoc_fk"] != idInst || docInfo["creatorUserDoc_fk"] != idUser {
				check = false
				return check
			}
			switch docInfo["authRoleStatus"] {
			case "5":
				if userInfo["roleAppUser_fk"] != "1" {
					check = false
					return check
				}
			case "6":
				if !(userInfo["roleAppUser_fk"] == "1" || userInfo["roleAppUser_fk"] == "2") {
					check = false
					return check
				}
			case "7":
				if !(userInfo["roleAppUser_fk"] == "1" || userInfo["roleAppUser_fk"] == "2" || userInfo["roleAppUser_fk"] == "3") {
					check = false
					return check
				}
			}
		}
	}
	return check
}
