//go:build hid

package hid

import "github.com/iantsysog/hwmon-go/internal/hid/source"

func withSource(src source.Source) Option { return func(c *config) { c.src = src } }
