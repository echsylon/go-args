package util_test

import (
	"fmt"
	"testing"

	"github.com/echsylon/go-args/internal/model"
	"github.com/echsylon/go-args/internal/util"
)

func Test_WhenComposingMainHelpSectionWithoutOptionsAndArguments_ThenOnlyTheAppNameAndDescriptionIsIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	expected := fmt.Sprintf("Usage: %s\n%s", appName, appDescr)
	actual := util.GetMainHelpSection(appName, appDescr, nil, nil)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}

func Test_WhenComposingMainHelpSectionWithSingleOption_ThenOptionNotationIsIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	options := []model.Option{model.NewOption("n", "name", "", "descr")}
	expected := fmt.Sprintf("Usage: %s [OPTION]\n%s", appName, appDescr)
	actual := util.GetMainHelpSection(appName, appDescr, &options, nil)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}

func Test_WhenComposingMainHelpSectionWithMultipleOptions_ThenOptionsNotationWithEllipsisIsIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	options := []model.Option{
		model.NewOption("n", "name", "", "descr"),
		model.NewOption("", "other", "", "descr")}
	expected := fmt.Sprintf("Usage: %s [OPTIONS...]\n%s", appName, appDescr)
	actual := util.GetMainHelpSection(appName, appDescr, &options, nil)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}

func Test_WhenComposingMainHelpSectionWithSingleValueArgument_ThenArgumentNameIsIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	argName := "ARG"
	arguments := []model.Argument{model.NewArgument(1, 1, argName, "", "descr")}
	expected := fmt.Sprintf("Usage: %s %s\n%s", appName, argName, appDescr)
	actual := util.GetMainHelpSection("app", "description", nil, &arguments)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}

func Test_WhenComposingMainHelpSectionWithMultiValueArgument_ThenArgumentNameWithEllipsisIsIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	argName := "ARG"
	arguments := []model.Argument{model.NewArgument(1, 2, argName, "", "descr")}
	expected := fmt.Sprintf("Usage: %s %s...\n%s", appName, argName, appDescr)
	actual := util.GetMainHelpSection(appName, appDescr, nil, &arguments)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}

func Test_WhenComposingMainHelpSectionWithTwoArguments_ThenBothArgumentNamesAreIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	argName1 := "ARG1"
	argName2 := "ARG2"
	arguments := []model.Argument{
		model.NewArgument(1, 1, argName1, "", "descr"),
		model.NewArgument(1, 1, argName2, "", "descr")}
	expected := fmt.Sprintf("Usage: %s %s %s\n%s", appName, argName1, argName2, appDescr)
	actual := util.GetMainHelpSection("app", "description", nil, &arguments)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}

func Test_WhenComposingMainHelpSectionWithOptionsAndArguments_ThenBothOptionNotationAndArgumentNameIsIncluded(t *testing.T) {
	appName := "app"
	appDescr := "description"
	argName := "ARG"
	options := []model.Option{model.NewOption("n", "name", "", "descr")}
	arguments := []model.Argument{model.NewArgument(1, 1, argName, "", "descr")}
	expected := fmt.Sprintf("Usage: %s [OPTION] %s\n%s", appName, argName, appDescr)
	actual := util.GetMainHelpSection("app", "description", &options, &arguments)
	if actual != expected {
		t.Errorf("Expected:\n<%s>, but got\n<%s>", expected, actual)
	}
}
