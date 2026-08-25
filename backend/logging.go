package backend

import (
	"fmt"
	"log/slog"
	"os"
)

func fatal(args ...any) {
	slog.Error(fmt.Sprint(args...))
	os.Exit(1)
}
