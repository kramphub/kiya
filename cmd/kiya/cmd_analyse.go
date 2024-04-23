package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kramphub/kiya/backend"

	passwordvalidator "github.com/wagslane/go-password-validator"
)

// commandAnalyse perform an analysis of the profile
func commandAnalyse(ctx context.Context, b backend.Backend, target *backend.Profile) {
	fmt.Printf("loading all secrets in profile [%s]...\n", target.Label)

	when := time.Now()
	kv, err := getAllItems(ctx, b, *target, "")
	if err != nil {
		log.Printf("error: failed to get all items, %s", err.Error())
		return
	}
	fmt.Printf("loaded all secrets [%d] in %v\n", len(kv), time.Since(when))

	fmt.Println("detecting weak secrets with entroopy < 50 ...")
	when = time.Now()
	count := 0
	for k, v := range kv {
		entropy := passwordvalidator.GetEntropy(string(v))
		if entropy < 50 {
			count++
			fmt.Printf("WARNING: kiya %s copy %s (entropy:%.2f length:%d)\n", target.Label, k, entropy, len(v))
		}
	}
	fmt.Printf("detected weak secrets [%d] in %v\n", count, time.Since(when))
}

// getAllItems returns all keys and values in store.
func getAllItems(ctx context.Context, b backend.Backend, target backend.Profile, filter string) (map[string][]byte, error) {
	items := make(map[string][]byte)
	keys := commandList(ctx, b, &target, filter)
	for _, key := range keys {
		buf, err := b.Get(ctx, &target, key.Name)
		if err != nil {
			fmt.Printf("error: get key '%s' failed, %s", key.Name, err.Error())
			continue
		}
		items[key.Name] = buf
	}
	return items, nil
}
