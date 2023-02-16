package gci

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
)

func init() {
	log.InitLogger()
	defer log.L().Sync()
}

var testFilesPath = "internal/testdata"

func isTestInputFile(_ string, file os.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), ".in.go")
}

func TestRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping test on Windows")
	}

	testFiles, err := io.FindFilesForPath(testFilesPath, isTestInputFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, testFile := range testFiles {
		fileBaseName := strings.TrimSuffix(testFile, ".in.go")
		t.Run("pkg/gci/"+testFile, func(t *testing.T) {
			t.Parallel()

			gciCfg, err := config.InitializeGciConfigFromYAML(fileBaseName + ".cfg.yaml")
			if err != nil {
				t.Fatal(err)
			}

			inputSrcFile := io.File{FilePath: fileBaseName + ".in.go"}
			inputSrc, err := inputSrcFile.Load()
			require.NoError(t, err)

			unmodifiedFile, formattedFile, err := LoadFormatGoFile(inputSrcFile, *gciCfg)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, inputSrc, unmodifiedFile)

			expectedOutput, err := os.ReadFile(fileBaseName + ".out.go")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, string(expectedOutput), string(formattedFile), "output")
			assert.NoError(t, err)
		})
	}
}

// func TestSkippingOverIncorrectlyFormattedFiles(t *testing.T) {
// 	cfg, err := config.YamlConfig{}.Parse()
// 	assert.NoError(t, err)

// 	var importUnclosedCtr, noImportCtr, validCtr int
// 	var files []io.FileObj
// 	files = append(files, TestFile{io.File{FilePath: "internal/skipTest/import-unclosed.testgo"}, &importUnclosedCtr})
// 	files = append(files, TestFile{io.File{FilePath: "internal/skipTest/no-import.testgo"}, &noImportCtr})
// 	files = append(files, TestFile{io.File{FilePath: "internal/skipTest/valid.testgo"}, &validCtr})

// 	validFileProcessedChan := make(chan bool, len(files))

// 	generatorFunc := func() ([]io.FileObj, error) {
// 		return files, nil
// 	}
// 	fileAccessTestFunc := func(filePath string, unmodifiedFile, formattedFile []byte) error {
// 		validFileProcessedChan <- true
// 		return nil
// 	}
// 	err = ProcessFiles(generatorFunc, *cfg, fileAccessTestFunc)

// 	assert.NoError(t, err)
// 	// check all files have been accessed
// 	assert.Equal(t, importUnclosedCtr, 1)
// 	assert.Equal(t, noImportCtr, 1)
// 	assert.Equal(t, validCtr, 1)
// 	// check that processing for the valid file was called
// 	assert.True(t, <-validFileProcessedChan)
// }

// type TestFile struct {
// 	wrappedFile   io.File
// 	accessCounter *int
// }

// func (t TestFile) Load() ([]byte, error) {
// 	*t.accessCounter++
// 	return t.wrappedFile.Load()
// }

// func (t TestFile) Path() string {
// 	return t.wrappedFile.Path()
// }
