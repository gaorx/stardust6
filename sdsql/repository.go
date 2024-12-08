package sdsql

import (
	"context"
)

// EntityId 表示repository的实体ID类型
type EntityId interface {
	~string | ~int | ~int64 | ~uint | ~uint64
}

// Repository 定义了一个通用的仓库接口
type Repository[T any, ID EntityId, Q any] interface {
	// CountAll 计算所有记录数
	CountAll(ctx context.Context) (int64, error)
	// CountBy 计算符合条件的记录数
	CountBy(ctx context.Context, q Q) (int64, error)
	// GetById 根据ID获取记录
	GetById(ctx context.Context, id ID, postProcs ...RowsProc[T]) (T, error)
	// GetBy 根据条件获取记录
	GetBy(ctx context.Context, q Q, postProcs ...RowsProc[T]) (T, error)
	// ExistsById 判断是否存在指定ID的记录
	ExistsById(ctx context.Context, id ID) (bool, error)
	// ExistsBy 判断是否存在符合条件的记录
	ExistsBy(ctx context.Context, q Q) (bool, error)
	// FindAll 获取所有记录
	FindAll(ctx context.Context, postProcs ...RowsProc[T]) ([]T, error)
	// FindByIds 根据ID列表获取记录
	FindByIds(ctx context.Context, ids []ID, postProcs ...RowsProc[T]) ([]T, error)
	// FindBy 根据条件获取记录
	FindBy(ctx context.Context, q Q, postProcs ...RowsProc[T]) ([]T, error)
	// PaginateAll 以分页形式获取所有数据
	PaginateAll(ctx context.Context, page Page, postProcs ...RowsProc[T]) (*PageData[T], error)
	// PaginateBy 以分页形式获取符合条件的数据
	PaginateBy(ctx context.Context, q Q, page Page, postProcs ...RowsProc[T]) (*PageData[T], error)
	// Create 创建记录
	Create(ctx context.Context, entity T) (Result, error)
	// CreateMany 创建多条记录
	CreateMany(ctx context.Context, entities []T) (Result, error)
	// CreateAndGet 创建记录并返回创建后的记录
	CreateAndGet(ctx context.Context, entity T, postProcs ...RowsProc[T]) (T, error)
	// Update 更新一条记录
	Update(ctx context.Context, entity T, cols []string) (Result, error)
	// UpdateMany 更新多条记录
	UpdateMany(ctx context.Context, entities []T, cols []string) (Result, error)
	// UpdateBy 根据条件更新记录
	UpdateBy(ctx context.Context, entity T, q Q, cols []string) (Result, error)
	// UpdateAndGet 更新一条记录并返回更新后的记录
	UpdateAndGet(ctx context.Context, entity T, cols []string, postProcs ...RowsProc[T]) (T, error)
	// DeleteAll 删除所有记录
	DeleteAll(ctx context.Context) (Result, error)
	// Delete 删除一条记录
	Delete(ctx context.Context, entity T) (Result, error)
	// DeleteMany 删除多条记录
	DeleteMany(ctx context.Context, entities []T) (Result, error)
	// DeleteById 根据ID删除记录
	DeleteById(ctx context.Context, id ID) (Result, error)
	// DeleteByIds 根据ID列表删除记录
	DeleteByIds(ctx context.Context, ids []ID) (Result, error)
	// DeleteBy 根据条件删除记录
	DeleteBy(ctx context.Context, q Q) (Result, error)
}
