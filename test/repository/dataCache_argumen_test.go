package repository_test

import (
	"strings"
	"testing"

	"github.com/echsylon/go-args/internal/repository"
)

func Test_WhenDefiningArgumentWithMinCountZeroAndMaxCountGreaterThanMinCount_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 0, 2, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithMinCountOneAndMaxCountGreaterThanMinCount_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 1, 2, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndEqualToMaxCount_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 2, 2, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndGreaterThanToMaxCount_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 3, 2, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentStartingWithUpperCaseLetterFollowedByNumbersAndUnderscore_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ABC123_", "description", 1, 1, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentContainingNonAlphaNumericLetters_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("AB?123_", "description", 1, 1, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 1, 1, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 1, 1, ".+")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	err := repo.DefineArgument("ARGUMENT", "description", 1, 1, "*")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningAlreadyExistingArgument_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT", "description", 1, 1, "")
	err := repo.DefineArgument("ARGUMENT", "description", 1, 1, "")
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
	repo.DefineArgument("ARGUMENT", "description", 1, 1, "[a-z0-9]+")
	repo.AddArgumentValue("value1")
	err := repo.AddArgumentValue("value2")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenAddingValueWithNonMatchingNonSaturatedPatternArgument_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT", "description", 1, 2, `\d`)
	err := repo.AddArgumentValue("text")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenAddingValueWithMatchingPatternNonSaturatedArgument_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT", "description", 1, 2, "[a-z]+")
	err := repo.AddArgumentValue("value")
	if err != nil {
		t.Errorf("Expeted <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAddingValueWithEmptyPatternNonSaturatedArgument_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT", "description", 1, 2, "")
	err := repo.AddArgumentValue("value")
	if err != nil {
		t.Errorf("Expeted <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAddingValueWithMatchingPatternSaturatedAndNonSaturatedArguments_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT1", "description", 1, 1, "[a-z]+")
	repo.DefineArgument("ARGUMENT2", "description", 1, 1, "")
	repo.AddArgumentValue("value")
	err := repo.AddArgumentValue("value")
	if err != nil {
		t.Errorf("Expeted <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAddingValueWithMultipleMatchingPatternNonSaturatedArguments_ThenItIsOnlyAddedToTheFirstArgument(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT1", "description", 1, 1, "[a-zA-Z]+")
	repo.DefineArgument("ARGUMENT2", "description", 1, 1, "[a-z]+")
	repo.AddArgumentValue("value")
	count1 := len(repo.GetArgumentValues("ARGUMENT1"))
	count2 := len(repo.GetArgumentValues("ARGUMENT2"))
	if count1 != 1 || count2 != 0 {
		t.Errorf("Expeted <1> and <0>, but got <%d> and <%d>", count1, count2)
	}
}

func Test_WhenRequestingValuesForArguments_ThenOnlyMatchingValuesAreReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT1", "description", 1, 2, "[a-z]+")
	repo.DefineArgument("ARGUMENT2", "description", 1, 2, "[0-9]+")
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
	repo.DefineArgument("ARGUMENT", "description", 1, 2, "")
	repo.AddArgumentValue("value")
	values := repo.GetArgumentValues("NON_EXISTING_ARGUMENT")
	actual := strings.Join(values, ", ")
	if actual != "" {
		t.Errorf("Expected <>, but got <%s>", actual)
	}
}

func Test_WhenAssertingValuesCountWithAllArgumentsSatisfied_ThenNoErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT1", "description", 1, 2, `[a-z]+`)
	repo.DefineArgument("ARGUMENT2", "description", 1, 2, `\d`)
	repo.AddArgumentValue("value")
	repo.AddArgumentValue("2")
	err := repo.AssertAllArgumentValuesProvided()
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenAssertingValuesCountWithUnsatisfiedArguments_ThenErrorIsReturned(t *testing.T) {
	repo := repository.NewDataCache()
	repo.DefineArgument("ARGUMENT1", "description", 1, 2, `[a-z]+`)
	repo.DefineArgument("ARGUMENT2", "description", 1, 2, `\d`)
	repo.AddArgumentValue("2")
	err := repo.AssertAllArgumentValuesProvided()
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}
