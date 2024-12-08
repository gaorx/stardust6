package sdbunmysql

import (
	"database/sql"
	"github.com/gaorx/stardust6/sdbun"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdtime"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
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
	// 连接最大生命周期
	ConnMaxLifeTimeMS int64 `json:"conn_max_lifetime" toml:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	// 连接最大空闲时间
	ConnMaxIdleTimeMS int64 `json:"conn_max_idle_time" toml:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	// 最大空闲连接数
	MaxIdleConns int `json:"max_idle_conns" toml:"max_idle_conns" yaml:"max_idle_conns"`
	// 最大打开连接数
	MaxOpenConns int `json:"max_open_conns" toml:"max_open_conns" yaml:"max_open_conns"`
}

// Dial 连接数据库
func Dial(addr Address, opts ...bun.DBOption) (*bun.DB, error) {
	sqlDB, err := sql.Open("mysql", addr.DSN)
	if err != nil {
		return nil, sderr.Wrapf(err, "open mysql error")
	}
	db := bun.NewDB(sqlDB, mysqldialect.New(), opts...)
	addr.apply(db)
	return db, nil
}

func (addr *Address) apply(db *bun.DB) {
	if addr.ConnMaxLifeTimeMS > 0 {
		db.SetConnMaxLifetime(sdtime.Milliseconds(addr.ConnMaxLifeTimeMS))
	}
	if addr.ConnMaxIdleTimeMS > 0 {
		db.SetConnMaxIdleTime(sdtime.Milliseconds(addr.ConnMaxIdleTimeMS))
	}
	if addr.MaxIdleConns > 0 {
		db.SetMaxIdleConns(addr.MaxIdleConns)
	}
	if addr.MaxOpenConns > 0 {
		db.SetMaxOpenConns(addr.MaxOpenConns)
	}
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
