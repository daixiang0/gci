package section

import (
	"github.com/daixiang0/gci/v2/pkg/parse"
	"github.com/daixiang0/gci/v2/pkg/specificity"
)

const DefaultType = "default"

type Default struct{}

func (d Default) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	return specificity.DefaultMatch{}
}

func (d Default) String() string {
	return DefaultType
}

func (d Default) Type() string {
	return DefaultType
}
