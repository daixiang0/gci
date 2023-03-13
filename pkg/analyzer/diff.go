package analyzer

import (
	"bytes"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/tools/go/analysis"
)

var hunkRE = regexp.MustCompile(`@@ -(\d+),(\d+) \+\d+,\d+ @@`)

func GetSuggestedFix(file *token.File, a, b []byte) (*analysis.SuggestedFix, error) {
	d := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(a)),
		B:       difflib.SplitLines(string(b)),
		Context: 1,
	}
	diff, err := difflib.GetUnifiedDiffString(d)
	if err != nil {
		return nil, err
	}
	if diff == "" {
		return nil, nil
	}
	var (
		fix   analysis.SuggestedFix
		found = false
		edit  analysis.TextEdit
		buf   bytes.Buffer
	)
	for _, line := range strings.Split(diff, "\n") {
		if hunk := hunkRE.FindStringSubmatch(line); len(hunk) > 0 {
			if found {
				edit.NewText = buf.Bytes()
				buf = bytes.Buffer{}
				fix.TextEdits = append(fix.TextEdits, edit)
				edit = analysis.TextEdit{}
			}
			found = true
			start, err := strconv.Atoi(hunk[1])
			if err != nil {
				return nil, err
			}
			lines, err := strconv.Atoi(hunk[2])
			if err != nil {
				return nil, err
			}
			edit.Pos = file.LineStart(start)
			end := start + lines
			if end > file.LineCount() {
				edit.End = token.Pos(file.Size())
			} else {
				edit.End = file.LineStart(end)
			}
			continue
		}
		// skip any lines until first hunk found
		if !found {
			continue
		}
		if line == "" {
			continue
		}
		switch line[0] {
		case '+':
			buf.WriteString(line[1:])
			buf.WriteRune('\n')
		case '-':
			// just skip
		case ' ':
			buf.WriteString(line[1:])
			buf.WriteRune('\n')
		}
	}
	edit.NewText = buf.Bytes()
	fix.TextEdits = append(fix.TextEdits, edit)

	return &fix, nil
}
