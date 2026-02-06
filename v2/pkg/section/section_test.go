package section

import (
	"fmt"
	"testing"

	"github.com/daixiang0/gci/v2/pkg/parse"
	"github.com/daixiang0/gci/v2/pkg/specificity"
)

type specificityTestData struct {
	path                string
	section             Section
	expectedSpecificity specificity.MatchSpecificity
}

func testSpecificity(t *testing.T, testCases []specificityTestData) {
	for _, test := range testCases {
		testName := fmt.Sprintf("%s:%v", test.path, test.section)
		t.Run(testName, testSpecificityCase(test))
	}
}

func testSpecificityCase(testData specificityTestData) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		detectedSpecificity := testData.section.MatchSpecificity(&parse.GciImports{Path: testData.path})
		if detectedSpecificity != testData.expectedSpecificity {
			t.Fatalf("specificity mismatch: got=%v want=%v", detectedSpecificity, testData.expectedSpecificity)
		}
	}
}
