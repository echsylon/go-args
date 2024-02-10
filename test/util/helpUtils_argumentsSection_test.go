package util_test

import (
	"strings"
	"testing"

	"github.com/echsylon/go-args/internal/model"
	"github.com/echsylon/go-args/internal/util"
)

func Test_WhenComposingArgumentsHelpSectionWithNilPointerArguments_ThenEmptyStringIsReturned(t *testing.T) {
	expected := ""
	actual := util.GetArgumentsHelpSection(nil)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingArgumentsHelpSectionWithEmptyArguments_ThenEmptyStringIsReturned(t *testing.T) {
	expected := ""
	actual := util.GetArgumentsHelpSection(&[]model.Argument{})
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingArgumentsHelpSectionWithMultipleArguments_ThenEachArgumentIsIncludedOnItsOwnRow(t *testing.T) {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("Arguments:\n")
	stringBuilder.WriteString("  ARGUMENT First argument description\n")
	stringBuilder.WriteString("  ARGS     Second argument description")
	expected := stringBuilder.String()
	arguments := []model.Argument{
		model.NewArgument("ARGUMENT", "First argument description", 1, 1, ""),
		model.NewArgument("ARGS", "Second argument description", 1, 2, ""),
	}
	actual := util.GetArgumentsHelpSection(&arguments)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}
