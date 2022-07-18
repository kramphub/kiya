package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kramphub/kiya/backend"
)

// commandPutPasteGenerate ...
func commandPutPasteGenerate(
	ctx context.Context,
	b backend.Backend,
	target *backend.Profile,
	command, key, value string,
	mustPrompt bool,
) {

	overwrite := false
	if exists, _ := b.CheckExists(ctx, target, key); exists {
		if mustPrompt && !promptForYes(fmt.Sprintf("Are you sure to overwrite [%s] from [%s] (y/N)? ", key, target.Label)) {
			log.Fatalln(command + " aborted")
			return
		}
		overwrite = true
	}

	if err := b.Put(ctx, target, key, value, overwrite); err != nil {
		log.Fatal(err)
	}
}
