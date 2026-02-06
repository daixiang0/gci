package section

import (
	"fmt"
	"strings"

	"github.com/daixiang0/gci/v2/pkg/parse"
	"github.com/daixiang0/gci/v2/pkg/specificity"
)

type Custom struct {
	Prefix string
}

const CustomSeparator = ","

const CustomType = "custom"

func (c Custom) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	for _, prefix := range strings.Split(c.Prefix, CustomSeparator) {
		prefix = strings.TrimSpace(prefix)
		if strings.HasPrefix(spec.Path, prefix) {
			return specificity.Match{Length: len(prefix)}
		}
	}

	return specificity.MisMatch{}
}

func (c Custom) String() string {
	return fmt.Sprintf("prefix(%s)", c.Prefix)
}

func (c Custom) Type() string {
	return CustomType
}
