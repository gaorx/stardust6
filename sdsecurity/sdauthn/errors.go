package sdauthn

import (
	"github.com/gaorx/stardust6/sderr"
)

var (
	// ErrPrincipalNotFound 没找到principal
	ErrPrincipalNotFound = sderr.Sentinel("principal not found")
	// ErrCredentialInvalid 找到了principal但是credential无效，例如密码错误
	ErrCredentialInvalid = sderr.Sentinel("credential invalid")
	// ErrPrincipalDisabled principal被禁用
	ErrPrincipalDisabled = sderr.Sentinel("principal disabled")
	// ErrPrincipalExpired principal过期
	ErrPrincipalExpired = sderr.Sentinel("principal expired")
	// ErrPrincipalLocked principal被锁定
	ErrPrincipalLocked = sderr.Sentinel("principal locked")
	// ErrCredentialExpired credential过期，例如user token超过有效期
	ErrCredentialExpired = sderr.Sentinel("credential expired")
)
