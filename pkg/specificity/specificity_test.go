package specificity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecificityOrder(t *testing.T) {
	testCases := testCasesInSpecificityOrder()
	for i := 1; i < len(testCases); i++ {
		t.Run(fmt.Sprintf("Specificity(%v)>Specificity(%v)", testCases[i], testCases[i-1]), func(t *testing.T) {
			assert.True(t, testCases[i].IsMoreSpecific(testCases[i-1]))
		})
	}
}

func TestSpecificityEquality(t *testing.T) {
	for _, testCase := range testCasesInSpecificityOrder() {
		t.Run(fmt.Sprintf("Specificity(%v)==Specificity(%v)", testCase, testCase), func(t *testing.T) {
			assert.True(t, testCase.Equal(testCase))
		})
	}
}

func testCasesInSpecificityOrder() []MatchSpecificity {
	return []MatchSpecificity{MisMatch{}, Default{}, StandardMatch{}, Match{0}, Match{1}}
}
