package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daixiang0/gci/pkg/section"
)

// the custom sections sort alphabetically as default.
func TestParseOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/daixiang0/gci)", "prefix(github/daixiang0/gai)"},
	}
	gciCfg, err := cfg.Parse()
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}, section.Custom{Prefix: "github/daixiang0/gai"}, section.Custom{Prefix: "github/daixiang0/gci"}}, gciCfg.Sections)
}

func TestParseCustomOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/daixiang0/gci)", "prefix(github/daixiang0/gai)"},
		Cfg: BoolConfig{
			CustomOrder: true,
		},
	}
	gciCfg, err := cfg.Parse()
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}, section.Custom{Prefix: "github/daixiang0/gci"}, section.Custom{Prefix: "github/daixiang0/gai"}}, gciCfg.Sections)
}

func TestParseNoLexOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"prefix(github/daixiang0/gci)", "prefix(github/daixiang0/gai)", "default"},
		Cfg: BoolConfig{
			NoLexOrder: true,
		},
	}

	gciCfg, err := cfg.Parse()
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}, section.Custom{Prefix: "github/daixiang0/gci"}, section.Custom{Prefix: "github/daixiang0/gai"}}, gciCfg.Sections)
}
