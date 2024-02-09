package model

type Option interface {
	IsParsed() bool
	SetParsed()
	GetShortName() string
	GetLongName() string
	GetPattern() string
	GetDescription() string
}

func NewOption(shortName string, longName string, pattern string, description string) Option {
	return &option{
		shortName:   shortName,
		longName:    longName,
		pattern:     pattern,
		description: description,
	}
}

type option struct {
	shortName   string
	longName    string
	pattern     string
	description string
	parsed      bool
}

func (o *option) IsParsed() bool         { return o.parsed }
func (o *option) SetParsed()             { o.parsed = true }
func (o *option) GetShortName() string   { return o.shortName }
func (o *option) GetLongName() string    { return o.longName }
func (o *option) GetPattern() string     { return o.pattern }
func (o *option) GetDescription() string { return o.description }
