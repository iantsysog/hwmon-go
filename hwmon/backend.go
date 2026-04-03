package hwmon

import "context"

type Backend interface {
	Name() string
	Collect(ctx context.Context, emit func(Reading)) error
}
