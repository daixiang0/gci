package gci

import (
	"fmt"
	"os"
	"path/filepath"
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

func TestRun(t *testing.T) {
	for i := range testCases {
		t.Run(fmt.Sprintf("run case: %s", testCases[i].name), func(t *testing.T) {
			config, err := config.ParseConfig(testCases[i].config)
			if err != nil {
				t.Fatal(err)
			}

			old, new, err := LoadFormat([]byte(testCases[i].in), "", *config)
			if err != nil {
				t.Fatal(err)
			}

			assert.NoError(t, err)
			assert.Equal(t, testCases[i].in, string(old))
			assert.Equal(t, testCases[i].out, string(new))
		})
	}
}

func chdir(t *testing.T, dir string) {
	oldWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))

	// change back at the end of the test
	t.Cleanup(func() { os.Chdir(oldWd) })
}

func readConfig(t *testing.T, configPath string) *config.Config {
	rawConfig, err := os.ReadFile(configPath)
	require.NoError(t, err)
	cfg, err := config.ParseConfig(string(rawConfig))
	require.NoError(t, err)

	return cfg
}

func TestRunWithLocalModule(t *testing.T) {
	tests := []struct {
		name      string
		moduleDir string
		// files with a corresponding '*.out.go' file containing the expected
		// result of formatting
		testedFiles []string
	}{
		{
			name:      `default module test case`,
			moduleDir: filepath.Join("testdata", "module"),
			testedFiles: []string{
				"main.go",
				filepath.Join("internal", "foo", "lib.go"),
			},
		},
		{
			name:      `canonical module without go sources in root dir`,
			moduleDir: filepath.Join("testdata", "module_canonical"),
			testedFiles: []string{
				filepath.Join("cmd", "client", "main.go"),
				filepath.Join("cmd", "server", "main.go"),
				filepath.Join("internal", "foo", "lib.go"),
			},
		},
		{
			name:      `non-canonical module without go sources in root dir`,
			moduleDir: filepath.Join("testdata", "module_noncanonical"),
			testedFiles: []string{
				filepath.Join("cmd", "client", "main.go"),
				filepath.Join("cmd", "server", "main.go"),
				filepath.Join("internal", "foo", "lib.go"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// run subtests for expected module loading behaviour
			chdir(t, tt.moduleDir)
			cfg := readConfig(t, "config.yaml")

			for _, path := range tt.testedFiles {
				t.Run(path, func(t *testing.T) {
					// *.go -> *.out.go
					expected, err := os.ReadFile(strings.TrimSuffix(path, ".go") + ".out.go")
					require.NoError(t, err)

					_, got, err := LoadFormatGoFile(io.File{path}, *cfg)

					require.NoError(t, err)
					require.Equal(t, string(expected), string(got))
				})
			}
		})
	}
}

func TestRunWithLocalModuleWithPackageLoadFailure(t *testing.T) {
	// just a directory with no Go modules
	dir := t.TempDir()
	configContent := "sections:\n  - LocalModule\n"

	chdir(t, dir)
	_, err := config.ParseConfig(configContent)
	require.ErrorContains(t, err, "go.mod: open go.mod:")
}
