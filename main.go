package main

import (
	"os"

	"github.com/shportix/blob-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
