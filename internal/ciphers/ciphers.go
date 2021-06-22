package ciphers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"io"
)

var secretKey string

func init() {
	secretKey = "secret"
}

// EncryptGCM cifra un contenido, recibe el texto a cifrar y la llave.
func Encrypt(textToEncrypt string) string {
	key := to32Bytes(secretKey)
	if len(key) != 32 {
		logger.Error.Printf("se intentó cifrar con EncryptGCM un contenido con una clave de un tamaño diferente a 32 bytes")
		return ""
	}

	c, err := aes.NewCipher(key[:])
	if err != nil {
		logger.Error.Printf("no se pudo crear el cifrado aes.NewCipher en EncryptGCM: %v", err)
		return ""
	}

	aesGCM, err := cipher.NewGCM(c)
	if err != nil {
		logger.Error.Printf("no se pudo crear el cifrado ciphers.NewGCM en EncryptGCM: %v", err)
		return ""
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error.Printf("no se pudo leer el contenido de rand.Reader al tratar de cifrar un contenido con EncryptGCM: %v", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(aesGCM.Seal(nonce, nonce, []byte(textToEncrypt), nil))
}

// DecryptGCM decifra un contenido con una clave de 32 bytes
func Decrypt(encryptedString string) ([]byte, error) {
	key := to32Bytes(secretKey)
	if len(key) != 32 {
		logger.Error.Printf("se intentó cifrar con EncryptGCM un contenido con una clave de un tamaño diferente a 32 bytes")
		return nil, errors.New("la clave de cifrado debe ser de 32 bytes")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		logger.Error.Printf("no se pudo hacer decode del texto cifrado a un slice de bytes base64Decode en DecryptGCM: %v", err)
		return nil, err
	}
	c, err := aes.NewCipher(key[:])
	if err != nil {
		logger.Error.Printf("no se pudo crear el cifrado aes.NewCipher en DecryptGCM: %v", err)
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger.Error.Printf("no se pudo crear el cifrado ciphers.NewGCM en DencryptGCM: %v", err)
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func to32Bytes(s string) []byte {
	return []byte(fmt.Sprintf("%x", md5.Sum([]byte(s))))
}

func GetSecret() string {
	return secretKey
}
