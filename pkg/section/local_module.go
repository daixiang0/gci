package section

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/daixiang0/gci/pkg/parse"
	"github.com/daixiang0/gci/pkg/specificity"
)

const LocalModuleType = "localmodule"

type LocalModule struct {
	Paths []string
}

func (m *LocalModule) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	for _, modPath := range m.Paths {
		// also check path etc.
		if strings.HasPrefix(spec.Path, modPath) {
			return specificity.LocalModule{}
		}
	}

	return specificity.MisMatch{}
}

func (m *LocalModule) String() string {
	return LocalModuleType
}

func (m *LocalModule) Type() string {
	return LocalModuleType
}

// Configure configures the module section by finding the module
// for the current path
func (m *LocalModule) Configure() error {
	modPaths, err := findLocalModules()
	if err != nil {
		return err
	}
	m.Paths = modPaths
	return nil
}

func findLocalModules() ([]string, error) {
	packages, err := packages.Load(
		// find the package in the current dir and load its module
		// NeedFiles so there is some more info in package errors
		&packages.Config{Mode: packages.NeedModule | packages.NeedFiles},
		".",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load local modules: %v", err)
	}

	uniqueModules := make(map[string]struct{})

	for _, pkg := range packages {
		if len(pkg.Errors) != 0 {
			return nil, fmt.Errorf("error reading local packages: %v", pkg.Errors)
		}
		if pkg.Module != nil {
			uniqueModules[pkg.Module.Path] = struct{}{}
		}
	}

	modPaths := make([]string, 0, len(uniqueModules))
	for mod := range uniqueModules {
		modPaths = append(modPaths, mod)
	}

	return modPaths, nil
}
