package section

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMatching(t *testing.T) {
	assert.True(t, errors.Is(MissingParameterClosingBracketsError, MissingParameterClosingBracketsError))
	assert.True(t, errors.Is(MoreThanOneOpeningQuotesError, MoreThanOneOpeningQuotesError))
	assert.True(t, errors.Is(SectionTypeDoesNotAcceptParametersError, SectionTypeDoesNotAcceptParametersError))
	assert.True(t, errors.Is(SectionTypeDoesNotAcceptPrefixError, SectionTypeDoesNotAcceptPrefixError))
	assert.True(t, errors.Is(SectionTypeDoesNotAcceptSuffixError, SectionTypeDoesNotAcceptSuffixError))
}

func TestErrorClass(t *testing.T) {
	subError := MissingParameterClosingBracketsError
	errorGroup := SectionParsingError{subError}
	assert.True(t, errors.Is(errorGroup, SectionParsingError{}))
	assert.True(t, errors.Is(errorGroup, subError))
	assert.True(t, errors.Is(errorGroup.Wrap("x"), SectionParsingError{}))
	assert.True(t, errors.Is(errorGroup.Wrap("x"), subError))
}
