package model

type Reading struct {
	Kind    Kind
	Name    string
	Unit    string
	Source  string
	Value   float64
	KeyOrID string
}
