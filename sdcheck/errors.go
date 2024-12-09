package sdcheck

import (
	"fmt"
	"github.com/gaorx/stardust6/sdlo"
)

type Error struct {
	Message  string
	Internal error
}

func (e Error) Error() string {
	msg := sdlo.EmptyOr(e.Message, "check failed")
	if e.Internal == nil {
		return msg
	}
	internalMessage := e.Internal.Error()
	return fmt.Sprintf("%s: %s", msg, internalMessage)
}

func newErr(message string, internal error) Error {
	return Error{Message: message, Internal: internal}
}

func errorOf(message any) error {
	switch v := message.(type) {
	case nil:
		return newErr("", nil)
	case string:
		return newErr(v, nil)
	case error:
		return newErr("", v)
	case func() error:
		var err0 error
		if v != nil {
			err0 = v()
		}
		return newErr("", err0)
	case func() string:
		var msg0 string
		if v != nil {
			msg0 = v()
		}
		return newErr(msg0, nil)
	case fmt.Stringer:
		var msg0 string
		if v != nil {
			msg0 = v.String()
		}
		return newErr(msg0, nil)
	default:
		return newErr("", nil)
	}
}
