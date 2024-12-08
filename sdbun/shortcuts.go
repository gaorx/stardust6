package sdbun

import (
	"context"
	"database/sql"
	"github.com/gaorx/stardust6/sdreflect"
	"github.com/gaorx/stardust6/sdsql"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"reflect"
)

// Tx 在事务中执行操作
func Tx(ctx context.Context, db bun.IDB, action func(context.Context, bun.Tx) error, opts *sql.TxOptions) error {
	return db.RunInTx(ctx, opts, action)
}

// TxFor 在事务中执行操作，并返回操作结果
func TxFor[R any](ctx context.Context, db bun.IDB, action func(context.Context, bun.Tx) (R, error), opts *sql.TxOptions) (R, error) {
	var r R
	err := db.RunInTx(ctx, opts, func(ctx context.Context, tx bun.Tx) error {
		r0, err := action(ctx, tx)
		if err != nil {
			return err
		}
		r = r0
		return nil
	})
	return r, err
}

// Insert 插入一条记录
func Insert[ROW any](ctx context.Context, db bun.IDB, row ROW, qfn func(query *bun.InsertQuery) *bun.InsertQuery) (sdsql.Result, error) {
	sr, err := db.NewInsert().Model(row).Apply(qfn).Exec(ctx)
	if err != nil {
		return sdsql.ResultOfRowsAffected(0), err
	}
	return sdsql.ResultOf(sr), nil
}

// InsertMany 插入多条记录
func InsertMany[ROW any](ctx context.Context, db bun.IDB, rows []ROW, qfn func(query *bun.InsertQuery) *bun.InsertQuery) (sdsql.Result, error) {
	sr, err := db.NewInsert().Model(&rows).Apply(qfn).Exec(ctx)
	if err != nil {
		return sdsql.ResultOfRowsAffected(0), err
	}
	return sdsql.ResultOf(sr), nil
}

// Update 更新一条记录
func Update[ROW any](ctx context.Context, db bun.IDB, row ROW, qfn func(query *bun.UpdateQuery) *bun.UpdateQuery) (sdsql.Result, error) {
	sr, err := db.NewUpdate().Model(row).Apply(qfn).Exec(ctx)
	if err != nil {
		return sdsql.ResultOfRowsAffected(0), err
	}
	return sdsql.ResultOf(sr), nil
}

// UpdateMany 更新多条记录
func UpdateMany[ROW any](ctx context.Context, db bun.IDB, rows []ROW, qfn func(query *bun.UpdateQuery) *bun.UpdateQuery) (sdsql.Result, error) {
	sr, err := db.NewUpdate().Model(rows).Apply(qfn).Exec(ctx)
	if err != nil {
		return sdsql.ResultOfRowsAffected(0), err
	}
	return sdsql.ResultOf(sr), nil
}

// Delete 删除记录
func Delete[ROW any](ctx context.Context, db bun.IDB, qfn func(query *bun.DeleteQuery) *bun.DeleteQuery) (sdsql.Result, error) {
	sr, err := db.NewDelete().Apply(qfn).Apply(modelApplier[*bun.DeleteQuery, ROW]()).Exec(ctx)
	if err != nil {
		return sdsql.ResultOfRowsAffected(0), err
	}
	return sdsql.ResultOf(sr), nil
}

// SelectMany 查询多条记录
func SelectMany[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery, postProcs ...sdsql.RowsProc[ROW]) ([]ROW, error) {
	var r []ROW
	err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Scan(ctx, &r)
	if err != nil {
		return nil, err
	}
	return sdsql.ProcRows(ctx, r, postProcs...)
}

// SelectFirst 查询第一条记录
func SelectFirst[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery, postProcs ...sdsql.RowsProc[ROW]) (ROW, error) {
	t := sdreflect.T[ROW]()
	if isPtrToStruct(t) {
		dest := reflect.New(t.Elem()).Interface()
		err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Scan(ctx, dest)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return sdsql.ProcRow(ctx, dest.(ROW), postProcs...)
	} else {
		var r ROW
		err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Scan(ctx, &r)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return sdsql.ProcRow(ctx, r, postProcs...)
	}
}

// SelectManyRaw 直接通过SQL语句查询多条记录
func SelectManyRaw[ROW any](ctx context.Context, db bun.IDB, q string, args []any, postProcs ...sdsql.RowsProc[ROW]) ([]ROW, error) {
	var r []ROW
	err := db.NewRaw(q, args...).Scan(ctx, &r)
	if err != nil {
		return nil, err
	}
	return sdsql.ProcRows(ctx, r, postProcs...)
}

// SelectFirstRaw 直接通过SQL语句查询第一条记录
func SelectFirstRaw[ROW any](ctx context.Context, db bun.IDB, q string, args []any, postProcs ...sdsql.RowsProc[ROW]) (ROW, error) {
	t := sdreflect.T[ROW]()
	if isPtrToStruct(t) {
		dest := reflect.New(t.Elem()).Interface()
		err := db.NewRaw(q, args...).Scan(ctx, dest)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return sdsql.ProcRow(ctx, dest.(ROW), postProcs...)
	} else {
		var r ROW
		err := db.NewRaw(q, args...).Scan(ctx, &r)
		if err != nil {
			var zero ROW
			return zero, err
		}
		return sdsql.ProcRow(ctx, r, postProcs...)
	}
}

// SelectOne 通过SQL语句查询一个值
func SelectOne[T any](ctx context.Context, db bun.IDB, q string, args []any) (T, error) {
	return SelectFirstRaw[T](ctx, db, q, args)
}

// Count 查询记录数
func Count[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) (int64, error) {
	n, err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Count(ctx)
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

// Exists 是否存在记录
func Exists[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery) (bool, error) {
	exists, err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// SelectManyAndCount 查询多条记录并返回符合条件的总行数
func SelectManyAndCount[ROW any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery, postProcs ...sdsql.RowsProc[ROW]) ([]ROW, int64, error) {
	var r []ROW
	n, err := db.NewSelect().Apply(qfn).Apply(modelApplier[*bun.SelectQuery, ROW]()).ScanAndCount(ctx, &r)
	if err != nil {
		return nil, 0, err
	}
	r, err = sdsql.ProcRows(ctx, r, postProcs...)
	if err != nil {
		return nil, 0, err
	}
	return r, int64(n), nil
}

// SelectMap 查询多条记录，以map形式返回
func SelectMap[ROW any, K comparable, V any](ctx context.Context, db bun.IDB, qfn func(*bun.SelectQuery) *bun.SelectQuery, transform func(ROW) (K, V), postProcs ...sdsql.RowsProc[ROW]) (map[K]V, error) {
	rows, err := SelectMany[ROW](ctx, db, qfn, postProcs...)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(rows, transform), nil
}

// SelectMapRaw 直接通过SQL语句查询多条记录，以map形式返回
func SelectMapRaw[ROW any, K comparable, V any](ctx context.Context, db bun.IDB, q string, args []any, transform func(ROW) (K, V), postProcs ...sdsql.RowsProc[ROW]) (map[K]V, error) {
	rows, err := SelectManyRaw[ROW](ctx, db, q, args, postProcs...)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(rows, transform), nil
}

type modelQuery[Q any] interface {
	GetTableName() string
	Model(model any) Q
}

func modelApplier[Q modelQuery[Q], ROW any]() func(Q) Q {
	return func(q Q) Q {
		if q.GetTableName() == "" {
			if model := modelOfTyped[ROW](); model != nil {
				q = q.Model(model)
			}
		}
		return q
	}
}

func modelOf(model any) any {
	if model == nil {
		return nil
	}
	t := reflect.TypeOf(model)
	k := t.Kind()
	if k == reflect.Struct {
		return reflect.Zero(reflect.PointerTo(t)).Interface()
	} else if k == reflect.Pointer {
		if t.Elem().Kind() == reflect.Struct {
			return reflect.Zero(t).Interface()
		} else {
			var getElem func(reflect.Type) reflect.Type
			getElem = func(t1 reflect.Type) reflect.Type {
				if t1.Kind() == reflect.Pointer {
					return getElem(t1.Elem())
				} else {
					return t1
				}
			}
			base := getElem(t)
			if base.Kind() == reflect.Struct {
				return reflect.Zero(reflect.PtrTo(base)).Interface()
			} else {
				return nil
			}
		}
	} else {
		return nil
	}
}

func modelOfTyped[T any]() any {
	var model T
	return modelOf(model)
}

func isPtrToStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct
}
