package model_test

import (
	"testing"

	"github.com/echsylon/go-args/internal/model"
)

func Test_WhenParsedFlagIsNotSet_ThenIsParsedReturnsFalse(t *testing.T) {
	opt := model.NewOption("n", "name", "description", "")
	is := opt.IsParsed()
	if is {
		t.Errorf("Expected <false>, but got <true>")
	}
}

func Test_WhenSettingTheParsedFlag_ThenIsParsedReturnsTrue(t *testing.T) {
	opt := model.NewOption("n", "name", "description", "")
	opt.SetParsed()
	is := opt.IsParsed()
	if !is {
		t.Errorf("Expected <true>, but got <false>")
	}
}

func Test_WhenCreatingNewOption_ThenItsShortNameCanBeRetrievedUndistorted(t *testing.T) {
	expected := "n"
	opt := model.NewOption(expected, "name", "description", "")
	actual := opt.GetShortName()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenCreatingNewOption_ThenItsLongNameCanBeRetrievedUndistorted(t *testing.T) {
	expected := "name"
	opt := model.NewOption("n", expected, "description", "")
	actual := opt.GetLongName()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenCreatingNewOption_ThenItsValuePatternCanBeRetrievedUndistorted(t *testing.T) {
	expected := "[a-z]{2}"
	option := model.NewOption("n", "name", "description", expected)
	actual := option.GetPattern()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenCreatingNewOption_ThenItsDescriptionCanBeRetrievedUndistorted(t *testing.T) {
	expected := "Some kind of description"
	opt := model.NewOption("n", "name", expected, "")
	actual := opt.GetDescription()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}
