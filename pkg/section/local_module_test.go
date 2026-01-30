package section

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daixiang0/gci/pkg/specificity"
)

func TestLocalModule_specificity(t *testing.T) {
	singleLocalModule := &LocalModule{Paths: []string{"example.com/hello"}}
	multipleLocalModules := &LocalModule{Paths: []string{
		"example.com/module1",
		"example.com/module2",
	}}

	testCases := []specificityTestData{
		{"example.com/hello", singleLocalModule, specificity.LocalModule{}},
		{"example.com/hello/world", singleLocalModule, specificity.LocalModule{}},
		{"example.com/hello-world", singleLocalModule, specificity.MisMatch{}},
		{"example.com/helloworld", singleLocalModule, specificity.MisMatch{}},
		{"example.com", singleLocalModule, specificity.MisMatch{}},

		{"example.com/module1", multipleLocalModules, specificity.LocalModule{}},
		{"example.com/module2", multipleLocalModules, specificity.LocalModule{}},
		{"example.com/module1/world", multipleLocalModules, specificity.LocalModule{}},
		{"example.com/module2/foo/bar", multipleLocalModules, specificity.LocalModule{}},
		{"example.com/module2butnotreally", multipleLocalModules, specificity.MisMatch{}},
		{"example.com", multipleLocalModules, specificity.MisMatch{}},
	}

	testSpecificity(t, testCases)
}

func TestLocalModule_findLocalModules(t *testing.T) {
	m := new(LocalModule)
	testdata := filepath.Join("./testdata", "local_module")

	for name, tt := range map[string]struct {
		testdataDir          string
		expectFailure        bool
		expectedModulesPaths []string
		expectedErrorMessage string
	}{
		"within root module": {
			testdataDir:          "mod_simple",
			expectedModulesPaths: []string{"fake.tld/example/simple"},
		},
		"within workspace": {
			testdataDir:          "work_simple",
			expectedModulesPaths: []string{"fake.tld/example/module1", "fake.tld/example/module2"},
		},
		"both files - go.work precedence": {
			testdataDir:          "both_files",
			expectedModulesPaths: []string{"fake.tld/example/work"},
		},
		"empty directory": {
			testdataDir:          "empty_dir",
			expectedModulesPaths: []string{},
		},
		"redundant paths deduplication": {
			testdataDir:          "work_redundant_paths",
			expectedModulesPaths: []string{"fake.tld/example/foo/bar", "fake.tld/other/project"},
		},
		"broken go.mod": {
			testdataDir:   "mod_malformed",
			expectFailure: true,
		},
		"broken go.work": {
			testdataDir:   "work_malformed",
			expectFailure: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Chdir(filepath.Join(testdata, tt.testdataDir))

			modPaths, err := m.findLocalModules()
			if tt.expectFailure {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.ElementsMatch(t, tt.expectedModulesPaths, modPaths)
		})
	}
}

func TestLocalModule_getModulePathFromRootMod(t *testing.T) {
	m := new(LocalModule)
	testdata := filepath.Join("./testdata", "local_module")

	t.Run("non existing file", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "empty_dir"))

		modPath, err := m.getModulePathFromRootMod()
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, modPath)
	})

	t.Run("invalid mod file", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "mod_malformed"))

		modPath, err := m.getModulePathFromRootMod()
		assert.ErrorContains(t, err, "no module path found")
		assert.NotErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, modPath)
	})

	t.Run("mod path found", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "mod_simple"))

		modPath, err := m.getModulePathFromRootMod()
		assert.NoError(t, err)
		assert.Equal(t, "fake.tld/example/simple", modPath)
	})

	t.Run("mod path found with GOMOD env", func(t *testing.T) {
		t.Setenv("GOMOD", filepath.Join(testdata, "mod_simple", "go.mod"))

		modPath, err := m.getModulePathFromRootMod()
		assert.NoError(t, err)
		assert.Equal(t, "fake.tld/example/simple", modPath)
	})
}

func TestLocalModule_getModulePathFromWorkspace(t *testing.T) {
	m := new(LocalModule)
	testdata := filepath.Join("./testdata", "local_module")

	t.Run("non existing file", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "empty_dir"))

		modsPath, err := m.getModulesPathFromWorkspace()
		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, modsPath)
	})

	t.Run("invalid work file", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "work_malformed"))

		modsPath, err := m.getModulesPathFromWorkspace()
		assert.ErrorContains(t, err, "unable to parse go.work file")
		assert.NotErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, modsPath)
	})

	t.Run("invalid use attributes", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "work_missing_referenced_mod"))

		modsPath, err := m.getModulesPathFromWorkspace()
		assert.ErrorContains(t, err, "unable to get mod file")
		assert.NotErrorIs(t, err, os.ErrNotExist)
		assert.Empty(t, modsPath)
	})

	t.Run("work file found and valid", func(t *testing.T) {
		t.Chdir(filepath.Join(testdata, "work_simple"))

		modsPath, err := m.getModulesPathFromWorkspace()
		assert.NoError(t, err)
		assert.Equal(t, []string{"fake.tld/example/module1", "fake.tld/example/module2"}, modsPath)
	})
}

func TestLocalModule_removeRedundantModulePaths(t *testing.T) {
	m := new(LocalModule)

	for name, tt := range map[string]struct {
		input    []string
		expected []string
	}{
		"empty input": {
			input:    []string{},
			expected: []string{},
		},
		"single input": {
			input:    []string{"github.com/foo"},
			expected: []string{"github.com/foo"},
		},
		"no redundancy": {
			input:    []string{"github.com/foo", "github.com/bar"},
			expected: []string{"github.com/foo", "github.com/bar"},
		},
		"multiple redundancies": {
			input:    []string{"github.com/foo/bar", "github.com/foo/bar/sub1", "github.com/foo/bar/sub2", "fake.tld/other"},
			expected: []string{"github.com/foo/bar", "fake.tld/other"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			result := m.removeRedundantModulePaths(tt.input)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
