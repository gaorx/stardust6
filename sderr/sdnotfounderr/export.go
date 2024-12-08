package sdnotfounderr

import (
	"database/sql"
	"github.com/gaorx/stardust6/sdcache"
	"github.com/gaorx/stardust6/sderr"
	"github.com/redis/go-redis/v9"
)

var notFoundErrs = []error{
	sql.ErrNoRows,
	sdcache.ErrNotFound,
	redis.Nil,
}

// Register 注册一个新的错误，用于判断是否为NotFound语义的错误
// 例如注册gorm.ErrRecordNotFound
func Register(err error) {
	if err == nil {
		return
	}
	for _, err0 := range notFoundErrs {
		if sderr.Is(err, err0) {
			return
		}
	}
	notFoundErrs = append(notFoundErrs, err)
}

// Is 判断是否为NotFound语义的错误
func Is(err error) bool {
	for _, notFoundErr := range notFoundErrs {
		if sderr.Is(err, notFoundErr) {
			return true
		}
	}
	return false
}
