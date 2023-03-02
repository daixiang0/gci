package section

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sectionTestData struct {
	input           []string
	expectedSection SectionList
	expectedError   error
}

func TestParse(t *testing.T) {
	testCases := []sectionTestData{
		{
			input:           []string{""},
			expectedSection: nil,
			expectedError:   nil,
		},
		{
			input:           []string{"prefix(go)"},
			expectedSection: SectionList{Custom{"go"}},
			expectedError:   nil,
		},
		{
			input:           []string{"prefix(go-UPPER-case)"},
			expectedSection: SectionList{Custom{"go-UPPER-case"}},
			expectedError:   nil,
		},
		{
			input:           []string{"PREFIX(go-UPPER-case)"},
			expectedSection: SectionList{Custom{"go-UPPER-case"}},
			expectedError:   nil,
		},
		{
			input:           []string{"prefix("},
			expectedSection: nil,
			expectedError:   errors.New("invalid params: prefix("),
		},
		{
			input:           []string{"prefix(domainA,domainB)"},
			expectedSection: SectionList{Custom{"domainA,domainB"}},
			expectedError:   nil,
		},
	}
	for _, test := range testCases {
		parsedSection, err := Parse(test.input)
		assert.Equal(t, test.expectedSection, parsedSection)
		assert.Equal(t, test.expectedError, err)
	}
}
