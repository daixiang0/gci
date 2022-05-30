package sections

import (
	"testing"

	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func TestNewLineSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, NewLine{}, specificity.MisMatch{}},
		{`"x"`, NewLine{}, specificity.MisMatch{}},
		{`"\n"`, NewLine{}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}

func TestNewLineParsing(t *testing.T) {
	testCases := []sectionTestData{
		{"nl", NewLine{}, nil},
		{"newLine", NewLine{}, nil},
		{"newLine:nl", nil, SectionTypeDoesNotAcceptPrefixError},
		{"NL(invalid)", nil, SectionTypeDoesNotAcceptParametersError},
	}
	testSectionParser(t, testCases)
}

func TestNewLineToString(t *testing.T) {
	testSectionToString(t, NewLine{})
}
