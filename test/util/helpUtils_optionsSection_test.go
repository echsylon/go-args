package util_test

import (
	"strings"
	"testing"

	"github.com/echsylon/go-args/internal/model"
	"github.com/echsylon/go-args/internal/util"
)

func Test_WhenComposingOptionsHelpSectionWithNilPointerOptions_ThenEmptyStringIsReturned(t *testing.T) {
	expected := ""
	actual := util.GetOptionsHelpSection(nil)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingOptionsHelpSectionWithEmptyOptions_ThenEmptyStringIsReturned(t *testing.T) {
	expected := ""
	actual := util.GetOptionsHelpSection(&[]model.Option{})
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingOptionsHelpSectionWithoutShortName_ThenShortNameColumnIsNotIncluded(t *testing.T) {
	expected := "Options:\n  --name  description"
	options := []model.Option{model.NewOption("", "name", "description", "")}
	actual := util.GetOptionsHelpSection(&options)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingOptionsHelpSectionWithoutLongName_ThenLongNameColumnIsNotIncluded(t *testing.T) {
	expected := "Options:\n  -n  description"
	options := []model.Option{model.NewOption("n", "", "description", "")}
	actual := util.GetOptionsHelpSection(&options)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingOptionsHelpSectionWithShortAndLongNames_ThenBothColumnsAreIncluded(t *testing.T) {
	expected := "Options:\n  -n, --name  description"
	options := []model.Option{model.NewOption("n", "name", "description", "")}
	actual := util.GetOptionsHelpSection(&options)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenComposingOptionsHelpSectionWithMultipleOptions_ThenEachOptionIsIncludedOnItsOwnRow(t *testing.T) {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("Options:\n")
	stringBuilder.WriteString("  -n, --name   Name description\n")
	stringBuilder.WriteString("      --value  Value description\n")
	stringBuilder.WriteString("  -s           Status description")
	expected := stringBuilder.String()
	options := []model.Option{
		model.NewOption("n", "name", "Name description", ""),
		model.NewOption("", "value", "Value description", ""),
		model.NewOption("s", "", "Status description", ""),
	}
	actual := util.GetOptionsHelpSection(&options)
	if actual != expected {
		t.Errorf("Expected: <%s>, but got <%s>", expected, actual)
	}
}
