package model

type Argument interface {
	Constrainable

	GetName() string
	GetDescription() string
	ExpectsMultipleValues() bool
}

func NewArgument(name string, description string, minCount int, maxCount int, pattern string) Argument {
	return &argument{
		minCount:    minCount,
		maxCount:    maxCount,
		pattern:     pattern,
		name:        name,
		description: description}
}

type argument struct {
	minCount    int
	maxCount    int
	pattern     string
	name        string
	description string
}

// Constrainable interface
func (a *argument) GetMaxValuesCount() int { return a.maxCount }
func (a *argument) GetMinValuesCount() int { return a.minCount }
func (a *argument) GetPattern() string     { return a.pattern }

// Argument interface
func (a *argument) GetName() string             { return a.name }
func (a *argument) GetDescription() string      { return a.description }
func (a *argument) ExpectsMultipleValues() bool { return a.maxCount > 1 }
