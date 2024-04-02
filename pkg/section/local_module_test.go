package section

import (
	"testing"

	"github.com/daixiang0/gci/pkg/specificity"
)

func TestLocalModule_specificity(t *testing.T) {
	testCases := []specificityTestData{
		{"example.com/hello", &LocalModule{Path: "example.com/hello"}, specificity.LocalModule{}},
		{"example.com/hello/world", &LocalModule{Path: "example.com/hello"}, specificity.LocalModule{}},
		{"example.com/hello-world", &LocalModule{Path: "example.com/hello"}, specificity.MisMatch{}},
		{"example.com/helloworld", &LocalModule{Path: "example.com/hello"}, specificity.MisMatch{}},
		{"example.com", &LocalModule{Path: "example.com/hello"}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
