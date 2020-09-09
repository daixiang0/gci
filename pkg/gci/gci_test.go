package gci

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestProcessSource(t *testing.T) {
	asGoFile := func(imports string) []byte {
		return []byte(fmt.Sprintf(`package example

%s

func main() {
	fmt.Println("hello world")
}
`, strings.TrimSpace(imports)))
	}

	type testCase struct {
		name string

		localFlag string
		doWrite   bool
		doDiff    bool

		sourceImports   string
		expectedImports string
	}

	testCases := []testCase{
		{
			name:      "simple case",
			localFlag: "github.com/daixiang0",
			sourceImports: `
import (
	"golang.org/x/tools"

	"fmt"

	"github.com/daixiang0/gci"
)
`,
			expectedImports: `
import (
	"fmt"

	"golang.org/x/tools"

	"github.com/daixiang0/gci"
)
`,
		},
		{
			name:      "with aliases",
			localFlag: "github.com/daixiang0",
			sourceImports: `
import (
	go "github.com/golang"
	"fmt"
	"github.com/daixiang0"
)
`,
			expectedImports: `
import (
	"fmt"

	go "github.com/golang"

	"github.com/daixiang0"
)
`,
		},
		{
			name:      "with multiple aliases",
			localFlag: "github.com/daixiang0",
			sourceImports: `
import (
	a "github.com/golang"
	b "github.com/golang"
	"fmt"
	"github.com/daixiang0"
)
`,
			expectedImports: `
import (
	"fmt"

	a "github.com/golang"
	b "github.com/golang"

	"github.com/daixiang0"
)
`,
		},
		{
			name:      "with comment and aliases",
			localFlag: "github.com/daixiang0",
			sourceImports: `
import (
	"fmt"
	_ "github.com/golang" // golang
	"github.com/daixiang0"
)
`,
			expectedImports: `
import (
	"fmt"

	// golang
	_ "github.com/golang"

	"github.com/daixiang0"
)
`,
		},
		{
			name:      "with above comment and aliases",
			localFlag: "github.com/daixiang0",
			sourceImports: `
import (
	"fmt"
	// golang
	_ "github.com/golang"
	"github.com/daixiang0"
)
`,
			expectedImports: `
import (
	"fmt"

	// golang
	_ "github.com/golang"

	"github.com/daixiang0"
)
`,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			filename := "test_filename.go"
			source := asGoFile(c.sourceImports)
			expected := asGoFile(c.expectedImports)

			var buffer bytes.Buffer
			if err := processSource(filename, source, &buffer, &FlagSet{
				LocalFlag: c.localFlag,
				DoWrite:   &c.doWrite,
				DoDiff:    &c.doDiff,
			}); err != nil {
				t.Errorf("failed to process source: %v", err)
				t.FailNow()
			}
			actual := buffer.Bytes()

			if bytes.Compare(expected, actual) != 0 {
				diff, _ := diff(expected, actual, "")
				t.Errorf("unexpected output %s", diff)
				t.FailNow()
			}
		})
	}
}
