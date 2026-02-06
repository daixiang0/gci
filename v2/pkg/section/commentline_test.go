package section

import (
	"testing"

	"github.com/daixiang0/gci/v2/pkg/specificity"
)

func TestCommentLineSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, CommentLine{""}, specificity.MisMatch{}},
		{`"x"`, CommentLine{""}, specificity.MisMatch{}},
		{`"//"`, CommentLine{""}, specificity.MisMatch{}},
		{`"/"`, CommentLine{""}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
