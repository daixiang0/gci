package config

import (
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/daixiang0/gci/v2/pkg/section"
)

var defaultOrder = map[string]int{
	section.StandardType:    0,
	section.DefaultType:     1,
	section.CustomType:      2,
	section.BlankType:       3,
	section.DotType:         4,
	section.AliasType:       5,
	section.LocalModuleType: 6,
}

type BoolConfig struct {
	NoInlineComments bool `yaml:"no-inlineComments"`
	NoPrefixComments bool `yaml:"no-prefixComments"`
	Debug            bool `yaml:"-"`
	SkipGenerated    bool `yaml:"skipGenerated"`
	SkipVendor       bool `yaml:"skipVendor"`
	CustomOrder      bool `yaml:"customOrder"`
	NoLexOrder       bool `yaml:"noLexOrder"`
}

type Config struct {
	BoolConfig
	Sections          section.SectionList
	SectionSeparators section.SectionList
}

type YamlConfig struct {
	Cfg                     BoolConfig `yaml:",inline"`
	SectionStrings          []string   `yaml:"sections"`
	SectionSeparatorStrings []string   `yaml:"sectionseparators"`

	ModPath string `yaml:"-"`
}

func (g YamlConfig) Parse() (*Config, error) {
	var err error

	sections, err := section.Parse(g.SectionStrings)
	if err != nil {
		return nil, err
	}
	if sections == nil {
		sections = section.DefaultSections()
	}
	if err := configureSections(sections, g.ModPath); err != nil {
		return nil, err
	}

	if !g.Cfg.CustomOrder {
		sort.Slice(sections, func(i, j int) bool {
			sectionI, sectionJ := sections[i].Type(), sections[j].Type()

			if g.Cfg.NoLexOrder || sectionI != sectionJ {
				return defaultOrder[sectionI] < defaultOrder[sectionJ]
			}

			return sections[i].String() < sections[j].String()
		})
	}

	sectionSeparators, err := section.Parse(g.SectionSeparatorStrings)
	if err != nil {
		return nil, err
	}
	if sectionSeparators == nil {
		sectionSeparators = section.DefaultSectionSeparators()
	}

	return &Config{g.Cfg, sections, sectionSeparators}, nil
}

func ParseConfig(in string) (*Config, error) {
	config := YamlConfig{}

	err := yaml.Unmarshal([]byte(in), &config)
	if err != nil {
		return nil, err
	}

	gciCfg, err := config.Parse()
	if err != nil {
		return nil, err
	}

	return gciCfg, nil
}

func configureSections(sections section.SectionList, path string) error {
	for _, sec := range sections {
		switch s := sec.(type) {
		case *section.LocalModule:
			if err := s.Configure(path); err != nil {
				return err
			}
		}
	}
	return nil
}
