package section

import (
	"reflect"
	"testing"
)

type sectionTestData struct {
	input           []string
	expectedSection SectionList
	expectedError   string
}

func TestParse(t *testing.T) {
	testCases := []sectionTestData{
		{
			input:           []string{""},
			expectedSection: nil,
			expectedError:   "",
		},
		{
			input:           []string{"prefix(go)"},
			expectedSection: SectionList{Custom{Prefix: "go"}},
			expectedError:   "",
		},
		{
			input:           []string{"prefix(go-UPPER-case)"},
			expectedSection: SectionList{Custom{Prefix: "go-UPPER-case"}},
			expectedError:   "",
		},
		{
			input:           []string{"PREFIX(go-UPPER-case)"},
			expectedSection: SectionList{Custom{Prefix: "go-UPPER-case"}},
			expectedError:   "",
		},
		{
			input:           []string{"prefix("},
			expectedSection: nil,
			expectedError:   "invalid params: prefix(",
		},
		{
			input:           []string{"prefix(domainA,domainB)"},
			expectedSection: SectionList{Custom{Prefix: "domainA,domainB"}},
			expectedError:   "",
		},
	}
	for _, test := range testCases {
		parsedSection, err := Parse(test.input)
		if test.expectedError == "" {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		} else {
			if err == nil {
				t.Fatalf("expected error: %s", test.expectedError)
			}
			if err.Error() != test.expectedError {
				t.Fatalf("error mismatch: got=%v want=%v", err, test.expectedError)
			}
		}
		if !reflect.DeepEqual(test.expectedSection, parsedSection) {
			t.Fatalf("section mismatch: got=%v want=%v", parsedSection, test.expectedSection)
		}
	}
}
