package domain

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/echsylon/go-args/internal/configuration"
	"github.com/echsylon/go-args/internal/data"
	"github.com/echsylon/go-args/internal/model"
)

type StateMachine interface {
	SetName(name string)
	GetName() string
	SetDescription(description string)
	GetDescription() string
	DefineOption(shortName string, longName string, description string, pattern string) error
	GetDefinedOptions() []model.Option
	GetOptionValue(name string) string
	DefineArgument(name string, description string, minCount int, maxCount int, pattern string) error
	GetDefinedArguments() []model.Argument
	GetArgumentValues(name string) []string
	Parse() error
	Reset()
}

func NewStateMachine(name string, description string, data data.Repository) StateMachine {
	return &stateMachine{name, description, data}
}

type stateMachine struct {
	name        string
	description string
	data        data.Repository
}

func (state *stateMachine) SetName(name string) {
	state.name = name
}

func (state *stateMachine) GetName() string {
	return state.name
}

func (state *stateMachine) SetDescription(description string) {
	state.description = description
}

func (state *stateMachine) GetDescription() string {
	return state.description
}

func (state *stateMachine) DefineOption(shortName string, longName string, description string, pattern string) error {
	var result error = nil
	if shortName == "" && longName == "" {
		result = fmt.Errorf("no name given for option")
	} else if shortName != "" && !isValidOptionShortName(shortName) {
		result = fmt.Errorf("unexpected short name: %s", shortName)
	} else if longName != "" && !isValidOptionLongName(longName) {
		result = fmt.Errorf("unexpected long name: %s", longName)
	} else if !isValidRegularExpression(pattern) {
		result = fmt.Errorf("unexpected option value pattern: %s", pattern)
	} else if isOptionAlreadyDefined(shortName, longName, state.data) {
		result = fmt.Errorf("option already defined: %s, %s", shortName, longName)
	} else {
		state.data.SaveOption(shortName, longName, description, pattern)
	}
	return result
}

func (state *stateMachine) GetDefinedOptions() []model.Option {
	return state.data.GetOptions()
}

func (state *stateMachine) GetOptionValue(name string) string {
	option := state.data.GetOption(name)
	value := state.data.GetOptionValue(name)
	if value == "" && option != nil && option.IsParsed() {
		value = "true"
	}
	return value
}

func (state *stateMachine) DefineArgument(name string, description string, minCount int, maxCount int, pattern string) error {
	var result error = nil
	if !isValidArgumentCountRange(minCount, maxCount) {
		result = fmt.Errorf("unexpected range: [%d..%d]", minCount, maxCount)
	} else if !isValidArgumentName(name) {
		result = fmt.Errorf("unexpected argument name: %s", name)
	} else if !isValidRegularExpression(pattern) {
		result = fmt.Errorf("unexpected argument value pattern: %s", pattern)
	} else if isArgumentAlreadyDefined(name, state.data) {
		result = fmt.Errorf("argument already defined: %s", name)
	} else {
		state.data.SaveArgument(name, description, minCount, maxCount, pattern)
	}

	return result
}

func (state *stateMachine) GetDefinedArguments() []model.Argument {
	return state.data.GetArguments()
}

func (state *stateMachine) GetArgumentValues(name string) []string {
	return state.data.GetArgumentValues(name)
}

func (state *stateMachine) Parse() error {
	var result error = nil
	var currentOptionName string = ""

	state.data.ClearValues()

	for _, data := range getInput() {
		if isExpectedOption(data, state.data) {
			currentOptionName = strings.Trim(data, "-")
			option := state.data.GetOption(currentOptionName)
			option.SetParsed()
		} else if isExpectedOptionValue(currentOptionName, data, state.data) {
			state.data.SaveOptionValue(currentOptionName, data)
			currentOptionName = ""
		} else if argument := findArgumentForValue(data, state.data); argument != nil {
			state.data.SaveArgumentValue(argument.GetName(), data)
			currentOptionName = ""
		} else {
			result = fmt.Errorf("unexpected input: %s", data)
			break
		}
	}

	if result == nil {
		missing := getUnsatisfiedArguments(state.data)
		if len(missing) > 0 {
			names := strings.Join(missing, ", ")
			result = fmt.Errorf("missing input for: %s", names)
		}
	}

	return result
}

func (state *stateMachine) Reset() {
	state.data.ClearAll()
}

func getInput() []string {
	return os.Args[1:]
}

func isValidOptionShortName(name string) bool {
	return configuration.OptionShortNamePattern.MatchString(name)
}

func isValidOptionLongName(name string) bool {
	return configuration.OptionLongNamePattern.MatchString(name)
}

func isOptionAlreadyDefined(shortName string, longName string, data data.Repository) bool {
	return data.GetOption(shortName) != nil || data.GetOption(longName) != nil
}

func isExpectedOption(input string, data data.Repository) bool {
	var result = false
	if strings.HasPrefix(input, "-") {
		name := strings.Trim(input, "-")
		if option := data.GetOption(name); option != nil {
			if value := data.GetOptionValue(name); value == "" {
				result = true
			}
		}
	}
	return result
}

func isExpectedOptionValue(name string, input string, data data.Repository) bool {
	var result = false
	if option := data.GetOption(name); option != nil {
		if value := data.GetOptionValue(name); value == "" {
			if test, err := regexp.Compile(option.GetPattern()); err == nil {
				if test.MatchString(value) {
					result = true
				}
			}
		}
	}
	return result
}

func isValidArgumentCountRange(min int, max int) bool {
	return min >= 1 && min <= max
}

func isValidArgumentName(name string) bool {
	return configuration.ArgumentNamePattern.MatchString(name)
}

func isArgumentAlreadyDefined(name string, data data.Repository) bool {
	return data.GetArgument(name) != nil
}

func findArgumentForValue(value string, data data.Repository) model.Argument {
	var result model.Argument = nil
	var arguments = data.GetArguments()
	for _, argument := range arguments {
		values := data.GetArgumentValues(argument.GetName())
		if len(values) < argument.GetMaxValuesCount() {
			if test, err := regexp.Compile(argument.GetPattern()); err == nil {
				if test.MatchString(value) {
					result = argument
					break
				}
			}
		}
	}
	return result
}

func getUnsatisfiedArguments(data data.Repository) []string {
	var missing []string
	var arguments = data.GetArguments()
	for _, argument := range arguments {
		name := argument.GetName()
		values := data.GetArgumentValues(name)
		if len(values) < argument.GetMinValuesCount() {
			missing = append(missing, name)
		}
	}
	return missing
}

func isValidRegularExpression(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}
