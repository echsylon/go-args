package data_test

import (
	"testing"

	"github.com/echsylon/go-args/internal/data"
)

func Test_WhenSavingAnOptionSuccessfully_ThenThatOptionCanBeRetrieved(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("n", "name", "description", ".*")
	option := repository.GetOption("n")
	if option == nil {
		t.Errorf("Expected <option>, but got <nil>")
	}
}

func Test_WhenSavingMultipleOptionsSuccessfully_ThenThoseOptionsCanAllBeRetrieved(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("o", "one", "description", "")
	repository.SaveOption("t", "two", "description", "")
	options := repository.GetOptions()
	count := len(options)
	if count != 2 {
		t.Errorf("Expected <2>, but got <%d>", count)
	}
}

func Test_WhenSavingAnOptionValueSuccessfully_ThenThatValueCanBeRetrieved(t *testing.T) {
	expected := "?value!"
	repository := data.NewRepository()
	repository.SaveOption("o", "option", "description", "")
	repository.SaveOptionValue("option", expected)
	actual := repository.GetOptionValue("o")
	if actual != expected {
		t.Errorf("Expected <%s>, but got <%s>", expected, actual)
	}
}

func Test_WhenSavingValueForNonExistingOption_ThenTheValueIsNotSaved(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("n", "name", "description", "")
	repository.SaveOptionValue("nonexisting", "some_value")
	actual := repository.GetOptionValue("nonexisting")
	if actual != "" {
		t.Errorf("Expected <>, but got <%s>", actual)
	}
}

func Test_WhenRequestingNonExistingValueForOption_ThenEmptyStringIsReturned(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("o", "opt", "description", "")
	actual := repository.GetOptionValue("o")
	if actual != "" {
		t.Errorf("Expected <>, but got <%s>", actual)
	}
}

func Test_WhenSavingAnArgumentSuccessfully_ThenThatArgumentCanBeRetrieved(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	argument := repository.GetArgument("ARG")
	if argument == nil {
		t.Errorf("Expected <argument>, but got <nil>")
	}
}

func Test_WhenSavingMultipleArgumentsSuccessfully_ThenThoseArgumentsCanBeRetrieved(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveArgument("ARG1", "description", 1, 1, ".*")
	repository.SaveArgument("ARG2", "description", 1, 2, "")
	arguments := repository.GetArguments()
	count := len(arguments)
	if count != 2 {
		t.Errorf("Expected <2>, but got <%d>", count)
	}
}

func Test_WhenSavingArgumentValuesSuccessfully_ThenAllThoseValuesCanBeRetrieved(t *testing.T) {
	expected := []string{"value-1", "value-2"}
	repository := data.NewRepository()
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	repository.SaveArgumentValue("ARG", expected[0])
	repository.SaveArgumentValue("ARG", expected[1])
	actual := repository.GetArgumentValues("ARG")
	if actual[0] != expected[0] || actual[1] != expected[1] {
		t.Errorf("Expected <%s> and <%s>, but got <%s> and <%s>", expected[0], expected[1], actual[0], actual[1])
	}
}

func Test_WhenSavingValueForNonExistingArgument_ThenTheValueIsNotSaved(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	repository.SaveArgumentValue("INVALID", "some_value")
	values := repository.GetArgumentValues("INVALID")
	count := len(values)
	if count != 0 {
		t.Errorf("Expected <0>, but got <%d>", count)
	}
}

func Test_WhenRequestingNonExistingValuesForArgument_ThenEmptySliceIsReturned(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	values := repository.GetArgumentValues("INVALID")
	count := len(values)
	if count != 0 {
		t.Errorf("Expected <0>, but got <%d>", count)
	}
}

func Test_WhenClearingAllValues_ThenAllOptionsAndArgumentsAreKept(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("o", "opt", "description", "")
	repository.SaveOptionValue("opt", "option_value")
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	repository.SaveArgumentValue("ARG", "argument_value")
	repository.ClearValues()
	option := repository.GetOption("o")
	argument := repository.GetArgument("ARG")
	if option == nil && argument == nil {
		t.Errorf("Expected <option> and <argument>, but got <nil> and <nil>")
	} else if option == nil {
		t.Errorf("Expected <option> and <argument>, but got <nil> and <argument>")
	} else if argument == nil {
		t.Errorf("Expected <option> and <argument>, but got <option> and <nil>")
	}
}

func Test_WhenClearingAllValues_ThenAllValuesAreDeleted(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("o", "opt", "description", "")
	repository.SaveOptionValue("opt", "optionValue")
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	repository.SaveArgumentValue("ARG", "argumentValue")
	repository.ClearValues()
	optionValue := repository.GetOptionValue("o")
	argumentValues := repository.GetArgumentValues("ARG")
	if optionValue != "" || len(argumentValues) != 0 {
		t.Errorf("Expected <> and <[]>, but got <%s> and <%v>", optionValue, argumentValues)
	}
}

func Test_WhenClearingAllData_ThenAllOptionsAndArgumentsAreDeleted(t *testing.T) {
	repository := data.NewRepository()
	repository.SaveOption("o", "opt", "description", "")
	repository.SaveOptionValue("opt", "option_value")
	repository.SaveArgument("ARG", "description", 1, 1, ".*")
	repository.SaveArgumentValue("ARG", "argument_value")
	repository.ClearAll()
	option := repository.GetOption("o")
	argument := repository.GetArgument("ARG")
	if option != nil && argument != nil {
		t.Errorf("Expected <nil> and <nil>, but got <%v> and <%v>", option, argument)
	}
}
