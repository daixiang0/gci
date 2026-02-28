package specificity

import (
	"testing"
)

func TestSpecificityOrder(t *testing.T) {
	testCases := testCasesInSpecificityOrder()
	for i := 1; i < len(testCases); i++ {
		if !testCases[i].IsMoreSpecific(testCases[i-1]) {
			t.Fatalf("expected %v to be more specific than %v", testCases[i], testCases[i-1])
		}
	}
}

func TestSpecificityEquality(t *testing.T) {
	for _, testCase := range testCasesInSpecificityOrder() {
		if !testCase.Equal(testCase) {
			t.Fatalf("expected %v to equal itself", testCase)
		}
	}
}

func testCasesInSpecificityOrder() []MatchSpecificity {
	return []MatchSpecificity{MisMatch{}, DefaultMatch{}, StandardMatch{}, Match{Length: 0}, Match{Length: 1}}
}
