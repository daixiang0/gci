package gci

import (
	"fmt"
	"testing"
)

func TestGetPkgType(b *testing.T) {
	testCases := []struct {
		Line           string
		LocalFlag      string
		ExpectedResult int
	}{
		{Line: `"foo/pkg/bar"`, LocalFlag: "foo", ExpectedResult: local},
		{Line: `"github.com/foo/bar"`, LocalFlag: "foo", ExpectedResult: remote},
		{Line: `"context"`, LocalFlag: "foo", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "foo", ExpectedResult: standard},
		{Line: `"foo/pkg/bar"`, LocalFlag: "bar", ExpectedResult: remote},
		{Line: `"github.com/foo/bar"`, LocalFlag: "bar", ExpectedResult: remote},
		{Line: `"context"`, LocalFlag: "bar", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "bar", ExpectedResult: standard},
		{Line: `"foo/pkg/bar"`, LocalFlag: "github.com/foo/bar", ExpectedResult: remote},
		{Line: `"github.com/foo/bar"`, LocalFlag: "github.com/foo/bar", ExpectedResult: local},
		{Line: `"context"`, LocalFlag: "github.com/foo/bar", ExpectedResult: standard},
		{Line: `"os/signal"`, LocalFlag: "github.com/foo/bar", ExpectedResult: standard},
	}

	for _, _tCase := range testCases {
		tCase := _tCase
		testFn := func(t *testing.T) {
			result := getPkgType(tCase.Line, tCase.LocalFlag)
			if got, want := result, tCase.ExpectedResult; got != want {
				t.Errorf("bad result: %d, expected: %d", got, want)
			}
		}
		b.Run(fmt.Sprintf("%s:%s", tCase.LocalFlag, tCase.Line), testFn)
	}
}
