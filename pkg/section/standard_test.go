package section

import (
	"testing"

	"github.com/daixiang0/gci/pkg/specificity"
)

func TestStandardPackageSpecificity(t *testing.T) {
	standard := NewStandard()
	testCases := []specificityTestData{
		{"context", standard, specificity.StandardMatch{}},
		{"contexts", standard, specificity.MisMatch{}},
		{"crypto", standard, specificity.StandardMatch{}},
		{"crypto1", standard, specificity.MisMatch{}},
		{"crypto/ae", standard, specificity.MisMatch{}},
		{"crypto/aes", standard, specificity.StandardMatch{}},
		{"crypto/aes2", standard, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
