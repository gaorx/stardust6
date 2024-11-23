package sdauthn

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"time"
)

func Load(ctx context.Context, req Request, loader Loader) (*Principal, error) {
	if req == nil {
		return nil, sderr.Newf("request is nil")
	}
	if loader == nil {
		loader = LoaderFunc(defaultPrincipalLoader)
	}
	pid := req.PrincipalId()
	if pid.IsZero() {
		return nil, ErrPrincipalNotFound
	}
	p, err := loader.LoadPrincipal(ctx, pid)
	if err != nil {
		return nil, sderr.Wrapf(err, "load principal failed")
	}
	return p, nil
}

func Authenticate(
	ctx context.Context,
	req Request,
	loader Loader,
	requestExpiration Expiration,
	now time.Time,
) (*Principal, error) {
	p, err := Load(ctx, req, loader)
	if err != nil {
		return nil, err
	}
	if requestExpiration == nil {
		requestExpiration = NoExpiration()
	}
	if err := req.ValidateRequest(ctx, p, requestExpiration, now); err != nil {
		return p, err
	}
	return p, nil
}

func defaultPrincipalLoader(_ context.Context, _ PrincipalId) (*Principal, error) {
	return nil, ErrPrincipalNotFound
}
