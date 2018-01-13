package main

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func command_put_paste_generate(kmsService *cloudkms.Service, storageService *cloudstore.Client,
	target profile, command, key, value string) {
	// check for exists
	_, err := loadSecret(storageService, target, key)
	if err == nil {
		if !promptForYes(fmt.Sprintf("Are you sure to overwrite [%s] from [%s] (y/N)? ", key, target.Label)) {
			fmt.Println(command + " aborted")
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
