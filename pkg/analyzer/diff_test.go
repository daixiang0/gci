package analyzer_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"

	"github.com/daixiang0/gci/pkg/analyzer"
)

const formattedFile = `package analyzer

import (
	"fmt"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
)
`

func TestGetSuggestedFix(t *testing.T) {
	for _, tt := range []struct {
		name            string
		unformattedFile string
		expectedFix     *analysis.SuggestedFix
		expectedErr     string
	}{
		{
			name:            "same files",
			unformattedFile: formattedFile,
		},
		{
			name: "one change",
			unformattedFile: `package analyzer

import (
	"fmt"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"

	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
)
`,
			expectedFix: &analysis.SuggestedFix{
				TextEdits: []analysis.TextEdit{
					{
						Pos: 133,
						End: 205,
						NewText: []byte(`	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
`,
						),
					},
				},
			},
		},
		{
			name: "multiple changes",
			unformattedFile: `package analyzer

import (
	"fmt"
	"go/token"

	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"

	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
)
`,
			expectedFix: &analysis.SuggestedFix{
				TextEdits: []analysis.TextEdit{
					{
						Pos: 35,
						End: 59,
						NewText: []byte(`	"go/token"
	"strings"
`,
						),
					},
					{
						Pos: 134,
						End: 206,
						NewText: []byte(`	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
`,
						),
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "analyzer.go", tt.unformattedFile, 0)
			assert.NoError(t, err)

			actualFix, err := analyzer.GetSuggestedFix(fset.File(f.Pos()), []byte(tt.unformattedFile), []byte(formattedFile))
			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedFix, actualFix)
		})
	}
}
