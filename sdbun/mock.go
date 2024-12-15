package sdbun

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaorx/stardust6/sderr"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func mockDB(d dialect.Name, version string) (*bun.DB, error) {
	switch d {
	case dialect.MySQL:
		return mockMysql(version), nil
	default:
		return nil, sderr.Newf("unsupported dialect: %s", d)
	}
}

func mockMysql(version string) *bun.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(sderr.Newf("mock mysql error"))
	}
	mock.ExpectQuery("SELECT version()").WillReturnRows(
		sqlmock.NewRows([]string{"version()"}).AddRow(version),
	)
	return bun.NewDB(db, mysqldialect.New())
}
