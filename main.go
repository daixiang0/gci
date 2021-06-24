package main

import (
	"flag"
	"fmt"
	"go/scanner"
	"os"

	"github.com/daixiang0/gci/pkg/gci"
)

var (
	doWrite = flag.Bool("w", false, "doWrite result to (source) file instead of stdout")
	doDiff  = flag.Bool("d", false, "display diffs instead of rewriting files")

	localFlag []string

	exitCode = 0
)

func report(err error) {
	if err == nil {
		return
	}
	scanner.PrintError(os.Stderr, err)
	exitCode = 1
}

func parseFlags() []string {
	var localFlagStr string
	flag.StringVar(&localFlagStr, "local", "", "put imports beginning with this string after 3rd-party packages; comma-separated list")

	flag.Parse()
	localFlag = gci.ParseLocalFlag(localFlagStr)
	return flag.Args()
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: gci [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	paths := parseFlags()

	flagSet := &gci.FlagSet{
		LocalFlag: localFlag,
		DoWrite:   doWrite,
		DoDiff:    doDiff,
	}

	for _, path := range paths {
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			report(gci.WalkDir(path, flagSet))
		default:
			if err := gci.ProcessFile(path, os.Stdout, flagSet); err != nil {
				report(err)
			}
		}
	}
	os.Exit(exitCode)
}
