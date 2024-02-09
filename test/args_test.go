package args

import (
	"os"
	"testing"

	"github.com/echsylon/go-args"
)

func Test_WhenRegisteringInvalidOption_ThenPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("<Expected <panic>, but got nothing")
		}
	}()

	args.DefineConstrainedOption("shortName", "shortName", "description", "")
}

func Test_WhenRegisteringInvalidArgument_ThenPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("<Expected <panic>, but got nothing")
		}
	}()

	args.DefineConstrainedArgument("ARG", "Description", "", 0, 1)
}

func Test_WhenGettingValueForNonRegisteredOption_ThenEmptyStringIsReturned(t *testing.T) {
	expected := ""
	actual := args.GetOptionValue("unregisterd")
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenGettingValueForNonRegisteredArgument_ThenEmptyResultSetIsReturned(t *testing.T) {
	expected := 0
	actual := len(args.GetArgumentValues("UNREGISTERED"))
	if actual != expected {
		t.Errorf("Expected <%d>, but got <%d>", expected, actual)
	}
}

func Test_WhenParsingRegisteredOptionWithNoValue_ThenBooleanTrueCanBeRetrievedForIt(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expected := "true"
	os.Args = []string{"appName", "--value"}

	args.Reset()
	args.DefineOption("value", "description")
	args.Parse()
	actual := args.GetOptionValue("value")

	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenParsingRegisteredOptionWithValue_ThenTheValueCanBeRetrievedUndistorted(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expected := "some_value"
	os.Args = []string{"appName", "-v", expected}

	args.Reset()
	args.DefineOption("v", "description")
	args.Parse()
	actual := args.GetOptionValue("v")

	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}
