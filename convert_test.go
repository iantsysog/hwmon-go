package smc

import (
	"encoding/binary"
	"testing"

	"github.com/shoenig/test/must"
)

func TestSMCValUintInt(t *testing.T) {
	v := SMCVal{Bytes: []byte{0x01, 0x00}}
	u, err := v.Uint(binary.LittleEndian)
	must.NoError(t, err)
	must.Eq(t, uint64(1), u)

	i, err := v.Int(binary.LittleEndian)
	must.NoError(t, err)
	must.Eq(t, int64(1), i)
}

func TestSMCValFloat32(t *testing.T) {
	v := SMCVal{Bytes: BytesFloat32(binary.LittleEndian, 12.5)}
	f, err := v.Float32LE()
	must.NoError(t, err)
	must.Eq(t, float32(12.5), f)
}

func TestSMCValCString(t *testing.T) {
	v := SMCVal{Bytes: []byte{'h', 'i', 0x00, 'x'}}
	must.Eq(t, "hi", v.CString())
	must.Eq(t, "hi\x00x", v.String())
}

func TestSMCValErrors(t *testing.T) {
	_, err := (SMCVal{}).Bool()
	must.Eq(t, ErrInsufficientBytes, err)

	_, err = (SMCVal{Bytes: []byte{0x01, 0x02, 0x03}}).Uint(binary.BigEndian)
	must.Eq(t, ErrUnsupportedSize, err)

	_, err = (SMCVal{Bytes: []byte{0x01, 0x02}}).Float32LE()
	must.Eq(t, ErrInsufficientBytes, err)
}
