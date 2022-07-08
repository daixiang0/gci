package section

import (
	"fmt"
	"testing"

	"github.com/daixiang0/gci/pkg/parse"
	"github.com/daixiang0/gci/pkg/specificity"
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
			t.Errorf("Specificity is %v and not %v", detectedSpecificity, testData.expectedSpecificity)
		}
	}
}
