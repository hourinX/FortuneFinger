package errorx

import (
	"errors"
	"fmt"
)

type Error struct {
	code Code
	msg  string
	err  error
}

func New(code Code, msg string) *Error {
	return newError(code, msg, nil)
}

func Wrap(code Code, msg string, err error) *Error {
	return newError(code, msg, err)
}

func newError(code Code, msg string, err error) *Error {
	resolvedCode, def := lookup(code)
	if msg == "" || resolvedCode != code {
		msg = def.message
	}
	return &Error{
		code: resolvedCode,
		msg:  msg,
		err:  err,
	}
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *Error) Code() Code {
	if e == nil {
		return InternalServer
	}
	return e.code
}

func (e *Error) Message() string {
	if e == nil {
		return definitions[InternalServer].message
	}
	return e.msg
}

func As(err error) (*Error, bool) {
	var target *Error
	if !errors.As(err, &target) {
		return nil, false
	}
	return target, true
}
