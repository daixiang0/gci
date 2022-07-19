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
	_ = gciCfg
	assert.Equal(t, section.DefaultSections(), gciCfg.Sections)
	assert.Equal(t, section.DefaultSectionSeparators(), gciCfg.SectionSeparators)
	assert.False(t, gciCfg.Debug)
	assert.False(t, gciCfg.NoInlineComments)
	assert.False(t, gciCfg.NoPrefixComments)
}

func TestInitGciConfigFromYAML(t *testing.T) {
	gciCfg, err := InitializeGciConfigFromYAML(path.Join(testFilesPath, "configTest.cfg.yaml"))
	assert.NoError(t, err)
	_ = gciCfg
	assert.Equal(t, section.SectionList{section.Default{}}, gciCfg.Sections)
	assert.False(t, gciCfg.Debug)
	assert.True(t, gciCfg.SkipGenerated)
}
