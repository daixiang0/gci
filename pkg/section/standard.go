package section

import (
	"golang.org/x/tools/go/packages"

	"github.com/daixiang0/gci/pkg/parse"
	"github.com/daixiang0/gci/pkg/specificity"
)

const StandardType = "standard"

type Standard struct {
	standardPackages map[string]struct{}
}

func NewStandard() Standard {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	standardPackages := make(map[string]struct{})
	for _, p := range pkgs {
		standardPackages[p.PkgPath] = struct{}{}
	}
	return Standard{standardPackages: standardPackages}
}

func (s Standard) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	if _, ok := s.standardPackages[spec.Path]; ok {
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
