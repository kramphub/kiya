package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func commandList(storageService *cloudstore.Client, target profile, filter string) {
	ctx := context.Background()
	bucket := storageService.Bucket(target.Bucket)
	query := &cloudstore.Query{}
	it := bucket.Objects(ctx, query)
	data := [][]string{}
	filteredCount := 0

	for {
		next, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatal(tre.New(err, "list failed"))
		}
		if len(filter) > 0 {
			if !caseInsensitiveContains(next.Name, filter) {
				filteredCount++
				continue
			}
		}
		data = append(data, []string{fmt.Sprintf("kiya %s copy %s", target.Label, next.Name), next.Created.Format(time.RFC822), next.Owner})
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
