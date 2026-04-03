//go:build darwin && cgo && iohid && !eventsystem

package source

func defaultEventSystemSource() Source { return nil }
