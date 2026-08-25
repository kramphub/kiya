package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func fatal(args ...any) {
	slog.Error(fmt.Sprint(args...))
	os.Exit(1)
}

func fatalf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func fatalln(args ...any) {
	slog.Error(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
	os.Exit(1)
}
