package main

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kramphub/kiya/backend"
)

// commandList lists keys in a specific profile
func commandList(ctx context.Context, b backend.Backend, target *backend.Profile, filter string) []backend.Key {
	keys, err := b.List(ctx, target)
	if err != nil {
		log.Fatal(err)
	}

	filteredKeys := make([]backend.Key, 0)
	for _, k := range keys {
		if len(filter) > 0 {
			if !caseInsensitiveContains(k.Name, filter) {
				continue
			}
		}

		filteredKeys = append(filteredKeys, k)
	}

	return filteredKeys
}

// writeTable writes a human-readable table with parameters info.
func writeTable(keys []backend.Key, target *backend.Profile, filter string) {
	filteredCount := 0

	data := make([][]string, 0)

	for _, k := range keys {
		if len(filter) > 0 {
			if !caseInsensitiveContains(k.Name, filter) {
				filteredCount++
				continue
			}
		}
		data = append(data, []string{fmt.Sprintf("kiya %s copy %s", target.Label, k.Name), k.CreatedAt.Format(time.RFC822), k.Info})
	}

	if len(filter) > 0 {
		fmt.Printf("Showing %d key(s) matching '%s', skipped %d key(s)\n", len(data), filter, filteredCount)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Copy to clipboard command", "Created", "Info"})
	table.AppendBulk(data)
	table.Render() // writes to stdout
}

func caseInsensitiveContains(key, filter string) bool {
	key, filter = strings.ToLower(key), strings.ToLower(filter)
	return strings.Contains(key, filter)
}
