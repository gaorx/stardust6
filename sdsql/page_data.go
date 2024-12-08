package sdsql

import (
	"slices"
)

type PageData[T any] struct {
	Data       []T   `json:"data"`
	TotalRows  int64 `json:"totalRows"`
	PageNum    int64 `json:"pageNum"`
	PageSize   int64 `json:"pageSize"`
	TotalPages int64 `json:"totalPages"`
}

func NewPageData[T any](data []T, page Page, totalRows int64) *PageData[T] {
	totalPages := totalRows / page.Size()
	if totalRows%page.Size() > 0 {
		totalPages++
	}
	if data == nil {
		data = []T{}
	}
	return &PageData[T]{
		Data:       data,
		TotalRows:  totalRows,
		PageNum:    page.Num(),
		PageSize:   page.Size(),
		TotalPages: totalPages,
	}
}

func (pd *PageData[T]) Len() int {
	return len(pd.Data)
}

func PaginateSlice[ROW any](rows []ROW, page Page) *PageData[ROW] {
	limit, offset := page.LimitAndOffset()
	start, end := offset, offset+limit
	numRows := int64(len(rows))
	if start > end {
		start, end = end, start
	}
	if start > numRows {
		start = numRows
	}
	if start < 0 {
		start = 0
	}

	if end > numRows {
		end = numRows
	}
	if end < 0 {
		end = 0
	}
	return &PageData[ROW]{
		Data:       slices.Clone(rows[start:end]),
		TotalRows:  numRows,
		PageSize:   limit,
		PageNum:    page.Num(),
		TotalPages: (numRows + limit - 1) / limit,
	}
}
