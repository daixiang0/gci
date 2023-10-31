package gci

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/log"
)

func init() {
	log.InitLogger()
	defer log.L().Sync()
}

func TestRun(t *testing.T) {
	// if runtime.GOOS == "windows" {
	// 	t.Skip("Skipping test on Windows")
	// }

	for i := range testCases {
		t.Run(fmt.Sprintf("run case: %s", testCases[i].name), func(t *testing.T) {
			config, err := config.ParseConfig(testCases[i].config)
			if err != nil {
				t.Fatal(err)
			}

			old, new, err := LoadFormat([]byte(testCases[i].in), "", *config)
			if err != nil {
				t.Fatal(err)
			}

			assert.NoError(t, err)
			assert.Equal(t, testCases[i].in, string(old))
			assert.Equal(t, testCases[i].out, string(new))
		})
	}
}
