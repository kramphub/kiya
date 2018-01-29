package main

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func commandDelete(kmsService *cloudkms.Service, storageService *cloudstore.Client, target profile, key string) {
	_, err := getValueByKey(kmsService, storageService, key, target)
	if err != nil {
		log.Fatal(tre.New(err, "delete failed", "key", key, "err", err))
	}
	if promptForYes(fmt.Sprintf("Are you sure to delete [%s] from [%s] (y/N)? ", key, target.Label)) {
		if err := deleteSecret(storageService, target, key); err != nil {
			fmt.Printf("failed to delete [%s] from [%s] because [%v]\n", key, target.Label, err)
		} else {
			fmt.Printf("successfully deleted [%s] from [%s]\n", key, target.Label)
		}
	} else {
		fmt.Println("delete aborted")
	}
}
