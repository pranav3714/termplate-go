package main

import (
	"os"

	"github.com/blacksilver/ever-so-powerful/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
