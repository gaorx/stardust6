package internal

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust6/sdbun"
	"github.com/gaorx/stardust6/sdrand"
	"github.com/gaorx/stardust6/sdsql"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"testing"
)

func CreateTables(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	models := []any{
		&todoEntity1{},
	}
	for i := len(models) - 1; i >= 0; i-- {
		_, err := db.NewDropTable().Model(models[i]).IfExists().Exec(context.Background())
		is.NoError(err)
	}
	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).Exec(context.Background())
		is.NoError(err)
	}
}

func TestRepositoryCount(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(_ []*todoEntity1) {
		// CountAll
		n, err := repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(10), n)

		// CountBy
		n, err = repo.CountBy(context.Background(), nil)
		is.NoError(err)
		is.Equal(int64(10), n)
		n, err = repo.CountBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title IN (?)", bun.In([]string{"xyz-2", "xyz-4", "xyz-not-exists"}))
		}))
		is.NoError(err)
		is.Equal(int64(2), n)
	})
}

func TestRepositoryGet(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(entities []*todoEntity1) {
		// GetById
		expected := entities[2]
		id := expected.Id
		actual, err := repo.GetById(context.Background(), id)
		is.NoError(err)
		is.True(expected.contentEqual(actual))
		actual, err = repo.GetById(context.Background(), "not-exists")
		is.True(sdbun.IsNotFoundErr(err))
		is.Nil(actual)

		// GetBy
		actual, err = repo.GetBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title=?", "xyz-3")
		}))
		is.NoError(err)
		is.Equal("xyz-3", actual.Title)
	})
}

func TestRepositoryExists(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(entities []*todoEntity1) {
		// ExistsById
		id3, id6 := entities[2].Id, entities[5].Id
		exists, err := repo.ExistsById(context.Background(), id3)
		is.NoError(err)
		is.True(exists)
		exists, err = repo.ExistsById(context.Background(), id6)
		is.NoError(err)
		is.True(exists)
		exists, err = repo.ExistsById(context.Background(), "not-exists")
		is.NoError(err)
		is.False(exists)

		// ExistsBy
		exists, err = repo.ExistsBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title=?", "xyz-7")
		}))
		is.NoError(err)
		is.True(exists)
		exists, err = repo.ExistsBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title=?", "xyz-not-exists")
		}))
		is.NoError(err)
		is.False(exists)
	})
}

func TestRepositoryFind(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(entities []*todoEntity1) {
		// FindAll
		actuals, err := repo.FindAll(context.Background())
		is.NoError(err)
		is.Len(actuals, 10)

		// FindByIds
		id3, id6 := entities[2].Id, entities[5].Id
		actuals, err = repo.FindByIds(context.Background(), []string{id3, id6, "not-exists"})
		is.NoError(err)
		is.Len(actuals, 2)
		is.True(todoEntity1ContainsById(entities, id3))
		is.True(todoEntity1ContainsById(entities, id6))

		// FindBy
		actuals, err = repo.FindBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title IN (?)", bun.In([]string{"xyz-3", "xyz-6", "xyz-not-exists"}))
		}))
		is.NoError(err)
		is.Len(actuals, 2)
		is.True(todoEntity1ContainsById(entities, id3))
		is.True(todoEntity1ContainsById(entities, id6))
	})
}

func TestRepositoryPaginate(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(entities []*todoEntity1) {
		// PaginateAll
		pageData, err := repo.PaginateAll(context.Background(), sdsql.Page1(3, 2))
		is.NoError(err)
		is.Equal(int64(10), pageData.TotalRows)
		is.Equal(int64(5), pageData.TotalPages)
		is.Equal(int64(3), pageData.PageNum)
		is.Equal(int64(2), pageData.PageSize)
		is.Equal(2, len(pageData.Data))
		is.True(todoEntity1ContainsById(pageData.Data, entities[4].Id))
		is.True(todoEntity1ContainsById(pageData.Data, entities[5].Id))
		pageData, err = repo.PaginateAll(context.Background(), sdsql.Page1(4, 3))
		is.NoError(err)
		is.Equal(int64(10), pageData.TotalRows)
		is.Equal(int64(4), pageData.TotalPages)
		is.Equal(int64(4), pageData.PageNum)
		is.Equal(int64(3), pageData.PageSize)
		is.Equal(1, len(pageData.Data))
		is.True(todoEntity1ContainsById(pageData.Data, entities[9].Id))
		pageData, err = repo.PaginateAll(context.Background(), sdsql.Page1(5, 3))
		is.NoError(err)
		is.Equal(int64(10), pageData.TotalRows)
		is.Equal(int64(4), pageData.TotalPages)
		is.Equal(int64(5), pageData.PageNum)
		is.Equal(int64(3), pageData.PageSize)
		is.Equal(0, len(pageData.Data))

		// PaginateBy
		pageData, err = repo.PaginateBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title IN (?)", bun.In([]string{"xyz-1", "xyz-3", "xyz-5", "xyz-7", "xyz-9"}))
		}), sdsql.Page1(2, 2))
		is.NoError(err)
		is.Equal(int64(5), pageData.TotalRows)
		is.Equal(int64(3), pageData.TotalPages)
		is.Equal(int64(2), pageData.PageNum)
		is.Equal(int64(2), pageData.PageSize)
		is.Equal(2, len(pageData.Data))
		is.True(todoEntity1ContainsByTitle(pageData.Data, "xyz-5"))
		is.True(todoEntity1ContainsByTitle(pageData.Data, "xyz-7"))
		pageData, err = repo.PaginateBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title IN (?)", bun.In([]string{"xyz-1", "xyz-3", "xyz-5", "xyz-7", "xyz-9"}))
		}), sdsql.Page1(3, 2))
		is.NoError(err)
		is.Equal(int64(5), pageData.TotalRows)
		is.Equal(int64(3), pageData.TotalPages)
		is.Equal(int64(3), pageData.PageNum)
		is.Equal(int64(2), pageData.PageSize)
		is.Equal(1, len(pageData.Data))
		is.True(todoEntity1ContainsByTitle(pageData.Data, "xyz-9"))
		pageData, err = repo.PaginateBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title IN (?)", bun.In([]string{"xyz-1", "xyz-3", "xyz-5", "xyz-7", "xyz-9"}))
		}), sdsql.Page1(4, 2))
		is.NoError(err)
		is.Equal(int64(5), pageData.TotalRows)
		is.Equal(int64(3), pageData.TotalPages)
		is.Equal(int64(4), pageData.PageNum)
		is.Equal(int64(2), pageData.PageSize)
		is.Equal(0, len(pageData.Data))
	})
}

func TestRepositoryCreate(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(_ []*todoEntity1) {
		// Create
		entity00 := todoEntity1{
			Id:          "new-id",
			Title:       "new-title",
			Description: "new-description",
			IsComplete:  true,
			Priority:    3,
		}
		entity01 := entity00
		sr, err := repo.Create(context.Background(), &entity00)
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		// 注释掉下面这行是因为 LastInsertIdOr 方法在不同数据库中的行为不一致
		// is.True(sr.LastInsertIdOr(0) > 0)
		actual, err := repo.GetById(context.Background(), "new-id")
		is.NoError(err)
		is.True(entity00.contentEqual(actual))
		sr, err = repo.Create(context.Background(), &entity01)
		is.Error(err)
		is.Equal(int64(0), sr.RowsAffectedOr(0))

		// CreateMany
		var newEntities00, newEntities01 []*todoEntity1
		for i := 0; i < 3; i++ {
			entity00 := todoEntity1{
				Id:          sdrand.String(10, sdrand.LowerCaseAlphanumericCharset),
				Title:       fmt.Sprintf("abc-%d", i+1),
				Description: fmt.Sprintf("abcabc-%d", i+1),
				IsComplete:  true,
				Priority:    2,
			}
			entity01 = entity00
			newEntities00 = append(newEntities00, &entity00)
			newEntities01 = append(newEntities01, &entity01)
		}
		sr, err = repo.CreateMany(context.Background(), newEntities00)
		is.NoError(err)
		is.Equal(int64(3), sr.RowsAffectedOr(0))
		// 注释掉下面这行是因为 LastInsertIdOr 方法在不同数据库中的行为不一致
		// is.True(sr.LastInsertIdOr(0) > 0)
		sr, err = repo.CreateMany(context.Background(), newEntities01)
		is.Error(err)
		is.Equal(int64(0), sr.RowsAffectedOr(0))

		// CreateAndGet
		creating := &todoEntity1{
			Id:          "new-id-2",
			Title:       "new-title-2",
			Description: "new-description-2",
			IsComplete:  true,
			Priority:    34,
		}
		created, err := repo.CreateAndGet(context.Background(), creating)
		is.NoError(err)
		is.NotNil(created)
		is.True(creating.contentEqual(created))
	})
}

func TestRepositoryUpdate(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(entities []*todoEntity1) {
		// Update
		id3 := entities[2].Id
		sr, err := repo.Update(context.Background(), &todoEntity1{
			Id:          id3,
			Title:       "new-title-3",
			Description: "new-desc-3",
		}, []string{"title"}) // 只更新 title
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		actual, err := repo.GetById(context.Background(), id3)
		is.NoError(err)
		is.Equal("new-title-3", actual.Title)
		is.Equal("xyzxyz-3", actual.Description)
		sr, err = repo.Update(context.Background(), &todoEntity1{
			Id:    id3,
			Title: "new-title-3",
		}, nil) // 不指定列，导致除了title外其他字段都被清零了
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		actual, err = repo.GetById(context.Background(), id3)
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		is.Equal("new-title-3", actual.Title)
		is.Equal("", actual.Description)
		is.Equal(false, actual.IsComplete)
		is.Equal(0, actual.Priority)
		sr, err = repo.Update(context.Background(), &todoEntity1{
			Id:    "not-exists",
			Title: "new-title-3",
		}, []string{"title"})
		is.NoError(err)
		is.Equal(int64(0), sr.RowsAffectedOr(0))

		// UpdateMany
		var updatingEntities []*todoEntity1
		for i := 3; i <= 5; i++ {
			updatingEntities = append(updatingEntities, &todoEntity1{
				Id:    entities[i].Id,
				Title: fmt.Sprintf("xxx-%d", i),
			})
		}
		sr, err = repo.UpdateMany(context.Background(), updatingEntities, []string{"title"})
		is.NoError(err)
		is.Equal(int64(3), sr.RowsAffectedOr(0))
		actuals, err := repo.FindByIds(context.Background(), []string{
			entities[3].Id,
			entities[4].Id,
			entities[5].Id},
		)
		is.NoError(err)
		is.True(todoEntity1ContainsByTitle(actuals, "xxx-3"))
		is.True(todoEntity1ContainsByTitle(actuals, "xxx-4"))
		is.True(todoEntity1ContainsByTitle(actuals, "xxx-5"))
		updatingEntities[0].Id = "not-exists"
		updatingEntities[1].Title = "yyy-4"
		updatingEntities[2].Title = "yyy-5"
		sr, err = repo.UpdateMany(context.Background(), updatingEntities, []string{"title"})
		is.NoError(err)
		is.Equal(int64(2), sr.RowsAffectedOr(0))
		actuals, err = repo.FindByIds(context.Background(), []string{
			updatingEntities[0].Id,
			updatingEntities[1].Id,
			updatingEntities[2].Id,
		})
		is.NoError(err)
		is.Equal(2, len(actuals))
		is.True(todoEntity1ContainsByTitle(actuals, "yyy-4"))
		is.True(todoEntity1ContainsByTitle(actuals, "yyy-5"))

		// UpdateBy
		sr, err = repo.UpdateBy(context.Background(), &todoEntity1{
			Title: "zzz",
		}, sdbun.Q(func(u *bun.UpdateQuery) *bun.UpdateQuery {
			return u.Where("title IN (?)", bun.In([]string{"xyz-7", "xyz-8", "not-exists"}))
		}), []string{"title"})
		is.NoError(err)
		is.Equal(int64(2), sr.RowsAffectedOr(0))
		actuals, err = repo.FindBy(context.Background(), sdbun.Q(func(s *bun.SelectQuery) *bun.SelectQuery {
			return s.Where("title=?", "zzz")
		}))
		is.NoError(err)
		is.Equal(2, len(actuals))

		// UpdateAndGet
		id10 := entities[9].Id
		updating := &todoEntity1{
			Id:          id10,
			Title:       "new-title-10",
			Description: "new-desc-10",
		}
		updated, err := repo.UpdateAndGet(context.Background(), updating, []string{"title"})
		is.NoError(err)
		is.NotNil(updated)
		is.Equal("new-title-10", updated.Title)
		is.Equal("xyzxyz-10", updated.Description)
		updated, err = repo.UpdateAndGet(context.Background(), &todoEntity1{
			Id:    "not-exists",
			Title: "nnn",
		}, []string{"title"})
		is.Error(err)
		is.True(sdbun.IsNotFoundErr(err))
	})
}

func TestRepositoryDelete(t *testing.T, db *bun.DB) {
	is := assert.New(t)
	repo := sdbun.RepoOf[*todoEntity1, string](db)
	useTodoEntities(t, db, 10, func(entities []*todoEntity1) {
		n, err := repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(10), n)

		// DeleteBy
		sr, err := repo.DeleteBy(context.Background(), sdbun.Q(func(d *bun.DeleteQuery) *bun.DeleteQuery {
			return d.Where("title IN (?)", bun.In([]string{"xyz-9"}))
		}))
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		n, err = repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(9), n)
		sr, err = repo.DeleteBy(context.Background(), sdbun.Q(func(d *bun.DeleteQuery) *bun.DeleteQuery {
			return d.Where("title IN (?)", bun.In([]string{"not-exits"}))
		}))
		is.NoError(err)
		is.Equal(int64(0), sr.RowsAffectedOr(0))

		// DeleteById
		id8 := entities[7].Id
		sr, err = repo.DeleteById(context.Background(), id8)
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		n, err = repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(8), n)
		sr, err = repo.DeleteById(context.Background(), "not-exists")
		is.NoError(err)
		is.Equal(int64(0), sr.RowsAffectedOr(0))

		// DeleteByIds
		id7, id6 := entities[6].Id, entities[5].Id
		sr, err = repo.DeleteByIds(context.Background(), []string{id7, id6, "not-exists"})
		is.NoError(err)
		is.Equal(int64(2), sr.RowsAffectedOr(0))
		n, err = repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(6), n)

		// Delete
		entity5 := entities[4]
		sr, err = repo.Delete(context.Background(), entity5)
		is.NoError(err)
		is.Equal(int64(1), sr.RowsAffectedOr(0))
		n, err = repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(5), n)
		sr, err = repo.Delete(context.Background(), &todoEntity1{Id: "not-exists"})
		is.NoError(err)
		is.Equal(int64(0), sr.RowsAffectedOr(0))

		// DeleteMany
		entity4, entity3 := entities[3], entities[2]
		sr, err = repo.DeleteMany(context.Background(), []*todoEntity1{entity4, entity3})
		is.NoError(err)
		is.Equal(int64(2), sr.RowsAffectedOr(0))
		n, err = repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(3), n)

		// DeleteAll
		sr, err = repo.DeleteAll(context.Background())
		is.NoError(err)
		is.Equal(int64(3), sr.RowsAffectedOr(0))
		n, err = repo.CountAll(context.Background())
		is.NoError(err)
		is.Equal(int64(0), n)
	})
}
