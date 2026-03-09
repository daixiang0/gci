package section

import (
	"testing"

	"github.com/daixiang0/gci/pkg/specificity"
)

func TestStandardPackageSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{"cmp", Standard{}, specificity.StandardMatch{}},
		{"context", Standard{}, specificity.StandardMatch{}},
		{"contexts", Standard{}, specificity.MisMatch{}},
		{"crypto", Standard{}, specificity.StandardMatch{}},
		{"crypto1", Standard{}, specificity.MisMatch{}},
		{"crypto/ae", Standard{}, specificity.MisMatch{}},
		{"crypto/aes", Standard{}, specificity.StandardMatch{}},
		{"crypto/aes2", Standard{}, specificity.MisMatch{}},
		{"crypto/hpke", Standard{}, specificity.StandardMatch{}},
		{"syscall/js", Standard{}, specificity.StandardMatch{}},
	}
	testSpecificity(t, testCases)
}
