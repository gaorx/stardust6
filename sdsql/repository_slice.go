package sdsql

import (
	"context"
	"database/sql"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
)

// SlicedEntity SliceRepository使用的实体接口
type SlicedEntity[T any, ID EntityId] interface {
	EntityId() ID
	Clone() T
}

// SliceRepository 使用切片实现的只读repository
type SliceRepository[T SlicedEntity[T, ID], ID EntityId] interface {
	Repository[T, ID, func(T) bool]
}

type sliceRepository[T SlicedEntity[T, ID], ID EntityId] struct {
	s []T
}

// SliceRepoOf 从一个切片创建一个SliceRepository
func SliceRepoOf[T SlicedEntity[T, ID], ID EntityId](s []T) SliceRepository[T, ID] {
	return sliceRepository[T, ID]{s: s}
}

// Slice 返回切片, 注意不要修改其中的内容
func (repo sliceRepository[T, ID]) Slice() []T {
	return repo.s
}

// CloneSlice 和 Slice 返回相同的切片，但是可以修改其中的内容
func (repo sliceRepository[T, ID]) CloneSlice() []T {
	return repo.cloneSlice(repo.s)
}

// Len 返回切片长度
func (repo sliceRepository[T, ID]) Len() int {
	return len(repo.s)
}

// IsEmpty 判断切片是否为空
func (repo sliceRepository[T, ID]) IsEmpty() bool {
	return len(repo.s) <= 0
}

// At 返回指定索引的元素
func (repo sliceRepository[T, ID]) At(index int) T {
	return repo.s[index]
}

// CountAll 实现 Repository.CountAll 接口
func (repo sliceRepository[T, ID]) CountAll(_ context.Context) (int64, error) {
	return int64(len(repo.s)), nil
}

// CountBy 实现 Repository.CountBy 接口
func (repo sliceRepository[T, ID]) CountBy(_ context.Context, q func(T) bool) (int64, error) {
	n := lo.CountBy(repo.s, q)
	return int64(n), nil
}

// GetById 实现 Repository.GetById 接口
func (repo sliceRepository[T, ID]) GetById(ctx context.Context, id ID, postProcs ...RowsProc[T]) (T, error) {
	entity, ok := lo.Find(repo.s, func(e T) bool {
		return repo.getId(e) == id
	})
	if !ok {
		var zero T
		return zero, sderr.Wrap(sql.ErrNoRows)
	}
	return ProcRow(ctx, repo.clone(entity), postProcs...)
}

// GetBy 实现 Repository.GetBy 接口
func (repo sliceRepository[T, ID]) GetBy(ctx context.Context, q func(T) bool, postProcs ...RowsProc[T]) (T, error) {
	if q == nil {
		q = func(T) bool { return true }
	}
	entity, ok := lo.Find(repo.s, q)
	if !ok {
		var zero T
		return zero, sderr.Wrap(sql.ErrNoRows)
	}
	return ProcRow(ctx, repo.clone(entity), postProcs...)
}

// ExistsById 实现 Repository.ExistsById 接口
func (repo sliceRepository[T, ID]) ExistsById(_ context.Context, id ID) (bool, error) {
	found := lo.ContainsBy(repo.s, func(e T) bool {
		return repo.getId(e) == id
	})
	return found, nil
}

// ExistsBy 实现 Repository.ExistsBy 接口
func (repo sliceRepository[T, ID]) ExistsBy(_ context.Context, q func(T) bool) (bool, error) {
	if q == nil {
		q = func(T) bool { return true }
	}
	found := lo.ContainsBy(repo.s, q)
	return found, nil
}

// FindAll 实现 Repository.FindAll 接口
func (repo sliceRepository[T, ID]) FindAll(ctx context.Context, postProcs ...RowsProc[T]) ([]T, error) {
	return ProcRows(ctx, repo.cloneSlice(repo.s), postProcs...)
}

// FindByIds 实现 Repository.FindByIds 接口
func (repo sliceRepository[T, ID]) FindByIds(ctx context.Context, ids []ID, postProcs ...RowsProc[T]) ([]T, error) {
	if len(ids) <= 0 {
		return []T{}, nil
	}
	filtered := lo.Filter(repo.s, func(e T, _ int) bool {
		id := repo.getId(e)
		return lo.Contains(ids, id)
	})
	return ProcRows(ctx, repo.cloneSlice(filtered), postProcs...)
}

// FindBy 实现 Repository.FindBy 接口
func (repo sliceRepository[T, ID]) FindBy(ctx context.Context, q func(T) bool, postProcs ...RowsProc[T]) ([]T, error) {
	if q == nil {
		q = func(T) bool { return true }
	}
	filtered := lo.Filter(repo.s, func(e T, _ int) bool {
		return q(e)
	})
	return ProcRows(ctx, repo.cloneSlice(filtered), postProcs...)
}

// PaginateAll 实现 Repository.PaginateAll 接口
func (repo sliceRepository[T, ID]) PaginateAll(ctx context.Context, page Page, postProcs ...RowsProc[T]) (*PageData[T], error) {
	return repo.PaginateBy(ctx, nil, page, postProcs...)
}

// PaginateBy 实现 Repository.PaginateBy 接口
func (repo sliceRepository[T, ID]) PaginateBy(ctx context.Context, q func(T) bool, page Page, postProcs ...RowsProc[T]) (*PageData[T], error) {
	if q == nil {
		q = func(T) bool { return true }
	}
	filtered := lo.Filter(repo.s, func(e T, _ int) bool {
		return q(e)
	})
	processed, err := ProcRows(ctx, repo.cloneSlice(filtered), postProcs...)
	if err != nil {
		return nil, err
	}
	return PaginateSlice(processed, page), nil
}

// Create 实现 Repository.Create 接口
func (repo sliceRepository[T, ID]) Create(_ context.Context, _ T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// CreateMany 实现 Repository.CreateMany 接口
func (repo sliceRepository[T, ID]) CreateMany(_ context.Context, _ []T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// CreateAndGet 实现 Repository.CreateAndGet 接口
func (repo sliceRepository[T, ID]) CreateAndGet(_ context.Context, _ T, _ ...RowsProc[T]) (T, error) {
	var zero T
	return zero, ErrReadonly
}

// Update 实现 Repository.Update 接口
func (repo sliceRepository[T, ID]) Update(_ context.Context, _ T, _ []string) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// UpdateMany 实现 Repository.UpdateMany 接口
func (repo sliceRepository[T, ID]) UpdateMany(_ context.Context, _ []T, _ []string) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// UpdateBy 实现 Repository.UpdateBy 接口
func (repo sliceRepository[T, ID]) UpdateBy(_ context.Context, _ T, _ func(T) bool, _ []string) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// UpdateAndGet 实现 Repository.UpdateAndGet 接口
func (repo sliceRepository[T, ID]) UpdateAndGet(_ context.Context, _ T, _ []string, _ ...RowsProc[T]) (T, error) {
	var zero T
	return zero, ErrReadonly
}

// DeleteAll 实现 Repository.DeleteAll 接口
func (repo sliceRepository[T, ID]) DeleteAll(_ context.Context) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// Delete 实现 Repository.Delete 接口
func (repo sliceRepository[T, ID]) Delete(_ context.Context, _ T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteMany 实现 Repository.DeleteMany 接口
func (repo sliceRepository[T, ID]) DeleteMany(_ context.Context, _ []T) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteById 实现 Repository.DeleteById 接口
func (repo sliceRepository[T, ID]) DeleteById(_ context.Context, _ ID) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteByIds 实现 Repository.DeleteByIds 接口
func (repo sliceRepository[T, ID]) DeleteByIds(_ context.Context, _ []ID) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

// DeleteBy 实现 Repository.DeleteBy 接口
func (repo sliceRepository[T, ID]) DeleteBy(_ context.Context, _ func(T) bool) (Result, error) {
	return ResultOfRowsAffected(0), ErrReadonly
}

func (repo sliceRepository[T, ID]) getId(e T) ID {
	return e.EntityId()
}

func (repo sliceRepository[T, ID]) clone(e T) T {
	return e.Clone()
}

func (repo sliceRepository[T, ID]) cloneSlice(slice []T) []T {
	return lo.Map(slice, func(e T, _ int) T {
		return e.Clone()
	})
}
