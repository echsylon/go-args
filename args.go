// Package args offers a more Unix style options and arguments handling
// alternative than the native "flags" implementation.
package args

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/echsylon/go-args/internal/data"
	"github.com/echsylon/go-args/internal/domain"
	"github.com/echsylon/go-args/internal/util"
)

var state = domain.NewStateMachine(filepath.Base(os.Args[0]), "", data.NewRepository())

// SetApplicationDescription takes a human readable description of the app.
// This text is only shown in the help output.
func SetApplicationDescription(text string) {
	state.SetDescription(text)
}

// DefineOption allows the developer to define a simple optional command line
// argument the caller can pass to the application. Only defined options will
// be accepted during the parsing phase.
//
// If the name is a single character it will serve as the shortName, else as
// the longName. The name is a mandatory field.
//
// The description, if given, is only shown in the help output.
//
// See DefineOptionStrict for more details.
func DefineOption(name string, description string) {
	if len(name) > 1 {
		DefineOptionStrict("", name, description, "")
	} else {
		DefineOptionStrict(name, "", description, "")
	}
}

// DefineOptionStrict allows the developer to define an optional command
// line argument the caller can pass to the application.
//
// At least a shortName or a longName must be given. Defining both is nice
// but not functionally required. The shortName must be a single alphabetic
// character matching the `^[a-zA-Z]{1}$` regular expression, while the
// longName must be at least two alphabetic+ characters, matching the
// `^[a-zA-Z-._]{2,}$` regular expression.
//
// The description, if given, is only shown in the help output.
//
// If a pattern regular expression is given, then any caller provided value
// for this option must match it for the value to be assigned to the option.
// Empty string patterns are allowed and will behave as "accept everything".
// The library will validate the pattern regular expression. If the validation
// fails the library will panic runtime.
//
// If a caller passes the defined (short or long) option name alone, without
// any corresponding value, the library will treat it as a boolean true flag
// and return "true" for it's value.
//
// If the caller provides multiple instances of the same option, the library
// will print a help text and exit the application gracefully.
//
// If the caller passes a value that doesn't match the given option pattern,
// then the library will try to match the value for any argument instead. If
// there is a defined argument that accepts the value it will be assigned to
// that argument, otherwise the library will print a help text and exit the
// application gracefully.
func DefineOptionStrict(shortName string, longName string, description string, pattern string) {
	err := state.DefineOption(shortName, longName, description, pattern)
	if err != nil {
		panic(err)
	}
}

// DefineOptionHelp allows the devleoper to define a graceful help trigger
// option.
//
// At least a shortName or a longName must be given. Defining both is nice
// but not functionally required. The shortName must be a single alphabetic
// character matching the `^[a-zA-Z]{1}$` regular expression, while the
// longName must be at least two alphabetic+ characters, matching the
// `^[a-zA-Z-._]{2,}$` regular expression.
//
// The description, if given, is only shown in the help output.
//
// If a caller passes this option, the library will immediately abort parsing
// and print the help text - without any error message (!) - and exit the
// application gracefully.
func DefineOptionHelp(shortName string, longName string, description string) {
	err := state.DefineHelpOption(shortName, longName, description)
	if err != nil {
		panic(err)
	}
}

// DefineArgument allows the developer to define a mandatory argument the
// caller must pass to the application. By default the defined argument will
// accept exactly one value of any shape and size.
//
// See DefineArgumentStrict for more details.
func DefineArgument(name string, description string) {
	DefineArgumentStrict(name, description, 1, 1, "")
}

// DefineArgumentStrict allows the developer to define more granular
// mandatory arguments the caller must pass to the application.
//
// If minCount and maxCount is given the number of caller provided values will
// be validated to be in that (inclusive) range. The minCount must be greater
// than or equal to 1 and the maxCount must be greater than or equal to minCount.
//
// If a pattern is given, then the caller provided input argument value will be
// matched against it. This allows the caller to mix the order of input values.
// The values will in the parsing phase be associated with the first argument
// definition that matches them. It is the developers responsibility to define
// non-overlapping patterns (if this is important) and to provide sufficient
// documentation in the argument descriptions for the caller to make an
// educated call statement.
//
// If a pattern is given, it will be validated, causing the library to panic
// runtime if it's invalid.
//
// If the caller fails to pass the constrained number of matching arguments,
// the library will print a help text and exit the application gracefully.
//
// If the caller provides an argument value that doesn't match any defined
// arguments, the library will print a help text and exit the application
// gracefully.
func DefineArgumentStrict(name string, description string, minCount int, maxCount int, pattern string) {
	err := state.DefineArgument(name, description, minCount, maxCount, pattern)
	if err != nil {
		panic(err)
	}
}

// Parse operates on the user provided command line arguments and matches them
// against the developer defined option and argument configurations. The parse
// function will validate the input and print the help text and exit gracefully
// if:
//
//   - An unknown option is parsed.
//
//   - A value doesn't match any option or argument patterns.
//
//   - An argument value would violate the defined argument value count limits.
//
//   - After all input is parsed, there are defined mandatory arguments that
//     hasn't received the minimum number of input values.
func Parse() {
	if err := state.Parse(); err != nil {
		exitWithHelpMessage(err, state)
	}
}

// GetOptionValue returns the parsed value for a defined option as a string.
// If there is no parsed value for the option, the fallback is returned
// instead.
func GetOptionValue(name string, fallback string) string {
	result := state.GetOptionValue(name)
	if result == "" {
		result = fallback
	}
	return result
}

// GetOptionIntValue returns the parsed value for a defined option as a 64 bit
// integer. If there is no parsed value for the option, the fallback is
// returned instead.
func GetOptionIntValue(name string, fallback int64) int64 {
	value := state.GetOptionValue(name)
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		result = fallback
	}
	return result
}

// GetOptionFloatValue returns the parsed value for a defined option as a
// 64 bit floating point number. If there is no parsed value for the option,
// the fallback is returned instead.
func GetOptionFloatValue(name string, fallback float64) float64 {
	value := state.GetOptionValue(name)
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		result = fallback
	}
	return result
}

// GetOptionBoolValue returns the parsed value for a defined option as a
// boolean. If there is no parsed value for the option, the fallback is
// returned instead.
func GetOptionBoolValue(name string, fallback bool) bool {
	value := state.GetOptionValue(name)
	result, err := strconv.ParseBool(value)
	if err != nil {
		result = fallback
	}
	return result
}

// GetArgumentValues returns all parsed mandatory values that matched the
// defined argument.
func GetArgumentValues(name string) []string {
	return state.GetArgumentValues(name)
}

// GetArgumentIntValues returns all parsed mandatory values that matched the
// defined argument and can be cast into a 64 bit integer. Values that can not
// be cast into a 64 bit integer are simply omitted from the result.
func GetArgumentIntValues(name string) []int64 {
	values := state.GetArgumentValues(name)
	result := []int64{}
	for _, value := range values {
		if number, err := strconv.ParseInt(value, 10, 64); err == nil {
			result = append(result, number)
		}
	}
	return result
}

// GetArgumentFloatValues returns all parsed mandatory values that matched the
// defined argument and can be cast into a 64 bit floating point number. Values
// that can not be cast into a 64 bit floating point number are simply omitted
// from the result.
func GetArgumentFloatValues(name string) []float64 {
	values := state.GetArgumentValues(name)
	result := []float64{}
	for _, value := range values {
		if number, err := strconv.ParseFloat(value, 64); err == nil {
			result = append(result, number)
		}
	}
	return result
}

// GetArgumentBoolValues returns all parsed mandatory values that matched the
// defined argument and can be cast into a boolean value. Values that can not
// be cast into a boolean value are simply omitted from the result.
func GetArgumentBoolValues(name string) []bool {
	values := state.GetArgumentValues(name)
	result := []bool{}
	for _, value := range values {
		if state, err := strconv.ParseBool(value); err == nil {
			result = append(result, state)
		}
	}
	return result
}

// Reset will delete all previously configured options and arguments and
// purge any corresponding parsed values.
func Reset() {
	state.Reset()
}

func exitWithHelpMessage(err error, state domain.StateMachine) {
	var stringBuilder strings.Builder
	var name = state.GetName()
	var description = state.GetDescription()
	var options = state.GetDefinedOptions()
	var arguments = state.GetDefinedArguments()

	var message = err.Error()
	if message != "" {
		stringBuilder.WriteString(message)
		stringBuilder.WriteString("\n\n")
	}

	var mainSection = util.GetMainHelpSection(name, description, &options, &arguments)
	if mainSection != "" {
		stringBuilder.WriteString(mainSection)
	}

	var argumentsSection = util.GetArgumentsHelpSection(&arguments)
	if argumentsSection != "" {
		stringBuilder.WriteString("\n\n")
		stringBuilder.WriteString(argumentsSection)
	}

	var optionsSection = util.GetOptionsHelpSection(&options)
	if optionsSection != "" {
		stringBuilder.WriteString("\n\n")
		stringBuilder.WriteString(optionsSection)
	}

	fmt.Println(stringBuilder.String())
	os.Exit(0)
}
