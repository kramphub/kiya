package main

import (
	"fmt"
	"log"
	"os"
	"time"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func commandList(storageService *cloudstore.Client, target profile) {
	ctx := context.Background()
	bucket := storageService.Bucket(target.Bucket)
	query := &cloudstore.Query{}
	it := bucket.Objects(ctx, query)
	data := [][]string{}
	for {
		next, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Fatal(tre.New(err, "list failed"))
		}
		data = append(data, []string{fmt.Sprintf("kiya %s copy %s", target.Label, next.Name), next.Created.Format(time.RFC822), next.Owner})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Copy to clipboard command", "Created", "Creator"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render() // writes to stdout
}
