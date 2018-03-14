package main

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"testing"
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

func TestURLEncodingFreeCharset(t *testing.T) {
	v := url.Values{}
	v.Set("test", defaultSecreteCharSet)
	if got, want := v.Encode(), "test="+defaultSecreteCharSet; got != want {
		t.Error("got [%s] want [%s]", got, want)
	}
}

func TestGenerateSecretDefaultChars(t *testing.T) {
	for _, each := range []int{0, 1, 50, 100000} {
		testGenerateSecret(t, each, "", defaultSecreteCharSet)
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
