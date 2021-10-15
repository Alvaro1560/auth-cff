package ciphers

import (
	"fmt"
	openssl "github.com/Luzifer/go-openssl/v4"
	"github.com/google/uuid"
	"sync"
)

var (
	secretKey string
	once     sync.Once
	secretKeyTemp string
)

func init() {
	//TODO get SecretKey
	secretKey = "204812730425442A472D2F423F452847"
	once.Do(func() {
		secretKeyTemp = uuid.New().String()
	})
}

func Encrypt(strToEncrypt string) string{

	o := openssl.New()

	enc, err := o.EncryptBytes(secretKey, []byte(strToEncrypt), openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		return ""
	}

	return string(enc)
}


func Decrypt(strToDecrypt string) string{
	o := openssl.New()
	dec, err := o.DecryptBytes(secretKey, []byte(strToDecrypt), openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		return ""
	}
	return  string(dec)
}

func EncryptKeyTemp(strToEncrypt string) string{

	o := openssl.New()

	enc, err := o.EncryptBytes(secretKey, []byte(strToEncrypt), openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		return ""
	}

	return string(enc)
}

func DecryptKeyTemp(strToDecrypt string) string{
	o := openssl.New()
	dec, err := o.DecryptBytes(secretKey, []byte(strToDecrypt), openssl.BytesToKeyMD5)
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		return ""
	}
	return  string(dec)
}


func GetSecret() string {
	return  secretKey
}

func GetSecretKeyTemp() string {
	return  secretKey
}
