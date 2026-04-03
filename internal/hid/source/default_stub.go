//go:build !darwin || !cgo || (!iohid && !eventsystem)

package source

func defaultSource() Source { return unsupportedSource{} }
