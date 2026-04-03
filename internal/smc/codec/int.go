package codec

import (
	"encoding/binary"
	"math"
)

func fixedIntDecoder(dataType string, order binary.ByteOrder, size int, signed bool) Decoder {
	return func(raw []byte) (any, error) {
		if len(raw) != size {
			return nil, ErrInvalidLength{DataType: dataType, Got: len(raw), Want: size}
		}

		switch size {
		case 1:
			if signed {
				return int8(raw[0]), nil
			}
			return uint8(raw[0]), nil
		case 2:
			u := order.Uint16(raw)
			if signed {
				return int16(u), nil
			}
			return u, nil
		case 4:
			u := order.Uint32(raw)
			if signed {
				return int32(u), nil
			}
			return u, nil
		case 8:
			u := order.Uint64(raw)
			if signed {
				return int64(u), nil
			}
			return u, nil
		default:
			return nil, ErrInvalidLength{DataType: dataType, Got: len(raw), Want: size}
		}
	}
}

func fixedIntEncoder(dataType string, order binary.ByteOrder, size int, signed bool) Encoder {
	return func(v any) ([]byte, error) {
		out := make([]byte, size)
		switch size {
		case 1:
			if signed {
				x, ok := asInt64(v)
				if !ok || x < math.MinInt8 || x > math.MaxInt8 {
					if ok {
						return nil, ErrOutOfRange{DataType: dataType, Value: v}
					}
					return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
				}
				out[0] = byte(int8(x))
				return out, nil
			}

			x, ok := asUint64(v)
			if !ok || x > math.MaxUint8 {
				if ok {
					return nil, ErrOutOfRange{DataType: dataType, Value: v}
				}
				return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
			}
			out[0] = byte(uint8(x))
			return out, nil
		case 2:
			if signed {
				x, ok := asInt64(v)
				if !ok || x < math.MinInt16 || x > math.MaxInt16 {
					if ok {
						return nil, ErrOutOfRange{DataType: dataType, Value: v}
					}
					return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
				}
				order.PutUint16(out, uint16(int16(x)))
				return out, nil
			}

			x, ok := asUint64(v)
			if !ok || x > math.MaxUint16 {
				if ok {
					return nil, ErrOutOfRange{DataType: dataType, Value: v}
				}
				return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
			}
			order.PutUint16(out, uint16(x))
			return out, nil
		case 4:
			if signed {
				x, ok := asInt64(v)
				if !ok || x < math.MinInt32 || x > math.MaxInt32 {
					if ok {
						return nil, ErrOutOfRange{DataType: dataType, Value: v}
					}
					return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
				}
				order.PutUint32(out, uint32(int32(x)))
				return out, nil
			}

			x, ok := asUint64(v)
			if !ok || x > math.MaxUint32 {
				if ok {
					return nil, ErrOutOfRange{DataType: dataType, Value: v}
				}
				return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
			}
			order.PutUint32(out, uint32(x))
			return out, nil
		case 8:
			if signed {
				x, ok := asInt64(v)
				if !ok {
					return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
				}
				order.PutUint64(out, uint64(x))
				return out, nil
			}

			x, ok := asUint64(v)
			if !ok {
				return nil, ErrUnsupportedValueType{DataType: dataType, Value: v}
			}
			order.PutUint64(out, x)
			return out, nil
		default:
			return nil, ErrUnsupportedType{DataType: dataType}
		}
	}
}

func asInt64(v any) (int64, bool) {
	switch x := v.(type) {
	case int:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	case int64:
		return x, true
	case uint:
		return int64(x), true
	case uint8:
		return int64(x), true
	case uint16:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		if x > math.MaxInt64 {
			return 0, false
		}
		return int64(x), true
	default:
		return 0, false
	}
}

func asUint64(v any) (uint64, bool) {
	switch x := v.(type) {
	case int:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int8:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int16:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int32:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int64:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case uint:
		return uint64(x), true
	case uint8:
		return uint64(x), true
	case uint16:
		return uint64(x), true
	case uint32:
		return uint64(x), true
	case uint64:
		return x, true
	default:
		return 0, false
	}
}
