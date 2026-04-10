package hwmon

type Reading struct {
	Value    any
	Kind     Kind
	Name     string
	Unit     string
	Source   string
	KeyOrID  string
	DataType string
	Raw      []byte
}
