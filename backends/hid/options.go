package hid

import "github.com/iantsysog/hwmon-go/internal/hid/source"

type config struct {
	src              source.Source
	emitUnclassified bool
}

type Option func(*config)

func WithEmitUnclassified(enabled bool) Option {
	return func(c *config) { c.emitUnclassified = enabled }
}
