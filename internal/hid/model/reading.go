package model

type Reading struct {
	Kind    Kind
	Name    string
	Unit    string
	Source  string
	KeyOrID string
	Value   float64
}
