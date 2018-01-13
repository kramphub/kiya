package main

import (
	"testing"
	"regexp"
	"fmt"
	"bytes"
	"unicode/utf8"
)

func TestGenerateSecret(t *testing.T) {
	var testCases = []struct {
		length      int
		secretChars string
	}{
		{0, "a"},
		{1, "a"},
		{10, "a"},
		{25, "ABCabc123\\+="},
		{25, "☠☑♨⛱"},
		{25, "   "},
		{25, "  + "},
		{25, "¯|_(ツ)_/¯"},
	}
	for _, tc := range testCases {
		testGenerateSecret(t, tc.length, tc.secretChars, tc.secretChars)
	}
}

func TestGenerateSecretDefaultChars(t *testing.T) {
	var testCases = []struct {
		length int
	}{
		{0},
		{1},
		{50},
		{100000},
	}
	for _, tc := range testCases {
		testGenerateSecret(t, tc.length, "", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-={}|[]\\:\"<>?,./")
	}
}

func testGenerateSecret(t *testing.T, length int, charsIn string, expectedChars string) {
	pattern, err := matchCharsRegex(expectedChars)
	got, err := GenerateSecret(length, []rune(charsIn))
	if err != nil {
		t.Errorf("want random string; got [%s] -> GenerateSecret(%d, %s)", err, length, charsIn)
	}
	if utf8.RuneCountInString(got) != length || !pattern.MatchString(got) {
		t.Errorf("want random string with {length: %d, chars: %s}\ngot %s", length, expectedChars, got)
	}
}

// create a regex pattern matching any random combination of the characters in the chars string
func matchCharsRegex(chars string) (*regexp.Regexp, error) {
	var regex bytes.Buffer
	regex.WriteString("^[")
	for _, char := range chars {
		regex.WriteString("\\x{" + fmt.Sprintf("%X", char) + "}")
	}
	regex.WriteString("]*$")
	return regexp.Compile(regex.String())
}
