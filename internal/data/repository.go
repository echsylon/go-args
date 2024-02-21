package data

import (
	"github.com/echsylon/go-args/internal/model"
)

type Repository interface {
	ClearAll()
	ClearValues()
	SaveOption(shortName string, longName string, description string, pattern string)
	GetOptions() []model.Option
	GetOption(name string) model.Option
	SaveOptionValue(name string, value string)
	GetOptionValue(name string) string
	SaveArgument(name string, description string, min int, max int, pattern string)
	GetArguments() []model.Argument
	GetArgument(name string) model.Argument
	SaveArgumentValue(name string, value string)
	GetArgumentValues(name string) []string
}

type any = interface{}

func NewRepository() Repository {
	return &repository{
		definitions: []any{},
		values:      make(map[any][]string)}
}

type repository struct {
	// !!! NOTE !!!
	// We're intentionally separating the values map keys in
	// to the definitions slice in order to maintain the order
	// they were added in. Go has not only undefined order of
	// keys in their map implementation, but they are also
	// actively randomizing (ish) them to prevent implicit
	// order dependencies.
	definitions []any
	values      map[any][]string
}

func (cache *repository) ClearAll() {
	cache.definitions = []any{}
	cache.values = make(map[any][]string)
}

func (cache *repository) ClearValues() {
	cache.values = make(map[any][]string)
}

func (cache *repository) SaveOption(shortName string, longName string, description string, pattern string) {
	cache.definitions = append(cache.definitions, model.NewOption(shortName, longName, description, pattern))
}

func (cache *repository) GetOptions() []model.Option {
	var result []model.Option
	for _, item := range cache.definitions {
		if option, isOption := item.(model.Option); isOption {
			result = append(result, option)
		}
	}
	return result
}

func (cache *repository) GetOption(name string) model.Option {
	return findOption(name, name, &cache.definitions)
}

func (cache *repository) SaveOptionValue(name string, value string) {
	if option := findOption(name, name, &cache.definitions); option != nil {
		cache.values[option] = []string{value}
	}
}

func (cache *repository) GetOptionValue(name string) string {
	var result = ""
	option := findOption(name, name, &cache.definitions)
	if option != nil {
		values, hasValue := cache.values[option]
		if hasValue {
			result = values[0]
		}
	}
	return result
}

func (cache *repository) SaveArgument(name string, description string, min int, max int, pattern string) {
	cache.definitions = append(cache.definitions, model.NewArgument(name, description, min, max, pattern))
}

func (cache *repository) GetArguments() []model.Argument {
	var result []model.Argument
	for _, item := range cache.definitions {
		if argument, isArgument := item.(model.Argument); isArgument {
			result = append(result, argument)
		}
	}
	return result
}

func (cache *repository) GetArgument(name string) model.Argument {
	return findArgument(name, &cache.definitions)
}

func (cache *repository) SaveArgumentValue(name string, value string) {
	if argument := findArgument(name, &cache.definitions); argument != nil {
		values, hasValues := cache.values[argument]
		if !hasValues {
			cache.values[argument] = []string{value}
		} else {
			values = append(values, value)
			cache.values[argument] = values
		}
	}
}

func (cache *repository) GetArgumentValues(name string) []string {
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

func findOption(shortName string, longName string, definitions *[]any) model.Option {
	var result model.Option = nil
	if shortName != "" || longName != "" {
		for _, item := range *definitions {
			if option, isOption := item.(model.Option); !isOption {
				continue
			} else if option.GetShortName() == shortName || option.GetLongName() == longName {
				result = option
				break
			}
		}
	}
	return result
}

func findArgument(name string, definitions *[]any) model.Argument {
	var result model.Argument = nil
	if name != "" {
		for _, item := range *definitions {
			if argument, isArgument := item.(model.Argument); !isArgument {
				continue
			} else if argument.GetName() == name {
				result = argument
				break
			}
		}
	}
	return result
}
