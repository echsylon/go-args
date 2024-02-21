package domain_test

import (
	"os"
	"testing"

	"github.com/echsylon/go-args/internal/data"
	"github.com/echsylon/go-args/internal/domain"
	"github.com/echsylon/go-args/internal/model"
)

type mockRepository struct {
	data.Repository
	argumentsProvider     func() []model.Argument
	argumentProvider      func() model.Argument
	argumentValueProvider func() []string
	argumentValueListener func(string, string)
	optionsProvider       func() []model.Option
	optionProvider        func() model.Option
	optionValueProvider   func() string
	optionValueListener   func(string, string)
}

func (mock *mockRepository) SaveArgument(string, string, int, int, string) {}
func (mock *mockRepository) GetArgument(string) model.Argument             { return mock.argumentProvider() }
func (mock *mockRepository) GetArguments() []model.Argument                { return mock.argumentsProvider() }
func (mock *mockRepository) GetArgumentValues(string) []string             { return mock.argumentValueProvider() }
func (mock *mockRepository) SaveArgumentValue(k string, v string)          { mock.argumentValueListener(k, v) }
func (mock *mockRepository) SaveOption(string, string, string, string)     {}
func (mock *mockRepository) GetOption(string) model.Option                 { return mock.optionProvider() }
func (mock *mockRepository) GetOptions() []model.Option                    { return mock.optionsProvider() }
func (mock *mockRepository) ClearValues()                                  {}
func (mock *mockRepository) SaveOptionValue(k string, v string)            { mock.optionValueListener(k, v) }
func (mock *mockRepository) GetOptionValue(string) string                  { return mock.optionValueProvider() }

func newEmptyMockRepository() *mockRepository {
	return &mockRepository{
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
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineOption("", "", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithTwoCharacterShortName_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineOption("sn", "longName", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithSingleCharacterLongName_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineOption("s", "l", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningOptionWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineOption("s", "short", "description", "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningOptionWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineOption("s", "short", "description", ".+")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningOptionWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineOption("s", "short", "description", "*")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningDuplicateOption_ThenErrorIsReturned(t *testing.T) {
	option := model.NewOption("n", "name", "description", "")
	provider := func() model.Option { return option }
	state := domain.NewStateMachine("", "", &mockRepository{optionProvider: provider})
	err := state.DefineOption("n", "name", "description", "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithInvalidCharacterInName_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("arg$", "description", 1, 1, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithMinCountZeroAndMaxCountGreaterThanMinCount_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 0, 2, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithMinCountOneAndMaxCountGreaterThanMinCount_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 1, 2, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndEqualToMaxCount_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 2, 2, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithMinCountGreaterThanZeroAndGreaterThanToMaxCount_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 3, 2, "")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningArgumentWithNoPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, "")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithValidRegexPattern_ThenNoErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, ".+")
	if err != nil {
		t.Errorf("Expected <nil>, but got <error>: %s", err.Error())
	}
}

func Test_WhenDefiningArgumentWithInvalidRegexPattern_ThenErrorIsReturned(t *testing.T) {
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, "*")
	if err == nil {
		t.Errorf("Expected <error>, but got <nil>")
	}
}

func Test_WhenDefiningDuplicateArgument_ThenErrorIsReturned(t *testing.T) {
	argument := model.NewArgument("ARGUMENT", "description", 1, 1, "")
	provider := func() model.Argument { return argument }
	state := domain.NewStateMachine("", "", &mockRepository{argumentProvider: provider})
	err := state.DefineArgument("ARGUMENT", "description", 1, 1, "")
	if err == nil {
		t.Errorf("Expeted <error>, but got <nil>")
	}
}

func Test_WhenParsingUndefinedOption_ThenErrorIsReturned(t *testing.T) {
	actualArgs := os.Args
	defer func() { os.Args = actualArgs }()

	os.Args = []string{"appName", "--undefined"}
	state := domain.NewStateMachine("", "", newEmptyMockRepository())
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

	mockRepository := newEmptyMockRepository()
	mockRepository.optionProvider = func() model.Option { return option }
	mockRepository.optionValueListener = func(k string, v string) { actualName = k; actualValue = v }

	state := domain.NewStateMachine("", "", mockRepository)
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

	mockRepository := newEmptyMockRepository()
	mockRepository.optionProvider = func() model.Option { return option }

	state := domain.NewStateMachine("", "", mockRepository)
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

	mockRepository := newEmptyMockRepository()
	mockRepository.optionProvider = func() model.Option { return option }
	mockRepository.argumentsProvider = func() []model.Argument { return []model.Argument{argument} }
	mockRepository.optionValueListener = func(string, v string) { actualOptionValue = v }
	mockRepository.argumentValueListener = func(string, v string) { actualArgumentValue = v }

	state := domain.NewStateMachine("", "", mockRepository)
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

	mockRepository := newEmptyMockRepository()
	mockRepository.optionProvider = func() model.Option { return option }

	state := domain.NewStateMachine("", "", mockRepository)
	actual := state.GetOptionValue("opt")
	expected := "true"

	if actual != expected {
		t.Errorf("Expeted <%s>, but got <%s>", expected, actual)
	}
}
