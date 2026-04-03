//go:build !darwin

package platform

func model() (string, error) { return "", ErrUnsupported }

func family() (string, error) { return "Unknown", ErrUnsupported }
