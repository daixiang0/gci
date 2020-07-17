package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)



var (
	write  = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff = flag.Bool("d", false, "display diffs instead of rewriting files")

	localFlag string

	exitCode = 0
)

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func parseFlags() []string {
	flag.StringVar(&localFlag, "local", "", "put imports beginning with this string after 3rd-party packages, only support one string")

	flag.Parse()
	return flag.Args()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gci [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func processFile(filename string, out io.Writer) error {
	var err error

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	ori := make([]byte, len(src))
	copy(ori, src)
	start := bytes.Index(src, importStartFlag)
	// in case no importStartFlag or importStartFlag exist in the commentFlag
	if start < 0 {
		fmt.Printf("skip file %s since no import\n", filename)
		return nil
	}
	end := bytes.Index(src[start:], importEndFlag) + start

	ret := bytes.Split(src[start+len(importStartFlag):end], []byte(linebreak))

	p := newPkg(ret, localFlag)

	res := append(src[:start+len(importStartFlag)], append(p.fmt(), src[end+1:]...)...)

	if !bytes.Equal(ori, res) {
		exitCode = 1

		if *write {
			// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
			var perms os.FileMode
			if fi, err := os.Stat(filename); err == nil {
				perms = fi.Mode() & os.ModePerm
			}
			err = ioutil.WriteFile(filename, res, perms)
			if err != nil {
				return err
			}
		}
		if *doDiff {
			data, err := diff(ori, res, filename)
			if err != nil {
				return fmt.Errorf("failed to diff: %v", err)
			}
			fmt.Printf("diff -u %s %s\n", filepath.ToSlash(filename+".orig"), filepath.ToSlash(filename))
			if _, err := out.Write(data); err != nil {
				return fmt.Errorf("failed to write: %v", err)
			}
		}
	}
	if !*write && !*doDiff {
		if _, err = out.Write(res); err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}
	}

	return err
}

func main() {
	flag.Usage = usage
	paths := parseFlags()
	for _, path := range paths {
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, os.Stdout); err != nil {
				report(err)
			}
		}
	}
	os.Exit(exitCode)
}
