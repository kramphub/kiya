package kiya

import (
	"crypto/rand"
	"log"
	"math/big"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"google.golang.org/api/cloudkms/v1"
)

// GetValueByKey is very self explanatory :P
func GetValueByKey(kmsService *cloudkms.Service, storageService *cloudstore.Client, key string, target Profile) (string, error) {
	encryptedValue, err := LoadSecret(storageService, target, key)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "key", key))
		return "", err
	}
	decryptedValue, err := GetDecryptedValue(kmsService, target, encryptedValue)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "cipherText", encryptedValue))
		return "", err
	}

	return decryptedValue, nil
}

// Generate_secret generates a random key
func Generate_secret(length int, chars string) (string, error) {
	if len(chars) == 0 {
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	}

	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		randomString[i] = chars[n.Int64()]
	}

	return string(randomString), nil
}
