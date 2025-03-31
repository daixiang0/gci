package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidNumberOfFiles(t *testing.T) {
	assert.ErrorIs(t, InvalidNumberOfFilesInAnalysis{1, 2}, InvalidNumberOfFilesInAnalysis{})
}
