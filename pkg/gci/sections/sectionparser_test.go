package sections

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sectionTestData struct {
	sectionDef      string
	expectedSection Section
	expectedError   error
}

func testSectionParser(t *testing.T, testCases []sectionTestData) {
	for _, test := range testCases {
		testName := fmt.Sprintf("%q-->(%v,%v)", test.sectionDef, test.expectedSection, test.expectedError)
		t.Run(testName, func(t *testing.T) {
			parsedSection, err := SectionParserInst.parseSectionString(test.sectionDef, true, true)
			assert.Equal(t, test.expectedSection, parsedSection)
			assert.True(t, errors.Is(err, test.expectedError))
		})
	}
}

func testSectionToString(t *testing.T, section Section) {
	testName := fmt.Sprintf("%#v", section)
	t.Run(testName, func(t *testing.T) {
		sectionStr := section.String()
		parsedSection, err := SectionParserInst.parseSectionString(sectionStr, true, true)
		assert.NoError(t, err)
		assert.Equal(t, section, parsedSection)
	})
}

func TestComplexParsingCases(t *testing.T) {
	testCases := []sectionTestData{
		{"Comment:defAult", DefaultSection{CommentLine{""}, nil}, nil},
		{":defAult:Comment(u)", DefaultSection{nil, CommentLine{"u"}}, nil},
	}
	testSectionParser(t, testCases)
}

func TestRegisterSectionAliasTwice(t *testing.T) {
	parser := SectionParser{}
	t1 := SectionType{
		aliases: []string{"a", "x"},
	}
	err := parser.RegisterSection(&t1)
	assert.NoError(t, err)
	err = parser.RegisterSection(&t1)
	assert.True(t, errors.Is(err, TypeAlreadyRegisteredError{}))
	t2 := SectionType{
		aliases: []string{"b", "x"},
	}
	err = parser.RegisterSection(&t2)
	assert.True(t, errors.Is(err, TypeAlreadyRegisteredError{}))
}
