//go:build !darwin || !cgo

package smc

type unsupportedConnection struct{}

func New() Connection { return unsupportedConnection{} }

func (unsupportedConnection) Open() error  { return ErrUnsupported }
func (unsupportedConnection) Close() error { return nil }
func (unsupportedConnection) Read(string) (SMCVal, error) {
	return SMCVal{}, ErrUnsupported
}
func (unsupportedConnection) Write(string, []byte) error { return ErrUnsupported }
