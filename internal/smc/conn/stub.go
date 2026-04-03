//go:build !darwin || !cgo

package conn

import "github.com/iantsysog/hwmon-go/internal/smc/model"

type unsupportedConnection struct{}

func New() Connection { return unsupportedConnection{} }

func (unsupportedConnection) Open() error  { return model.ErrUnsupported }
func (unsupportedConnection) Close() error { return nil }
func (unsupportedConnection) Read(string) (model.SMCVal, error) {
	return model.SMCVal{}, model.ErrUnsupported
}
func (unsupportedConnection) Write(string, []byte) error { return model.ErrUnsupported }
