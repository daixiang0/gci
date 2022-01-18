package sections

import (
	"testing"

	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func TestDefaultSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, DefaultSection{}, specificity.Default{}},
		{`"x"`, DefaultSection{}, specificity.Default{}},
	}
	testSpecificity(t, testCases)
}

func TestDefaultSectionParsing(t *testing.T) {
	testCases := []sectionTestData{
		{"def", DefaultSection{}, nil},
		{"defAult", DefaultSection{nil, nil}, nil},
		{"defAult(invalid)", nil, SectionTypeDoesNotAcceptParametersError},
	}
	testSectionParser(t, testCases)
}

func TestDefaultSectionToString(t *testing.T) {
	testSectionToString(t, DefaultSection{})
	testSectionToString(t, DefaultSection{nil, nil})
	testSectionToString(t, DefaultSection{nil, NewLine{}})
	testSectionToString(t, DefaultSection{CommentLine{"a"}, CommentLine{"b"}})
}
