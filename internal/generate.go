package main

import (
	"bytes"
	"go/format"
	"os"
	"runtime"
	"strings"
	"text/template"

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

func generate() error {
	all, err := packages.Load(nil, "std")
	if err != nil {
		return err
	}

	var pkgs []string

	// go list std | grep -v vendor | grep -v internal
	for _, pkg := range all {
		if !strings.Contains(pkg.PkgPath, "internal") && !strings.Contains(pkg.PkgPath, "vendor") {
			pkgs = append(pkgs, pkg.PkgPath)
		}
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

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
