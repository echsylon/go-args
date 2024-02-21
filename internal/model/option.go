package model

type Option interface {
	Constrainable

	IsParsed() bool
	SetParsed()
	IsHelpTrigger() bool
	GetShortName() string
	GetLongName() string
	GetDescription() string
}

func NewOption(shortName string, longName string, description string, pattern string) Option {
	return &option{
		shortName:   shortName,
		longName:    longName,
		pattern:     pattern,
		description: description,
		parsed:      false,
		help:        false,
	}
}

func NewHelpOption(shortName string, longName string, description string) Option {
	return &option{
		shortName:   shortName,
		longName:    longName,
		pattern:     "",
		description: description,
		parsed:      false,
		help:        true,
	}
}

type option struct {
	shortName   string
	longName    string
	pattern     string
	description string
	parsed      bool
	help        bool
}

// Constrainable interface
func (o *option) GetMinValuesCount() int { return 0 }
func (o *option) GetMaxValuesCount() int { return 1 }
func (o *option) GetPattern() string     { return o.pattern }

// Option interface
func (o *option) IsParsed() bool         { return o.parsed }
func (o *option) SetParsed()             { o.parsed = true }
func (o *option) IsHelpTrigger() bool    { return o.help }
func (o *option) GetShortName() string   { return o.shortName }
func (o *option) GetLongName() string    { return o.longName }
func (o *option) GetDescription() string { return o.description }
