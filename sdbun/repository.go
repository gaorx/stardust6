package sdbun

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdsql"
	"github.com/uptrace/bun"
	"reflect"
)

// Repository 一个用表示bun实现的 sdsql.Repository 接口
type Repository[T any, ID sdsql.EntityId] interface {
	sdsql.Repository[T, ID, *RepositoryQuery]
}

type repository[T any, ID sdsql.EntityId] struct {
	db   bun.IDB
	info repoInfo
}

// RepoOf 通过db创建一个Repository
func RepoOf[T any, ID sdsql.EntityId](db bun.IDB) Repository[T, ID] {
	var model T
	info := getRepoInfo(db, reflect.TypeOf(model))
	if info == nil {
		panic(sderr.Newf("get model info failed"))
	}
	return repository[T, ID]{db: db, info: *info}
}

// CountAll 实现 Repository.CountAll 接口
func (repo repository[T, ID]) CountAll(ctx context.Context) (int64, error) {
	return repo.CountBy(ctx, nil)
}

// CountBy 实现 Repository.CountBy 接口
func (repo repository[T, ID]) CountBy(ctx context.Context, q *RepositoryQuery) (int64, error) {
	return Count[int64](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Apply(q.applySelect)
	})
}

// GetById 实现 Repository.GetById 接口
func (repo repository[T, ID]) GetById(ctx context.Context, id ID, postProcs ...sdsql.RowsProc[T]) (T, error) {
	idFieldSqlName := repo.idFieldSqlName()
	return SelectFirst[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Where("?=?", bun.Ident(idFieldSqlName), id)
	}, postProcs...)
}

// GetBy 实现 Repository.GetBy 接口
func (repo repository[T, ID]) GetBy(ctx context.Context, q *RepositoryQuery, postProcs ...sdsql.RowsProc[T]) (T, error) {
	return SelectFirst[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Apply(q.applySelect)
	}, postProcs...)
}

// ExistsById 实现 Repository.ExistsById 接口
func (repo repository[T, ID]) ExistsById(ctx context.Context, id ID) (bool, error) {
	idFieldSqlName := repo.idFieldSqlName()
	return Exists[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Where("?=?", bun.Ident(idFieldSqlName), id)
	})
}

// ExistsBy 实现 Repository.ExistsBy 接口
func (repo repository[T, ID]) ExistsBy(ctx context.Context, q *RepositoryQuery) (bool, error) {
	return Exists[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Apply(q.applySelect)
	})
}

// FindAll 实现 Repository.FindAll 接口
func (repo repository[T, ID]) FindAll(ctx context.Context, postProcs ...sdsql.RowsProc[T]) ([]T, error) {
	return repo.FindBy(ctx, nil, postProcs...)
}

// FindByIds 实现 Repository.FindByIds 接口
func (repo repository[T, ID]) FindByIds(ctx context.Context, ids []ID, postProcs ...sdsql.RowsProc[T]) ([]T, error) {
	idFieldSqlName := repo.idFieldSqlName()
	if len(ids) <= 0 {
		return []T{}, nil
	}
	rows, err := SelectMany[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Where("? IN (?)", bun.Ident(idFieldSqlName), bun.In(ids))
	}, postProcs...)
	return ensureRows(rows), err
}

// FindBy 实现 Repository.FindBy 接口
func (repo repository[T, ID]) FindBy(ctx context.Context, q *RepositoryQuery, postProcs ...sdsql.RowsProc[T]) ([]T, error) {
	rows, err := SelectMany[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		var model T
		return s.Model(model).Apply(q.applySelect)
	}, postProcs...)
	return ensureRows(rows), err
}

// PaginateAll 实现 Repository.PaginateAll 接口
func (repo repository[T, ID]) PaginateAll(ctx context.Context, page sdsql.Page, postProcs ...sdsql.RowsProc[T]) (*sdsql.PageData[T], error) {
	return repo.PaginateBy(ctx, nil, page, postProcs...)
}

// PaginateBy 实现 Repository.PaginateBy 接口
func (repo repository[T, ID]) PaginateBy(ctx context.Context, q *RepositoryQuery, page sdsql.Page, postProcs ...sdsql.RowsProc[T]) (*sdsql.PageData[T], error) {
	rows, numRowsTotal, err := SelectManyAndCount[T](ctx, repo.db, func(s *bun.SelectQuery) *bun.SelectQuery {
		limit, offset := page.LimitAndOffset()
		return s.Apply(q.applySelect).Limit(int(limit)).Offset(int(offset))
	}, postProcs...)
	if err != nil {
		return nil, err
	}
	return sdsql.NewPageData(ensureRows(rows), page, numRowsTotal), nil
}

// Create 实现 Repository.Create 接口
func (repo repository[T, ID]) Create(ctx context.Context, entity T) (sdsql.Result, error) {
	return Insert(ctx, repo.db, entity, nil)
}

// CreateMany 实现 Repository.CreateMany 接口
func (repo repository[T, ID]) CreateMany(ctx context.Context, entities []T) (sdsql.Result, error) {
	return InsertMany(ctx, repo.db, entities, nil)
}

// CreateAndGet 实现 Repository.CreateAndGet 接口
func (repo repository[T, ID]) CreateAndGet(ctx context.Context, entity T, postProcs ...sdsql.RowsProc[T]) (T, error) {
	_, err := repo.Create(ctx, entity)
	if err != nil {
		var zero T
		return zero, err
	}
	return repo.GetById(ctx, repo.getEntityId(entity), postProcs...)
}

// Update 实现 Repository.Update 接口
func (repo repository[T, ID]) Update(ctx context.Context, entity T, cols []string) (sdsql.Result, error) {
	return Update[T](ctx, repo.db, entity, func(u *bun.UpdateQuery) *bun.UpdateQuery {
		id := repo.getEntityId(entity)
		if len(cols) > 0 {
			u = u.Column(cols...)
		}
		return u.Where("?=?", bun.Ident(repo.idFieldSqlName()), id)
	})
}

// UpdateMany 实现 Repository.UpdateMany 接口
func (repo repository[T, ID]) UpdateMany(ctx context.Context, entities []T, cols []string) (sdsql.Result, error) {
	if len(entities) <= 0 {
		return sdsql.ResultNoRows, nil
	}
	rowsAffected := int64(0)
	for _, entity := range entities {
		sr0, err := repo.Update(ctx, entity, cols)
		if err != nil {
			return sdsql.Result{}, err
		}
		rowsAffected += sr0.RowsAffectedOr(0)
	}
	return sdsql.ResultOfRowsAffected(rowsAffected), nil
}

// UpdateAndGet 实现 Repository.UpdateAndGet 接口
func (repo repository[T, ID]) UpdateAndGet(ctx context.Context, entity T, cols []string, postProcs ...sdsql.RowsProc[T]) (T, error) {
	_, err := repo.Update(ctx, entity, cols)
	if err != nil {
		var zero T
		return zero, err
	}
	return repo.GetById(ctx, repo.getEntityId(entity), postProcs...)
}

// UpdateBy 实现 Repository.UpdateBy 接口
func (repo repository[T, ID]) UpdateBy(ctx context.Context, entity T, q *RepositoryQuery, cols []string) (sdsql.Result, error) {
	return Update[T](ctx, repo.db, entity, func(u *bun.UpdateQuery) *bun.UpdateQuery {
		if len(cols) > 0 {
			u = u.Column(cols...)
		}
		return u.Apply(q.applyUpdate)
	})
}

// DeleteAll 实现 Repository.DeleteAll 接口
func (repo repository[T, ID]) DeleteAll(ctx context.Context) (sdsql.Result, error) {
	return Delete[T](ctx, repo.db, func(d *bun.DeleteQuery) *bun.DeleteQuery {
		return d.Where("1=1")
	})
}

// Delete 实现 Repository.Delete 接口
func (repo repository[T, ID]) Delete(ctx context.Context, entity T) (sdsql.Result, error) {
	return Delete[T](ctx, repo.db, func(d *bun.DeleteQuery) *bun.DeleteQuery {
		id := repo.getEntityId(entity)
		return d.Where("?=?", bun.Ident(repo.idFieldSqlName()), id)
	})
}

// DeleteMany 实现 Repository.DeleteMany 接口
func (repo repository[T, ID]) DeleteMany(ctx context.Context, entities []T) (sdsql.Result, error) {
	if len(entities) <= 0 {
		return sdsql.ResultNoRows, nil
	}
	ids := repo.getEntityIds(entities)
	return Delete[T](ctx, repo.db, func(d *bun.DeleteQuery) *bun.DeleteQuery {
		return d.Where("? IN (?)", bun.Ident(repo.idFieldSqlName()), bun.In(ids))
	})
}

// DeleteById 实现 Repository.DeleteById 接口
func (repo repository[T, ID]) DeleteById(ctx context.Context, id ID) (sdsql.Result, error) {
	return Delete[T](ctx, repo.db, func(d *bun.DeleteQuery) *bun.DeleteQuery {
		return d.Where("?=?", bun.Ident(repo.idFieldSqlName()), id)
	})
}

// DeleteByIds 实现 Repository.DeleteByIds 接口
func (repo repository[T, ID]) DeleteByIds(ctx context.Context, ids []ID) (sdsql.Result, error) {
	return Delete[T](ctx, repo.db, func(d *bun.DeleteQuery) *bun.DeleteQuery {
		if len(ids) <= 0 {
			return d.Where("1!=1")
		}
		return d.Where("? IN (?)", bun.Ident(repo.idFieldSqlName()), bun.In(ids))
	})
}

// DeleteBy 实现 Repository.DeleteBy 接口
func (repo repository[T, ID]) DeleteBy(ctx context.Context, q *RepositoryQuery) (sdsql.Result, error) {
	return Delete[T](ctx, repo.db, func(d *bun.DeleteQuery) *bun.DeleteQuery {
		return d.Apply(q.applyDelete)
	})
}

func (repo repository[T, ID]) idFieldSqlName() string {
	return repo.info.idFieldSqlName
}

func (repo repository[T, ID]) getEntityId(entity T) ID {
	return getEntityId[ID](entity, &repo.info)
}

func (repo repository[T, ID]) getEntityIds(entities []T) []ID {
	ids := make([]ID, len(entities))
	for i, entity := range entities {
		ids[i] = repo.getEntityId(entity)
	}
	return ids
}
