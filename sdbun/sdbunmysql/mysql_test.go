package sdbunmysql

import (
	"github.com/gaorx/stardust6/sdbun/internal"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestRepository(t *testing.T) {
	dsn := os.Getenv("TEST_MYSQL")
	if dsn == "" {
		return
	}

	is := assert.New(t)
	db, err := Dial(Address{
		DSN:  dsn,
		Slog: slog.Default(),
	})
	is.NoError(err)

	internal.CreateTables(t, db)
	internal.TestRepositoryCount(t, db)
	internal.TestRepositoryGet(t, db)
	internal.TestRepositoryExists(t, db)
	internal.TestRepositoryFind(t, db)
	internal.TestRepositoryPaginate(t, db)
	internal.TestRepositoryCreate(t, db)
	internal.TestRepositoryUpdate(t, db)
	internal.TestRepositoryDelete(t, db)
}
