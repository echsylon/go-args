package domain_test

import (
	"os"
	"testing"

	"github.com/echsylon/go-args/internal/domain"
	"github.com/echsylon/go-args/internal/model"
	"github.com/echsylon/go-args/internal/repository"
)

type mockDataCache struct {
	repository.DataCache
	argumentsProvider     func() []model.Argument
	argumentProvider      func() model.Argument
	argumentValueProvider func() []string
	argumentValueListener func(string, string)
	optionsProvider       func() []model.Option
	optionProvider        func() model.Option
	optionValueProvider   func() string
	optionValueListener   func(string, string)
}

func (mock *mockDataCache) DefineArgument(string, string, int, int, string) {}
func (mock *mockDataCache) GetArgument(string) model.Argument               { return mock.argumentProvider() }
func (mock *mockDataCache) GetArguments() []model.Argument                  { return mock.argumentsProvider() }
func (mock *mockDataCache) GetArgumentValues(string) []string               { return mock.argumentValueProvider() }
func (mock *mockDataCache) AddArgumentValue(k string, v string)             { mock.argumentValueListener(k, v) }
func (mock *mockDataCache) DefineOption(string, string, string, string)     {}
func (mock *mockDataCache) GetOption(string) model.Option                   { return mock.optionProvider() }
func (mock *mockDataCache) GetOptions() []model.Option                      { return mock.optionsProvider() }
func (mock *mockDataCache) ClearValues()                                    {}
func (mock *mockDataCache) SetOptionValue(k string, v string)               { mock.optionValueListener(k, v) }
func (mock *mockDataCache) GetOptionValue(string) string                    { return mock.optionValueProvider() }

func newEmptyMockDataCache() *mockDataCache {
	return &mockDataCache{
		argumentsProvider:     func() []model.Argument { return []model.Argument{} },
		argumentProvider:      func() model.Argument { return nil },
		argumentValueProvider: func() []string { return []string{} },
		argumentValueListener: func(string, string) {},
		optionsProvider:       func() []model.Option { return []model.Option{} },
		optionProvider:        func() model.Option { return nil },
		optionValueProvider:   func() string { return "" },
		optionValueListener:   func(string, string) {},
	}
}

func Test_WhenDefiningOptionWithEmptyNames_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineOption("", "", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithTwoCharacterShortName_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineOption("sn", "longName", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithSingleCharacterLongName_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineOption("s", "l", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineOption("s", "short", "description", "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningOptionWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineOption("s", "short", "description", ".+")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningOptionWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineOption("s", "short", "description", "*")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningDuplicateOption_ThenErrorIsReturned(t *testing.T) {
	option := model.NewOption("n", "name", "description", "")
	provider := func() model.Option { return option }
	state := domain.NewStateMachine("", "", &mockDataCache{optionProvider: provider})
	err := state.DefineOption("n", "name", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithInvalidCharacterInName_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("arg$", "description", 1, 1, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithMinCountZeroAndMaxCountGreaterThanMinCount_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 0, 2, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithMinCountOneAndMaxCountGreaterThanMinCount_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 1, 2, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndEqualToMaxCount_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 2, 2, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndGreaterThanToMaxCount_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 3, 2, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, ".+")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, "*")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningDuplicateArgument_ThenErrorIsReturned(t *testing.T) {
	argument := model.NewArgument("ARGUMENT", "description", 1, 1, "")
	provider := func() model.Argument { return argument }
	state := domain.NewStateMachine("", "", &mockDataCache{argumentProvider: provider})
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, "")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenParsingUndefinedOption_ThenErrorIsReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "--undefined"}
	state := domain.NewStateMachine("", "", newEmptyMockDataCache())
	err := state.Parse()

	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenParsingDefinedOptionValue_ThenThatValueIsSaved(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	expectedName := "verbose"
	expectedValue := "t"
	os.Args = []string{"appName", "--" + expectedName, expectedValue}

	option := model.NewOption("v", "verbose", "description", "")
	actualName := ""
	actualValue := ""

	mockDataCache := newEmptyMockDataCache()
	mockDataCache.optionProvider = func() model.Option { return option }
	mockDataCache.optionValueListener = func(k string, v string) { actualName = k; actualValue = v }

	state := domain.NewStateMachine("", "", mockDataCache)
	state.Parse()

	if actualName != expectedName || actualValue != expectedValue {
		t.Errorf("Expected <%s> and <%s>, but got <%s> and <%s>", expectedName, expectedValue, actualName, actualValue)
	}
}

func Test_WhenParsingNonMatchingOptionValueWithoutArguments_ThenErrorIsReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "--option", "12"}
	option := model.NewOption("", "option", "description", `^(true|false)$`)

	mockDataCache := newEmptyMockDataCache()
	mockDataCache.optionProvider = func() model.Option { return option }

	state := domain.NewStateMachine("", "", mockDataCache)
	err := state.Parse()

	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenParsingNonMatchingOptionValueWithArguments_ThenFirstMatchingArgumentReceivesTheValue(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "--option", "12"}
	option := model.NewOption("", "option", "description", `^(true|false)$`)
	argument := model.NewArgument("ARG", "description", 1, 1, ``)

	expectedArgumentValue := "12"
	expectedOptionValue := ""
	actualOptionValue := ""
	actualArgumentValue := ""

	mockDataCache := newEmptyMockDataCache()
	mockDataCache.optionProvider = func() model.Option { return option }
	mockDataCache.argumentsProvider = func() []model.Argument { return []model.Argument{argument} }
	mockDataCache.optionValueListener = func(string, v string) { actualOptionValue = v }
	mockDataCache.argumentValueListener = func(string, v string) { actualArgumentValue = v }

	state := domain.NewStateMachine("", "", mockDataCache)
	state.Parse()

	if actualOptionValue != expectedOptionValue || actualArgumentValue != expectedArgumentValue {
		t.Errorf("Expected <%s> and <%s>, but got <%s> and <%s>",
			expectedOptionValue,
			expectedArgumentValue,
			actualOptionValue,
			actualArgumentValue,
		)
	}
}

func Test_WhenRequestingValueForParsedOptionWithNoGivenValue_ThenBooleanTrueIsReturned(t *testing.T) {
	option := model.NewOption("", "opt", "description", "")
	option.SetParsed()

	mockDataCache := newEmptyMockDataCache()
	mockDataCache.optionProvider = func() model.Option { return option }

	state := domain.NewStateMachine("", "", mockDataCache)
	actual := state.GetOptionValue("opt")
	expected := "true"

	if actual != expected {
		t.Errorf("Expeted <%s>, but got <%s>", expected, actual)
	}
}
