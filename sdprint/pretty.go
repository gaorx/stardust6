package sdprint

import (
	"github.com/k0kubun/pp/v3"
)

func init() {
	defaultPP := pp.Default
	defaultPP.SetColoringEnabled(false)
	defaultPP.SetExportedOnly(false)
	defaultPP.SetOmitEmpty(false)
}

// SetPrettyColorEnabled 设置全局打印pretty格式时，是否启用带有颜色的输出
func SetPrettyColorEnabled(enabled bool) {
	pp.Default.SetColoringEnabled(enabled)
}
