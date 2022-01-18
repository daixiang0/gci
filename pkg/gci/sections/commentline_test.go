package sections

import (
	"testing"

	"github.com/daixiang0/gci/pkg/gci/specificity"
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

func TestCommentLineParsing(t *testing.T) {
	testCases := []sectionTestData{
		{"commentline", CommentLine{""}, nil},
		{"Commentline(abc)", CommentLine{"abc"}, nil},
		{"cOmMenT(x)", CommentLine{"x"}, nil},
		{"Comment:Comment", nil, SectionTypeDoesNotAcceptPrefixError},
		{"Comment:Comment:Comment()", nil, SectionTypeDoesNotAcceptPrefixError},
	}
	testSectionParser(t, testCases)
}

func TestCommentLineToString(t *testing.T) {
	testSectionToString(t, CommentLine{""})
	testSectionToString(t, CommentLine{"abc"})
}
