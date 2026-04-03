//go:build darwin && cgo && eventsystem && !iohid

package source

func defaultIOHIDSource() Source { return nil }
