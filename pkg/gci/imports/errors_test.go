package imports

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMatching(t *testing.T) {
	assert.True(t, errors.Is(InvalidCharacterError{'a', "abc"}, InvalidCharacterError{}))
	assert.True(t, errors.Is(ValidationError{MissingOpeningQuotesError}, ValidationError{}))
	assert.True(t, errors.Is(MissingOpeningQuotesError, MissingOpeningQuotesError))
	assert.True(t, errors.Is(MissingClosingQuotesError, MissingClosingQuotesError))
}

func TestErrorClass(t *testing.T) {
	subError := errors.New("test")
	errorGroup := ValidationError{subError}
	assert.True(t, errors.Is(errorGroup, ValidationError{}))
	assert.True(t, errors.Is(errorGroup, subError))
	assert.True(t, errors.Is(MissingOpeningQuotesError, ValidationError{}))
	assert.True(t, errors.Is(MissingClosingQuotesError, ValidationError{}))
}
