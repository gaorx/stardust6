package sdsql

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
)

// Aggregator 将多行数据所需要聚合的部分进行统一处理，然后将聚合结果回填到行中
type Aggregator[ROW any, ON comparable, COMPLEMENT any] struct {
	disabled        bool
	Collect         func(context.Context, ROW) []ON
	Fetch           func(context.Context, []ON) (map[ON]COMPLEMENT, error)
	Complete        func(context.Context, ROW, map[ON]COMPLEMENT) (ROW, error)
	CompleteInplace func(context.Context, ROW, map[ON]COMPLEMENT) error
}

var _ RowsProc[struct{}] = Aggregator[struct{}, string, struct{}]{}

// If 设置此聚合器否启用
func (a Aggregator[ROW, ON, COMPLEMENT]) If(enabled bool) Aggregator[ROW, ON, COMPLEMENT] {
	a1 := a
	a1.disabled = !enabled
	return a1
}

// ProcRows 处理多行数据
func (a Aggregator[ROW, ON, COMPLEMENT]) ProcRows(ctx context.Context, rows []ROW) ([]ROW, error) {
	if a.disabled {
		return rows, nil
	}

	if a.Collect == nil || a.Fetch == nil || (a.Complete == nil && a.CompleteInplace == nil) {
		return nil, sderr.Newf("invalid aggregator")
	}

	if len(rows) <= 0 {
		return rows, nil
	}
	collected := make(map[ON]struct{})
	for _, row := range rows {
		ons := a.Collect(ctx, row)
		if len(ons) <= 0 {
			// do nothing
		} else if len(ons) == 1 {
			collected[ons[0]] = struct{}{}
		} else {
			for _, on := range ons {
				collected[on] = struct{}{}
			}
		}
	}
	complements, err := a.Fetch(ctx, lo.Keys(collected))
	if err != nil {
		return nil, err
	}
	if a.Complete != nil {
		newRows := make([]ROW, len(rows))
		for _, row := range rows {
			newRow, err := a.Complete(ctx, row, complements)
			if err != nil {
				return nil, err
			}
			newRows = append(newRows, newRow)
		}
		return newRows, nil
	} else if a.CompleteInplace != nil {
		for _, row := range rows {
			err := a.CompleteInplace(ctx, row, complements)
			if err != nil {
				return nil, err
			}
		}
		return rows, nil
	} else {
		return nil, sderr.Newf("aggregator.complete is nil")
	}
}
