package kiya

import (
	"log"

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
