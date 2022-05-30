package gci

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/daixiang0/gci/pkg/configuration"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"
)

type GciConfiguration struct {
	configuration.FormatterConfiguration
	Sections          SectionList
	SectionSeparators SectionList
}

type GciStringConfiguration struct {
	Cfg                     configuration.FormatterConfiguration `yaml:",inline"`
	SectionStrings          []string                             `yaml:"sections"`
	SectionSeparatorStrings []string                             `yaml:"sectionseparators"`
}

func (g GciStringConfiguration) Parse() (*GciConfiguration, error) {
	sections := DefaultSections()
	var err error
	if len(g.SectionStrings) > 0 {
		sections, err = sectionsPkg.SectionParserInst.ParseSectionStrings(g.SectionStrings, true, true)
		if err != nil {
			return nil, err
		}
	}
	sectionSeparators := DefaultSectionSeparators()
	if len(g.SectionSeparatorStrings) > 0 {
		sectionSeparators, err = sectionsPkg.SectionParserInst.ParseSectionStrings(g.SectionSeparatorStrings, false, false)
		if err != nil {
			return nil, err
		}
	}
	return &GciConfiguration{g.Cfg, sections, sectionSeparators}, nil
}

// InitializeModules collects and remembers Go module names for the given
// files, by traversing the file system.
//
// This method requires that g.Sections contains the Module section,
// otherwise InitializeModules does nothing. This also implies that
// this method should be called after changes to g.Sections, for example
// right after (*GciStringConfiguration).Parse().
func (g *GciConfiguration) InitializeModules(files []string) error {
	var moduleSection *sectionsPkg.Module
	for _, section := range g.Sections {
		if m, ok := section.(sectionsPkg.Module); ok {
			moduleSection = &m
			break
		}
	}
	if moduleSection == nil {
		// skip collecting Go modules when not needed
		return nil
	}

	resolver := make(moduleResolver)
	knownModulePaths := map[string]struct{}{} // unique list of Go modules
	for _, file := range files {
		path, err := resolver.Lookup(file)
		if err != nil {
			return err
		}
		if path != "" {
			knownModulePaths[path] = struct{}{}
		}
	}
	modulePaths := make([]string, 0, len(knownModulePaths))
	for path := range knownModulePaths {
		modulePaths = append(modulePaths, path)
	}
	moduleSection.SetModulePaths(modulePaths)
	return nil
}

func initializeGciConfigFromYAML(filePath string) (*GciConfiguration, error) {
	yamlCfg := GciStringConfiguration{}
	yamlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlData, &yamlCfg)
	if err != nil {
		return nil, err
	}
	gciCfg, err := yamlCfg.Parse()
	if err != nil {
		return nil, err
	}
	return gciCfg, nil
}
