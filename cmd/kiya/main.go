package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/atotto/clipboard"

	"github.com/kramphub/kiya"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"golang.org/x/net/context"
	"google.golang.org/api/cloudkms/v1"
)

var version = "build-" + time.Now().String()

const (
	doPrompt    = true
	doNotPrompt = false
)

func main() {
	flag.Parse()
	if *oVersion {
		fmt.Println("kiya version", version)
		os.Exit(0)
	}
	kiya.LoadConfiguration(*oConfigFilename)
	if len(flag.Args()) < 2 {
		fmt.Println("kiya [flags] [profile] [get|put|delete|list|template|copy|paste|move|generate] [|parent/key] [|value] [|template-filename] [|secret-length]")
		fmt.Println("    if value, template-filename or secret length is needed, but missing, it is read from stdin")
		flag.PrintDefaults()
		os.Exit(0)
	}
	// Create the KMS client.
	kmsService, err := cloudkms.New(kiya.NewAuthenticatedClient(*oAuthLocation))
	if err != nil {
		log.Fatal(err)
	}
	// Create the Bucket client
	storageService, err := cloudstore.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create client [%v]", err)
	}
	profileName := flag.Arg(0)
	target, ok := kiya.Profiles[profileName]
	if !ok {
		log.Fatalf("no such profile [%s] please check your .kiya file", profileName)
	}
	// what command?
	switch flag.Arg(1) {

	case "put":
		key := flag.Arg(2)
		value := flag.Arg(3)
		if len(value) != 0 {
			commandPutPasteGenerate(kmsService, storageService, target, "put", key, value, doPrompt)
		} else {
			value = readFromStdIn()
			commandPutPasteGenerate(kmsService, storageService, target, "put", key, value, doNotPrompt)
		}

	case "paste":
		key := flag.Arg(2)
		value, err := clipboard.ReadAll()

		if err != nil {
			log.Fatal(tre.New(err, "clipboard read failed", "key", key))
		}
		commandPutPasteGenerate(kmsService, storageService, target, "paste", key, value, doPrompt)

	case "generate":
		key := flag.Arg(2)
		value := flag.Arg(3)

		var length string
		var mustPrompt bool
		if len(value) != 0 {
			length = value
			mustPrompt = true
		} else {
			length = readFromStdIn()
			mustPrompt = false
		}

		secretLength, err := strconv.Atoi(length)
		if err != nil {
			log.Fatal(tre.New(err, "generate failed", "key", key, "err", err))
		}
		secret, err := kiya.GenerateSecret(secretLength, target.SecretRunes)
		if err != nil {
			log.Fatal(tre.New(err, "generate failed", "key", key, "err", err))
		}
		commandPutPasteGenerate(kmsService, storageService, target, "generate", key, secret, mustPrompt)
		// make it available on the clipboard, ignore error
		clipboard.WriteAll(secret)

	case "copy":
		key := flag.Arg(2)
		value, err := kiya.GetValueByKey(kmsService, storageService, key, target)
		if err != nil {
			log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
		}
		if err := clipboard.WriteAll(value); err != nil {
			log.Fatal(tre.New(err, "copy failed", "key", key, "err", err))
		}

	case "get":
		key := flag.Arg(2)
		value, err := kiya.GetValueByKey(kmsService, storageService, key, target)
		if err != nil {
			log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
		}
		if len(*oOutputFilename) > 0 {
			if err := ioutil.WriteFile(*oOutputFilename, []byte(value), os.ModePerm); err != nil {
				log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
			}
			return
		}
		fmt.Println(value)

	case "delete":
		key := flag.Arg(2)
		commandDelete(kmsService, storageService, target, key)
	case "list":
		// kiya [profile] list [|filter-term]
		filter := flag.Arg(2)
		commandList(storageService, target, filter)
	case "template":
		commandTemplate(kmsService, storageService, target, *oOutputFilename)
	case "move":
		// kiya [source] move [source-key] [target] [|target-key]
		sourceProfile := kiya.Profiles[flag.Arg(0)]
		sourceKey := flag.Arg(2)
		targetProfile := kiya.Profiles[flag.Arg(3)]
		targetKey := sourceKey
		if len(flag.Args()) == 5 {
			targetKey = flag.Arg(4)
		}
		commandMove(kmsService, storageService, sourceProfile, sourceKey, targetProfile, targetKey)
	default:
		commandList(storageService, target, flag.Arg(1))
	}
}
