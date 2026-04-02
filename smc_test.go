//go:build darwin && cgo

package smc

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestSMC(t *testing.T) {
	c := New()
	if err := c.Open(); err != nil {
		t.Skipf("SMC open failed: %v", err)
	}
	defer c.Close()

	if err := c.Write("CH0B", []byte{0x0}); err != nil {
		t.Skipf("SMC write failed: %v", err)
	}

	v, err := c.Read("CH0B")
	must.NoError(t, err)
	must.Eq(t, "CH0B", v.Key)
	must.Eq(t, "hex_", v.DataType)
	must.Eq(t, []uint8{0x0}, v.Bytes)
}
