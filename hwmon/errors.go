package hwmon

import "errors"

var (
	ErrNoBackends     = errors.New("no backends are enabled")
	ErrUnknownBackend = errors.New("unknown backend")
)
