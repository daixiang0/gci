package gci

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/daixiang0/gci/pkg/gci/sections"
	"github.com/daixiang0/gci/pkg/io"

	"github.com/stretchr/testify/assert"
)

var testFilesPath = "internal/testdata"

func isTestInputFile(file os.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), ".in.go")
}

func TestRun(t *testing.T) {
	testFiles, err := io.FindFilesForPath(testFilesPath, isTestInputFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, testFile := range testFiles {
		fileBaseName := strings.TrimSuffix(testFile, ".in.go")
		t.Run(fileBaseName, func(t *testing.T) {
			t.Parallel()

			gciCfg, err := initializeGciConfigFromYAML(fileBaseName + ".cfg.yaml")
			if err != nil {
				t.Fatal(err)
			}

			_, formattedFile, err := LoadFormatGoFile(io.File{fileBaseName + ".in.go"}, *gciCfg)
			if err != nil {
				t.Fatal(err)
			}
			expectedOutput, err := ioutil.ReadFile(fileBaseName + ".out.go")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, string(expectedOutput), string(formattedFile), "output")
			assert.NoError(t, err)
		})
	}
}

func TestInitGciConfigFromEmptyYAML(t *testing.T) {
	gciCfg, err := initializeGciConfigFromYAML(path.Join(testFilesPath, "defaultValues.cfg.yaml"))
	assert.NoError(t, err)
	_ = gciCfg
	assert.Equal(t, DefaultSections(), gciCfg.Sections)
	assert.Equal(t, DefaultSectionSeparators(), gciCfg.SectionSeparators)
	assert.False(t, gciCfg.Debug)
	assert.False(t, gciCfg.NoInlineComments)
	assert.False(t, gciCfg.NoPrefixComments)
}

func TestInitGciConfigFromYAML(t *testing.T) {
	gciCfg, err := initializeGciConfigFromYAML(path.Join(testFilesPath, "configTest.cfg.yaml"))
	assert.NoError(t, err)
	_ = gciCfg
	assert.Equal(t, SectionList{sections.DefaultSection{}}, gciCfg.Sections)
	assert.Equal(t, SectionList{sections.CommentLine{"---"}}, gciCfg.SectionSeparators)
	assert.False(t, gciCfg.Debug)
	assert.True(t, gciCfg.NoInlineComments)
	assert.True(t, gciCfg.NoPrefixComments)
}
