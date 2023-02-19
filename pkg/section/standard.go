package section

import (
	"strings"

	"github.com/daixiang0/gci/pkg/parse"
	"github.com/daixiang0/gci/pkg/specificity"
)

const StandardType = "standard"

type Standard struct{}

func (s Standard) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	if isStandard(spec.Path) {
		return specificity.StandardMatch{}
	}
	return specificity.MisMatch{}
}

func (s Standard) String() string {
	return StandardType
}

func (s Standard) Type() string {
	return StandardType
}

func isStandard(pkg string) bool {
	if index := strings.Index(pkg, "/"); index != -1 {
		pkg = pkg[:index]
	}
	return !strings.Contains(pkg, ".")
}
