package repository_test

import (
	"strings"
	"testing"

	"github.com/echsylon/go-args/internal/repository"
)

func Test_WhenDefiningArgumentWithMinCountZeroAndMaxCountGreaterThanMinCount_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(0, 2, "ARGUMENT", "", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithMinCountOneAndMaxCountGreaterThanMinCount_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(1, 2, "ARGUMENT", "", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndEqualToMaxCount_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(2, 2, "ARGUMENT", "", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndGreaterThanToMaxCount_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(3, 2, "ARGUMENT", "", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentStartingWithUpperCaseLetterFollowedByNumbersAndUnderscore_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(1, 1, "ABC123_", "", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentContainingNonAlphaNumericLetters_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(1, 1, "AB?123_", "", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(1, 1, "ARGUMENT", ".+", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(1, 1, "ARGUMENT", "", "description")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument(1, 1, "ARGUMENT", "*", "description")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningAlreadyExistingArgument_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 1, "ARGUMENT", "", "description")
	err := repo.DefineArgument(1, 1, "ARGUMENT", "", "description")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenAddingValueWithNoArgumentsDefined_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.AddArgumentValue("value")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenAddingValueWithMatchingPatternSaturatedArgument_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 1, "ARGUMENT", "[a-z0-9]+", "description")
	repo.AddArgumentValue("value1")
	err := repo.AddArgumentValue("value2")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenAddingValueWithNonMatchingNonSaturatedPatternArgument_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT", `\d`, "description")
	err := repo.AddArgumentValue("text")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenAddingValueWithMatchingPatternNonSaturatedArgument_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT", "[a-z]+", "description")
	err := repo.AddArgumentValue("value")
	if err != nil {
		t.Errorf("Expeted <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAddingValueWithEmptyPatternNonSaturatedArgument_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT", "", "description")
	err := repo.AddArgumentValue("value")
	if err != nil {
		t.Errorf("Expeted <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAddingValueWithMatchingPatternSaturatedAndNonSaturatedArguments_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 1, "ARGUMENT1", "[a-z]+", "description")
	repo.DefineArgument(1, 1, "ARGUMENT2", "", "description")
	repo.AddArgumentValue("value")
	err := repo.AddArgumentValue("value")
	if err != nil {
		t.Errorf("Expeted <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAddingValueWithMultipleMatchingPatternNonSaturatedArguments_ThenItIsOnlyAddedToTheFirstArgument(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 1, "ARGUMENT1", "[a-zA-Z]+", "description")
	repo.DefineArgument(1, 1, "ARGUMENT2", "[a-z]+", "description")
	repo.AddArgumentValue("value")
	count1 := len(repo.GetArgumentValues("ARGUMENT1"))
	count2 := len(repo.GetArgumentValues("ARGUMENT2"))
	if count1 != 1 || count2 != 0 {
		t.Errorf("Expeted <1> and <0>, but got <%d> and <%d>", count1, count2)
	}
}

func Test_WhenRequestingValuesForArguments_ThenOnlyMatchingValuesAreReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT1", "[a-z]+", "description")
	repo.DefineArgument(1, 2, "ARGUMENT2", "[0-9]+", "description")
	repo.AddArgumentValue("value")
	repo.AddArgumentValue("123")
	values := repo.GetArgumentValues("ARGUMENT1")
	actual := strings.Join(values, ", ")
	if actual != "value" {
		t.Errorf("Expected <value>, but got <%s>", actual)
	}
}

func Test_WhenRequestingValuesForNonExistingArgument_ThenEmptyResultSetIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT", "", "description")
	repo.AddArgumentValue("value")
	values := repo.GetArgumentValues("NON_EXISTING_ARGUMENT")
	actual := strings.Join(values, ", ")
	if actual != "" {
		t.Errorf("Expected <>, but got <%s>", actual)
	}
}

func Test_WhenAssertingValuesCountWithAllArgumentsSatisfied_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT1", `[a-z]+`, "description")
	repo.DefineArgument(1, 2, "ARGUMENT2", `\d`, "description")
	repo.AddArgumentValue("value")
	repo.AddArgumentValue("2")
	err := repo.AssertAllArgumentValuesProvided()
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAssertingValuesCountWithUnsatisfiedArguments_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument(1, 2, "ARGUMENT1", `[a-z]+`, "description")
	repo.DefineArgument(1, 2, "ARGUMENT2", `\d`, "description")
	repo.AddArgumentValue("2")
	err := repo.AssertAllArgumentValuesProvided()
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}
