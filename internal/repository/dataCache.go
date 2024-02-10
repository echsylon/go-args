package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/echsylon/go-args/internal/configuration"
	"github.com/echsylon/go-args/internal/model"
)

type DataCache interface {
	ClearValues()
	DefineOption(shortName string, longName string, description string, pattern string) error
	IsValidOptionValue(name string, value string) bool
	GetOptions() []model.Option
	SetOptionParsed(name string) error
	SetOptionValue(name string, value string) error
	GetOptionValue(name string) string
	DefineArgument(name string, description string, min int, max int, pattern string) error
	IsValidArgumentValue(value string) bool
	GetArguments() []model.Argument
	AddArgumentValue(value string) error
	GetArgumentValues(name string) []string
	AssertAllArgumentValuesProvided() error
}

func NewDataCache() DataCache {
	return &dataCache{
		definitions: []any{},
		values:      make(map[any][]string)}
}

type dataCache struct {
	definitions []any
	values      map[any][]string
}

func (cache *dataCache) ClearValues() {
	clear(cache.values)
}
func (cache *dataCache) GetOptions() []model.Option {
	var result []model.Option
	for _, item := range cache.definitions {
		if option, isOption := item.(model.Option); isOption {
			result = append(result, option)
		}
	}
	return result
}

func (cache *dataCache) GetArguments() []model.Argument {
	var result []model.Argument
	for _, item := range cache.definitions {
		if argument, isArgument := item.(model.Argument); isArgument {
			result = append(result, argument)
		}
	}
	return result
}

func (cache *dataCache) DefineOption(shortName string, longName string, description string, pattern string) error {
	var result error = nil
	if shortName == "" && longName == "" {
		result = fmt.Errorf("no name given for option")
	} else if shortName != "" && !isValidOptionShortName(shortName) {
		result = fmt.Errorf("unexpected short name: %s", shortName)
	} else if longName != "" && !isValidOptionLongName(longName) {
		result = fmt.Errorf("unexpected long name: %s", longName)
	} else if !isValidRegularExpression(pattern) {
		result = fmt.Errorf("unexpected option value pattern: %s", pattern)
	} else if isOptionAlreadyDefined(shortName, longName, &cache.definitions) {
		result = fmt.Errorf("option already defined: %s, %s", shortName, longName)
	} else {
		cache.definitions = append(cache.definitions, model.NewOption(shortName, longName, description, pattern))
	}

	return result
}

func (cache *dataCache) IsValidOptionValue(name string, value string) bool {
	var result = false
	if !isValidOptionShortName(name) && !isValidOptionLongName(name) {
		result = false
	} else if option := findOption(name, &cache.definitions); option != nil {
		result = canSaveOptionValue(option, value, &cache.values)
	}
	return result
}

func (cache *dataCache) SetOptionParsed(name string) error {
	var result error = nil
	option := findOption(name, &cache.definitions)
	if option == nil {
		result = fmt.Errorf("unknown option: %s", name)
	} else {
		option.SetParsed()
	}
	return result
}

func (cache *dataCache) SetOptionValue(name string, value string) error {
	var result error = nil
	option := findOption(name, &cache.definitions)
	if option == nil {
		result = fmt.Errorf("unknown option: %s", name)
	} else if !canSaveOptionValue(option, value, &cache.values) {
		result = fmt.Errorf("unexpected option value: %s", value)
	} else {
		option.SetParsed()
		cache.values[option] = []string{value}
	}
	return result
}

func (cache *dataCache) GetOptionValue(name string) string {
	var result = ""
	option := findOption(name, &cache.definitions)
	if option != nil {
		values, hasValue := cache.values[option]
		if hasValue {
			result = values[0]
		}
		if !hasValue && option.IsParsed() {
			result = "true"
		}
	}
	return result
}

func (cache *dataCache) DefineArgument(name string, description string, min int, max int, pattern string) error {
	var result error = nil
	if !isValidArgumentCountRange(min, max) {
		result = fmt.Errorf("unexpected range: [%d..%d]", min, max)
	} else if !isValidArgumentName(name) {
		result = fmt.Errorf("unexpected argument name: %s", name)
	} else if !isValidRegularExpression(pattern) {
		result = fmt.Errorf("unexpected argument value pattern: %s", pattern)
	} else if isArgumentAlreadyDefined(name, &cache.definitions) {
		result = fmt.Errorf("argument already defined: %s", name)
	} else {
		cache.definitions = append(cache.definitions, model.NewArgument(name, description, min, max, pattern))
	}

	return result
}

func (cache *dataCache) IsValidArgumentValue(value string) bool {
	var result = false
	for _, item := range cache.definitions {
		if argument, isArgument := item.(model.Argument); !isArgument {
			continue
		} else if canSaveArgumentValue(argument, value, &cache.values) {
			result = true
			break
		}
	}
	return result
}

func (cache *dataCache) AddArgumentValue(value string) error {
	var result = fmt.Errorf("unexpected argument value: %s", value)
	for _, item := range cache.definitions {
		if argument, isArgument := item.(model.Argument); !isArgument {
			continue
		} else if canSaveArgumentValue(argument, value, &cache.values) {
			values, hasValues := cache.values[argument]
			if !hasValues {
				cache.values[argument] = []string{value}
			} else {
				values = append(values, value)
				cache.values[argument] = values
			}
			result = nil
			break
		}
	}
	return result
}

func (cache *dataCache) GetArgumentValues(name string) []string {
	var result []string
	argument := findArgument(name, &cache.definitions)
	if argument != nil {
		values, hasValues := cache.values[argument]
		if hasValues {
			result = values
		}
	}
	return result
}

func (cache *dataCache) AssertAllArgumentValuesProvided() error {
	var result error = nil
	var missing = []string{}
	for _, item := range cache.definitions {
		if argument, isArgument := item.(model.Argument); isArgument {
			values := cache.values[argument]
			if len(values) < argument.GetMinValuesCount() {
				missing = append(missing, argument.GetName())
			}
		}
	}
	if len(missing) > 0 {
		var labels = strings.Join(missing, ", ")
		result = fmt.Errorf("missing values for: %s", labels)
	}
	return result
}

func findOption(name string, definitions *[]any) model.Option {
	return findOptionAny(name, name, definitions)
}

func findOptionAny(shortName string, longName string, definitions *[]any) model.Option {
	var result model.Option = nil
	shortName = strings.Trim(shortName, " -")
	longName = strings.Trim(longName, " -")
	for _, item := range *definitions {
		if option, isOption := item.(model.Option); !isOption {
			continue
		} else if option.GetShortName() == shortName || option.GetLongName() == longName {
			result = option
			break
		}
	}
	return result
}

func findArgument(name string, definitions *[]any) model.Argument {
	var result model.Argument = nil
	name = strings.Trim(name, " -")
	for _, item := range *definitions {
		if argument, isArgument := item.(model.Argument); !isArgument {
			continue
		} else if argument.GetName() == name {
			result = argument
			break
		}
	}
	return result
}

func isValidArgumentName(name string) bool {
	return configuration.ArgumentNamePattern.MatchString(name)
}

func isValidOptionShortName(name string) bool {
	return configuration.OptionShortNamePattern.MatchString(name)
}

func isValidOptionLongName(name string) bool {
	return configuration.OptionLongNamePattern.MatchString(name)
}

func isOptionAlreadyDefined(shortName string, longName string, definitions *[]any) bool {
	return findOptionAny(shortName, longName, definitions) != nil
}

func isValidArgumentCountRange(min int, max int) bool {
	return min >= 1 && min <= max
}

func isValidRegularExpression(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}

func isArgumentAlreadyDefined(name string, definitions *[]any) bool {
	return findArgument(name, definitions) != nil
}

func canSaveOptionValue(option model.Option, value string, values *map[any][]string) bool {
	pattern := option.GetPattern()
	test, err := regexp.Compile(pattern)

	var canSave bool
	if err != nil { // "This should never happen"
		canSave = false
	} else if test.MatchString(value) {
		optionValues, hasValue := (*values)[option]
		canSave = !hasValue || len(optionValues) == 0 || optionValues[0] == value
	} else {
		canSave = false
	}

	return canSave
}

func canSaveArgumentValue(argument model.Argument, value string, values *map[any][]string) bool {
	pattern := argument.GetPattern()
	test, err := regexp.Compile(pattern)
	argValues := (*values)[argument]
	hasRoomForMore := len(argValues) < argument.GetMaxValuesCount()

	var canSave bool
	if err != nil { // "This should never happen"
		canSave = false
	} else if hasRoomForMore && test.MatchString(value) {
		canSave = true
	} else {
		canSave = false
	}

	return canSave
}
