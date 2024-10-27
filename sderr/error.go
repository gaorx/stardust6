package sderr

import (
	"maps"
)

type Error struct {
	err   error
	msg   string
	stack *Stack
	attrs map[string]any
}

var _ error = &Error{}

func Newf(msg string, a ...any) error {
	return newBuilder().Newf(msg, a...)
}

func Wrap(err error) error {
	return newBuilder().Wrap(err)
}

func Wrapf(err error, msg string, a ...any) error {
	return newBuilder().Wrapf(err, msg, a...)
}

func Wrap2[T1 any](a T1, err error) (T1, error) {
	return a, Wrap(err)
}

func Wrap3[T1, T2 any](a T1, b T2, err error) (T1, T2, error) {
	return a, b, Wrap(err)
}

func Wrap4[T1, T2, T3 any](a T1, b T2, c T3, err error) (T1, T2, T3, error) {
	return a, b, c, Wrap(err)
}

func Recover(f func()) error {
	return newBuilder().Recover(f)
}

func Recoverf(f func(), msg string, a ...any) error {
	return newBuilder().Recoverf(f, msg, a...)
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}
	if e.msg == "" {
		return e.err.Error()
	}
	return e.msg + ": " + e.err.Error()
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.msg
}

func (e *Error) Stack() *Stack {
	if e == nil {
		return nil
	}
	return e.stack
}

func (e *Error) OwnAttrs() map[string]any {
	if e == nil {
		return nil
	}
	return maps.Clone(e.attrs)
}
