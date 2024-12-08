package sdsql

import (
	"context"
)

// RowsProc 处理记录行数据的处理器，用于补充记录行中的内容
type RowsProc[ROW any] interface {
	// ProcRows 处理记录行数据，返回处理过的记录行数据，如出错返回error
	ProcRows(ctx context.Context, rows []ROW) ([]ROW, error)
}

// ProcRows 使用一组处理器处理记录行数据
func ProcRows[ROW any](ctx context.Context, rows []ROW, procs ...RowsProc[ROW]) ([]ROW, error) {
	if len(procs) <= 0 {
		return rows, nil
	}
	for _, proc := range procs {
		if proc == nil {
			continue
		}
		var err error
		rows, err = proc.ProcRows(ctx, rows)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

// ProcRow 使用一组处理器处理单条记录行数据
func ProcRow[ROW any](ctx context.Context, row ROW, procs ...RowsProc[ROW]) (ROW, error) {
	if len(procs) <= 0 {
		return row, nil
	}
	rows, err := ProcRows(ctx, []ROW{row}, procs...)
	if err != nil {
		var zero ROW
		return zero, err
	}
	return rows[0], nil
}

// RowsProcFunc 将一个函数转换为 RowsProc 接口
type RowsProcFunc[ROW any] func(ctx context.Context, rows []ROW) ([]ROW, error)

// ProcRows 实现 RowsProc.ProcRows 接口
func (f RowsProcFunc[ROW]) ProcRows(ctx context.Context, rows []ROW) ([]ROW, error) {
	if f == nil {
		return rows, nil
	}
	return f(ctx, rows)
}
