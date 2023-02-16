package section

import (
	"testing"

	"github.com/daixiang0/gci/pkg/specificity"
)

func TestStandardPackageSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{"context", NewStandard(), specificity.StandardMatch{}},
		{"contexts", NewStandard(), specificity.MisMatch{}},
		{"crypto", NewStandard(), specificity.StandardMatch{}},
		{"crypto1", NewStandard(), specificity.MisMatch{}},
		{"crypto/ae", NewStandard(), specificity.MisMatch{}},
		{"crypto/aes", NewStandard(), specificity.StandardMatch{}},
		{"crypto/aes2", NewStandard(), specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
