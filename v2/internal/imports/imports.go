package imports

import (
	"bufio"
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/daixiang0/gci/v2/pkg/config"
)

type Options struct {
	Config *config.Config

	Fragment  bool
	AllErrors bool

	Comments  bool
	TabIndent bool
	TabWidth  int

	FormatOnly bool
}

func Process(filename string, src []byte, opt *Options) (formatted []byte, err error) {
	fileSet := token.NewFileSet()
	file, adjust, err := parse(fileSet, filename, src, opt)
	if err != nil {
		return nil, err
	}

	return formatFile(fileSet, file, src, adjust, opt)
}

func formatFile(fset *token.FileSet, file *ast.File, src []byte, adjust func(orig []byte, src []byte) []byte, opt *Options) ([]byte, error) {
	moveCgoDeclsToTop(file)
	mergeImports(file)
	sortImports(opt.Config, fset.File(file.Pos()), file)

	var spacesBefore []string
	for _, impSection := range astutil.Imports(fset, file) {
		lastGroup := -1
		for _, importSpec := range impSection {
			importPath, _ := strconv.Unquote(importSpec.Path.Value)
			groupNum := importGroup(opt.Config, importPath, importSpec)
			if groupNum != lastGroup && lastGroup != -1 {
				spacesBefore = append(spacesBefore, importPath)
			}
			lastGroup = groupNum
		}
	}

	printerMode := printer.UseSpaces
	if opt.TabIndent {
		printerMode |= printer.TabIndent
	}
	printConfig := &printer.Config{Mode: printerMode, Tabwidth: opt.TabWidth}

	var buf bytes.Buffer
	err := printConfig.Fprint(&buf, fset, file)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	if len(spacesBefore) > 0 {
		out, err = addImportSpaces(bytes.NewReader(out), spacesBefore)
		if err != nil {
			return nil, err
		}
	}

	out = normalizePackageImportSpacing(src, len(file.Imports), out)
	out = normalizeImportDeclSpacing(out)
	out = restoreInlineImportCommentSpacing(src, out)
	out = normalizeInlineImportCommentSpacing(src, out)
	out = fixImportCommentLayout(opt.Config, src, out)
	out = normalizeCgoCommentIndentation(out)

	if opt.FormatOnly {
		return out, nil
	}

	out, err = format.Source(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func normalizePackageImportSpacing(src []byte, importCount int, out []byte) []byte {
	needsBlank := importCount > 1 || hasBlankLineBetweenPackageAndImport(src)
	lines := bytes.Split(out, []byte("\n"))
	for i := 0; i+2 < len(lines); i++ {
		if bytes.HasPrefix(bytes.TrimSpace(lines[i]), []byte("package ")) {
			hasBlank := len(bytes.TrimSpace(lines[i+1])) == 0
			isImportNext := bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("import"))
			isImportAfterBlank := bytes.HasPrefix(bytes.TrimSpace(lines[i+2]), []byte("import"))
			if hasBlank && isImportAfterBlank && !needsBlank {
				merged := append([][]byte{}, lines[:i+1]...)
				merged = append(merged, lines[i+2:]...)
				return bytes.Join(merged, []byte("\n"))
			}
			if !hasBlank && isImportNext && needsBlank {
				merged := append([][]byte{}, lines[:i+1]...)
				merged = append(merged, []byte(""))
				merged = append(merged, lines[i+1:]...)
				return bytes.Join(merged, []byte("\n"))
			}
			break
		}
	}
	return out
}

func hasBlankLineBetweenPackageAndImport(src []byte) bool {
	lines := bytes.Split(src, []byte("\n"))
	seenPackage := false
	for i := 0; i < len(lines); i++ {
		trimmed := bytes.TrimSpace(lines[i])
		if !seenPackage {
			if bytes.HasPrefix(trimmed, []byte("package ")) {
				seenPackage = true
			}
			continue
		}
		if bytes.HasPrefix(trimmed, []byte("import")) {
			return false
		}
		if len(trimmed) == 0 {
			return true
		}
	}
	return false
}

func normalizeImportDeclSpacing(out []byte) []byte {
	lines := bytes.Split(out, []byte("\n"))
	var result [][]byte
	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmed, []byte("import")) {
			result = append(result, line)
			i++
			continue
		}
		start := i
		end := i + 1
		if bytes.HasPrefix(trimmed, []byte("import (")) {
			for end < len(lines) && !bytes.HasPrefix(bytes.TrimSpace(lines[end]), []byte(")")) {
				end++
			}
			if end < len(lines) {
				end++
			}
		}
		result = append(result, lines[start:end]...)
		i = end
		j := i
		for j < len(lines) && len(bytes.TrimSpace(lines[j])) == 0 {
			j++
		}
		if j < len(lines) && bytes.HasPrefix(bytes.TrimSpace(lines[j]), []byte("import")) {
			if len(result) == 0 || len(bytes.TrimSpace(result[len(result)-1])) != 0 {
				result = append(result, []byte(""))
			}
			i = j
		}
	}
	return bytes.Join(result, []byte("\n"))
}

func restoreInlineImportCommentSpacing(src []byte, out []byte) []byte {
	paths := extractNoSpaceImportPaths(src)
	if len(paths) == 0 {
		return out
	}
	lines := bytes.Split(out, []byte("\n"))
	for i, line := range lines {
		idx := bytes.Index(line, []byte("//"))
		if idx <= 0 || line[idx-1] != ' ' {
			continue
		}
		path := extractImportPath(line)
		if path == "" {
			continue
		}
		if _, ok := paths[path]; !ok {
			continue
		}
		lines[i] = append(line[:idx-1], line[idx:]...)
	}
	return bytes.Join(lines, []byte("\n"))
}

func fixImportCommentLayout(cfg *config.Config, src []byte, out []byte) []byte {
	docPaths := extractDocImportPaths(src)
	if len(docPaths) == 0 {
		return out
	}
	trailingPaths := extractTrailingImportPaths(src)
	lines := bytes.Split(out, []byte("\n"))
	inImports := false
	lastGroup := -1
	for i := 0; i < len(lines); i++ {
		trimmed := bytes.TrimSpace(lines[i])
		if !inImports && bytes.HasPrefix(trimmed, []byte("import")) {
			inImports = true
		}
		if inImports && bytes.HasPrefix(trimmed, []byte(")")) {
			inImports = false
			continue
		}
		if !inImports {
			continue
		}
		m := impLine.FindStringSubmatch(string(lines[i]))
		if m == nil {
			continue
		}
		path := m[1]
		group := importGroup(cfg, path, buildImportSpecFromLine(lines[i], path))
		needsBreak := lastGroup != -1 && group != lastGroup
		if _, ok := docPaths[path]; !ok {
			lastGroup = group
			continue
		}
		if idx := bytes.Index(lines[i], []byte("//")); idx >= 0 {
			if _, ok := trailingPaths[path]; ok {
				lastGroup = group
				continue
			}
			indent := leadingWhitespace(lines[i])
			comment := bytes.TrimSpace(lines[i][idx:])
			line := bytes.TrimRight(lines[i][:idx], " \t")
			lines[i] = line
			newLine := append(append([]byte{}, indent...), comment...)
			lines = append(lines[:i], append([][]byte{newLine}, lines[i:]...)...)
			i++
		}
		start := i - 1
		for start >= 0 {
			t := bytes.TrimSpace(lines[start])
			if bytes.HasPrefix(t, []byte("//")) || bytes.HasPrefix(t, []byte("/*")) || bytes.HasPrefix(t, []byte("*")) || bytes.HasPrefix(t, []byte("*/")) {
				start--
				continue
			}
			break
		}
		blockStart := start + 1
		if needsBreak {
			if blockStart > 0 && len(bytes.TrimSpace(lines[blockStart-1])) != 0 {
				lines = append(lines[:blockStart], append([][]byte{[]byte("")}, lines[blockStart:]...)...)
				i++
			}
		} else {
			if blockStart > 0 && len(bytes.TrimSpace(lines[blockStart-1])) == 0 {
				lines = append(lines[:blockStart-1], lines[blockStart:]...)
				i--
			}
		}
		lastGroup = group
	}
	return bytes.Join(lines, []byte("\n"))
}

func normalizeCgoCommentIndentation(out []byte) []byte {
	lines := bytes.Split(out, []byte("\n"))
	for i := 0; i < len(lines); i++ {
		trimmed := bytes.TrimSpace(lines[i])
		if !bytes.HasPrefix(trimmed, []byte("import (")) {
			continue
		}
		end := i + 1
		for end < len(lines) && !bytes.HasPrefix(bytes.TrimSpace(lines[end]), []byte(")")) {
			end++
		}
		if end >= len(lines) {
			continue
		}
		hasC := false
		for k := i + 1; k < end; k++ {
			if bytes.Contains(lines[k], []byte(`"C"`)) {
				hasC = true
				break
			}
		}
		if !hasC {
			i = end
			continue
		}
		inComment := false
		var prefix []byte
		for k := i + 1; k < end; k++ {
			line := lines[k]
			if !inComment {
				start := bytes.Index(line, []byte("/*"))
				if start >= 0 {
					endIdx := bytes.Index(line, []byte("*/"))
					if endIdx == -1 || endIdx < start {
						prefix = leadingWhitespace(line)
						inComment = true
					}
				}
				continue
			}
			trimmedLine := bytes.TrimLeft(line, " \t")
			lines[k] = append(append([]byte{}, prefix...), trimmedLine...)
			if bytes.Contains(trimmedLine, []byte("*/")) {
				inComment = false
			}
		}
		i = end
	}
	return bytes.Join(lines, []byte("\n"))
}

func extractDocImportPaths(src []byte) map[string]struct{} {
	lines := bytes.Split(src, []byte("\n"))
	paths := map[string]struct{}{}
	inImports := false
	var pending bool
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if !inImports && bytes.HasPrefix(trimmed, []byte("import")) {
			inImports = true
		}
		if inImports && bytes.HasPrefix(trimmed, []byte(")")) {
			inImports = false
			continue
		}
		if !inImports {
			continue
		}
		if len(trimmed) == 0 {
			pending = false
			continue
		}
		if bytes.HasPrefix(trimmed, []byte("//")) || bytes.HasPrefix(trimmed, []byte("/*")) || bytes.HasPrefix(trimmed, []byte("*")) || bytes.HasPrefix(trimmed, []byte("*/")) {
			pending = true
			continue
		}
		m := impLine.FindStringSubmatch(string(line))
		if m == nil {
			pending = false
			continue
		}
		if pending {
			paths[m[1]] = struct{}{}
		}
		pending = false
	}
	return paths
}

func extractTrailingImportPaths(src []byte) map[string]struct{} {
	lines := bytes.Split(src, []byte("\n"))
	paths := map[string]struct{}{}
	inImports := false
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if !inImports && bytes.HasPrefix(trimmed, []byte("import")) {
			inImports = true
		}
		if inImports && bytes.HasPrefix(trimmed, []byte(")")) {
			inImports = false
			continue
		}
		if !inImports {
			continue
		}
		if bytes.Index(line, []byte("//")) == -1 {
			continue
		}
		path := extractImportPath(line)
		if path == "" {
			continue
		}
		paths[path] = struct{}{}
	}
	return paths
}

func leadingWhitespace(line []byte) []byte {
	i := 0
	for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
		i++
	}
	return line[:i]
}

func buildImportSpecFromLine(line []byte, path string) *ast.ImportSpec {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)}}
	}
	if trimmed[0] == '"' {
		return &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)}}
	}
	parts := strings.Fields(string(trimmed))
	if len(parts) >= 2 {
		return &ast.ImportSpec{
			Name: ast.NewIdent(parts[0]),
			Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)},
		}
	}
	return &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)}}
}

func normalizeInlineImportCommentSpacing(src []byte, out []byte) []byte {
	noSpace := extractNoSpaceImportPaths(src)
	lines := bytes.Split(out, []byte("\n"))
	for i, line := range lines {
		idx := bytes.Index(line, []byte("//"))
		if idx <= 0 || line[idx-1] != ' ' {
			continue
		}
		path := extractImportPath(line)
		if path == "" {
			continue
		}
		if _, ok := noSpace[path]; ok {
			continue
		}
		j := idx - 1
		for j-1 >= 0 && line[j-1] == ' ' {
			j--
		}
		if j < idx-1 {
			lines[i] = append(append([]byte{}, line[:j+1]...), line[idx:]...)
		}
	}
	return bytes.Join(lines, []byte("\n"))
}

func extractNoSpaceImportPaths(src []byte) map[string]struct{} {
	lines := bytes.Split(src, []byte("\n"))
	paths := map[string]struct{}{}
	for _, line := range lines {
		idx := bytes.Index(line, []byte("//"))
		if idx <= 0 {
			continue
		}
		if line[idx-1] == ' ' || line[idx-1] == '\t' {
			continue
		}
		path := extractImportPath(line)
		if path == "" {
			continue
		}
		paths[path] = struct{}{}
	}
	return paths
}

func extractImportPath(line []byte) string {
	first := bytes.IndexByte(line, '"')
	if first == -1 {
		return ""
	}
	second := bytes.IndexByte(line[first+1:], '"')
	if second == -1 {
		return ""
	}
	second += first + 1
	return string(line[first+1 : second])
}

func parse(fset *token.FileSet, filename string, src []byte, opt *Options) (*ast.File, func(orig, src []byte) []byte, error) {
	parserMode := parser.Mode(0)
	if opt.Comments {
		parserMode |= parser.ParseComments
	}
	if opt.AllErrors {
		parserMode |= parser.AllErrors
	}

	file, err := parser.ParseFile(fset, filename, src, parserMode)
	if err == nil {
		return file, nil, nil
	}
	if !opt.Fragment || !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	const prefix = "package main;"
	psrc := append([]byte(prefix), src...)
	file, err = parser.ParseFile(fset, filename, psrc, parserMode)
	if err == nil {
		psrc[len(prefix)-1] = '\n'
		fset.File(file.Package).SetLinesForContent(psrc)

		if containsMainFunc(file) {
			return file, nil, nil
		}

		adjust := func(orig, src []byte) []byte {
			src = src[len(prefix):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}

	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err = parser.ParseFile(fset, filename, fsrc, parserMode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			src = bytes.ReplaceAll(src, []byte("\n\t"), []byte("\n"))
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}

	return nil, nil, err
}

func containsMainFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				continue
			}

			if len(f.Type.Params.List) != 0 {
				continue
			}

			if f.Type.Results != nil && len(f.Type.Results.List) != 0 {
				continue
			}

			return true
		}
	}

	return false
}

func cutSpace(b []byte) (before, middle, after []byte) {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n') {
		i++
	}
	j := len(b)
	for j > 0 && (b[j-1] == ' ' || b[j-1] == '\t' || b[j-1] == '\n') {
		j--
	}
	if i <= j {
		return b[:i], b[i:j], b[j:]
	}
	return nil, nil, b[j:]
}

func matchSpace(orig []byte, src []byte) []byte {
	before, _, after := cutSpace(orig)
	i := bytes.LastIndex(before, []byte{'\n'})
	before, indent := before[:i+1], before[i+1:]

	_, src, _ = cutSpace(src)

	var b bytes.Buffer
	b.Write(before)
	for len(src) > 0 {
		line := src
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, src = line[:i+1], src[i+1:]
		} else {
			src = nil
		}
		if len(line) > 0 && line[0] != '\n' {
			b.Write(indent)
		}
		b.Write(line)
	}
	b.Write(after)
	return b.Bytes()
}

var impLine = regexp.MustCompile(`^\s+(?:[\w\.]+\s+)?"(.+?)"`)

func addImportSpaces(r io.Reader, breaks []string) ([]byte, error) {
	var out bytes.Buffer
	in := bufio.NewReader(r)
	inImports := false
	done := false
	var pendingComments []string
	for {
		s, err := in.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if !inImports && !done && strings.HasPrefix(s, "import") {
			inImports = true
		}
		if inImports && (strings.HasPrefix(s, "var") ||
			strings.HasPrefix(s, "func") ||
			strings.HasPrefix(s, "const") ||
			strings.HasPrefix(s, "type")) {
			done = true
			inImports = false
		}
		if inImports && isCommentOnlyLine(s) {
			pendingComments = append(pendingComments, s)
			continue
		}
		if inImports && len(breaks) > 0 {
			if m := impLine.FindStringSubmatch(s); m != nil {
				if m[1] == breaks[0] {
					out.WriteByte('\n')
					breaks = breaks[1:]
				}
			}
		}
		if len(pendingComments) > 0 {
			for _, c := range pendingComments {
				out.WriteString(c)
			}
			pendingComments = nil
		}

		out.WriteString(s)
	}
	if len(pendingComments) > 0 {
		for _, c := range pendingComments {
			out.WriteString(c)
		}
	}
	return out.Bytes(), nil
}

func isCommentOnlyLine(s string) bool {
	trimmed := strings.TrimSpace(s)
	return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "*/")
}
