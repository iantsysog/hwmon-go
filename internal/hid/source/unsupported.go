package source

import "github.com/iantsysog/hwmon-go/internal/hid/model"

type unsupportedSource struct{}

func (unsupportedSource) Open() error                        { return ErrUnsupported }
func (unsupportedSource) Close() error                       { return nil }
func (unsupportedSource) Readings() ([]model.Reading, error) { return nil, ErrUnsupported }
