package sdauthn

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"time"
)

// Load 通过request加载一个principal,如果没找到返回ErrPrincipalNotFound。
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

// Authenticate 验证请求，如果验证通过返回principal，否则返回错误，error可以携带错误原因；
// 注意如果找到了principal但是验证失败，即使error不为nil也会返回principal。
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
