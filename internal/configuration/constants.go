package configuration

import "regexp"

// Treated as internal constants
var ArgumentNamePattern = regexp.MustCompile(`^[a-zA-Z0-9-._]+$`)
var OptionShortNamePattern = regexp.MustCompile(`^[a-zA-Z]{1}$`)
var OptionLongNamePattern = regexp.MustCompile(`^[a-zA-Z-._]{2,}$`)
var OptionNamePattern = regexp.MustCompile(`^(-[a-zA-Z]{1}$ | --[a-zA-Z-._]{2,})$`)
