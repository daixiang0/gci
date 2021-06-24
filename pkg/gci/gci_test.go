package gci

import (
	"fmt"
	"testing"
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
