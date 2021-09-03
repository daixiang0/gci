package gci

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPkgType(t *testing.T) {
	testCases := []struct {
		Line           string
		LocalFlag      string
		ExpectedResult int
	}{
		{Line: `"foo/pkg/bar"`, LocalFlag: "", ExpectedResult: remote},
		{Line: `"foo/pkg/bar"`, LocalFlag: "foo", ExpectedResult: local},
		{Line: `"foo/pkg/bar"`, LocalFlag: "bar", ExpectedResult: remote},
		{Line: `"foo/pkg/bar"`, LocalFlag: "github.com/foo/bar", ExpectedResult: remote},
		{Line: `"foo/pkg/bar"`, LocalFlag: "github.com/foo", ExpectedResult: remote},
		{Line: `"foo/pkg/bar"`, LocalFlag: "github.com/bar", ExpectedResult: remote},
		{Line: `"foo/pkg/bar"`, LocalFlag: "github.com/foo,github.com/bar", ExpectedResult: remote},
		{Line: `"foo/pkg/bar"`, LocalFlag: "github.com/foo,,github.com/bar", ExpectedResult: remote},

		{Line: `"github.com/foo/bar"`, LocalFlag: "", ExpectedResult: remote},
		{Line: `"github.com/foo/bar"`, LocalFlag: "foo", ExpectedResult: remote},
		{Line: `"github.com/foo/bar"`, LocalFlag: "bar", ExpectedResult: remote},
		{Line: `"github.com/foo/bar"`, LocalFlag: "github.com/foo/bar", ExpectedResult: local},
		{Line: `"github.com/foo/bar"`, LocalFlag: "github.com/foo", ExpectedResult: local},
		{Line: `"github.com/foo/bar"`, LocalFlag: "github.com/bar", ExpectedResult: remote},
		{Line: `"github.com/foo/bar"`, LocalFlag: "github.com/foo,github.com/bar", ExpectedResult: local},
		{Line: `"github.com/foo/bar"`, LocalFlag: "github.com/foo,,github.com/bar", ExpectedResult: local},

		{Line: `"context"`, LocalFlag: "", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "context", ExpectedResult: local},
		{Line: `"context"`, LocalFlag: "foo", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "bar", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "github.com/foo/bar", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "github.com/foo", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "github.com/bar", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "github.com/foo,github.com/bar", ExpectedResult: standard},
		{Line: `"context"`, LocalFlag: "github.com/foo,,github.com/bar", ExpectedResult: standard},

		{Line: `"os/signal"`, LocalFlag: "", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "os/signal", ExpectedResult: local},
		{Line: `"os/signal"`, LocalFlag: "foo", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "bar", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "github.com/foo/bar", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "github.com/foo", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "github.com/bar", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "github.com/foo,github.com/bar", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "github.com/foo,,github.com/bar", ExpectedResult: standard},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%s:%s", tc.Line, tc.LocalFlag), func(t *testing.T) {
			t.Parallel()

			result := getPkgType(tc.Line, ParseLocalFlag(tc.LocalFlag))
			if got, want := result, tc.ExpectedResult; got != want {
				t.Errorf("bad result: %d, expected: %d", got, want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	fileinfos, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, fileinfo := range fileinfos {
		inname := fileinfo.Name()
		if strings.HasPrefix(inname, ".") || !strings.HasSuffix(inname, ".in.go") {
			continue
		}
		name := strings.TrimSuffix(inname, ".in.go")
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			expectedInput, err := ioutil.ReadFile(filepath.Join("testdata", inname))
			if err != nil {
				t.Fatal(err)
			}
			outname := name + ".out.go"
			expectedOutput, err := ioutil.ReadFile(filepath.Join("testdata", outname))
			if err != nil {
				t.Fatal(err)
			}

			flagSet := &FlagSet{
				LocalFlag: []string{
					"github.com/daixiang0",
					"github.com/local",
				},
			}

			actualInput, actualOutput, err := Run(filepath.Join("testdata", inname), flagSet)
			assert.Equal(t, string(expectedInput), string(actualInput), "input")
			if bytes.Equal(expectedInput, expectedOutput) {
				assert.Nil(t, actualOutput, "output")

			}
			if actualOutput == nil {
				actualOutput = actualInput
			}
			assert.Equal(t, string(expectedOutput), string(actualOutput), "output")
			assert.NoError(t, err)
		})
	}
}
