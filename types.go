package smc

import (
	"errors"
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
	if e.Key == "" {
		return e.Op + ": ret=" + strconv.Itoa(e.Code)
	}
	return e.Op + " " + `"` + e.Key + `"` + ": ret=" + strconv.Itoa(e.Code)
}

func (e *SMCError) Unwrap() error { return ErrSMCFailure }
