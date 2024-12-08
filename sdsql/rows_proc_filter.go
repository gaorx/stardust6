package sdsql

import (
	"context"
)

// Filter 过滤记录行数据的处理器，只保留符合条件的记录行
type Filter[ROW any] func(context.Context, ROW) bool

var _ RowsProc[struct{}] = (Filter[struct{}])(nil)

// ProcRows 处理记录行数据
func (f Filter[ROW]) ProcRows(ctx context.Context, rows []ROW) ([]ROW, error) {
	if f == nil {
		return rows, nil
	}
	if len(rows) <= 0 {
		return rows, nil
	}
	newRows := make([]ROW, 0)
	for _, row := range rows {
		if f(ctx, row) {
			newRows = append(newRows, row)
		}
	}
	return newRows, nil
}
