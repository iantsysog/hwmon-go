package hid

type config struct {
	emitUnclassified bool
}

type Option func(*config)

func WithEmitUnclassified(enabled bool) Option {
	return func(c *config) { c.emitUnclassified = enabled }
}
