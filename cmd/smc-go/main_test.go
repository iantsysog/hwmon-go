package main

import (
	"bytes"
	"testing"

	smc "github.com/iantsysog/smc-go"
	"github.com/shoenig/test/must"
)

func TestParseArgs(t *testing.T) {
	opts, err := parseArgs([]string{"-k", "ABCD"})
	must.NoError(t, err)
	must.Eq(t, "ABCD", opts.key)
	must.Eq(t, "", opts.valueHex)
}

func TestRunRead(t *testing.T) {
	conn := smc.NewMockConnection()
	must.NoError(t, conn.Write("ABCD", []byte{0x01, 0x02}))

	var out bytes.Buffer
	err := run(options{key: "ABCD"}, conn, &out)
	must.NoError(t, err)
	must.Eq(t, "0102\n", out.String())
}

func TestRunWrite(t *testing.T) {
	conn := smc.NewMockConnection()

	err := run(options{key: "ABCD", valueHex: "0a0b"}, conn, &bytes.Buffer{})
	must.NoError(t, err)

	v, err := conn.Read("ABCD")
	must.NoError(t, err)
	must.Eq(t, []byte{0x0a, 0x0b}, v.Bytes)
}
