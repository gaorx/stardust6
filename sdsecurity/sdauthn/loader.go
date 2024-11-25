package sdauthn

import (
	"context"
)

// Loader 通过PrincipalId加载Principal
type Loader interface {
	LoadPrincipal(ctx context.Context, pid PrincipalId) (*Principal, error)
}

// LoaderFunc 通过函数实现Loader
type LoaderFunc func(ctx context.Context, pid PrincipalId) (*Principal, error)

func (f LoaderFunc) LoadPrincipal(ctx context.Context, pid PrincipalId) (*Principal, error) {
	return f(ctx, pid)
}

func LoaderOf(p *Principal) LoaderFunc {
	return LoaderFunc(func(_ context.Context, _ PrincipalId) (*Principal, error) {
		return p, nil
	})
}
