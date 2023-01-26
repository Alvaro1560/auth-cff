package password

import (
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

func Compare(id string, hashedPassword, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(p))
	if err != nil {
		logger.Warning.Printf("la contraseña de %s no es válida: %v", id, err)
		return false
	}
	return true
}

func Encrypt(password string) string {
	bp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error.Printf("generando el hash del password: %v", err)
	}
	return string(bp)
}
