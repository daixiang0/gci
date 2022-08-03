package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daixiang0/gci/pkg/section"
)

var testFilesPath = "../gci/internal/testdata"

func TestInitGciConfigFromEmptyYAML(t *testing.T) {
	gciCfg, err := InitializeGciConfigFromYAML(path.Join(testFilesPath, "defaultValues.cfg.yaml"))
	assert.NoError(t, err)
	assert.Equal(t, section.DefaultSections(), gciCfg.Sections)
	assert.Equal(t, section.DefaultSectionSeparators(), gciCfg.SectionSeparators)
	assert.False(t, gciCfg.Debug)
	assert.False(t, gciCfg.NoInlineComments)
	assert.False(t, gciCfg.NoPrefixComments)
}

func TestInitGciConfigFromYAML(t *testing.T) {
	gciCfg, err := InitializeGciConfigFromYAML(path.Join(testFilesPath, "configTest.cfg.yaml"))
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}}, gciCfg.Sections)
	assert.False(t, gciCfg.Debug)
	assert.True(t, gciCfg.SkipGenerated)
	assert.False(t, gciCfg.CustomOrder)
}

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
