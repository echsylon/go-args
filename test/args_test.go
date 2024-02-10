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

	args.DefineOptionStrict("shortName", "shortName", "description", "")
}

func Test_WhenRegisteringInvalidArgument_ThenPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("<Expected <panic>, but got nothing")
		}
	}()

	args.DefineArgumentStrict("ARG", "Description", 0, 1, "")
}

func Test_WhenGettingStringValueForNonRegisteredOption_ThenFallbackIsReturned(t *testing.T) {
	expected := "fallback"
	actual := args.GetOptionValue("unregisterd", expected)
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenGettingIntValueForNonIntOption_ThenFallbackIsReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expected := int64(12)
	os.Args = []string{"appName", "--value", "string"}

	args.Reset()
	args.DefineOption("value", "description")
	args.Parse()
	actual := args.GetOptionIntValue("value", expected)

	if actual != expected {
		t.Errorf("Expected <%d>, but got <%d>", expected, actual)
	}
}

func Test_WhenGettingFloatValueForNonFloatOption_ThenFallbackIsReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expected := float64(11.9)
	os.Args = []string{"appName", "--value", "true"}

	args.Reset()
	args.DefineOption("value", "description")
	args.Parse()
	actual := args.GetOptionFloatValue("value", expected)

	if actual != expected {
		t.Errorf("Expected <%f>, but got <%f>", expected, actual)
	}
}

func Test_WhenGettingBoolValueForNonBoolOption_ThenFallbackIsReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expected := true
	os.Args = []string{"appName", "--value", "13.1"}

	args.Reset()
	args.DefineOption("value", "description")
	args.Parse()
	actual := args.GetOptionBoolValue("value", expected)

	if actual != expected {
		t.Errorf("Expected <%t>, but got <%t>", expected, actual)
	}
}

func Test_WhenGettingValueForNonRegisteredArgument_ThenEmptyResultSetIsReturned(t *testing.T) {
	expected := 0
	actual := len(args.GetArgumentValues("UNREGISTERED"))
	if actual != expected {
		t.Errorf("Expected <%d>, but got <%d>", expected, actual)
	}
}

func Test_WhenGettingIntValuesForPartiallyIntArguments_ThenOnlyIntValuesAreReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	e := []int64{12, 14}
	os.Args = []string{"appName", "12", "13.1", "14"}

	args.Reset()
	args.DefineArgumentStrict("arg", "description", 1, 3, "")
	args.Parse()
	a := args.GetArgumentIntValues("arg")

	if len(a) != len(e) || a[0] != e[0] || a[1] != e[1] {
		t.Errorf("Expected <%v>, but got <%v>", e, a)
	}
}

func Test_WhenGettingFloatValuesForPartiallyFloatArguments_ThenOnlyFloatValuesAreReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	e := []float64{3.1, 1.4}
	os.Args = []string{"appName", "true", "3.1", "1.4"}

	args.Reset()
	args.DefineArgumentStrict("arg", "description", 1, 3, "")
	args.Parse()
	a := args.GetArgumentFloatValues("arg")

	if len(a) != len(e) || a[0] != e[0] || a[1] != e[1] {
		t.Errorf("Expected <%v>, but got <%v>", e, a)
	}
}

func Test_WhenGettingBoolValuesForPartiallyBoolArguments_ThenOnlyBoolValuesAreReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	e := []bool{true, true}
	os.Args = []string{"appName", "True", "1", "1.4"}

	args.Reset()
	args.DefineArgumentStrict("arg", "description", 1, 3, "")
	args.Parse()
	a := args.GetArgumentBoolValues("arg")

	if len(a) != len(e) || a[0] != e[0] || a[1] != e[1] {
		t.Errorf("Expected <%v>, but got <%v>", e, a)
	}
}

func Test_WhenParsingRegisteredOptionWithNoValue_ThenBooleanTrueCanBeRetrievedForIt(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expected := true
	os.Args = []string{"appName", "--value"}

	args.Reset()
	args.DefineOption("value", "description")
	args.Parse()
	actual := args.GetOptionBoolValue("value", false)

	if actual != expected {
		t.Errorf("Expected <%t>, but got <%t>", expected, actual)
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
	actual := args.GetOptionValue("v", "fallback")

	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenDefiningArgumentsWithDifferentConstraints_ThenTheValuesDoNotGetMixedUp(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "file1.txt", "2000", "file2.txt"}

	args.Reset()
	args.DefineArgumentStrict("TIMEOUT", "description", 1, 1, `^\d+$`)
	args.DefineArgumentStrict("FILES", "description", 1, 2, `^*\.txt$`)
	args.Parse()
	n := args.GetArgumentIntValues("TIMEOUT")
	s := args.GetArgumentValues("FILES")

	if len(n) != 1 || n[0] != 2000 || len(s) != 2 || s[0] != "file1.txt" || s[1] != "file2.txt" {
		t.Errorf("Expected <[2000]> and <[file1.txt file2.txt]>, but got <%v> and <%v>", n, s)
	}
}

func Test_WhenHavingMultipleMatchingArguments_ThenTheValuesAreDistributedInOrderOfDefinition(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "A", "b"}

	args.Reset()
	args.DefineArgument("LOWER-CASE", "description")
	args.DefineArgument("UPPER-CASE", "description")
	args.Parse()
	upper := args.GetArgumentValues("UPPER-CASE")[0]
	lower := args.GetArgumentValues("LOWER-CASE")[0]

	if lower != "A" || upper != "b" {
		t.Errorf("Expected <[A]> and <[b]>, but got <%v> and <%v>", lower, upper)
	}
}

func Test_WhenHavingMultipleMatchingOptionsAndArguments_ThenValuesAreValidatedInOrderOfDefinition(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "a", "-s", "b", "c"}

	args.Reset()
	args.DefineOption("s", "description")
	args.DefineArgument("arg1", "description")
	args.DefineArgument("arg2", "description")
	args.Parse()
	opt := args.GetOptionValue("s", "empty")
	arg1 := args.GetArgumentValues("arg1")[0]
	arg2 := args.GetArgumentValues("arg2")[0]

	if opt != "b" || arg1 != "a" || arg2 != "c" {
		t.Errorf("Expected <b> and <a> and <c>, but got <%s> and <%s> and <%s>", opt, arg1, arg2)
	}
}
