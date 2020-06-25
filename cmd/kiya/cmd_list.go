package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cloudstore "cloud.google.com/go/storage"
	"github.com/olekukonko/tablewriter"

	"github.com/kramphub/kiya"
)

// commandList lists keys in a specific profile
func commandList(storageService *cloudstore.Client, target kiya.Profile, filter string) {
	keys, err := kiya.List(storageService, target)
	if err != nil {
		log.Fatal(err)
	}

	var data [][]string
	filteredCount := 0

	for _, k := range keys {
		if len(filter) > 0 {
			if !caseInsensitiveContains(k.Name, filter) {
				filteredCount++
				continue
			}
		}
		data = append(data, []string{fmt.Sprintf("kiya %s copy %s", target.Label, k.Name), k.CreatedAt.Format(time.RFC822), k.Owner})
	}

	if len(filter) > 0 {
		fmt.Printf("Showing %d key(s) matching '%s', skipped %d key(s)\n", len(data), filter, filteredCount)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Copy to clipboard command", "Created", "Creator"})
	table.AppendBulk(data)
	table.Render() // writes to stdout
}

func caseInsensitiveContains(key, filter string) bool {
	key, filter = strings.ToLower(key), strings.ToLower(filter)
	return strings.Contains(key, filter)
}
