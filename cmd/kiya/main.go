//nolint:gomnd
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	cloudstore "cloud.google.com/go/storage"
	"github.com/atotto/clipboard"
	"github.com/emicklei/tre"
	"github.com/kramphub/kiya"
	"github.com/kramphub/kiya/backend"
	"golang.org/x/net/context"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/option"
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

	b, err := getBackend(ctx, &target)
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

		if shouldPromptForPassword(b) {
			pass := promptForPassword()
			b.SetParameter("masterPassword", pass)
		}

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

		if shouldPromptForPassword(b) {
			pass := promptForPassword()
			b.SetParameter("masterPassword", pass)
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

		if shouldPromptForPassword(b) {
			pass := promptForPassword()
			b.SetParameter("masterPassword", pass)
		}

		commandPutPasteGenerate(ctx, b, &target, "generate", key, secret, mustPrompt)

		// make it available on the clipboard, ignore error
		err = clipboard.WriteAll(secret)
		if err != nil {
			log.Printf("[WARN] cannot copy public key to clipboard, %s", err.Error())
		}

	case "copy":
		key := flag.Arg(2)
		copySecretToClipboard(ctx, b, target, key)

	case "get":
		key := flag.Arg(2)

		if shouldPromptForPassword(b) {
			pass := promptForPassword()
			b.SetParameter("masterPassword", pass)
		}

		bytes, err := b.Get(ctx, &target, key)
		if err != nil {
			log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
		}

		if len(*oOutputFilename) > 0 {
			if err := os.WriteFile(*oOutputFilename, bytes, os.ModePerm); err != nil {
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
		listMatchingKeys(ctx, b, target, flag.Arg(2))

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

		if shouldPromptForPassword(b) {
			pass := promptForPassword()
			b.SetParameter("masterPassword", pass)
		}
		commandMove(ctx, b, &sourceProfile, sourceKey, &targetProfile, targetKey)

	case "backup":
		filter := flag.Arg(2)

		if *oBackupPath == "" {
			log.Fatalln("--path not specified")
		}

		fmt.Printf("Backup profile '%s', filter: '%s' to %s\n", profileName, filter, *oBackupPath)
		if *oEncryptBackup {
			fmt.Printf("Backup will be encrypted. Public key path: '%s', public key location: '%s'\n", *oBackupKey, *oBackupKeyStore)
		}

		if shouldPromptForPassword(b) {
			pass := promptForPassword()
			b.SetParameter("masterPassword", pass)
		}

		backup, err := commandBackup(ctx, b, target, filter)
		if err != nil {
			log.Fatalln(err.Error())
		}

		file, err := os.Create(*oBackupPath)
		if err != nil {
			log.Fatalf("create file '%s' failed, %s", *oBackupPath, err.Error())
		}

		if *oEncryptBackup {
			pub, err := getPublicKey(ctx, b, target, *oBackupKeyStore, *oBackupKey)
			if err != nil {
				log.Fatalf("[FATAL] get public key failed, %s", err.Error())
			}

			backup.Secret = generateSecret()

			buf, err := encrypt(backup.Data, backup.SecretAsBytes())
			if err != nil {
				log.Fatalf("[FATAL] encrypt items failed, %s", err.Error())
			}

			backup.Data = buf
			encryptedSecret, err := encryptSecret(backup.Secret, pub)
			if err != nil {
				log.Fatalf("[FATAL] encrypt secret failed, %s", err.Error())
			}
			backup.Encrypted = true
			backup.Secret = encryptedSecret
		}

		_, err = file.Write([]byte(backup.String()))

		if err != nil {
			log.Fatalf("save file '%s' failed, %s", *oBackupPath, err.Error())
		}
	case "restore":
		fmt.Printf("Restore profile '%s' from %s\n", profileName, *oBackupPath)

		buf, err := os.ReadFile(*oBackupPath)
		if err != nil {
			log.Fatalf("read '%s' failed, %s", *oBackupPath, err.Error())
		}

		backup := Backup{}
		backup.FromString(string(buf))
		var items map[string][]byte

		fmt.Printf("Backend '%s', restoring keys...\n", target.Backend)

		if backup.Encrypted || *oEncryptBackup {
			fmt.Println("Backup is encrypted.")

			buf, err := os.ReadFile(*oBackupKey)
			if err != nil {
				log.Fatalf("[FATAL] read private key '%s' failed, %s", *oBackupKey, err.Error())
			}

			privKey := exportPrivateKeyFromPEMString(buf)
			if err != nil {
				log.Fatalf("[FATAL] export private key '%s' failed, %s", *oBackupKey, err.Error())
			}

			secret, err := decryptSecret(backup.Secret, privKey)
			if err != nil {
				log.Fatalf("[FATAL] cannot decrypt secret, %s", err.Error())
			}

			buf, err = decrypt(backup.Data, secret)
			if err != nil {
				log.Fatalf("[FATAL] decrypt items failed, %s", err.Error())
			}

			fmt.Println("Backup decrypted, decode from JSON")
			items = decodeJson[map[string][]byte](buf)
		} else {
			items = decodeJson[map[string][]byte](backup.Data)
		}

		fmt.Printf("\rBackend '%s', restoring %d key(s)\n", target.Backend, len(items))

		if items == nil {
			log.Fatalln("no items found")
		}

		for k, v := range items {
			err := b.Put(ctx, &target, k, string(v), *oBackupRestoreOverwrite)
			if err != nil {
				log.Printf("[ERROR] put key '%s' failed - %s", k, err.Error())
			}
		}

	case "keygen":
		priv, pub, err := generateKeyPair()
		if err != nil {
			log.Fatal(err)
		}

		path := flag.Arg(2)
		if path == "" {
			path = "kiya_backupkey_rsa"
		}

		pubKeyStr := exportPublicKeyAsPEM(pub)
		privKeyStr := exportPrivateKeyAsPEM(priv)

		err = saveKeyToFile(pubKeyStr, fmt.Sprintf("%s_pub", path))
		if err != nil {
			log.Fatal(err)
		}

		err = saveKeyToFile(privKeyStr, path)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Key '%s', '%s_pub' saved\n", path, path)
		if err := clipboard.WriteAll(pubKeyStr); err != nil {
			log.Fatal(tre.New(err, "copy failed", err))
		}
		fmt.Println("Public key copied to clipboard")

	default:
		listMatchingKeys(ctx, b, target, flag.Arg(1))
	}
}

// getBackend returns a backend based on the profile
func getBackend(ctx context.Context, p *backend.Profile) (backend.Backend, error) {
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
			log.Fatalf("failed to create client [%v]", err)
		}
		return backend.NewAKV(client), nil
	case "file":
		return backend.NewFileStore(p.Location, p.ProjectID), nil
	case "kms":
		fallthrough
	default:
		// Create the KMS client
		kmsService, err := cloudkms.NewService(ctx, option.WithHTTPClient(kiya.NewAuthenticatedClient(*oAuthLocation)))
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

func copySecretToClipboard(ctx context.Context, be backend.Backend, target backend.Profile, key string) {
	if shouldPromptForPassword(be) {
		pass := promptForPassword()
		be.SetParameter("masterPassword", pass)
	}
	value, err := be.Get(ctx, &target, key)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "key", key, "err", err))
	}
	if err := clipboard.WriteAll(string(value)); err != nil {
		log.Fatal(tre.New(err, "copy failed", "key", key, "err", err))
	}
}

func listMatchingKeys(ctx context.Context, be backend.Backend, target backend.Profile, filter string) {
	keys := commandList(ctx, be, &target, filter)
	writeTable(keys, &target, filter)
	if len(keys) == 0 {
		return
	}
	// if there is only one match and AutoCopy is enabled
	// then copy the secret to clipboard
	if len(keys) == 1 && target.AutoCopyEnabled {
		copySecretToClipboard(ctx, be, target, keys[0].Name)
		fmt.Printf("... copied secret [%s] to clipboard.\n", keys[0].Name)
	}
	// more than one match
	if target.PromptForSecretLine {
		fmt.Println("Enter the number of the key to copy to clipboard")
		var n int
		_, err := fmt.Scan(&n)
		if err != nil {
			return
		}
		if n < 1 || n > len(keys) {
			fmt.Printf("No such line number %d\n", n)
			return
		}
		copySecretToClipboard(ctx, be, target, keys[n-1].Name)
		fmt.Printf("... copied secret [%s] to clipboard.\n", keys[n-1].Name)
	}
}
