package main

import (
	"bytes"
	"testing"
)

var (
	data = []string{
		`
	"fmt"
	"github.com/golang"
	"github.com/daixiang0"`,

		`
	"fmt"
	go "github.com/golang"
	"github.com/daixiang0"`,

		`
	"fmt"
	_ "github.com/golang" // golang
	"github.com/daixiang0"`,
	}

	result = []string{
`	"fmt"

	"github.com/golang"

	"github.com/daixiang0"
`,

`	"fmt"

	go "github.com/golang"

	"github.com/daixiang0"
`,

`	"fmt"

	// golang
	_ "github.com/golang"

	"github.com/daixiang0"
`,
	}
)

func Test(t *testing.T) {
	for i := 0; i < len(data); i++ {
		p := newPkg(bytes.Split([]byte(data[i]), []byte(linebreak)), "github.com/daixiang0")
		ret := p.fmt()
		if !bytes.Equal([]byte(result[i]), ret) {
			t.Fatalf("Test %d Failed\n====want===\n%s\n====get=====\n%s", i, result[i], ret)
		}
	}
}
