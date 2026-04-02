package smc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

var (
	ErrInsufficientBytes = errors.New("insufficient bytes for conversion")
	ErrUnsupportedSize   = errors.New("unsupported byte length for conversion")
)

func (v SMCVal) Bool() (bool, error) {
	if len(v.Bytes) < 1 {
		return false, ErrInsufficientBytes
	}
	return v.Bytes[0] != 0, nil
}

func (v SMCVal) Uint(order binary.ByteOrder) (uint64, error) {
	switch len(v.Bytes) {
	case 1:
		return uint64(v.Bytes[0]), nil
	case 2:
		return uint64(order.Uint16(v.Bytes)), nil
	case 4:
		return uint64(order.Uint32(v.Bytes)), nil
	case 8:
		return order.Uint64(v.Bytes), nil
	default:
		return 0, ErrUnsupportedSize
	}
}

func (v SMCVal) Int(order binary.ByteOrder) (int64, error) {
	u, err := v.Uint(order)
	if err != nil {
		return 0, err
	}

	switch len(v.Bytes) {
	case 1:
		return int64(int8(u)), nil
	case 2:
		return int64(int16(u)), nil
	case 4:
		return int64(int32(u)), nil
	case 8:
		return int64(u), nil
	default:
		return 0, ErrUnsupportedSize
	}
}

func (v SMCVal) UintLE() (uint64, error) { return v.Uint(binary.LittleEndian) }
func (v SMCVal) UintBE() (uint64, error) { return v.Uint(binary.BigEndian) }
func (v SMCVal) IntLE() (int64, error)   { return v.Int(binary.LittleEndian) }
func (v SMCVal) IntBE() (int64, error)   { return v.Int(binary.BigEndian) }

func (v SMCVal) Float32(order binary.ByteOrder) (float32, error) {
	if len(v.Bytes) < 4 {
		return 0, ErrInsufficientBytes
	}
	return math.Float32frombits(order.Uint32(v.Bytes[:4])), nil
}

func (v SMCVal) Float64(order binary.ByteOrder) (float64, error) {
	if len(v.Bytes) < 8 {
		return 0, ErrInsufficientBytes
	}
	return math.Float64frombits(order.Uint64(v.Bytes[:8])), nil
}

func (v SMCVal) Float32LE() (float32, error) { return v.Float32(binary.LittleEndian) }
func (v SMCVal) Float32BE() (float32, error) { return v.Float32(binary.BigEndian) }
func (v SMCVal) Float64LE() (float64, error) { return v.Float64(binary.LittleEndian) }
func (v SMCVal) Float64BE() (float64, error) { return v.Float64(binary.BigEndian) }

func (v SMCVal) String() string {
	return string(v.Bytes)
}

func (v SMCVal) CString() string {
	if len(v.Bytes) == 0 {
		return ""
	}
	if i := bytes.IndexByte(v.Bytes, 0); i >= 0 {
		return string(v.Bytes[:i])
	}
	return string(v.Bytes)
}

func BytesUint16(order binary.ByteOrder, n uint16) []byte {
	b := make([]byte, 2)
	order.PutUint16(b, n)
	return b
}

func BytesUint32(order binary.ByteOrder, n uint32) []byte {
	b := make([]byte, 4)
	order.PutUint32(b, n)
	return b
}

func BytesUint64(order binary.ByteOrder, n uint64) []byte {
	b := make([]byte, 8)
	order.PutUint64(b, n)
	return b
}

func BytesInt16(order binary.ByteOrder, n int16) []byte { return BytesUint16(order, uint16(n)) }
func BytesInt32(order binary.ByteOrder, n int32) []byte { return BytesUint32(order, uint32(n)) }
func BytesInt64(order binary.ByteOrder, n int64) []byte { return BytesUint64(order, uint64(n)) }
func BytesFloat32(order binary.ByteOrder, f float32) []byte {
	return BytesUint32(order, math.Float32bits(f))
}

func BytesFloat64(order binary.ByteOrder, f float64) []byte {
	return BytesUint64(order, math.Float64bits(f))
}
