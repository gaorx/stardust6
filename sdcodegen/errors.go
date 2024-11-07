package sdcodegen

import (
	"github.com/gaorx/stardust6/sderr"
	"io/fs"
)

// IsNotExistErr 判断ReadFile返回的error是不是因为文件不存在
func IsNotExistErr(err error) bool {
	return sderr.Is(err, fs.ErrNotExist)
}
