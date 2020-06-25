package main

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"google.golang.org/api/cloudkms/v1"

	"github.com/kramphub/kiya"
)

// commandPutPasteGenerate ...
func commandPutPasteGenerate(kmsService *cloudkms.Service, storageService *cloudstore.Client,
	target kiya.Profile, command, key, value string, mustPrompt bool) {

	if kiya.CheckSecretExists(storageService, target, key) {
		if mustPrompt && !promptForYes(fmt.Sprintf("Are you sure to overwrite [%s] from [%s] (y/N)? ", key, target.Label)) {
			log.Fatalln(command + " aborted")
			return
		}
	}
	if err := kiya.PutSecret(kmsService, storageService, target, key, value); err != nil {
		log.Fatal(err)
	}
}
