package sdhttpua

import (
	"github.com/samber/lo"
)

// All 所有的 User-Agent
var All = lo.Must(FromLines(allRaw, false))
