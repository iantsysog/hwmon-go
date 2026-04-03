package main

import (
	"bytes"
	"errors"
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

func TestParseArgsErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "missing key", args: []string{}},
		{name: "bad key length", args: []string{"-k", "ABC"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseArgs(tc.args)
			must.NotEq(t, nil, err)
		})
	}
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

func TestRunInvalidHex(t *testing.T) {
	conn := smc.NewMockConnection()
	err := run(options{key: "ABCD", valueHex: "x"}, conn, &bytes.Buffer{})
	must.NotEq(t, nil, err)
}

type closeErrorConn struct {
	smc.Connection
	closeErr error
}

func (c closeErrorConn) Close() error {
	if c.closeErr != nil {
		return c.closeErr
	}
	return c.Connection.Close()
}

func TestRunReturnsCloseError(t *testing.T) {
	base := smc.NewMockConnection()
	must.NoError(t, base.Write("ABCD", []byte{0x01}))

	expected := errors.New("close failed")
	conn := closeErrorConn{
		Connection: base,
		closeErr:   expected,
	}

	err := run(options{key: "ABCD"}, conn, &bytes.Buffer{})
	must.Eq(t, expected, err)
}
