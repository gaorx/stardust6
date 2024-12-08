package sdbunsqlite

import (
	"database/sql"
	"github.com/gaorx/stardust6/sdbun"
	"github.com/gaorx/stardust6/sderr"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"log/slog"
)

// Address 数据库地址
type Address struct {
	// DataSourceName
	DSN string `json:"dsn" toml:"dsn" yaml:"dsn"`
	// 如果不为nil的话，使用此slog作为logger实现
	Slog *slog.Logger `json:"-" toml:"-" yaml:"-"`
	// 如果Slog为nil的话，这个字符串代表logger实现名称
	Logger string `json:"logger" toml:"logger"`
}

// Open 打开数据库
func Open(addr Address, opts ...bun.DBOption) (*bun.DB, error) {
	sqlDB, err := sql.Open(sqliteshim.ShimName, addr.DSN)
	if err != nil {
		return nil, sderr.Wrapf(err, "open sqlite error")
	}
	db := bun.NewDB(sqlDB, sqlitedialect.New(), opts...)
	addr.apply(db)
	return db, nil
}

func (addr *Address) apply(db *bun.DB) {
	var logger bun.QueryHook
	if addr.Slog != nil {
		logger = sdbun.Slog(addr.Slog)
	} else {
		logger = sdbun.LoggerOf(addr.Logger)
	}
	if logger != nil {
		db.AddQueryHook(logger)
	}
}
