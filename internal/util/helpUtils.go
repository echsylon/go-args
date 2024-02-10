package util

import (
	"fmt"
	"strings"

	"github.com/echsylon/go-args/internal/model"
)

func GetMainHelpSection(name string, description string, options *[]model.Option, arguments *[]model.Argument) string {
	var stringBuilder strings.Builder
	stringBuilder.WriteString("Usage: ")
	stringBuilder.WriteString(name)

	if options != nil {
		optionsCount := len(*options)
		if optionsCount > 1 {
			stringBuilder.WriteString(" [OPTIONS...]")
		} else if optionsCount > 0 {
			stringBuilder.WriteString(" [OPTION]")
		}
	}

	if arguments != nil {
		for _, argument := range *arguments {
			stringBuilder.WriteString(" ")
			stringBuilder.WriteString(argument.GetName())
			if argument.ExpectsMultipleValues() {
				stringBuilder.WriteString("...")
			}
		}
	}

	if description != "" {
		stringBuilder.WriteString("\n")
		stringBuilder.WriteString(description)
	}

	return stringBuilder.String()
}

func GetOptionsHelpSection(options *[]model.Option) string {
	var stringBuilder strings.Builder
	if options != nil && len(*options) > 0 {
		shortColumnWidth, longColumnWidth := calculateOptionNamesColumnWidth(options)
		stringBuilder.WriteString("Options:")

		for _, option := range *options {
			stringBuilder.WriteString("\n")

			shortName := option.GetShortName()
			shortText := buildOptionShortNameColumn(shortName, shortColumnWidth)
			stringBuilder.WriteString(shortText)

			longName := option.GetLongName()
			longText := buildOptionLongNameColumn(shortName, longName, longColumnWidth)
			stringBuilder.WriteString(longText)

			description := option.GetDescription()
			stringBuilder.WriteString("  " + description)
		}
	}

	return stringBuilder.String()
}

func GetArgumentsHelpSection(arguments *[]model.Argument) string {
	var stringBuilder strings.Builder
	if arguments != nil && len(*arguments) > 0 {
		columnWidth := calculateArgumentNameColumnWidth(arguments)
		stringBuilder.WriteString("Arguments:")

		for _, argument := range *arguments {
			stringBuilder.WriteString("\n")

			name := argument.GetName()
			text := buildArgumentNameColumn(name, columnWidth)
			stringBuilder.WriteString(text)

			description := argument.GetDescription()
			stringBuilder.WriteString("  " + description)
		}
	}

	return stringBuilder.String()
}

func calculateOptionNamesColumnWidth(options *[]model.Option) (int, int) {
	shortWidth := 0
	longWidth := 0
	if options != nil {
		for _, option := range *options {
			shortNameWidth := len(option.GetShortName())
			if shortNameWidth > shortWidth {
				shortWidth = shortNameWidth
			}

			longNameWidth := len(option.GetLongName())
			if longNameWidth > longWidth {
				longWidth = longNameWidth
			}
		}
	}
	return shortWidth, longWidth
}

func buildOptionShortNameColumn(shortName string, columnWidth int) string {
	result := ""
	if columnWidth > 0 {
		text := ""
		if shortName != "" {
			text = "  -" + shortName
		}
		width := columnWidth + 3 // prefix
		result = fmt.Sprintf("%*s", width, text)
	}
	return result
}

func buildOptionLongNameColumn(shortName string, longName string, columnWidth int) string {
	result := ""
	if columnWidth > 0 {
		text := ""
		if longName != "" && shortName == "" {
			text = "  --" + longName
		} else if longName != "" && shortName != "" {
			text = ", --" + longName
		}
		width := columnWidth + 4 // prefix
		result = fmt.Sprintf("%*s", -width, text)
	}
	return result
}

func calculateArgumentNameColumnWidth(arguments *[]model.Argument) int {
	widestWidth := 0
	if arguments != nil {
		for _, argument := range *arguments {
			nameWidth := len(argument.GetName())
			if nameWidth > widestWidth {
				widestWidth = nameWidth
			}
		}
	}
	return widestWidth
}
func buildArgumentNameColumn(name string, columnWidth int) string {
	result := ""
	if columnWidth > 0 {
		text := "  " + name
		size := columnWidth + 2 // 2: prefix
		result = fmt.Sprintf("%-*s", size, text)
	}
	return result
}
