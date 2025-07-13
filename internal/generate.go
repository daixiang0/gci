package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"
	"text/template"

	"golang.org/x/sync/errgroup"
	"golang.org/x/tools/go/packages"
)

//go:generate go run .

const outputFile = "../pkg/section/standard_list.go"

const stdTemplate = `
package section

// Code generated based on {{ .Version }}. DO NOT EDIT.

var standardPackages = map[string]struct{}{
{{- range $pkg := .Packages }}
		"{{ $pkg }}":  {},
{{- end}}
}

`

func main() {
	err := generate()
	if err != nil {
		panic(err)
	}
}

// update from https://go.dev/doc/install/source#environment
var list = `aix	ppc64
android	386
android	amd64
android	arm
android	arm64
darwin	amd64
darwin	arm64
dragonfly	amd64
freebsd	386
freebsd	amd64
freebsd	arm
illumos	amd64
ios	arm64
js	wasm
linux	386
linux	amd64
linux	arm
linux	arm64
linux	loong64
linux	mips
linux	mipsle
linux	mips64
linux	mips64le
linux	ppc64
linux	ppc64le
linux	riscv64
linux	s390x
netbsd	386
netbsd	amd64
netbsd	arm
openbsd	386
openbsd	amd64
openbsd	arm
openbsd	arm64
plan9	386
plan9	amd64
plan9	arm
solaris	amd64
wasip1	wasm
windows	386
windows	amd64
windows	arm
windows	arm64`

func generate() error {
	var all []*packages.Package

	writeLock := sync.Mutex{}

	g, _ := errgroup.WithContext(context.Background())
	for _, pair := range strings.Split(list, "\n") {
		pair := pair
		g.Go(func() error {
			goos, goarch, found := strings.Cut(pair, "\t")
			if !found {
				return nil
			}

			pkgs, err := packages.Load(&packages.Config{
				Mode: packages.NeedName,
				Env:  append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch, "GOEXPERIMENT=arenas,boringcrypto,synctest,jsonv2"),
			}, "std")
			if err != nil {
				return err
			}
			fmt.Println("loaded", goos, goarch, len(pkgs))

			writeLock.Lock()
			defer writeLock.Unlock()

			all = append(all, pkgs...)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	uniquePkgs := make(map[string]struct{})

	// go list std | grep -v vendor | grep -v internal
	for _, pkg := range all {
		if !strings.Contains(pkg.PkgPath, "internal") && !strings.Contains(pkg.PkgPath, "vendor") {
			uniquePkgs[pkg.PkgPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(uniquePkgs))
	for pkg := range uniquePkgs {
		pkgs = append(pkgs, pkg)
	}

	slices.Sort(pkgs)

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	models := map[string]interface{}{
		"Packages": pkgs,
		"Version":  runtime.Version(),
	}

	tlt, err := template.New("std-packages").Parse(stdTemplate)
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}

	err = tlt.Execute(b, models)
	if err != nil {
		return err
	}

	// gofmt
	source, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	_, err = file.Write(source)
	if err != nil {
		return err
	}

	return nil
}
