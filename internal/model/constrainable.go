package model

type Constrainable interface {
	GetMinValuesCount() int
	GetMaxValuesCount() int
	GetPattern() string
}
