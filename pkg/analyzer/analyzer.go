package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gci",
	Doc:      "A tool that control golang package import order and make it always deterministic.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}
