//go:build linux && darwin

package main

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/term"
)

func promptForPassword() []byte {

	log.Print("[INFO]: Make sure you use a secure and strong master password.")

	fmt.Println("Enter master password: ")
	password, err := term.ReadPassword(syscall.Stdin)

	if err != nil {
		log.Fatal("Error while reading password from standard in", err)
	}

	if len(password) == 0 {
		log.Fatal("Password should have at least one character.")
	}
	return password
}
