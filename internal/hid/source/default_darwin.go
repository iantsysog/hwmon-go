//go:build darwin && cgo && (iohid || eventsystem)

package source

func defaultSource() Source {
	if s := defaultEventSystemSource(); s != nil {
		return s
	}
	if s := defaultIOHIDSource(); s != nil {
		return s
	}
	return unsupportedSource{}
}
