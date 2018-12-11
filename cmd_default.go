package main

import (
	"fmt"
	"log"
	"strconv"

	cloudstore "cloud.google.com/go/storage"
	"github.com/atotto/clipboard"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func commandDefault(
	kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	source profile,
	cmdOrIndex string) {
	// check for list number
	listIndex, err := strconv.Atoi(cmdOrIndex)
	if err != nil {
		fmt.Println("unknown command", cmdOrIndex)
		return
	}
	// resolve list number to key
	key := keyAtListIndex(storageService, source, listIndex)
	if len(key) == 0 {
		fmt.Println("no key at index", cmdOrIndex)
		return
	}
	// get the decrypted value
	value, err := getValueByKey(kmsService, storageService, key, source)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
	}
	// write to clipboard
	if err := clipboard.WriteAll(value); err != nil {
		log.Fatal(tre.New(err, "copy failed", "key", key, "err", err))
	}
}
