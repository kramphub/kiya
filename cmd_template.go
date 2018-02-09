package main

import (
	"encoding/base64"
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

func commandTemplate(kmsService *cloudkms.Service, storageService *cloudstore.Client, target profile) {
	funcMap := template.FuncMap{
		"kiya": templateFunction(kmsService, storageService, target),
		"base64": func(value string) string {
			return base64.StdEncoding.EncodeToString([]byte(value))
		},
	}
	processor := template.New("base").Funcs(funcMap)
	templateName := "base"

	filename := flag.Arg(2)
	if len(filename) > 0 {
		t, err := processor.ParseFiles(filename)
		if err != nil {
			log.Fatal(tre.New(err, "templating failed", "filename", filename))
		}
		processor = t
		templateName = filepath.Base(filename)
	} else {
		templateContent := readFromStdIn()
		t, err := processor.Parse(templateContent)
		if err != nil {
			log.Fatal("templating failed", err)
		}
		processor = t
	}
	writer := os.Stdout
	// change writer is output was specified
	if len(*oOutputFilename) > 0 {
		out, err := os.Create(*oOutputFilename)
		if err != nil {
			log.Fatal("unable to create output", err)
		}
		writer = out
	}
	defer writer.Close()
	processor.ExecuteTemplate(writer, templateName, "")
}

func templateFunction(kmsService *cloudkms.Service, storageService *cloudstore.Client, target profile) func(string) string {
	return func(key string) string {
		value, err := getValueByKey(kmsService, storageService, key, target)
		if err != nil {
			log.Fatal(tre.New(err, "templating failed", "key", key))
			return ""
		}
		return value
	}
}
