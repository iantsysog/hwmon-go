package hwmon

type Reading struct {
	Kind    Kind
	Name    string
	Unit    string
	Source  string
	KeyOrID string
	DataType string
	Raw []byte
	Value any
}
