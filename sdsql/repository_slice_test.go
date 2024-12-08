package sdsql

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestSliceRepository(t *testing.T) {
	is := assert.New(t)
	repo := testNewSliceRepo(10)
	// CountAll
	n, err := repo.CountAll(context.Background())
	is.NoError(err)
	is.Equal(int64(10), n)

	// Count
	n, err = repo.CountBy(context.Background(), func(e *todoEntity1) bool {
		return e.IsComplete
	})
	is.NoError(err)
	is.Equal(int64(5), n)

	// GetById
	actual, err := repo.GetById(context.Background(), "id-2")
	is.NoError(err)
	is.True("title-2" == actual.Title && "desc-2" == actual.Description)
	actual, err = repo.GetById(context.Background(), "not-exists")
	is.Error(err)
	is.True(IsNotFoundErr(err))

	// GetBy
	actual, err = repo.GetBy(context.Background(), func(e *todoEntity1) bool {
		return e.Title == "title-3"
	})
	is.NoError(err)
	is.True("title-3" == actual.Title && "desc-3" == actual.Description)
	actual, err = repo.GetBy(context.Background(), func(e *todoEntity1) bool {
		return e.Title == "not-exists-title"
	})
	is.Error(err)
	is.True(IsNotFoundErr(err))

	// FindAll
	actuals, err := repo.FindAll(context.Background())
	is.NoError(err)
	is.Equal(10, len(actuals))
	is.True("title-1" == actuals[0].Title && "desc-1" == actuals[0].Description)
	is.True("title-10" == actuals[9].Title && "desc-10" == actuals[9].Description)

	// FindByIds
	actuals, err = repo.FindByIds(context.Background(), []string{"id-3", "id-7", "not-exists"})
	is.NoError(err)
	is.Equal(2, len(actuals))
	is.True("title-3" == actuals[0].Title && "desc-3" == actuals[0].Description)
	is.True("title-7" == actuals[1].Title && "desc-7" == actuals[1].Description)

	// FindBy
	actuals, err = repo.FindBy(context.Background(), func(e *todoEntity1) bool {
		return slices.Contains([]string{"title-2", "title-4", "title-not-exists"}, e.Title)
	})
	is.NoError(err)
	is.Equal(2, len(actuals))
	is.True("title-2" == actuals[0].Title && "desc-2" == actuals[0].Description)
	is.True("title-4" == actuals[1].Title && "desc-4" == actuals[1].Description)

	// PaginateAll
	actualPage, err := repo.PaginateAll(context.Background(), Page1(2, 3))
	is.NoError(err)
	is.Equal(int64(10), actualPage.TotalRows)
	is.Equal(int64(4), actualPage.TotalPages)
	is.Equal(int64(2), actualPage.PageNum)
	is.Equal(int64(3), actualPage.PageSize)
	is.Equal(3, len(actualPage.Data))
	is.True("title-4" == actualPage.Data[0].Title && "desc-4" == actualPage.Data[0].Description)
	is.True("title-5" == actualPage.Data[1].Title && "desc-5" == actualPage.Data[1].Description)
	is.True("title-6" == actualPage.Data[2].Title && "desc-6" == actualPage.Data[2].Description)
	actualPage, err = repo.PaginateAll(context.Background(), Page1(4, 3))
	is.NoError(err)
	is.Equal(int64(10), actualPage.TotalRows)
	is.Equal(int64(4), actualPage.TotalPages)
	is.Equal(int64(4), actualPage.PageNum)
	is.Equal(int64(3), actualPage.PageSize)
	is.Equal(1, len(actualPage.Data))
	is.True("title-10" == actualPage.Data[0].Title && "desc-10" == actualPage.Data[0].Description)
	actualPage, err = repo.PaginateAll(context.Background(), Page1(5, 3))
	is.NoError(err)
	is.Equal(int64(10), actualPage.TotalRows)
	is.Equal(int64(4), actualPage.TotalPages)
	is.Equal(int64(5), actualPage.PageNum)
	is.Equal(int64(3), actualPage.PageSize)
	is.Equal(0, len(actualPage.Data))

	// PaginateBy
	actualPage, err = repo.PaginateBy(context.Background(), func(e *todoEntity1) bool {
		return e.IsComplete
	}, Page1(2, 2))
	is.NoError(err)
	is.Equal(int64(5), actualPage.TotalRows)
	is.Equal(int64(3), actualPage.TotalPages)
	is.Equal(int64(2), actualPage.PageNum)
	is.Equal(int64(2), actualPage.PageSize)
	is.Equal(2, len(actualPage.Data))
	is.True("title-5" == actualPage.Data[0].Title && "desc-5" == actualPage.Data[0].Description)
	is.True("title-7" == actualPage.Data[1].Title && "desc-7" == actualPage.Data[1].Description)
	actualPage, err = repo.PaginateBy(context.Background(), func(e *todoEntity1) bool {
		return e.IsComplete
	}, Page1(3, 2))
	is.NoError(err)
	is.Equal(int64(5), actualPage.TotalRows)
	is.Equal(int64(3), actualPage.TotalPages)
	is.Equal(int64(3), actualPage.PageNum)
	is.Equal(int64(2), actualPage.PageSize)
	is.Equal(1, len(actualPage.Data))
	is.True("title-9" == actualPage.Data[0].Title && "desc-9" == actualPage.Data[0].Description)
	actualPage, err = repo.PaginateBy(context.Background(), func(e *todoEntity1) bool {
		return e.IsComplete
	}, Page1(4, 2))
	is.NoError(err)
	is.Equal(int64(5), actualPage.TotalRows)
	is.Equal(int64(3), actualPage.TotalPages)
	is.Equal(int64(4), actualPage.PageNum)
	is.Equal(int64(2), actualPage.PageSize)
	is.Equal(0, len(actualPage.Data))
}

func testNewSliceRepo(n int) SliceRepository[*todoEntity1, string] {
	var entities []*todoEntity1
	for i := 0; i < n; i++ {
		entity := &todoEntity1{
			Id:          fmt.Sprintf("id-%d", i+1),
			Title:       fmt.Sprintf("title-%d", i+1),
			Description: fmt.Sprintf("desc-%d", i+1),
			IsComplete:  i%2 == 0,
			Priority:    1,
		}
		entities = append(entities, entity)
	}
	return SliceRepoOf(entities)
}

type todoEntity1 struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsComplete  bool   `json:"is_complete"`
	Priority    int    `json:"priority"`
}

func (e *todoEntity1) EntityId() string {
	return e.Id
}

func (e *todoEntity1) Clone() *todoEntity1 {
	if e == nil {
		return nil
	}
	e1 := *e
	return &e1
}
