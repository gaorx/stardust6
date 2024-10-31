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

// Newf 构造一个带有stacktrace和一个可读消息的error
func Newf(msg string, a ...any) error {
	return newBuilder().Newf(msg, a...)
}

// Wrap wrap一个error，附加上stacktrace
func Wrap(err error) error {
	return newBuilder().Wrap(err)
}

// Wrapf wrap一个error，附加上stacktrace和一个可读消息
func Wrapf(err error, msg string, a ...any) error {
	return newBuilder().Wrapf(err, msg, a...)
}

// Wrap2 用于wrap一个函数的返回值
func Wrap2[T1 any](a T1, err error) (T1, error) {
	return a, Wrap(err)
}

// Wrap3 用于wrap一个函数的返回值
func Wrap3[T1, T2 any](a T1, b T2, err error) (T1, T2, error) {
	return a, b, Wrap(err)
}

// Wrap4 用于wrap一个函数的返回值
func Wrap4[T1, T2, T3 any](a T1, b T2, c T3, err error) (T1, T2, T3, error) {
	return a, b, c, Wrap(err)
}

// Recover 构建一个error，通过一个函数的panic信息
func Recover(f func()) error {
	return newBuilder().Recover(f)
}

// Recoverf 构建一个error，通过一个函数的panic信息和一个可读消息
func Recoverf(f func(), msg string, a ...any) error {
	return newBuilder().Recoverf(f, msg, a...)
}

// Error 实现error.Error
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	causeMsg := func() string {
		if e.err == nil {
			return ""
		}
		return e.err.Error()
	}
	if e.msg == "" {
		return makeMsgWithAttrs(causeMsg(), e.attrs)
	} else {
		msg1, msg2 := makeMsgWithAttrs(e.msg, e.attrs), causeMsg()
		if msg1 != "" && msg2 != "" {
			return msg1 + ": " + msg2
		} else if msg1 != "" {
			return msg1
		} else {
			return msg2
		}
	}
}

// Unwrap 兼容errors包中的Unwrap协议
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// Msg 返回此error的可读性信息
func (e *Error) Msg() string {
	if e == nil {
		return ""
	}
	return e.msg
}

// Stack 返回此error的stacktrace(仅含有此层的stack信息)
func (e *Error) Stack() *Stack {
	if e == nil {
		return nil
	}
	return e.stack
}

// OwnAttrs 返回此层error的Attrs
func (e *Error) OwnAttrs() map[string]any {
	if e == nil {
		return nil
	}
	return maps.Clone(e.attrs)
}
