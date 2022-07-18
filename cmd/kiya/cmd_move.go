package main

import (
	"context"
	"fmt"
	"log"

	"github.com/emicklei/tre"
	"github.com/kramphub/kiya/backend"
)

// commandMove transfers a secret from a source to a target profile.
func commandMove(
	ctx context.Context,
	b backend.Backend,
	source *backend.Profile,
	sourceKey string,
	target *backend.Profile,
	targetKey string,
) {

	if promptForYes(fmt.Sprintf("Are you sure you want to move [%s] from [%s] (y/N)", sourceKey, target.Label)) {
		if err := move(ctx, b, source, sourceKey, target, targetKey); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully moved [%s] to [%s]\n", sourceKey, target.Label)
	} else {
		log.Fatalln("delete aborted")
	}
}

func move(
	ctx context.Context,
	b backend.Backend,
	source *backend.Profile,
	sourceKey string,
	target *backend.Profile,
	targetKey string) error {

	// fetch value for key from source
	sourceValue, err := b.Get(ctx, source, sourceKey)
	if err != nil {
		return tre.New(err, "get source key failed", "key", sourceKey)
	}

	exists, _ := b.CheckExists(ctx, target, targetKey)
	if err := b.Put(ctx, target, targetKey, string(sourceValue), exists); err != nil {
		return tre.New(err, "save key failed", targetKey)
	}
	// delete key from source
	err = b.Delete(ctx, source, sourceKey)
	return tre.New(err, "could not delete key", targetKey)
}
