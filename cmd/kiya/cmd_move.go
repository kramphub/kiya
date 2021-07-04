package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kramphub/kiya"
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
		if err := kiya.Move(ctx, b, source, sourceKey, target, targetKey); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully moved [%s] to [%s]\n", sourceKey, target.Label)
	} else {
		log.Fatalln("delete aborted")
	}
}
