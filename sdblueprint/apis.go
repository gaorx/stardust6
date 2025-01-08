package sdblueprint

import (
	"github.com/samber/lo"
	"slices"
	"strings"
)

type APIDecls []*APIDecl

func (apis APIDecls) ForEach(f func(decl *APIDecl)) {
	if f == nil {
		return
	}
	for _, a := range apis {
		if a != nil {
			f(a)
		}
	}
}

func (apis APIDecls) GetById(id string) *APIDecl {
	if id == "" {
		return nil
	}
	for _, a := range apis {
		if a.id == id {
			return a
		}
	}
	return nil
}

func (apis APIDecls) GetByPath(method, path string) *APIDecl {
	if path == "" {
		return nil
	}
	api, _ := lo.Find(apis, func(v *APIDecl) bool {
		return v.httpMethod == method && v.httpPath == trimAPIPath(path)
	})
	return api
}

func (apis APIDecls) FindById(ids ...string) APIDecls {
	if len(ids) <= 0 {
		return nil
	}
	return lo.Filter(apis, func(v *APIDecl, _ int) bool {
		return lo.Contains(ids, v.id)
	})
}

func (apis APIDecls) FindByPath(paths ...string) APIDecls {
	if len(paths) <= 0 {
		return nil
	}
	trimmedPaths := lo.Map(paths, func(v string, _ int) string {
		return trimAPIPath(v)
	})
	return lo.Filter(apis, func(v *APIDecl, _ int) bool {
		return lo.Contains(trimmedPaths, v.httpPath)
	})
}

func (apis APIDecls) FindByPathPrefix(prefix string) APIDecls {
	if prefix == "" {
		return slices.Clone(apis)
	}

	return lo.Filter(apis, func(v *APIDecl, _ int) bool {
		return strings.HasPrefix(v.httpPath, trimAPIPath(prefix))
	})
}

func (apis APIDecls) FindByCategory(category string) APIDecls {
	if category == "" {
		return nil
	}
	return lo.Filter(apis, func(v *APIDecl, _ int) bool {
		return v.category == category
	})
}
