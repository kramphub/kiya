package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kramphub/kiya/backend"
	"golang.org/x/term"
)

func readFromStdIn() string {
	buffer, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error while reading from standard in", err)
	}

	// remove newline added to std in from command execution
	if buffer[len(buffer)-1] == '\n' {
		buffer = buffer[:len(buffer)-1]
	}

	return string(buffer)
}

// PromptForYes prompts for a yes or no in a CMD environment.
func promptForYes(message string) bool {

	// Don't prompt for confirmation if the quiet flag is enabled
	if *oQuiet {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	yn, _ := reader.ReadString('\n')
	return strings.HasPrefix(yn, "Y") || strings.HasPrefix(yn, "y")
}

func shouldPromptForPassword(b backend.Backend) bool {
	switch b.(type) {
	case *backend.FileStore:
		return true
	default:
		return false
	}
}

func promptForPassword() []byte {
	log.Print("[INFO]: Make sure you use a secure and strong master password.")

	fmt.Println("Enter master password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		log.Fatal("Error while reading password from standard in", err)
	}

	if len(password) == 0 {
		log.Fatal("Password should have at least one character.")
	}
	return password
}

// encodeToJson encodes the given object to JSON.
func encodeToJson(v interface{}) []byte {
	buf, err := json.Marshal(v)

	if err != nil {
		log.Fatalf("[FATAL] encode struct to JSON failed: %s", err.Error())
	}

	return buf
}

// decodeJson decodes the given JSON to the given object.
func decodeJson[T interface{}](data []byte) T {
	var obj T
	err := json.Unmarshal(data, &obj)

	if err != nil {
		log.Fatalf("[FATAL] decode JSON failed: %s", err.Error())
	}

	return obj
}
