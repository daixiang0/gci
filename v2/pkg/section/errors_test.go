package section

import (
	"errors"
	"testing"
)

func TestErrorMatching(t *testing.T) {
	if !errors.Is(MissingParameterClosingBracketsError, MissingParameterClosingBracketsError) {
		t.Fatal("expected error match for MissingParameterClosingBracketsError")
	}
	if !errors.Is(MoreThanOneOpeningQuotesError, MoreThanOneOpeningQuotesError) {
		t.Fatal("expected error match for MoreThanOneOpeningQuotesError")
	}
	if !errors.Is(SectionTypeDoesNotAcceptParametersError, SectionTypeDoesNotAcceptParametersError) {
		t.Fatal("expected error match for SectionTypeDoesNotAcceptParametersError")
	}
	if !errors.Is(SectionTypeDoesNotAcceptPrefixError, SectionTypeDoesNotAcceptPrefixError) {
		t.Fatal("expected error match for SectionTypeDoesNotAcceptPrefixError")
	}
	if !errors.Is(SectionTypeDoesNotAcceptSuffixError, SectionTypeDoesNotAcceptSuffixError) {
		t.Fatal("expected error match for SectionTypeDoesNotAcceptSuffixError")
	}
}

func TestErrorClass(t *testing.T) {
	subError := MissingParameterClosingBracketsError
	errorGroup := SectionParsingError{subError}
	if !errors.Is(errorGroup, SectionParsingError{}) {
		t.Fatal("expected SectionParsingError match")
	}
	if !errors.Is(errorGroup, subError) {
		t.Fatal("expected wrapped sub error match")
	}
	if !errors.Is(errorGroup.Wrap("x"), SectionParsingError{}) {
		t.Fatal("expected wrapped SectionParsingError match")
	}
	if !errors.Is(errorGroup.Wrap("x"), subError) {
		t.Fatal("expected wrapped sub error match after Wrap")
	}
}
