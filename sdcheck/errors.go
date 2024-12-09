package sdcheck

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
)

var Error = sderr.Sentinel("check failed")

func errorOf(message any) error {
	switch v := message.(type) {
	case nil:
		return sderr.Wrap(Error)
	case string:
		return sderr.Wrapf(Error, v)
	case error:
		return v
	case func() error:
		return v()
	case func() string:
		return sderr.Wrapf(Error, v())
	case fmt.Stringer:
		return sderr.Wrapf(Error, v.String())
	default:
		return sderr.Wrap(Error)
	}
}
