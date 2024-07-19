package grpt

import (
	"errors"
	"fmt"
)

type ErrorCode int

type Error interface {
	Code() ErrorCode
	Cause() error
	Wrap(err error) error
	error
}

var (
	ErrInvalidArgument   = newError(1, "invalid argument")
	ErrElementOverflow   = newError(2, "element doesn't fit its parent")
	ErrInvalidAxis       = newError(3, "invalid axis")
	ErrInvalidSize       = newError(3, "invalid size")
	ErrInvalidOffset     = newError(4, "invalid offset")
	ErrInvalidBorderSide = newError(5, "invalid border side")
	ErrElementRender     = newError(6, "element can't be renderized")
)

type baseError struct {
	code    ErrorCode
	message string
	cause   error
}

func newError(code ErrorCode, message string) Error {
	return &baseError{
		code:    code,
		message: message,
	}
}

func (e baseError) Code() ErrorCode {
	return 0
}

func (e baseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s", e.message, e.cause.Error())
	}
	return e.message
}

func (e baseError) Cause() error {
	return e.cause
}

func (e baseError) Wrap(err error) error {
	e.cause = err
	return e
}

func (e baseError) Unwrap() error {
	return e.cause
}

func (e baseError) Is(err error) bool {
	var target Error
	if !errors.As(err, &target) {
		return false
	}
	return e.code == target.Code()
}
