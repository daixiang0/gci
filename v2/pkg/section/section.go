package section

import (
	"github.com/daixiang0/gci/v2/pkg/parse"
	"github.com/daixiang0/gci/v2/pkg/specificity"
)

type Section interface {
	MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity
	String() string
	Type() string
}

type SectionList []Section

func (list SectionList) String() []string {
	var output []string
	for _, section := range list {
		output = append(output, section.String())
	}
	return output
}

func DefaultSections() SectionList {
	return SectionList{Standard{}, Default{}}
}

func DefaultSectionSeparators() SectionList {
	return SectionList{NewLine{}}
}
