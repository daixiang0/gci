package imports

import (
	"go/ast"
	"go/token"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/daixiang0/gci/v2/pkg/config"
	gciParse "github.com/daixiang0/gci/v2/pkg/parse"
	"github.com/daixiang0/gci/v2/pkg/specificity"
)

func sortImports(cfg *config.Config, tokFile *token.File, f *ast.File) {
	for i, d := range f.Decls {
		d, ok := d.(*ast.GenDecl)
		if !ok || d.Tok != token.IMPORT {
			break
		}

		if len(d.Specs) == 0 {
			f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		}

		if !d.Lparen.IsValid() {
			continue
		}

		// Check if this import block contains CGO import "C"
		// If so, skip sorting to preserve CGO structure
		hasCgo := false
		for _, spec := range d.Specs {
			if impSpec := spec.(*ast.ImportSpec); importPath(impSpec) == "C" {
				hasCgo = true
				break
			}
		}

		if hasCgo {
			// For CGO blocks, preserve the original structure
			// Only sort if there are multiple non-CGO imports
			continue
		}

		// Sort all specs together based on section configuration
		d.Specs = sortSpecs(cfg, tokFile, f, d.Specs)

		if len(d.Specs) > 0 {
			lastSpec := d.Specs[len(d.Specs)-1]
			lastLine := tokFile.Line(lastSpec.Pos())
			if rParenLine := tokFile.Line(d.Rparen); rParenLine > lastLine+1 {
				tokFile.MergeLine(rParenLine - 1)
			}
		}
	}
}

func mergeImports(f *ast.File) {
	if len(f.Decls) <= 1 {
		return
	}

	var first *ast.GenDecl
	for i := 0; i < len(f.Decls); i++ {
		decl := f.Decls[i]
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT || declImports(gen, "C") {
			continue
		}
		if first == nil {
			first = gen
			continue
		}
		first.Lparen = first.Pos()
		for _, spec := range gen.Specs {
			spec.(*ast.ImportSpec).Path.ValuePos = first.Pos()
			first.Specs = append(first.Specs, spec)
		}
		f.Decls = append(f.Decls[:i], f.Decls[i+1:]...)
		i--
	}
}

func moveCgoDeclsToTop(f *ast.File) {
	if len(f.Decls) == 0 {
		return
	}
	importEnd := 0
	for importEnd < len(f.Decls) {
		gen, ok := f.Decls[importEnd].(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			break
		}
		importEnd++
	}
	if importEnd <= 1 {
		return
	}
	var cgoDecls []ast.Decl
	var otherDecls []ast.Decl
	for _, decl := range f.Decls[:importEnd] {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT && declImports(gen, "C") {
			cgoDecls = append(cgoDecls, decl)
			continue
		}
		otherDecls = append(otherDecls, decl)
	}
	if len(cgoDecls) == 0 {
		return
	}
	reordered := append([]ast.Decl{}, cgoDecls...)
	reordered = append(reordered, otherDecls...)
	f.Decls = append(reordered, f.Decls[importEnd:]...)
}

func declImports(gen *ast.GenDecl, path string) bool {
	if gen.Tok != token.IMPORT {
		return false
	}
	for _, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)
		if importPath(impspec) == path {
			return true
		}
	}
	return false
}

func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err == nil {
		return t
	}
	return ""
}

func importName(s *ast.ImportSpec) string {
	n := s.Name
	if n == nil {
		return ""
	}
	return n.Name
}

func importComment(s *ast.ImportSpec) string {
	c := s.Comment
	if c == nil {
		return ""
	}
	return c.Text()
}

func collapse(prev, next *ast.ImportSpec) bool {
	if importPath(next) != importPath(prev) || importName(next) != importName(prev) {
		return false
	}
	return prev.Comment == nil
}

type posSpan struct {
	Start token.Pos
	End   token.Pos
}

func sortSpecs(cfg *config.Config, tokFile *token.File, f *ast.File, specs []ast.Spec) []ast.Spec {
	if len(specs) <= 1 {
		return specs
	}

	pos := make([]posSpan, len(specs))
	for i, s := range specs {
		pos[i] = posSpan{s.Pos(), s.End()}
	}

	lastLine := tokFile.Line(pos[len(pos)-1].End)
	cstart := len(f.Comments)
	cend := len(f.Comments)
	for i, g := range f.Comments {
		if g.Pos() < pos[0].Start {
			continue
		}
		if i < cstart {
			cstart = i
		}
		if tokFile.Line(g.End()) > lastLine {
			cend = i
			break
		}
	}
	comments := f.Comments[cstart:cend]
	specLines := map[int]struct{}{}
	for _, s := range specs {
		specLines[tokFile.Line(s.Pos())] = struct{}{}
	}
	var expandedComments []*ast.CommentGroup
	for _, g := range comments {
		var current *ast.CommentGroup
		for _, c := range g.List {
			_, trailing := specLines[tokFile.Line(c.Slash)]
			if trailing {
				if current != nil && len(current.List) > 0 {
					expandedComments = append(expandedComments, current)
					current = nil
				}
				expandedComments = append(expandedComments, &ast.CommentGroup{List: []*ast.Comment{c}})
				continue
			}
			if current == nil {
				current = &ast.CommentGroup{}
			}
			current.List = append(current.List, c)
		}
		if current != nil && len(current.List) > 0 {
			expandedComments = append(expandedComments, current)
		}
	}
	if len(expandedComments) > 0 {
		head := append([]*ast.CommentGroup{}, f.Comments[:cstart]...)
		head = append(head, expandedComments...)
		head = append(head, f.Comments[cend:]...)
		f.Comments = head
		comments = f.Comments[cstart : cstart+len(expandedComments)]
	}

	docGroups := map[*ast.ImportSpec][]*ast.CommentGroup{}
	trailingGroups := map[*ast.ImportSpec][]*ast.CommentGroup{}
	type commentGroupOffset struct {
		offsets   []int
		columns   []int
		maxOffset int
	}
	commentOffsets := map[*ast.CommentGroup]commentGroupOffset{}
	for _, g := range comments {
		cLine := tokFile.Line(g.Pos())
		bestIdx := len(specs) - 1
		bestDist := int(^uint(0) >> 1)
		for i, spec := range specs {
			specLine := tokFile.Line(spec.Pos())
			if cLine <= specLine {
				dist := specLine - cLine
				if dist < bestDist {
					bestDist = dist
					bestIdx = i
				}
			}
		}
		s := specs[bestIdx].(*ast.ImportSpec)
		lines := make([]int, len(g.List))
		columns := make([]int, len(g.List))
		minLine := int(^uint(0) >> 1)
		for i, c := range g.List {
			pos := tokFile.Position(c.Slash)
			lines[i] = pos.Line
			columns[i] = pos.Column
			if pos.Line < minLine {
				minLine = pos.Line
			}
		}
		offsets := make([]int, len(lines))
		maxOffset := 0
		for i, line := range lines {
			offset := line - minLine
			offsets[i] = offset
			if offset > maxOffset {
				maxOffset = offset
			}
		}
		commentOffsets[g] = commentGroupOffset{
			offsets:   offsets,
			columns:   columns,
			maxOffset: maxOffset,
		}
		if cLine < tokFile.Line(s.Pos()) {
			docGroups[s] = append(docGroups[s], g)
		} else {
			trailingGroups[s] = append(trailingGroups[s], g)
		}
	}

	sort.Sort(byImportSpec{cfg, specs})

	deduped := specs[:0]
	for i, s := range specs {
		if i == len(specs)-1 || !collapse(s.(*ast.ImportSpec), specs[i+1].(*ast.ImportSpec)) {
			deduped = append(deduped, s)
		} else {
			p := s.Pos()
			tokFile.MergeLine(tokFile.Line(p))
		}
	}
	specs = deduped

	for i, s := range specs {
		s := s.(*ast.ImportSpec)
		if i == 0 && len(docGroups[s]) > 0 {
			maxOffset := 0
			for _, g := range docGroups[s] {
				if meta := commentOffsets[g]; meta.maxOffset > maxOffset {
					maxOffset = meta.maxOffset
				}
			}
			shift := maxOffset + 1
			startLine := tokFile.Line(pos[i].Start)
			startOffset := pos[i].Start - tokFile.LineStart(startLine)
			pos[i].Start = tokFile.LineStart(startLine+shift) + startOffset
			endLine := tokFile.Line(pos[i].End)
			endOffset := pos[i].End - tokFile.LineStart(endLine)
			pos[i].End = tokFile.LineStart(endLine+shift) + endOffset
		}
		if s.Name != nil {
			s.Name.NamePos = pos[i].Start
		}
		s.Path.ValuePos = pos[i].Start
		s.EndPos = pos[i].End

		importLine := tokFile.Line(pos[i].Start)
		for _, g := range docGroups[s] {
			meta := commentOffsets[g]
			startLine := importLine - 1 - meta.maxOffset
			if i > 0 {
				prevLine := tokFile.Line(pos[i-1].End)
				if startLine < prevLine+1 {
					startLine = prevLine + 1
				}
			}
			if i == 0 {
				startLine = importLine - meta.maxOffset
			}
			if startLine < 1 {
				startLine = 1
			}
			for j, c := range g.List {
				line := startLine + meta.offsets[j]
				if line < 1 {
					line = 1
				}
				c.Slash = tokFile.LineStart(line) + token.Pos(meta.columns[j]-1)
			}
		}
		for _, g := range trailingGroups[s] {
			meta := commentOffsets[g]
			startLine := importLine
			for j, c := range g.List {
				line := startLine + meta.offsets[j]
				if line < 1 {
					line = 1
				}
				c.Slash = tokFile.LineStart(line) + token.Pos(meta.columns[j]-1)
			}
		}
	}

	sort.Sort(byCommentPos(comments))


	if len(comments) == 0 {
		firstSpecLine := tokFile.Line(specs[0].Pos())
		for _, s := range specs[1:] {
			p := s.Pos()
			line := tokFile.Line(p)
			for previousLine := line - 1; previousLine >= firstSpecLine; {
				if previousLine > 0 && previousLine < tokFile.LineCount() {
					tokFile.MergeLine(previousLine)
					previousLine--
				} else {
					req := "Please report what the imports section of your go file looked like."
					log.Printf("panic avoided: first:%d line:%d previous:%d max:%d. %s",
						firstSpecLine, line, previousLine, tokFile.LineCount(), req)
				}
			}
		}
	}
	return specs
}

type byImportSpec struct {
	cfg   *config.Config
	specs []ast.Spec
}

func (x byImportSpec) Len() int      { return len(x.specs) }
func (x byImportSpec) Swap(i, j int) { x.specs[i], x.specs[j] = x.specs[j], x.specs[i] }
func (x byImportSpec) Less(i, j int) bool {
	ipath := importPath(x.specs[i].(*ast.ImportSpec))
	jpath := importPath(x.specs[j].(*ast.ImportSpec))

	igroup := importGroup(x.cfg, ipath, x.specs[i].(*ast.ImportSpec))
	jgroup := importGroup(x.cfg, jpath, x.specs[j].(*ast.ImportSpec))
	if igroup != jgroup {
		return igroup < jgroup
	}

	if ipath != jpath {
		return ipath < jpath
	}
	iname := importName(x.specs[i].(*ast.ImportSpec))
	jname := importName(x.specs[j].(*ast.ImportSpec))

	if iname != jname {
		return iname < jname
	}
	return importComment(x.specs[i].(*ast.ImportSpec)) < importComment(x.specs[j].(*ast.ImportSpec))
}

func importGroup(cfg *config.Config, importPath string, spec *ast.ImportSpec) int {
	if cfg == nil || len(cfg.Sections) == 0 {
		return defaultImportGroup(importPath)
	}

	gciImport := &gciParse.GciImports{
		Path: importPath,
		Name: importName(spec),
	}

	// Find the most specific match
	bestMatchIdx := -1
	var bestMatchSpec specificity.MatchSpecificity = specificity.MisMatch{}

	for i, sec := range cfg.Sections {
		matchSpec := sec.MatchSpecificity(gciImport)
		if matchSpec.IsMoreSpecific(bestMatchSpec) {
			bestMatchIdx = i
			bestMatchSpec = matchSpec
		}
	}

	if bestMatchIdx >= 0 {
		return bestMatchIdx
	}

	return len(cfg.Sections)
}

func defaultImportGroup(importPath string) int {
	if strings.HasPrefix(importPath, "appengine") {
		return 2
	}
	firstComponent := strings.Split(importPath, "/")[0]
	if strings.Contains(firstComponent, ".") {
		return 1
	}
	return 0
}

type byCommentPos []*ast.CommentGroup

func (x byCommentPos) Len() int           { return len(x) }
func (x byCommentPos) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x byCommentPos) Less(i, j int) bool { return x[i].Pos() < x[j].Pos() }
