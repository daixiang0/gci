package section

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorMatching(t *testing.T) {
	require.ErrorIs(t, MissingParameterClosingBracketsError, MissingParameterClosingBracketsError)
	require.ErrorIs(t, MoreThanOneOpeningQuotesError, MoreThanOneOpeningQuotesError)
	require.ErrorIs(t, SectionTypeDoesNotAcceptParametersError, SectionTypeDoesNotAcceptParametersError)
	require.ErrorIs(t, SectionTypeDoesNotAcceptPrefixError, SectionTypeDoesNotAcceptPrefixError)
	require.ErrorIs(t, SectionTypeDoesNotAcceptSuffixError, SectionTypeDoesNotAcceptSuffixError)
}

func TestErrorClass(t *testing.T) {
	subError := MissingParameterClosingBracketsError
	errorGroup := SectionParsingError{subError}
	require.ErrorIs(t, errorGroup, SectionParsingError{})
	require.ErrorIs(t, errorGroup, subError)
	require.ErrorIs(t, errorGroup.Wrap("x"), SectionParsingError{})
	require.ErrorIs(t, errorGroup.Wrap("x"), subError)
}
