package sdgogen

import (
	"github.com/gaorx/stardust6/sdcodegen"
)

func Formatter() sdcodegen.Formatter {
	return sdcodegen.FormatByCmd("gofmt", nil)
}
