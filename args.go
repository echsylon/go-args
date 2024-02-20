// Package args offers a more Unix style options and arguments handling
// alternative than the native "flags" implementation.
package args

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/echsylon/go-args/internal/data"
	"github.com/echsylon/go-args/internal/domain"
	"github.com/echsylon/go-args/internal/util"
)

var state = domain.NewStateMachine(os.Args[0], "", data.NewRepository())

func SetApplicationDescription(text string) {
	state.SetDescription(text)
}

// DefineOption allows the developer to define a default optional command line
// argument the caller can pass to the application.
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
// line argument the caller can pass to the application. Only defined options
// will be accepted during the parsing phase.
//
// At least a shortName or a longName must be given. Defining both is nice but
// not functionally required. The shortName must be a single alphabetic
// character (`^[a-zA-Z]{1}$`), while the longName must be at least two
// alphabetic+ characters (`^[a-zA-Z-._]{2,}$`).
//
// The description, if given, is only shown in the help output.
//
// If a pattern is given, then any user provided value for this option must
// match it. Empty string patterns are allowed and will not constrain the input.
//
// If the caller passes the option name alone (short or long), without a giving
// corresponding value, the library will treat it the same way as would the
// value "true" have been passed along.
//
// If the caller passes multiple instances of the same option with different
// values, the library will print a help text and exit the application
// gracefully.
//
// If the caller passes a value that doesn't match the given pattern, the
// library will print a help text and exit the application gracefully.
func DefineOptionStrict(shortName string, longName string, description string, pattern string) {
	err := state.DefineOption(shortName, longName, description, pattern)
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
// than or equal to 1.
//
// If a pattern is given, then the caller provided input argument data will be
// matched against it. This allows the caller to mix the order of input values.
// The values will in the parsing phase be associated with the first argument
// definition that matches them. It is the developers responsibility to define
// non-overlapping patterns (if this is important) and to provide sufficient
// documentation in the argument descriptions for the caller to make an
// educated call statement.
//
// If a pattern is given, it will be validated, causing the library to panic if
// it's invalid.
//
// If the caller fails to pass the constrained number of matching arguments,
// the library will print a help text and exit the application gracefully.
func DefineArgumentStrict(name string, description string, minCount int, maxCount int, pattern string) {
	err := state.DefineArgument(name, description, minCount, maxCount, pattern)
	if err != nil {
		panic(err)
	}
}

// Parse operates on the user provided command line arguments and matches them
// agains the developer defined option- and argument configurations. The parser
// will validate the input and print the help text and exit if:
//
// - An unknown option is parsed, or it's corresponding value doesn't match the
// defined regular expression.
//
// - A mandatory argument is parsed that doesn't match the defined regular
// expression or violates the defined input count limitations.
//
// - After all input is parsed and there are defined mandatory arguments that
// hasn't received the minimum number of input values.
func Parse() {
	if err := state.Parse(); err != nil {
		exitWithHelpMessage(err, state)
	}
}

// GetOptionValue returns the parsed value for a defined option as a
// string. If there is no value for the option, the fallback is returned
// instead.
func GetOptionValue(name string, fallback string) string {
	result := state.GetOptionValue(name)
	if result == "" {
		result = fallback
	}
	return result
}

// GetOptionIntValue returns the parsed value for a defined option as a
// 64 bit integer. If there is no value for the option, the fallback is
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
// 64 bit floating point number. If there is no value for the option, the
// fallback is returned instead.
func GetOptionFloatValue(name string, fallback float64) float64 {
	value := state.GetOptionValue(name)
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		result = fallback
	}
	return result
}

// GetOptionBoolValue returns the parsed value for a defined option as a
// boolean. If there is no value for the option, the fallback is returned
// instead.
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

// GetArgumentIntValues returns all parsed mandatory values that matched
// the defined argument and can be parsed into a 64 bit integer. Values
// that can not be parsed are simply omitted.
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

// GetArgumentFloatValues returns all parsed mandatory values that matched
// the defined argument and can be parsed into a 64 bit floating point
// number. Values that can not be parsed are simply omitted.
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

// GetArgumentBoolValues returns all parsed mandatory values that matched
// the defined argument and can be parsed into a boolean value. Values
// that can not be parsed are simply omitted.
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

func exitWithHelpMessage(message error, state domain.StateMachine) {
	var stringBuilder strings.Builder
	var name = state.GetName()
	var description = state.GetDescription()
	var options = state.GetDefinedOptions()
	var arguments = state.GetDefinedArguments()

	if message != nil {
		stringBuilder.WriteString(message.Error())
	}

	var mainSection = util.GetMainHelpSection(name, description, &options, &arguments)
	if mainSection != "" {
		stringBuilder.WriteString("\n\n")
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
