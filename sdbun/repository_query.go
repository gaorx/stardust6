package sdbun

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/uptrace/bun"
)

// RepositoryQuery 仓库查询容器
type RepositoryQuery struct {
	s func(s *bun.SelectQuery) *bun.SelectQuery
	u func(u *bun.UpdateQuery) *bun.UpdateQuery
	d func(d *bun.DeleteQuery) *bun.DeleteQuery
}

// Q 创建一个查询容器
func Q(fn any) *RepositoryQuery {
	switch fn1 := fn.(type) {
	case nil:
		return nil
	case func(*bun.SelectQuery) *bun.SelectQuery:
		return &RepositoryQuery{s: fn1}
	case func(*bun.UpdateQuery) *bun.UpdateQuery:
		return &RepositoryQuery{u: fn1}
	case func(*bun.DeleteQuery) *bun.DeleteQuery:
		return &RepositoryQuery{d: fn1}
	default:
		panic(sderr.Newf("invalid query func type"))
	}
}

// IsZero 判断是否为零值
func (q *RepositoryQuery) IsZero() bool {
	return q == nil || (q.s == nil && q.u == nil && q.d == nil)
}

func (q *RepositoryQuery) applySelect(s *bun.SelectQuery) *bun.SelectQuery {
	if q == nil || q.s == nil {
		return s
	}
	return q.s(s)
}

func (q *RepositoryQuery) applyUpdate(u *bun.UpdateQuery) *bun.UpdateQuery {
	if q == nil || q.u == nil {
		return u
	}
	return q.u(u)
}

func (q *RepositoryQuery) applyDelete(d *bun.DeleteQuery) *bun.DeleteQuery {
	if q == nil || q.d == nil {
		return d
	}
	return q.d(d)
}
