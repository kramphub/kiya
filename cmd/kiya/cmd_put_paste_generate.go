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
	if exists, _ := b.CheckExists(ctx, target, key); exists {
		if mustPrompt && !promptForYes(fmt.Sprintf("Are you sure to overwrite [%s] from [%s] (y/N)? ", key, target.Label)) {
			log.Fatalln(command + " aborted")
			return
		}
	}

	if shouldPromptForPassword(b) {
		pass := promptForPassword()
		b.SetMasterPassword(pass)
	}

	if err := b.Put(ctx, target, key, value); err != nil {
		log.Fatal(err)
	}
}
