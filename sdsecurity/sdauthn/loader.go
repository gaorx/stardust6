package sdauthn

import (
	"context"
)

type Loader interface {
	LoadPrincipal(ctx context.Context, pid PrincipalId) (*Principal, error)
}

type LoaderFunc func(ctx context.Context, pid PrincipalId) (*Principal, error)

func (f LoaderFunc) LoadPrincipal(ctx context.Context, pid PrincipalId) (*Principal, error) {
	return f(ctx, pid)
}

func LoaderOf(p *Principal) LoaderFunc {
	return LoaderFunc(func(_ context.Context, _ PrincipalId) (*Principal, error) {
		return p, nil
	})
}
