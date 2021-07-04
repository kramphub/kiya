package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kramphub/kiya/backend"
)

// commandDelete deletes a stored key
func commandDelete(ctx context.Context, b backend.Backend, target *backend.Profile, key string) {
	if promptForYes(fmt.Sprintf("Are you sure to delete [%s] from [%s] (y/N)? ", key, target.Label)) {
		//if err := kiya.DeleteSecret(kmsService, storageService, target, key); err != nil {
		if err := b.Delete(ctx, target, key); err != nil {
			fmt.Printf("failed to delete [%s] from [%s] because [%v]\n", key, target.Label, err)
		} else {
			fmt.Printf("Successfully deleted [%s] from [%s]\n", key, target.Label)
		}
	} else {
		log.Fatalln("delete aborted")
	}
}
