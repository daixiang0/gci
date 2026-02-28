package section

import (
	"github.com/daixiang0/gci/v2/pkg/parse"
	"github.com/daixiang0/gci/v2/pkg/specificity"
)

type NewLine struct{}

const NewLineType = "newline"

func (n NewLine) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	return specificity.MisMatch{}
}

func (n NewLine) String() string {
	return ""
}

func (n NewLine) Type() string {
	return NewLineType
}
