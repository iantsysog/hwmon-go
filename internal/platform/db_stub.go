//go:build !platform

package platform

func lookupFamily(_ string) (string, bool) { return "", false }
