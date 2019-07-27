package kiya

import (
	"fmt"
	"log"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// commandMove transfers a secret from a source to a target profile.
func commandMove(
	kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	source profile,
	sourceKey string,
	target profile,
	targetKey string) {

	// fetch value for key from source
	sourceValue, err := getValueByKey(kmsService, storageService, sourceKey, source)
	if err != nil {
		log.Fatal(tre.New(err, "get source key failed", "key", sourceKey))
	}
	// store value for key to target
	commandPutPasteGenerate(kmsService, storageService, target, "put", targetKey, sourceValue, true)
	fmt.Printf("Successfully copied [%s] to [%s]\n", sourceKey, target.Label)
	// delete key from source
	commandDelete(kmsService, storageService, source, sourceKey)
}
