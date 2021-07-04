package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/emicklei/tre"
	"github.com/kramphub/kiya/backend"
)

func commandTemplate(ctx context.Context, b backend.Backend, target *backend.Profile, outputFilename string) {
	funcMap := template.FuncMap{
		"kiya": templateFunction(ctx, b, target),
		"base64": func(value string) string {
			return base64.StdEncoding.EncodeToString([]byte(value))
		},
		"env": func(value string) string {
			return os.Getenv(value)
		},
	}
	processor := template.New("base").Funcs(funcMap)
	templateName := "base"

	filename := flag.Arg(2)
	if len(filename) > 0 {
		t, err := processor.ParseFiles(filename)
		if err != nil {
			wd, _ := os.Getwd()
			log.Fatal(tre.New(err, "templating failed", "filename", filename, "current workdirectory", wd))
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
	if len(outputFilename) > 0 {
		out, err := os.Create(outputFilename)
		if err != nil {
			log.Fatal("unable to create output", err)
		}
		writer = out
	}
	defer writer.Close()
	processor.ExecuteTemplate(writer, templateName, "")
}

func templateFunction(ctx context.Context, b backend.Backend, target *backend.Profile) func(string) string {
	return func(key string) string {
		value, err := b.Get(ctx, target, key, "latest")
		if err != nil {
			log.Fatal(tre.New(err, "templating failed", "key", key))
			return ""
		}
		return string(value)
	}
}
