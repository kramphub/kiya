//go:build windows

package main

import (
	"fmt"
	"log"

	"golang.org/x/term"
)

// See also: https://github.com/golang/go/issues/11914
func promptForPassword() []byte {

	log.Print("[INFO]: Make sure you use a secure and strong master password.")

	fmt.Println("Enter master password: ")
	password, err := term.ReadPassword(0)

	if err != nil {
		log.Fatal("Error while reading password from standard in", err)
	}

	if len(password) == 0 {
		log.Fatal("Password should have at least one character.")
	}
	return password
}
