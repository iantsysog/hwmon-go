package source

import (
	"errors"

	"github.com/iantsysog/hwmon-go/internal/hid/model"
)

var ErrUnsupported = errors.New("HID sensors are only supported on darwin")

type Source interface {
	Open() error
	Close() error
	Readings() ([]model.Reading, error)
}

func Default() Source { return defaultSource() }
