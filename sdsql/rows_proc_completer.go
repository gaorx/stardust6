package sdsql

import (
	"context"
)

// Completer 补全记录行数据的处理器
type Completer[ROW any] func(context.Context, ROW) (ROW, error)

// InplaceCompleter 原地补全记录行数据的处理器，无须返回新行，在记录行是指针类型下，用这个更好
type InplaceCompleter[ROW any] func(context.Context, ROW) error

var (
	_ RowsProc[struct{}] = (Completer[struct{}])(nil)
	_ RowsProc[struct{}] = (InplaceCompleter[struct{}])(nil)
)

// ProcRows 处理记录行数据
func (c Completer[ROW]) ProcRows(ctx context.Context, rows []ROW) ([]ROW, error) {
	if c == nil {
		return rows, nil
	}
	if len(rows) <= 0 {
		return rows, nil
	}
	newRows := make([]ROW, 0, len(rows))
	for _, row := range rows {
		newRow, err := c(ctx, row)
		if err != nil {
			return nil, err
		}
		newRows = append(newRows, newRow)
	}
	return newRows, nil
}

// ProcRows 处理记录行数据
func (c InplaceCompleter[ROW]) ProcRows(ctx context.Context, rows []ROW) ([]ROW, error) {
	if c == nil {
		return rows, nil
	}
	if len(rows) <= 0 {
		return rows, nil
	}
	for _, row := range rows {
		err := c(ctx, row)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}
