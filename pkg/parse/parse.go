package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
)

const C = "\"C\""

type GciImports struct {
	// original index of import group, include doc, name, path and comment
	Start, End int
	Name, Path string
}
type ImportList []*GciImports

func (l ImportList) Len() int {
	return len(l)
}

func (l ImportList) Less(i, j int) bool {
	if strings.Compare(l[i].Path, l[j].Path) == 0 {
		return strings.Compare(l[i].Name, l[j].Name) < 0
	}

	return strings.Compare(l[i].Path, l[j].Path) < 0
}

func (l ImportList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

/*
 * AST considers a import block as below:
 * ```
 * Doc
 * Name Path Comment
 * ```
 * An example is like below:
 * ```
 * // test
 * test "fmt" // test
 * ```
 * getImports return a import block with name, start and end index
 */
func getImports(imp *ast.ImportSpec) (start, end int, name string) {
	if imp.Doc != nil {
		// doc poc need minus one to get the first index of comment
		start = int(imp.Doc.Pos()) - 1
	} else {
		if imp.Name != nil {
			// name pos need minus one too
			start = int(imp.Name.Pos()) - 1
		} else {
			// path pos start without quote, need minus one for it
			start = int(imp.Path.Pos()) - 1
		}
	}

	if imp.Name != nil {
		name = imp.Name.Name
	}

	if imp.Comment != nil {
		end = int(imp.Comment.End())
	} else {
		end = int(imp.Path.End())
	}
	return
}

func ParseFile(src []byte, filename string) (ImportList, int, int, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return nil, 0, 0, err
	}

	if len(f.Imports) == 0 {
		return nil, 0, 0, NoImportError{}
	}

	var (
		// headEnd means the start of import block
		headEnd int
		// tailStart means the end + 1 of import block
		tailStart int
		// lastImportStart means the start of last import block
		lastImportStart int
		data            ImportList
	)

	for i, d := range f.Decls {
		switch d.(type) {
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)
			if dd.Tok == token.IMPORT {
				// there are two cases, both end with linebreak:
				// 1.
				// import (
				//	 "xxxx"
				// )
				// 2.
				// import "xxx"
				if headEnd == 0 {
					headEnd = int(d.Pos()) - 1
				}
				tailStart = int(d.End())
				lastImportStart = i
			}
		}
	}

	if len(f.Decls) > lastImportStart+1 {
		tailStart = int(f.Decls[lastImportStart+1].Pos() - 1)
	}

	for _, imp := range f.Imports {
		if imp.Path.Value == C {
			if imp.Comment != nil {
				headEnd = int(imp.Comment.End())
			} else {
				headEnd = int(imp.Path.End())
			}
			continue
		}

		start, end, name := getImports(imp)

		data = append(data, &GciImports{
			Start: start,
			End:   end,
			Name:  name,
			Path:  strings.Trim(imp.Path.Value, `"`),
		})
	}

	sort.Sort(data)
	return data, headEnd, tailStart, nil
}

// IsGeneratedFileByComment reports whether the source file is generated code.
// Using a bit laxer rules than https://golang.org/s/generatedcode to
// match more generated code.
// Taken from https://github.com/golangci/golangci-lint.
func IsGeneratedFileByComment(in string) bool {
	const (
		genCodeGenerated = "code generated"
		genDoNotEdit     = "do not edit"
		genAutoFile      = "autogenerated file" // easyjson
	)

	markers := []string{genCodeGenerated, genDoNotEdit, genAutoFile}
	in = strings.ToLower(in)
	for _, marker := range markers {
		if strings.Contains(in, marker) {
			return true
		}
	}

	return false
}

type NoImportError struct{}

func (n NoImportError) Error() string {
	return "No imports"
}

func (i NoImportError) Is(err error) bool {
	_, ok := err.(NoImportError)
	return ok
}
