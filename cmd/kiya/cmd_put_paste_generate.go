package main

import (
	"fmt"
	"log"

	cloudkms "cloud.google.com/go/kms/apiv1"
	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"

	"github.com/kramphub/kiya"
)

// commandPutPasteGenerate ...
func commandPutPasteGenerate(kmsService *cloudkms.KeyManagementClient, storageService *cloudstore.Client,
	target kiya.Profile, command, key, value string, mustPrompt bool) {
	// check for exists
	_, err := kiya.LoadSecret(storageService, target, key)
	if err == nil {
		if mustPrompt && !promptForYes(fmt.Sprintf("Are you sure to overwrite [%s] from [%s] (y/N)? ", key, target.Label)) {
			log.Fatalln(command + " aborted")
			return
		}
	}
	encryptedValue, err := kiya.GetEncryptedValue(kmsService, target, value)
	if err != nil {
		log.Fatal(tre.New(err, command+" failed", "key", key))
	}
	if err := kiya.StoreSecret(storageService, target, key, encryptedValue); err != nil {
		log.Fatal(tre.New(err, command+" failed", "key", key, "encryptedValue", encryptedValue))
	}
}
