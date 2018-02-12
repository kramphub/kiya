package main

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"google.golang.org/api/cloudkms/v1"
)

func commandPutPasteGenerate(kmsService *cloudkms.Service, storageService *cloudstore.Client,
	target profile, command, key, value string, mustPrompt bool) {
	// check for exists
	_, err := loadSecret(storageService, target, key)
	if err == nil {
		if mustPrompt && !promptForYes(fmt.Sprintf("Are you sure to overwrite [%s] from [%s] (y/N)? ", key, target.Label)) {
			fmt.Println(command + " aborted")
			return
		}
	}
	encryptedValue, err := getEncryptedValue(kmsService, target, value)
	if err != nil {
		log.Fatal(tre.New(err, command+" failed", "key", key, "value", value))
	}
	if err := storeSecret(storageService, target, key, encryptedValue); err != nil {
		log.Fatal(tre.New(err, command+" failed", "key", key, "value", value, "encryptedValue", encryptedValue))
	}
}
