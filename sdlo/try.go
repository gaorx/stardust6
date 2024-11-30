package sdlo

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
)

func TryWithPanicError(callback func()) error {
	err, ok := lo.TryWithErrorValue(func() error {
		callback()
		return nil
	})
	if !ok {
		return sderr.Ensure(err)
	}
	return nil
}
