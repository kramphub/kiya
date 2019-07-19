package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"google.golang.org/api/cloudkms/v1"
	"math/big"
	"crypto/rand"
)

func getValueByKey(kmsService *cloudkms.Service, storageService *cloudstore.Client, key string, target profile) (string, error) {
	encryptedValue, err := loadSecret(storageService, target, key)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "key", key))
		return "", err
	}
	decryptedValue, err := getDecryptedValue(kmsService, target, encryptedValue)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "cipherText", encryptedValue))
		return "", err
	}

	return decryptedValue, nil
}

// readFromStdIn tries to read tries to read required input from standard in
func readFromStdIn() string {
	buffer, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error while reading from standard in", err)
	}

	// remove newline added to std in from command execution
	if buffer[len(buffer)-1] == '\n' {
		buffer = buffer[:len(buffer)-1]
	}

	return string(buffer)
}

func promptForYes(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	yn, _ := reader.ReadString('\n')
	return strings.HasPrefix(yn, "Y") || strings.HasPrefix(yn, "y")
}

func generate_secret(length int, chars string) (string, error) {
	if len(chars) == 0 {
		chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
	}

	randomString := make([]byte, length)
	for i := 0; i < length; i ++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		randomString[i] = chars[n.Int64()]
	}

	return string(randomString), nil
}
