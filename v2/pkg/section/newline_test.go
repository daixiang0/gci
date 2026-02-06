package section

import (
	"testing"

	"github.com/daixiang0/gci/v2/pkg/specificity"
)

func TestNewLineSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, NewLine{}, specificity.MisMatch{}},
		{`"x"`, NewLine{}, specificity.MisMatch{}},
		{`"\n"`, NewLine{}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
