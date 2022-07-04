package sections

import (
	"strings"

	"github.com/daixiang0/gci/pkg/configuration"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func init() {
	prefixType := SectionType{
		generatorFun: func(parameter string, sectionPrefix, sectionSuffix Section) (Section, error) {
			return Module{}, nil
		},
		aliases:     []string{"Module", "Mod"},
		description: "Groups all imports of the corresponding Go module",
	}.StandAloneSection().WithoutParameter()
	SectionParserInst.registerSectionWithoutErr(&prefixType)
}

type Module struct {
	// modulePaths contains all known Go module path names.
	//
	// This must be a pointer, because gci.formatImportBlock() will create
	// mapping between sections and imports, and slices are unhashable.
	modulePaths *[]string
}

func (m Module) MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity {
	if m.modulePaths == nil {
		return specificity.MisMatch{}
	}

	importPath := spec.Path()
	for _, path := range *m.modulePaths {
		if strings.HasPrefix(importPath, path) {
			return specificity.Module{}
		}
	}
	return specificity.MisMatch{}
}

func (m Module) Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	return inorderSectionFormat(m, imports, cfg)
}

func (Module) sectionPrefix() Section { return nil }
func (Module) sectionSuffix() Section { return nil }

func (Module) String() string {
	return "Module"
}

func (m *Module) SetModulePaths(paths []string) {
	dup := make([]string, len(paths), len(paths))
	copy(dup, paths)

	m.modulePaths = &dup
}
