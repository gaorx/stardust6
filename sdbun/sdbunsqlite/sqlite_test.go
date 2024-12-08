package sdbunsqlite

import (
	"github.com/gaorx/stardust6/sdbun/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository(t *testing.T) {
	is := assert.New(t)
	db, err := Open(Address{
		DSN: "file::memory:?cache=shared",
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
