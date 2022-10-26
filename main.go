package main

import (
	"os"

	"github.com/daixiang0/gci/cmd/gci"
)

var version = "0.8.2"

func main() {
	e := gci.NewExecutor(version)

	err := e.Execute()
	if err != nil {
		os.Exit(1)
	}
}
