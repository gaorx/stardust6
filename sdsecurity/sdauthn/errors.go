package sdauthn

import (
	"github.com/gaorx/stardust6/sderr"
)

var (
	ErrPrincipalNotFound = sderr.Sentinel("principal not found")
	ErrCredentialInvalid = sderr.Sentinel("credential invalid")
	ErrPrincipalDisabled = sderr.Sentinel("principal disabled")
	ErrPrincipalExpired  = sderr.Sentinel("principal expired")
	ErrPrincipalLocked   = sderr.Sentinel("principal locked")
	ErrCredentialExpired = sderr.Sentinel("credential expired")
)
