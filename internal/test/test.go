package test

import (
	"errors"
	"reflect"
)

type TB interface {
	Helper()
	Fatalf(format string, args ...any)
	Logf(format string, args ...any)
}

func NoError(t TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Error(t TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func ErrorIs(t TB, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected errors.Is(err, %v) to be true; err=%v", target, err)
	}
}

func True(t TB, v bool) {
	t.Helper()
	if !v {
		t.Fatalf("expected condition to be true")
	}
}

func Len[T any](t TB, s []T, want int) {
	t.Helper()
	if got := len(s); got != want {
		t.Fatalf("expected length %d, got %d", want, got)
	}
}

func Eq(t TB, want, got any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Logf("want: %#v", want)
		t.Logf("got:  %#v", got)
		t.Fatalf("values are not equal")
	}
}
