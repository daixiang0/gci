package gci

import (
	"errors"
	"testing"

	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"

	"github.com/stretchr/testify/assert"
)

func TestErrorMatching(t *testing.T) {
	section := sectionsPkg.DefaultSection{}
	importDef := importPkg.ImportDef{"abc", "abc.com", []string{}, ""}
	assert.True(t, errors.Is(EqualSpecificityMatchError{importDef, section, section}, EqualSpecificityMatchError{}))
	assert.True(t, errors.Is(NoMatchingSectionForImportError{importDef}, NoMatchingSectionForImportError{}))
	assert.True(t, errors.Is(InvalidImportSplitError{[]string{"a"}}, InvalidImportSplitError{}))
	assert.True(t, errors.Is(InvalidAliasSplitError{[]string{"a"}}, InvalidAliasSplitError{}))
	assert.True(t, errors.Is(MissingImportStatementError, MissingImportStatementError))
	assert.True(t, errors.Is(ImportStatementNotClosedError, ImportStatementNotClosedError))
}

func TestErrorClass(t *testing.T) {
	subError := errors.New("test")
	errorGroup := FileParsingError{subError}
	assert.True(t, errors.Is(errorGroup, FileParsingError{}))
	assert.True(t, errors.Is(errorGroup, subError))
	assert.True(t, errors.Is(errorGroup, MissingImportStatementError))
	assert.True(t, errors.Is(errorGroup, ImportStatementNotClosedError))
	// unavoidable with the current implementation
	assert.True(t, errors.Is(MissingImportStatementError, ImportStatementNotClosedError))
}
