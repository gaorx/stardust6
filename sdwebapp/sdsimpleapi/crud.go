package sdsimpleapi

import (
	"context"
	"github.com/gaorx/stardust6/sderr/sdnotfounderr"
	"github.com/gaorx/stardust6/sdsql"
	"github.com/gaorx/stardust6/sdurl"
	"github.com/gaorx/stardust6/sdwebapp"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

// RequestCreate 创建实体的请求
type RequestCreate[T any] struct {
	Data T `json:"data"`
}

// RequestUpdate 更新实体的请求
type RequestUpdate[T any] struct {
	Data    T        `json:"data"`
	Columns []string `json:"columns"`
}

// RequestDelete 删除实体的请求
type RequestDelete[ID sdsql.EntityId] struct {
	Id ID `json:"id"`
}

// RequestGet 获取实体的请求
type RequestGet[ID sdsql.EntityId] struct {
	Id ID `json:"id"`
}

// RequestList 列出实体列表请求
type RequestList[F any] struct {
	Filter F `json:"filter"`
}

// RequestFind 分页查找实体请求
type RequestFind[F any] struct {
	Filter F   `json:"filter"`
	Page   int `json:"page"`
	Size   int `json:"size"`
}

// NullQuery 一个永远返回零值的查询
func NullQuery[F any, Q any](_ F) Q {
	var zero Q
	return zero
}

// CrudOptions CRUD选项
type CrudOptions struct {
	Enable      int // 使用EnableXXX常量组合
	GuardRead   sdwebapp.RouteGuard
	GuardWrite  sdwebapp.RouteGuard
	GuardCreate sdwebapp.RouteGuard
	GuardUpdate sdwebapp.RouteGuard
	GuardDelete sdwebapp.RouteGuard
	GuardGet    sdwebapp.RouteGuard
	GuardList   sdwebapp.RouteGuard
	GuardFind   sdwebapp.RouteGuard
	DocCreate   *sdwebapp.Doc
	DocUpdate   *sdwebapp.Doc
	DocDelete   *sdwebapp.Doc
	DocGet      *sdwebapp.Doc
	DocList     *sdwebapp.Doc
	DocFind     *sdwebapp.Doc
	Middlewares []echo.MiddlewareFunc
}

const (
	EnableCreate = 1 << iota
	EnableUpdate
	EnableDelete
	EnableGet
	EnableList
	EnableFind
	EnableRead  = EnableGet | EnableFind // 默认的读并不包括list，因为容易数据量过大
	EnableWrite = EnableCreate | EnableUpdate | EnableDelete
	EnableAll   = EnableRead | EnableWrite | EnableList
)

// Crud 创建一组CRUD路由，path参数是CRUD的根路径，factory是Repo的工厂函数，filter是查询条件生成函数
func Crud[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any, F any](
	path string,
	factory func(context2 echo.Context) TRepo,
	filter func(F) Q,
	opts *CrudOptions,
) sdwebapp.Routes {
	opts1 := lo.FromPtr(opts)
	var apis []*API
	if opts1.Enable&EnableCreate != 0 {
		api := Create(
			getPath(path, "create"),
			factory,
			selectGuard(opts1.GuardCreate, opts1.GuardWrite, sdwebapp.RejectAll()),
		).SetDoc(opts1.DocCreate).AddMiddlewares(opts1.Middlewares...)
		apis = append(apis, api)
	}
	if opts1.Enable&EnableUpdate != 0 {
		api := Update(
			getPath(path, "update"),
			factory,
			selectGuard(opts1.GuardUpdate, opts1.GuardWrite, sdwebapp.RejectAll()),
		).SetDoc(opts1.DocUpdate).AddMiddlewares(opts1.Middlewares...)
		apis = append(apis, api)
	}
	if opts1.Enable&EnableDelete != 0 {
		api := Delete(
			getPath(path, "delete"),
			factory,
			selectGuard(opts1.GuardDelete, opts1.GuardWrite, sdwebapp.RejectAll()),
		).SetDoc(opts1.DocDelete).AddMiddlewares(opts1.Middlewares...)
		apis = append(apis, api)
	}
	if opts1.Enable&EnableGet != 0 {
		api := Get(
			getPath(path, "get"),
			factory,
			selectGuard(opts1.GuardGet, opts1.GuardRead, sdwebapp.PermitAll()),
		).SetDoc(opts1.DocGet).AddMiddlewares(opts1.Middlewares...)
		apis = append(apis, api)
	}
	if opts1.Enable&EnableList != 0 {
		api := List(
			getPath(path, "list"),
			factory,
			filter,
			selectGuard(opts1.GuardList, opts1.GuardRead, sdwebapp.PermitAll()),
		).SetDoc(opts1.DocList).AddMiddlewares(opts1.Middlewares...)
		apis = append(apis, api)
	}
	if opts1.Enable&EnableFind != 0 {
		api := Find(
			getPath(path, "find"),
			factory,
			filter,
			selectGuard(opts1.GuardFind, opts1.GuardRead, sdwebapp.PermitAll()),
		).SetDoc(opts1.DocFind).AddMiddlewares(opts1.Middlewares...)
		apis = append(apis, api)
	}
	return lo.Map(apis, func(api *API, _ int) sdwebapp.Routable {
		return api
	})
}

// Create 创建实体API
func Create[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any](
	path string,
	factory func(echo.Context) TRepo,
	guard sdwebapp.RouteGuard,
) *API {
	h := func(ctx echo.Context, req RequestCreate[T]) *Result {
		repo := factory(ctx)
		created, err := repo.CreateAndGet(context.Background(), req.Data)
		if err != nil {
			return Err(makeCrudErr(err))
		}
		return OK(created)
	}
	return R(path, h, guard)
}

// Update 更新实体API
func Update[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any](
	path string,
	factory func(echo.Context) TRepo,
	guard sdwebapp.RouteGuard,
) *API {
	h := func(ctx echo.Context, req RequestUpdate[T]) *Result {
		repo := factory(ctx)
		updated, err := repo.UpdateAndGet(context.Background(), req.Data, req.Columns)
		if err != nil {
			return Err(makeCrudErr(err))
		}
		return OK(updated)
	}
	return R(path, h, guard)
}

// Delete 删除实体API
func Delete[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any](
	path string,
	factory func(echo.Context) TRepo,
	guard sdwebapp.RouteGuard,
) *API {
	h := func(ctx echo.Context, req RequestDelete[ID]) *Result {
		repo := factory(ctx)
		updated, err := repo.DeleteById(context.Background(), req.Id)
		if err != nil {
			return Err(makeCrudErr(err))
		}
		return OK(updated)
	}
	return R(path, h, guard)
}

// Get 通过ID获取实体API
func Get[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any](
	path string,
	factory func(echo.Context) TRepo,
	guard sdwebapp.RouteGuard,
) *API {
	h := func(ctx echo.Context, req RequestGet[ID]) *Result {
		repo := factory(ctx)
		entity, err := repo.GetById(context.Background(), req.Id)
		if err != nil {
			return Err(makeCrudErr(err))
		}
		return OK(entity)
	}
	return R(path, h, guard)
}

// List 列出符合条件的实体列表API
func List[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any, F any](
	path string,
	factory func(context2 echo.Context) TRepo,
	filter func(F) Q,
	guard sdwebapp.RouteGuard,
) *API {
	h := func(ctx echo.Context, req RequestList[F]) *Result {
		repo := factory(ctx)
		q := filter(req.Filter)
		list, err := repo.FindBy(context.Background(), q)
		if err != nil {
			return Err(makeCrudErr(err))
		}
		return OK(list)
	}
	return R(path, h, guard)
}

// Find 分页查找实体API
func Find[TRepo sdsql.Repository[T, ID, Q], T any, ID sdsql.EntityId, Q any, F any](
	path string,
	factory func(context2 echo.Context) TRepo,
	filter func(F) Q,
	guard sdwebapp.RouteGuard,
) *API {
	h := func(ctx echo.Context, req RequestFind[F]) *Result {
		repo := factory(ctx)
		q := filter(req.Filter)
		paged, err := repo.PaginateBy(context.Background(), q, sdsql.Page1(int64(req.Page), int64(req.Size)))
		if err != nil {
			return Err(makeCrudErr(err))
		}
		return OK(paged.Data).SetMetas(map[string]any{
			"page":       paged.PageNum,
			"size":       paged.PageSize,
			"totalPages": paged.TotalPages,
			"totalRows":  paged.TotalRows,
		})
	}
	return R(path, h, guard)
}

func selectGuard(preferred, secondary, def sdwebapp.RouteGuard) sdwebapp.RouteGuard {
	if preferred != nil {
		return preferred
	}
	if secondary != nil {
		return secondary
	}
	return def
}

func getPath(path string, action string) string {
	return sdurl.JoinPath(path, action)
}

func makeCrudErr(err any) (string, any) {
	if err == nil {
		return CodeOK, nil
	}
	switch err1 := err.(type) {
	case error:
		if sdnotfounderr.Is(err1) {
			return CodeNotFound, err1
		}
		return CodeDataError, err1
	default:
		return CodeDataError, err1
	}
}
