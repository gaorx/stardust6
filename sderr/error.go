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

var _ error = Error{}

func New(msg string) error {
	return newBuilder().New(msg)
}

func Newf(msg string, a ...any) error {
	return newBuilder().Newf(msg, a...)
}

func Wrap(err error, msg string) error {
	return newBuilder().Wrap(err, msg)
}

func Wrapf(err error, msg string, a ...any) error {
	return newBuilder().Wrapf(err, msg, a...)
}

func (e Error) Error() string {
	if e.err == nil {
		return e.msg
	} else {
		if e.msg == "" {
			return e.err.Error()
		} else {
			return e.msg + ": " + e.err.Error()
		}
	}
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Is(err error) bool {
	return e.err == err
}

func (e Error) Message() string {
	return e.msg
}

func (e Error) Attrs() map[string]any {
	merged := map[string]any{}
	mergeAttrs(merged, e)
	return merged
}

func (e Error) Stack() *Stack {
	return e.stack
}

func mergeAttrs(dst map[string]any, e Error) {
	ee := e.err
	if ee != nil {
		if ee2, ok := ee.(Error); ok {
			mergeAttrs(dst, ee2)
		}
	}
	maps.Copy(dst, e.attrs)
}
