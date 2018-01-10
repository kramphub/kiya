package main

import (
	"log"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// command_move transfers a secret from a source to a target profile.
func command_move(
	kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	source profile,
	sourceKey string,
	target profile,
	targetKey string) {

	// fetch value for key from source
	sourceValue, err := getValueByKey(kmsService, storageService, sourceKey, source)
	if err != nil {
		log.Fatal(tre.New(err, "get source key failed", "key", sourceKey, "err", err))
	}
	// store value for key to target
	command_put_paste(kmsService, storageService, target, "put", targetKey, sourceValue)
	// delete key from source
	command_delete(kmsService, storageService, source, sourceKey)
}
