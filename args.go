// Package args offers a more Unix style options and arguments handling
// alternative than the native "flags" implementation.
package args

import (
	"fmt"
	"os"
	"strings"

	"github.com/echsylon/go-args/internal/repository"
	"github.com/echsylon/go-args/internal/util"
)

var dataCache repository.DataCache = nil

// Name of the application. Defaults to the first fraction in the OS provided
// arguments list. This field will only be shown in the help text.
var Name string = os.Args[0]

// Description of the application. Defaults to empty-string. This field will
// only be shown in the help text.
var Description string

// DefineOption allows the developer to define a default optional command line
// argument the caller can pass to the application.
//
// If the name is a single character it will serve as the shortName, else as
// the longName. The name is a mandatory field.
//
// The description, if given, is only shown in the help output.
//
// See DefineConstrainedOption for more details.
func DefineOption(name string, description string) {
	if len(name) > 1 {
		DefineConstrainedOption("", name, description, "")
	} else {
		DefineConstrainedOption(name, "", description, "")
	}
}

// DefineConstrainedOption allows the developer to define an optional command
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
// Passing the option name alone, without a corresponding value, is treated the
// same way as would the value "true" have been passed along.
//
// If the caller passes multiple instances of the same option with different
// values, the library will print a help text and exit the application
// gracefully.
//
// If the caller passes a value that doesn't match the given pattern, the
// library will print a help text and exit the application gracefully.
func DefineConstrainedOption(shortName string, longName string, description string, pattern string) {
	err := getRepository().DefineOption(shortName, longName, pattern, description)
	if err != nil {
		panic(err.Error())
	}
}

// DefineArgument allows the developer to define a mandatory argument the
// caller must pass to the application. By default the defined argument will
// accept exactly one value of any shape and size.
//
// See DefineConstrainedArgument for more details.
func DefineArgument(name string, description string) {
	DefineConstrainedArgument(name, description, "", 1, 1)
}

// DefineConstrainedArgument allows the developer to define more granular
// mandatory arguments the caller must pass to the application.
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
// If minCount and maxCount is given the number of caller provided values will
// be validated to be in that (inclusive) range. The minCount must be greater
// than or equal to 1.
//
// If the caller fails to pass the constrained number of matching arguments,
// the library will print a help text and exit the application gracefully.
func DefineConstrainedArgument(name string, description string, pattern string, minCount int, maxCount int) {
	err := getRepository().DefineArgument(minCount, maxCount, name, pattern, description)
	if err != nil {
		panic(err.Error())
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
	var currentOptionName string = ""
	var cache = getRepository()
	var input = getInput()

	cache.ClearValues()

	for _, data := range input {
		if isPotentialOptionName(data) {
			name := strings.Trim(data, " -")
			err := cache.SetOptionParsed(name)
			if err != nil {
				exitWithHelpMessage(err.Error())
			} else {
				currentOptionName = data
			}
		} else if isPotentialOptionValue(currentOptionName, data) {
			err := cache.SetOptionValue(currentOptionName, data)
			if err != nil {
				exitWithHelpMessage(err.Error())
			} else {
				currentOptionName = ""
			}
		} else if isPotentialArgumentValue(currentOptionName, data) {
			err := cache.AddArgumentValue(data)
			if err != nil {
				exitWithHelpMessage(err.Error())
			}
		} else {
			exitWithHelpMessage("Unexpected input: " + data)
		}
	}

	err := cache.AssertAllArgumentValuesProvided()
	if err != nil {
		exitWithHelpMessage(err.Error())
	}
}

// GetOptionValue returns the parsed value for a defined option. Note that
// options are optional and that this function therefore may return an empty
// string. The function accepts both the shortName and the longName of the
// defined option.
func GetOptionValue(name string) string {
	return getRepository().GetOptionValue(name)
}

// GetArgumentValues returns all parsed mandatory values that matched the
// defined argument.
func GetArgumentValues(name string) []string {
	return getRepository().GetArgumentValues(name)
}

// Reset will delete all previously configured options and arguments and
// purge any corresponding parsed values.
func Reset() {
	dataCache = nil
}

func getRepository() repository.DataCache {
	if dataCache == nil {
		dataCache = repository.NewDataCache()
	}
	return dataCache
}

func getInput() []string {
	return os.Args[1:]
}

func isPotentialOptionName(data string) bool {
	return strings.HasPrefix(data, "-")
}

func isPotentialOptionValue(optionName string, value string) bool {
	return isPotentialOptionName(optionName) && value != "" && !isPotentialOptionName(value)
}

func isPotentialArgumentValue(optionName string, value string) bool {
	return optionName == "" && value != ""
}

func exitWithHelpMessage(message string) {
	var stringBuilder strings.Builder
	var cache = dataCache
	var options = cache.GetOptions()
	var arguments = cache.GetArguments()

	if message != "" {
		stringBuilder.WriteString(message)
	}

	var mainSection = util.GetMainHelpSection(Name, Description, &options, &arguments)
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
