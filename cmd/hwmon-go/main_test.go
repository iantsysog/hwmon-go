package main

import "testing"

func TestParseArgs(t *testing.T) {
	_, err := parseArgs([]string{"-backends", "smc,hid"})
	if err != nil {
		t.Fatalf("parseArgs() returned an error: %v", err)
	}
}

func TestParseArgsEmptyBackendsAllowed(t *testing.T) {
	_, err := parseArgs([]string{"-backends", ""})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
