package sections

import (
	"testing"

	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func TestPrefixSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`"foo/pkg/bar"`, Prefix{"", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"foo", nil, nil}, specificity.Match{3}},
		{`"foo/pkg/bar"`, Prefix{"bar", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"github.com/foo/bar", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"github.com/foo", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"github.com/bar", nil, nil}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}

func TestPrefixParsing(t *testing.T) {
	testCases := []sectionTestData{
		{"pkgPREFIX", Prefix{"", nil, nil}, nil},
		{"prefix(test.com)", Prefix{"test.com", nil, nil}, nil},
	}
	testSectionParser(t, testCases)
}

func TestPrefixToString(t *testing.T) {
	testSectionToString(t, Prefix{})
	testSectionToString(t, Prefix{"", nil, nil})
	testSectionToString(t, Prefix{"abc.org", nil, nil})
	testSectionToString(t, Prefix{"abc.org", nil, CommentLine{"a"}})
	testSectionToString(t, Prefix{"abc.org", CommentLine{"a"}, NewLine{}})
}
