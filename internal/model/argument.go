package model

type Argument interface {
	GetMaxValuesCount() int
	GetMinValuesCount() int
	ExpectsMultipleValues() bool
	GetName() string
	GetPattern() string
	GetDescription() string
}

func NewArgument(minCount int, maxCount int, name string, pattern string, description string) Argument {
	return &argument{
		minCount:    minCount,
		maxCount:    maxCount,
		name:        name,
		pattern:     pattern,
		description: description}
}

type argument struct {
	minCount    int
	maxCount    int
	name        string
	pattern     string
	description string
}

func (a *argument) GetMaxValuesCount() int      { return a.maxCount }
func (a *argument) GetMinValuesCount() int      { return a.minCount }
func (a *argument) ExpectsMultipleValues() bool { return a.maxCount > 1 }
func (a *argument) GetName() string             { return a.name }
func (a *argument) GetPattern() string          { return a.pattern }
func (a *argument) GetDescription() string      { return a.description }
