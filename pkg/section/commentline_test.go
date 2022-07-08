package section

import (
	"testing"

	"github.com/daixiang0/gci/pkg/specificity"
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

// func TestCommentLineToString(t *testing.T) {
// 	testSectionToString(t, CommentLine{""})
// 	testSectionToString(t, CommentLine{"abc"})
// }
