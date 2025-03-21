package main

import (
	"github.com/daixiang0/gci/cmd/gci"
	"os"
	"strings"
)

var Version = "0.0.0"

func main() {
	e := gci.NewExecutor(strings.TrimPrefix(Version, "v"))

	err := e.Execute()
	if err != nil {
		os.Exit(1)
	}
}
