package kiya

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// CommandMove transfers a secret from a source to a target profile.
func CommandMove(
	kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	source Profile,
	sourceKey string,
	target Profile,
	targetKey string) {

	// fetch value for key from source
	sourceValue, err := GetValueByKey(kmsService, storageService, sourceKey, source)
	if err != nil {
		log.Fatal(tre.New(err, "get source key failed", "key", sourceKey))
	}
	// store value for key to target
	CommandPutPasteGenerate(kmsService, storageService, target, "put", targetKey, sourceValue, true)
	fmt.Printf("Successfully copied [%s] to [%s]\n", sourceKey, target.Label)
	// delete key from source
	CommandDelete(kmsService, storageService, source, sourceKey)
}
