package main

import (
	"flag"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	cloudstore "cloud.google.com/go/storage"
	"github.com/atotto/clipboard"
	"github.com/emicklei/tre"
	"github.com/kramphub/kiya"
	"github.com/kramphub/kiya/backend"
	"golang.org/x/net/context"
	"google.golang.org/api/cloudkms/v1"
)

var version = "build-" + time.Now().String()

const (
	doPrompt    = true
	doNotPrompt = false
)

func main() {
	ctx := context.Background()

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

	profileName := flag.Arg(0)
	target, ok := kiya.Profiles[profileName]
	if !ok {
		log.Fatalf("no such profile [%s] please check your .kiya file", profileName)
	}

	b, err := getBackend(ctx, &target, *oMasterPassword)
	if err != nil {
		log.Fatalf("failed to intialize the secret provider backend, %s", err.Error())
	}
	defer func() {
		if err := b.Close(); err != nil {
			log.Fatalf("failed to close the secret provider backend, %s", err.Error())
		}
	}()

	// what command?
	switch flag.Arg(1) {

	case "put":
		key := flag.Arg(2)
		value := flag.Arg(3)
		if len(value) != 0 {
			commandPutPasteGenerate(ctx, b, &target, "put", key, value, doPrompt)
		} else {
			value = readFromStdIn()
			commandPutPasteGenerate(ctx, b, &target, "put", key, value, doNotPrompt)
		}

	case "paste":
		key := flag.Arg(2)
		value, err := clipboard.ReadAll()

		if err != nil {
			log.Fatal(tre.New(err, "clipboard read failed", "key", key))
		}
		commandPutPasteGenerate(ctx, b, &target, "paste", key, value, doPrompt)

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
		commandPutPasteGenerate(ctx, b, &target, "generate", key, secret, mustPrompt)
		// make it available on the clipboard, ignore error
		clipboard.WriteAll(secret)

	case "copy":
		key := flag.Arg(2)
		value, err := b.Get(ctx, &target, key)
		if err != nil {
			log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
		}
		if err := clipboard.WriteAll(string(value)); err != nil {
			log.Fatal(tre.New(err, "copy failed", "key", key, "err", err))
		}

	case "get":
		key := flag.Arg(2)

		bytes, err := b.Get(ctx, &target, key)
		if err != nil {
			log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
		}

		if len(*oOutputFilename) > 0 {
			if err := ioutil.WriteFile(*oOutputFilename, bytes, os.ModePerm); err != nil {
				log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
			}
			return
		}

		fmt.Println(string(bytes))

	case "delete":
		key := flag.Arg(2)
		commandDelete(ctx, b, &target, key)
	case "list":
		// kiya [profile] list [|filter-term]
		filter := flag.Arg(2)
		commandList(ctx, b, &target, filter)
	case "template":
		commandTemplate(ctx, b, &target, *oOutputFilename)
	case "move":
		// kiya [source] move [source-key] [target] [|target-key]
		sourceProfile := kiya.Profiles[flag.Arg(0)]
		sourceKey := flag.Arg(2)
		targetProfile := kiya.Profiles[flag.Arg(3)]
		targetKey := sourceKey
		if len(flag.Args()) == 5 {
			targetKey = flag.Arg(4)
		}
		commandMove(ctx, b, &sourceProfile, sourceKey, &targetProfile, targetKey)
	default:
		commandList(ctx, b, &target, flag.Arg(1))
	}
}

func getBackend(ctx context.Context, p *backend.Profile, masterPassword string) (backend.Backend, error) {
	switch p.Backend {
	case "ssm":
		return backend.NewAWSParameterStore(ctx, p)
	case "gsm":
		// Create GSM client
		gsmClient, err := secretmanager.NewClient(ctx)
		if err != nil {
			log.Fatalf("failed to setup client: %v", err)
		}

		return backend.NewGSM(gsmClient), nil
	case "akv":
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			log.Fatal(err)
		}
		client, err := azsecrets.NewClient(p.VaultUrl, cred, nil)
		if err != nil {
			log.Fatal(err)
		}
		return backend.NewAKV(client), nil
	case "file":
		return backend.NewFileStore(p.Location, p.ProjectID, masterPassword), nil
	case "kms":
		fallthrough
	default:
		// Create the KMS client
		kmsService, err := cloudkms.New(kiya.NewAuthenticatedClient(*oAuthLocation))
		if err != nil {
			log.Fatal(err)
		}
		// Create the Bucket client
		storageService, err := cloudstore.NewClient(ctx)
		if err != nil {
			log.Fatalf("failed to create client [%v]", err)
		}

		return backend.NewKMS(kmsService, storageService), nil
	}
}
