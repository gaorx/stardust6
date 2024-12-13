package sdsimpleapi

import (
	"context"
	"github.com/gaorx/stardust6/sdbun"
	"github.com/gaorx/stardust6/sdbun/sdbunsqlite"
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"testing"
)

func TestCrud(t *testing.T) {
	is := assert.New(t)

	// 创建数据库并创建表
	db, err := sdbunsqlite.Open(sdbunsqlite.Address{
		DSN: "file::memory:?cache=shared",
	})
	is.NoError(err)
	_, err = db.NewCreateTable().Model(&todoEntity{}).Exec(context.Background())
	is.NoError(err)

	app := sdwebapp.New(nil)
	app.MustInstall(
		sdwebapp.Inject(sdwebapp.State(db)),
		Crud("/api/todo", newTodoRepo, todoQuery, &CrudOptions{
			Enable:     EnableAll,
			GuardRead:  sdwebapp.PermitAll(),
			GuardWrite: sdwebapp.PermitAll(),
		}),
	)

	// Create
	NewTestRequest("/api/todo/create", RequestCreate[*todoEntity]{
		Data: &todoEntity{
			Id:          "todo1",
			Title:       "title1",
			Description: "desc1",
			IsComplete:  true,
			Priority:    3,
		},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotNil(r.Data)
		is.Equal("todo1", r.Data.Id)
		is.Equal("title1", r.Data.Title)
	})

	// Update
	NewTestRequest("/api/todo/update", RequestUpdate[*todoEntity]{
		Data: &todoEntity{
			Id:          "todo1",
			Title:       "title1_update",
			Description: "desc1_update",
			IsComplete:  true,
			Priority:    4,
		},
		Columns: []string{"priority"},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotNil(r.Data)
		is.Equal("todo1", r.Data.Id)
		is.Equal("title1", r.Data.Title)
		is.Equal(4, r.Data.Priority)
	})
	NewTestRequest("/api/todo/update", RequestUpdate[*todoEntity]{
		Data: &todoEntity{
			Id:          "id-not-exists",
			Title:       "title1_update",
			Description: "desc1_update",
			IsComplete:  true,
			Priority:    5,
		},
		Columns: []string{"priority"},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeNotFound, r.Code)
	})

	// Get
	NewTestRequest("/api/todo/get", RequestGet[string]{
		Id: "todo1",
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotNil(r.Data)
	})
	NewTestRequest("/api/todo/get", RequestGet[string]{
		Id: "id-not-exists",
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeNotFound, r.Code)
		is.Nil(r.Data)
	})

	// List
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: nil,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotEmpty(r.Data)
		is.Equal(1, len(r.Data))
		is.Equal("todo1", r.Data[0].Id)
		is.Equal("title1", r.Data[0].Title)
	})
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: &todoFilter{
			Title: "title1",
		},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotEmpty(r.Data)
		is.Equal(1, len(r.Data))
		is.Equal("todo1", r.Data[0].Id)
		is.Equal("title1", r.Data[0].Title)
	})
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: &todoFilter{
			Title: "title-not-exists",
		},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Empty(r.Data)
	})
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: &todoFilter{
			Title:       "title1",
			Description: "desc1",
		},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotEmpty(r.Data)
		is.Equal(1, len(r.Data))
		is.Equal("todo1", r.Data[0].Id)
		is.Equal("title1", r.Data[0].Title)
		is.Equal("desc1", r.Data[0].Description)
	})
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: &todoFilter{
			Title:       "title1",
			Description: "desc-not-exists",
		},
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.Empty(r.Data)
	})

	// Find
	NewTestRequest("/api/todo/find", RequestFind[*todoFilter]{
		Filter: nil,
		Page:   1,
		Size:   10,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotEmpty(r.Data)
		is.Equal(1, len(r.Data))
		is.Equal("todo1", r.Data[0].Id)
		is.Equal("title1", r.Data[0].Title)
		is.Equal(1, r.Meta.Get("page").AsInt())
		is.Equal(10, r.Meta.Get("size").AsInt())
		is.Equal(1, r.Meta.Get("totalRows").AsInt())
		is.Equal(1, r.Meta.Get("totalPages").AsInt())
	})
	NewTestRequest("/api/todo/find", RequestFind[*todoFilter]{
		Filter: &todoFilter{
			Title: "title1",
		},
		Page: 1,
		Size: 10,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.NotEmpty(r.Data)
		is.Equal(CodeOK, r.Code)
		is.Equal(1, len(r.Data))
		is.Equal("todo1", r.Data[0].Id)
		is.Equal("title1", r.Data[0].Title)
		is.Equal(1, r.Meta.Get("page").AsInt())
		is.Equal(10, r.Meta.Get("size").AsInt())
		is.Equal(1, r.Meta.Get("totalRows").AsInt())
		is.Equal(1, r.Meta.Get("totalPages").AsInt())
	})
	NewTestRequest("/api/todo/find", RequestFind[*todoFilter]{
		Filter: &todoFilter{
			Title:       "title1",
			Description: "desc1",
		},
		Page: 1,
		Size: 10,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.NotEmpty(r.Data)
		is.Equal(1, len(r.Data))
		is.Equal("todo1", r.Data[0].Id)
		is.Equal("title1", r.Data[0].Title)
		is.Equal(1, r.Meta.Get("page").AsInt())
		is.Equal(10, r.Meta.Get("size").AsInt())
		is.Equal(1, r.Meta.Get("totalRows").AsInt())
		is.Equal(1, r.Meta.Get("totalPages").AsInt())
	})
	NewTestRequest("/api/todo/find", RequestFind[*todoFilter]{
		Filter: &todoFilter{
			Title: "title-not-exists",
		},
		Page: 1,
		Size: 10,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.Empty(r.Data)
		is.Equal(1, r.Meta.Get("page").AsInt())
		is.Equal(10, r.Meta.Get("size").AsInt())
		is.Equal(0, r.Meta.Get("totalRows").AsInt())
		is.Equal(0, r.Meta.Get("totalPages").AsInt())
	})
	NewTestRequest("/api/todo/find", RequestFind[*todoFilter]{
		Filter: &todoFilter{
			Title:       "title1",
			Description: "desc-not-exists",
		},
		Page: 1,
		Size: 10,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.Empty(r.Data)
		is.Equal(1, r.Meta.Get("page").AsInt())
		is.Equal(10, r.Meta.Get("size").AsInt())
		is.Equal(0, r.Meta.Get("totalRows").AsInt())
		is.Equal(0, r.Meta.Get("totalPages").AsInt())
	})

	// Delete
	NewTestRequest("/api/todo/delete", RequestDelete[string]{
		Id: "id-not-exists",
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[struct{}]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
	})
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: nil,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.Equal(1, len(r.Data))
	})
	NewTestRequest("/api/todo/delete", RequestDelete[string]{
		Id: "todo1",
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[struct{}]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
	})
	NewTestRequest("/api/todo/list", RequestList[*todoFilter]{
		Filter: nil,
	}).Call(app).Let(func(resp *sdwebapp.TestResponse) {
		var r ResultT[[]*todoEntity]
		is.Equal(200, resp.Code)
		is.NotPanics(func() {
			resp.BodyJson(&r)
		})
		is.Equal(CodeOK, r.Code)
		is.Empty(r.Data)
	})
}

type todoFilter struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type todoEntity struct {
	bun.BaseModel `json:"-" bun:"table:t_todo,alias:u"`
	Id_           int64  `json:"id_" bun:"id_,pk,autoincrement"`
	Id            string `json:"id" bun:"id,unique" repoid:"true"`
	Title         string `json:"title" bun:"title"`
	Description   string `json:"description" bun:"description"`
	IsComplete    bool   `json:"is_complete" bun:"is_complete"`
	Priority      int    `json:"priority" bun:"priority"`
}

type todoRepo sdbun.Repository[*todoEntity, string]

func newTodoRepo(c echo.Context) todoRepo {
	return sdbun.RepoOf[*todoEntity, string](sdwebapp.C(c).State().(*bun.DB))
}

func todoQuery(filter todoFilter) *sdbun.RepositoryQuery {
	q := func(s *bun.SelectQuery) *bun.SelectQuery {
		if lo.IsEmpty(filter) {
			return s
		}
		if filter.Title != "" {
			s.Where("title=?", filter.Title)
		}
		if filter.Description != "" {
			s.Where("description=?", filter.Description)
		}
		return s
	}
	return sdbun.Q(q)
}
