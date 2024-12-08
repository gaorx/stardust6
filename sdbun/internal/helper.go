package internal

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust6/sdbun"
	"github.com/gaorx/stardust6/sdrand"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"testing"
)

func useTodoEntities(t *testing.T, db bun.IDB, n int, f func([]*todoEntity1)) {
	clearTodoEntities(t, db)
	entities := createTodoEntities(t, db, n)
	defer clearTodoEntities(t, db)
	f(entities)
}

func createTodoEntities(t *testing.T, db bun.IDB, n int) []*todoEntity1 {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	var entities []*todoEntity1
	for i := 0; i < n; i++ {
		todoId := sdrand.String(10, sdrand.LowerCaseAlphanumericCharset)
		entity0 := todoEntity1{
			Id:          todoId,
			Title:       fmt.Sprintf("xyz-%d", i+1),
			Description: fmt.Sprintf("xyzxyz-%d", i+1),
			IsComplete:  true,
			Priority:    3,
		}
		entity1 := entity0
		entities = append(entities, &entity0)
		sr, err := repo.Create(context.Background(), &entity1)
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
	}
	return entities
}

func clearTodoEntities(t *testing.T, db bun.IDB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	_, err := repo.DeleteAll(context.Background())
	is.NoError(err)
}
