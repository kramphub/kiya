package main

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"google.golang.org/api/cloudkms/v1"

	"github.com/kramphub/kiya"
)

// commandMove transfers a secret from a source to a target profile.
func commandMove(
	kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	source kiya.Profile,
	sourceKey string,
	target kiya.Profile,
	targetKey string) {

	if promptForYes(fmt.Sprintf("Are you sure you want to move [%s] from [%s] (y/N)", sourceKey, target.Label)) {
		if err := kiya.Move(kmsService, storageService, source, sourceKey, target, targetKey); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully copied [%s] to [%s]\n", sourceKey, target.Label)
	} else {
		log.Fatalln("delete aborted")
	}
}
