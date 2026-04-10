package codec

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"
	"sync"
)

type Value struct {
	Decoded  any
	Key      string
	DataType string
	Raw      []byte
}

type Decoder func(raw []byte) (any, error)
type Encoder func(v any) ([]byte, error)

type codecEntry struct {
	dec Decoder
	enc Encoder
}

var (
	regMu sync.RWMutex
	reg   = make(map[string]codecEntry, 64)
)

func Register(dataType string, dec Decoder, enc Encoder) {
	if dataType == "" {
		return
	}
	if dec == nil && enc == nil {
		return
	}
	regMu.Lock()
	reg[dataType] = codecEntry{dec: dec, enc: enc}
	regMu.Unlock()
}

func lookup(dataType string) (codecEntry, bool) {
	regMu.RLock()
	e, ok := reg[dataType]
	regMu.RUnlock()
	return e, ok
}

func Decode(dataType string, raw []byte) (any, error) {
	if strings.HasPrefix(dataType, "ch8") {
		return string(raw), nil
	}

	if e, ok := lookup(dataType); ok && e.dec != nil {
		return e.dec(raw)
	}
	return nil, ErrUnsupportedType{DataType: dataType}
}

func Encode(dataType string, v any) ([]byte, error) {
	if strings.HasPrefix(dataType, "ch8") {
		switch x := v.(type) {
		case string:
			return []byte(x), nil
		case []byte:
			out := make([]byte, len(x))
			copy(out, x)
			return out, nil
		default:
			return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
		}
	}

	if e, ok := lookup(dataType); ok && e.enc != nil {
		return e.enc(v)
	}
	return nil, ErrUnsupportedType{DataType: dataType}
}

func init() {
	Register("ui8 ", fixedIntDecoder("ui8 ", binary.BigEndian, 1, false), fixedIntEncoder("ui8 ", binary.BigEndian, 1, false))
	Register("ui16", fixedIntDecoder("ui16", binary.BigEndian, 2, false), fixedIntEncoder("ui16", binary.BigEndian, 2, false))
	Register("ui32", fixedIntDecoder("ui32", binary.BigEndian, 4, false), fixedIntEncoder("ui32", binary.BigEndian, 4, false))
	Register("ui64", fixedIntDecoder("ui64", binary.BigEndian, 8, false), fixedIntEncoder("ui64", binary.BigEndian, 8, false))

	Register("si8 ", fixedIntDecoder("si8 ", binary.BigEndian, 1, true), fixedIntEncoder("si8 ", binary.BigEndian, 1, true))
	Register("si16", fixedIntDecoder("si16", binary.BigEndian, 2, true), fixedIntEncoder("si16", binary.BigEndian, 2, true))
	Register("si32", fixedIntDecoder("si32", binary.BigEndian, 4, true), fixedIntEncoder("si32", binary.BigEndian, 4, true))
	Register("si64", fixedIntDecoder("si64", binary.BigEndian, 8, true), fixedIntEncoder("si64", binary.BigEndian, 8, true))

	Register("flag", func(raw []byte) (any, error) {
		if len(raw) != 1 {
			return nil, ErrInvalidLength{DataType: "flag", Got: len(raw), Want: 1}
		}
		return raw[0] != 0, nil
	}, func(v any) ([]byte, error) {
		switch x := v.(type) {
		case bool:
			if x {
				return []byte{1}, nil
			}
			return []byte{0}, nil
		case uint8:
			return []byte{x}, nil
		case int8:
			return []byte{byte(x)}, nil
		default:
			return nil, ErrUnsupportedValueType{DataType: "flag", Value: v}
		}
	})

	Register("flt ", func(raw []byte) (any, error) {
		if len(raw) != 4 {
			return nil, ErrInvalidLength{DataType: "flt ", Got: len(raw), Want: 4}
		}
		u := binary.LittleEndian.Uint32(raw)
		return math.Float32frombits(u), nil
	}, func(v any) ([]byte, error) {
		var f float32
		switch x := v.(type) {
		case float32:
			f = x
		case float64:
			f = float32(x)
		default:
			return nil, ErrUnsupportedValueType{DataType: "flt ", Value: v}
		}
		out := make([]byte, 4)
		binary.LittleEndian.PutUint32(out, math.Float32bits(f))
		return out, nil
	})

	Register("flt", func(raw []byte) (any, error) { return Decode("flt ", raw) }, func(v any) ([]byte, error) { return Encode("flt ", v) })

	Register("fltL", func(raw []byte) (any, error) { return Decode("flt ", raw) }, func(v any) ([]byte, error) { return Encode("flt ", v) })

	Register("fltB", func(raw []byte) (any, error) {
		if len(raw) != 4 {
			return nil, ErrInvalidLength{DataType: "fltB", Got: len(raw), Want: 4}
		}
		u := binary.BigEndian.Uint32(raw)
		return math.Float32frombits(u), nil
	}, func(v any) ([]byte, error) {
		var f float32
		switch x := v.(type) {
		case float32:
			f = x
		case float64:
			f = float32(x)
		default:
			return nil, ErrUnsupportedValueType{DataType: "fltB", Value: v}
		}
		out := make([]byte, 4)
		binary.BigEndian.PutUint32(out, math.Float32bits(f))
		return out, nil
	})

	Register("cstr", func(raw []byte) (any, error) {
		if len(raw) == 0 {
			return "", nil
		}
		if i := bytes.IndexByte(raw, 0); i >= 0 {
			return string(raw[:i]), nil
		}
		return string(raw), nil
	}, func(v any) ([]byte, error) {
		switch x := v.(type) {
		case string:
			return []byte(x), nil
		case []byte:
			out := make([]byte, len(x))
			copy(out, x)
			return out, nil
		default:
			return nil, ErrUnsupportedValueType{DataType: "cstr", Value: v}
		}
	})

	registerAppleFixed16()
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int8:
		return float64(x), true
	case int16:
		return float64(x), true
	case int32:
		return float64(x), true
	case int64:
		return float64(x), true
	case uint:
		return float64(x), true
	case uint8:
		return float64(x), true
	case uint16:
		return float64(x), true
	case uint32:
		return float64(x), true
	case uint64:
		return float64(x), true
	default:
		return 0, false
	}
}
