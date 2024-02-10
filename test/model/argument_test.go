package model_test

import (
	"testing"

	"github.com/echsylon/go-args/internal/model"
)

func Test_WhenCreatingNewArgument_ThenItsMinValuesCountCanBeRetrievedUndistorted(t *testing.T) {
	expected := 1
	arg := model.NewArgument("name", "description", expected, 2, "")
	actual := arg.GetMinValuesCount()
	if actual != expected {
		t.Errorf("Expected <%d>, but got <%d>", expected, actual)
	}
}

func Test_WhenCreatingNewArgument_ThenItsMaxValuesCountCanBeRetrievedUndistorted(t *testing.T) {
	expected := 2
	arg := model.NewArgument("name", "description", 1, expected, "")
	actual := arg.GetMaxValuesCount()
	if actual != expected {
		t.Errorf("Expected <%d>, but got <%d>", expected, actual)
	}
}

func Test_WhenMaxCountIsGreaterThanOne_ThenExpectsMultipleValuesReturnsTrue(t *testing.T) {
	arg := model.NewArgument("ARG", "description", 1, 2, "")
	expected := arg.ExpectsMultipleValues()
	if !expected {
		t.Errorf("Expected <true>, but got <false>")
	}
}

func Test_WhenMaxCountIsExactlyOne_ThenExpectsMultipleValuesReturnsFalse(t *testing.T) {
	arg := model.NewArgument("ARG", "description", 1, 1, "")
	expected := arg.ExpectsMultipleValues()
	if expected {
		t.Errorf("Expected <false>, but got <true>")
	}
}

func Test_WhenCreatingNewArgument_ThenItsNameCanBeRetrievedUndistorted(t *testing.T) {
	expected := "ARG"
	arg := model.NewArgument(expected, "description", 1, 2, "")
	actual := arg.GetName()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenCreatingNewArgument_ThenItsValuePatternCanBeRetrievedUndistorted(t *testing.T) {
	expected := "[a-z]{2}"
	arg := model.NewArgument("ARG", "description", 1, 2, expected)
	actual := arg.GetPattern()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenCreatingNewArgument_ThenItsDescriptionCanBeRetrievedUndistorted(t *testing.T) {
	expected := "Some description text"
	arg := model.NewArgument("ARG", expected, 1, 2, "")
	actual := arg.GetDescription()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}
