package conn

import "github.com/iantsysog/hwmon-go/internal/smc/model"

type Connection interface {
	Open() error
	Close() error
	Read(key string) (model.SMCVal, error)
	Write(key string, value []byte) error
}
