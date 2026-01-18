package main

import (
	"os"

	"github.com/blacksilver/termplate-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
