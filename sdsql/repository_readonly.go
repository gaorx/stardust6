package sdsql

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
)

type readonlyRepository[T any, ID EntityId, Q any] struct {
	repo Repository[T, ID, Q]
}

// ReadonlyRepoOf 返回一个repository只读的代理，这个代理只能执行查询操作
func ReadonlyRepoOf[T any, ID EntityId, Q any](repo Repository[T, ID, Q]) Repository[T, ID, Q] {
	if repo == nil {
		panic(sderr.Newf("repo is nil"))
	}
	if rr, ok := repo.(readonlyRepository[T, ID, Q]); ok {
		return rr
	}
	return &readonlyRepository[T, ID, Q]{repo: repo}
}

// CountAll 实现 Repository.CountAll 接口
func (r readonlyRepository[T, ID, Q]) CountAll(ctx context.Context) (int64, error) {
	return r.repo.CountAll(ctx)
}

// CountBy 实现 Repository.CountBy 接口
func (r readonlyRepository[T, ID, Q]) CountBy(ctx context.Context, q Q) (int64, error) {
	return r.repo.CountBy(ctx, q)
}

// GetById 实现 Repository.GetById 接口
func (r readonlyRepository[T, ID, Q]) GetById(ctx context.Context, id ID, postProcs ...RowsProc[T]) (T, error) {
	return r.repo.GetById(ctx, id, postProcs...)
}

// GetBy 实现 Repository.GetBy 接口
func (r readonlyRepository[T, ID, Q]) GetBy(ctx context.Context, q Q, postProcs ...RowsProc[T]) (T, error) {
	return r.repo.GetBy(ctx, q, postProcs...)
}

// ExistsById 实现 Repository.ExistsById 接口
func (r readonlyRepository[T, ID, Q]) ExistsById(ctx context.Context, id ID) (bool, error) {
	return r.repo.ExistsById(ctx, id)
}

// ExistsBy 实现 Repository.ExistsBy 接口
func (r readonlyRepository[T, ID, Q]) ExistsBy(ctx context.Context, q Q) (bool, error) {
	return r.repo.ExistsBy(ctx, q)
}

// FindAll 实现 Repository.FindAll 接口
func (r readonlyRepository[T, ID, Q]) FindAll(ctx context.Context, postProcs ...RowsProc[T]) ([]T, error) {
	return r.repo.FindAll(ctx, postProcs...)
}

// FindByIds 实现 Repository.FindByIds 接口
func (r readonlyRepository[T, ID, Q]) FindByIds(ctx context.Context, ids []ID, postProcs ...RowsProc[T]) ([]T, error) {
	return r.repo.FindByIds(ctx, ids, postProcs...)
}

// FindBy 实现 Repository.FindBy 接口
func (r readonlyRepository[T, ID, Q]) FindBy(ctx context.Context, q Q, postProcs ...RowsProc[T]) ([]T, error) {
	return r.repo.FindBy(ctx, q, postProcs...)
}

// PaginateAll 实现 Repository.PaginateAll 接口
func (r readonlyRepository[T, ID, Q]) PaginateAll(ctx context.Context, page Page, postProcs ...RowsProc[T]) (*PageData[T], error) {
	return r.repo.PaginateAll(ctx, page, postProcs...)
}

// PaginateBy 实现 Repository.PaginateBy 接口
func (r readonlyRepository[T, ID, Q]) PaginateBy(ctx context.Context, q Q, page Page, postProcs ...RowsProc[T]) (*PageData[T], error) {
	return r.repo.PaginateBy(ctx, q, page, postProcs...)
}

// Create 实现 Repository.Create 接口
func (r readonlyRepository[T, ID, Q]) Create(_ context.Context, _ T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// CreateMany 实现 Repository.CreateMany 接口
func (r readonlyRepository[T, ID, Q]) CreateMany(_ context.Context, _ []T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// CreateAndGet 实现 Repository.CreateAndGet 接口
func (r readonlyRepository[T, ID, Q]) CreateAndGet(_ context.Context, _ T, _ ...RowsProc[T]) (T, error) {
	var zero T
	return zero, ErrReadonly
}

// Update 实现 Repository.Update 接口
func (r readonlyRepository[T, ID, Q]) Update(_ context.Context, _ T, _ []string) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// UpdateMany 实现 Repository.UpdateMany 接口
func (r readonlyRepository[T, ID, Q]) UpdateMany(_ context.Context, _ []T, _ []string) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// UpdateBy 实现 Repository.UpdateBy 接口
func (r readonlyRepository[T, ID, Q]) UpdateBy(_ context.Context, _ T, _ Q, _ []string) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// UpdateAndGet 实现 Repository.UpdateAndGet 接口
func (r readonlyRepository[T, ID, Q]) UpdateAndGet(_ context.Context, _ T, _ []string, _ ...RowsProc[T]) (T, error) {
	var zero T
	return zero, ErrReadonly
}

// DeleteAll 实现 Repository.DeleteAll 接口
func (r readonlyRepository[T, ID, Q]) DeleteAll(_ context.Context) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// Delete 实现 Repository.Delete 接口
func (r readonlyRepository[T, ID, Q]) Delete(_ context.Context, _ T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteMany 实现 Repository.DeleteMany 接口
func (r readonlyRepository[T, ID, Q]) DeleteMany(_ context.Context, _ []T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteById 实现 Repository.DeleteById 接口
func (r readonlyRepository[T, ID, Q]) DeleteById(_ context.Context, _ ID) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteByIds 实现 Repository.DeleteByIds 接口
func (r readonlyRepository[T, ID, Q]) DeleteByIds(_ context.Context, _ []ID) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteBy 实现 Repository.DeleteBy 接口
func (r readonlyRepository[T, ID, Q]) DeleteBy(_ context.Context, _ Q) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}
