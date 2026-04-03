package model

import (
	"errors"
	"fmt"
	"strconv"
)

type SMCVal struct {
	Key      string
	DataType string
	Bytes    []byte
}

var (
	ErrInvalidKey       = errors.New("SMC key must be exactly 4 characters")
	ErrNoDataForKey     = errors.New("no data returned for SMC key")
	ErrConnectionClosed = errors.New("SMC connection is not open")
	ErrInvalidDataSize  = errors.New("invalid SMC data size")
	ErrUnsupported      = errors.New("SMC is only supported on darwin with cgo enabled")
	ErrSMCFailure       = errors.New("SMC operation failed")
)

type SMCError struct {
	Op   string
	Key  string
	Code int
}

func (e *SMCError) Error() string {
	if e == nil {
		return ErrSMCFailure.Error()
	}
	code := strconv.Itoa(e.Code)
	if e.Key == "" {
		return e.Op + ": ret=" + code
	}
	return fmt.Sprintf("%s %q: ret=%s", e.Op, e.Key, code)
}

func (e *SMCError) Unwrap() error { return ErrSMCFailure }
