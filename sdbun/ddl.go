package sdbun

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"github.com/uptrace/bun/dialect"
)

type SqlCreateTableOptions struct {
	Version     string
	IfNotExists bool
}

func SqlCreateTable(d dialect.Name, model any, opts *SqlCreateTableOptions) (string, error) {
	opts1 := lo.FromPtr(opts)
	if opts1.Version == "" {
		if d == dialect.MySQL {
			opts1.Version = "8.4.0"
		}
	}
	db, err := mockDB(d, opts1.Version)
	if err != nil {
		return "", sderr.Wrap(err)
	}
	defer func() { _ = db.Close() }()
	q := db.NewCreateTable().Model(model)
	if opts1.IfNotExists {
		q.IfNotExists()
	}
	r, err := q.AppendQuery(db.Formatter(), nil)
	if err != nil {
		return "", sderr.Wrap(err)
	}
	return string(r), nil
}

func MustSqlCreateTable(typ dialect.Name, model any, opts *SqlCreateTableOptions) string {
	r, err := SqlCreateTable(typ, model, opts)
	if err != nil {
		panic(err)
	}
	return r
}
