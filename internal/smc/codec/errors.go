package codec

import "fmt"

type ErrUnsupportedType struct {
	DataType string
}

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported SMC data type: %q", e.DataType)
}

type ErrInvalidLength struct {
	DataType string
	Got      int
	Want     int
}

func (e ErrInvalidLength) Error() string {
	return fmt.Sprintf("invalid length for %q: got=%d want=%d", e.DataType, e.Got, e.Want)
}

type ErrOutOfRange struct {
	DataType string
	Value    any
}

func (e ErrOutOfRange) Error() string {
	return fmt.Sprintf("value out of range for %q: %v", e.DataType, e.Value)
}

type ErrUnsupportedValueType struct {
	DataType string
	Value    any
}

func (e ErrUnsupportedValueType) Error() string {
	return fmt.Sprintf("unsupported value type for %q: %T", e.DataType, e.Value)
}
