//go:build darwin

package platform

import (
	"strings"
	"syscall"
)

func model() (string, error) {
	m, err := syscall.Sysctl("hw.model")
	if err != nil {
		return "", ErrUnsupported
	}
	m = strings.TrimRight(m, "\x00")
	if m == "" {
		return "", ErrUnsupported
	}
	return m, nil
}

func family() (string, error) {
	m, err := model()
	if err != nil || m == "" {
		return "Unknown", err
	}

	if f, ok := lookupFamily(m); ok {
		return f, nil
	}
	return "Unknown", nil
}
