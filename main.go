package main

import (
	"os"

	"github.com/daixiang0/gci/cmd/gci"
)

var Version string

func main() {
	e := gci.NewExecutor(Version)

	err := e.Execute()
	if err != nil {
		os.Exit(1)
	}
}
