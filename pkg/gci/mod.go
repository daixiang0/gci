package gci

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// moduleResolver looksup the module path for a given (Go) file.
// To improve performance, the file paths and module paths are
// cached.
//
// Given the following directory structure:
//
//	/path/to/example
//	+-- go.mod               (module example)
//	+-- cmd/sample/main.go   (package main, imports example/util)
//	+-- util/util.go         (package util)
//
// After looking up main.go and util.go, the internal cache will contain:
//
//	"/path/to/foobar/": "example"
//
// For more complex module structures (i.e. sub-modules), the cache
// might look like this:
//
//	"/path/to/example/":            "example"
//	"/path/to/example/cmd/sample/": "go.example.com/historic/path"
//
// When matching files against this cache, the resolver will select the
// entry with the most specific path (so that, in this example, the file
// cmd/sample/main.go will resolve to go.example.com/historic/path).
type moduleResolver map[string]string

func (m moduleResolver) Lookup(file string) (string, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return "", fmt.Errorf("could not make path absolute: %w", err)
	}

	var bestMatch string
	for path := range m {
		if strings.HasPrefix(abs, path) && len(path) > len(bestMatch) {
			bestMatch = path
		}
	}
	if bestMatch != "" {
		return m[bestMatch], nil
	}

	return m.findRecursively(filepath.Dir(abs))
}

func (m moduleResolver) findRecursively(dir string) (string, error) {
	// When going up the directory tree, we might never find a go.mod
	// file. In this case remember where we started, so that the next
	// time we can short circuit the recursive ascent.
	stop := dir

	for {
		gomod := filepath.Join(dir, "go.mod")
		_, err := os.Stat(gomod)
		if errors.Is(err, os.ErrNotExist) {
			// go.mod doesn't exist at current location
			next := filepath.Dir(dir)
			if next == dir {
				// we're at the top of the filesystem
				m[stop] = ""
				return "", nil
			}
			// go one level up
			dir = next
			continue
		} else if err != nil {
			// other error (likely EPERM)
			return "", fmt.Errorf("module lookup failed: %w", err)
		}

		// we found a go.mod
		mod, err := ioutil.ReadFile(gomod)
		if err != nil {
			return "", fmt.Errorf("reading module failed: %w", err)
		}

		// store module path at m[dir]. add path separator to avoid
		// false-positive (think of /foo and /foobar).
		mpath := modfile.ModulePath(mod)
		if dir != "/" {
			// add trailing path sep, but not for *nix root directory
			dir += string(os.PathListSeparator)
		}
		m[dir] = mpath
		return mpath, nil
	}
}
