package encryp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"gitlab.com/e-capture/ecatch-bpm/ecatch-auth/internal/logger"
	"io"
)

var KeyString string

// EncryptGCM cifra un contenido, recibe el texto a cifrar y la llave.
func EncryptGCM(plaintext string, keyString string) (string, error) {
	key := to32Bytes(keyString)
	if len(key) != 32 {
		logger.Error.Printf("se intentó cifrar con EncryptGCM un contenido con una clave de un tamaño diferente a 32 bytes")
		return "", errors.New("la clave de cifrado debe ser de 32 bytes")
	}

	c, err := aes.NewCipher(key[:])
	if err != nil {
		logger.Error.Printf("no se pudo crear el cifrado aes.NewCipher en EncryptGCM: %v", err)
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger.Error.Printf("no se pudo crear el cifrado cipher.NewGCM en EncryptGCM: %v", err)
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error.Printf("no se pudo leer el contenido de rand.Reader al tratar de cifrar un contenido con EncryptGCM: %v", err)
		return "", err
	}

	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil)), nil
}

// DecryptGCM decifra un contenido con una clave de 32 bytes
func DecryptGCM(cipherToDecrypt string, keyString string) ([]byte, error) {
	key := to32Bytes(keyString)
	if len(key) != 32 {
		logger.Error.Printf("se intentó cifrar con EncryptGCM un contenido con una clave de un tamaño diferente a 32 bytes")
		return nil, errors.New("la clave de cifrado debe ser de 32 bytes")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cipherToDecrypt)
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
		logger.Error.Printf("no se pudo crear el cifrado cipher.NewGCM en DencryptGCM: %v", err)
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

func ExampleNewGCMEncrypter(text string) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key := []byte("AES256Key-32Characters1234567890")
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	fmt.Printf("%x\n", ciphertext)
}

func ExampleNewGCMDecrypter(text string) {
	// The key argument should be the AES key, either 16 or 32 bytes
	// to select AES-128 or AES-256.
	key := []byte("AES256Key-32Characters1234567890")
	ciphertext, _ := hex.DecodeString("f90fbef747e7212ad7410d0eee2d965de7e890471695cddd2a5bc0ef5da1d04ad8147b62141ad6e4914aee8c512f64fba9037603d41de0d50b718bd665f019cdcd")

	nonce, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", string(plaintext))
}
