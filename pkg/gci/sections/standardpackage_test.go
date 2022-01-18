package sections

import (
	"testing"

	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func TestStandardPackageSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`"context"`, StandardPackage{}, specificity.StandardPackageMatch{}},
		{`"contexts"`, StandardPackage{}, specificity.MisMatch{}},
		{`"crypto"`, StandardPackage{}, specificity.StandardPackageMatch{}},
		{`"crypto1"`, StandardPackage{}, specificity.MisMatch{}},
		{`"crypto/ae"`, StandardPackage{}, specificity.MisMatch{}},
		{`"crypto/aes"`, StandardPackage{}, specificity.StandardPackageMatch{}},
		{`"crypto/aes2"`, StandardPackage{}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}

func TestStandardPackageParsing(t *testing.T) {
	testCases := []sectionTestData{
		{"sTd", StandardPackage{}, nil},
		{"STANDARD", StandardPackage{}, nil},
		{"Std(i)", nil, SectionTypeDoesNotAcceptParametersError},
	}
	testSectionParser(t, testCases)
}

func TestStandardPackageToString(t *testing.T) {
	testSectionToString(t, StandardPackage{})
	testSectionToString(t, StandardPackage{nil, CommentLine{"a"}})
	testSectionToString(t, StandardPackage{CommentLine{"a"}, NewLine{}})
}
