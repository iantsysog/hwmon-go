package platform

import "errors"

var ErrUnsupported = errors.New("platform information is only supported on darwin")

func Model() (string, error) { return model() }

func Family() (string, error) { return family() }
