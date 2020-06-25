package main

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"google.golang.org/api/cloudkms/v1"

	"github.com/kramphub/kiya"
)

// commandDelete deletes a stored key
func commandDelete(kmsService *cloudkms.Service, storageService *cloudstore.Client, target kiya.Profile, key string) {
	if promptForYes(fmt.Sprintf("Are you sure to delete [%s] from [%s] (y/N)? ", key, target.Label)) {
		if err := kiya.DeleteSecret(kmsService, storageService, target, key); err != nil {
			fmt.Printf("failed to delete [%s] from [%s] because [%v]\n", key, target.Label, err)
		} else {
			fmt.Printf("Successfully deleted [%s] from [%s]\n", key, target.Label)
		}
	} else {
		log.Fatalln("delete aborted")
	}
}
