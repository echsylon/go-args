package repository_test

import (
	"testing"

	"github.com/echsylon/go-args/internal/repository"
)

func Test_WhenDefiningOptionWithEmptyNames_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineOption("", "", "", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithTwoCharacterShortName_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineOption("sn", "longName", "", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithSingleCharacterLongName_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineOption("s", "l", "", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineOption("s", "short", "", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningOptionWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineOption("s", "short", ".+", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningOptionWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineOption("s", "short", "*", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningDuplicateOption_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err1 := repo.DefineOption("s", "longName", "", "description")
	err2 := repo.DefineOption("s", "longName", "", "description")
	if err1 != nil {
		t.Errorf("Couldn't create first option. Expected <nil> got <error>")
	}
	if err2 == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithShortName_ThenOptionWithThatShortNameCanBeRetreived(t *testing.T) {
	expected := "s"
	repo := repository.NewDataCache()
	repo.DefineOption(expected, "", "", "description")
	options := repo.GetOptions()
	actual := options[0].GetShortName()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenDefiningOptionWithLongName_ThenOptionWithThatLongNameCanBeRetreived(t *testing.T) {
	expected := "longName"
	repo := repository.NewDataCache()
	repo.DefineOption("", expected, "", "description")
	options := repo.GetOptions()
	actual := options[0].GetLongName()
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenSettingParsedFlagOfNonDefinedOption_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	err := repo.SetOptionParsed("u")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenSettingParsedFlagOfDefinedButPreviouslyUnparsedOption_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	err := repo.SetOptionParsed("test")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>")
	}
}

func Test_WhenSettingParsedFlagOfDefinedAndPreviouslyParsedOption_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	repo.SetOptionParsed("t")
	err := repo.SetOptionParsed("test")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>")
	}
}

func Test_WhenSettingValueOfNonDefinedOption_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	err := repo.SetOptionValue("u", "value")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenSettingDifferentValueOnDefinedOption_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	repo.SetOptionValue("t", "value")
	err := repo.SetOptionValue("test", "differentValue")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenSettingSameValueOnDefinedOption_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	repo.SetOptionValue("test", "value")
	err := repo.SetOptionValue("test", "value")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenRequestingNotSetValueFromParsedOption_ThenBooleanTrueIsReturnedAsString(t *testing.T) {
	expected := "true"
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	repo.SetOptionParsed("t")
	actual := repo.GetOptionValue("test")
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenRequestingSetValueFromExistingOption_ThenTheValueIsReturnedAsString(t *testing.T) {
	expected := "value"
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	repo.SetOptionValue("t", expected)
	actual := repo.GetOptionValue("test")
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenRequestingValueForNonExistingOption_ThenEmptyStringIsReturned(t *testing.T) {
	expected := ""
	repo := repository.NewDataCache()
	repo.DefineOption("t", "test", "", "description")
	repo.SetOptionValue("t", "value")
	actual := repo.GetOptionValue("nonExistingOption")
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}
