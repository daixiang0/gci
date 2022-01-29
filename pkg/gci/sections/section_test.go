package sections

import (
	"fmt"
	"testing"

	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
)

type specificityTestData struct {
	importPath          string
	section             Section
	expectedSpecificity specificity.MatchSpecificity
}

func testSpecificity(t *testing.T, testCases []specificityTestData) {
	for _, test := range testCases {
		testName := fmt.Sprintf("%s:%v", test.importPath, test.section)
		t.Run(testName, testSpecificityCase(test))
	}
}

func testSpecificityCase(testData specificityTestData) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		detectedSpecificity := testData.section.MatchSpecificity(importPkg.ImportDef{QuotedPath: testData.importPath})
		if detectedSpecificity != testData.expectedSpecificity {
			t.Errorf("Specificity is %v and not %v", detectedSpecificity, testData.expectedSpecificity)
		}
	}
}
