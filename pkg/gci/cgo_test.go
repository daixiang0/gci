package gci

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/daixiang0/gci/pkg/config"
)

func TestGroupedCgoImports(t *testing.T) {
	const cgoBlock = `import (
	// C preamble.
	"C"
	"fmt"

	_ "golang.org/x/example/hello/reverse"
)
`
	const otherImports = `import (
	"strings"
	"os"
)
`
	const sortedOtherImports = `import (
	"os"
	"strings"
)
`
	const body = `
func init() {
	var x C.int
	fmt.Print(x)
}
`

	for _, tt := range []struct {
		name string
		in   string
		want string
	}{
		{
			name: "C first",
			in:   cgoBlock,
			want: cgoBlock,
		},
		{
			name: "C last",
			in: `import (
	"fmt"

	_ "golang.org/x/example/hello/reverse"
	// C preamble.
	"C"
)
`,
			want: `import (
	"fmt"

	_ "golang.org/x/example/hello/reverse"
	// C preamble.
	"C"
)
`,
		},
		{
			name: "other imports before C block",
			in:   otherImports + "\n" + cgoBlock,
			want: cgoBlock + "\n" + sortedOtherImports,
		},
		{
			name: "other imports after C block",
			in:   cgoBlock + "\n" + otherImports,
			want: cgoBlock + "\n" + sortedOtherImports,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.ParseConfig(commonConfig)
			require.NoError(t, err)

			src := []byte("package main\n\n" + tt.in + body)
			want := "package main\n\n" + tt.want + body
			_, got, err := LoadFormat(src, "main.go", *cfg)
			require.NoError(t, err)
			require.Equal(t, want, string(got))

			_, again, err := LoadFormat(got, "main.go", *cfg)
			require.NoError(t, err)
			require.Equal(t, string(got), string(again), "formatting must be idempotent")
		})
	}
}
