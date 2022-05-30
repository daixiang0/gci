package gci

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/daixiang0/gci/pkg/gci/sections"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.InitLogger()
	defer log.L().Sync()
}

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

func TestSkippingOverIncorrectlyFormattedFiles(t *testing.T) {
	cfg, err := GciStringConfiguration{}.Parse()
	assert.NoError(t, err)

	var importUnclosedCtr, noImportCtr, validCtr int
	var files []io.FileObj
	files = append(files, TestFile{io.File{"internal/skipTest/import-unclosed.testgo"}, &importUnclosedCtr})
	files = append(files, TestFile{io.File{"internal/skipTest/no-import.testgo"}, &noImportCtr})
	files = append(files, TestFile{io.File{"internal/skipTest/valid.testgo"}, &validCtr})

	validFileProcessedChan := make(chan bool, len(files))

	generatorFunc := func() ([]io.FileObj, error) {
		return files, nil
	}
	fileAccessTestFunc := func(filePath string, unmodifiedFile, formattedFile []byte) error {
		validFileProcessedChan <- true
		return nil
	}
	err = processFiles(generatorFunc, *cfg, fileAccessTestFunc)

	assert.NoError(t, err)
	// check all files have been accessed
	assert.Equal(t, importUnclosedCtr, 1)
	assert.Equal(t, noImportCtr, 1)
	assert.Equal(t, validCtr, 1)
	// check that processing for the valid file was called
	assert.True(t, <-validFileProcessedChan)
}

type TestFile struct {
	wrappedFile   io.File
	accessCounter *int
}

func (t TestFile) Load() ([]byte, error) {
	*t.accessCounter++
	return t.wrappedFile.Load()
}

func (t TestFile) Path() string {
	return t.wrappedFile.Path()
}
