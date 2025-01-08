package sdblueprint

type APIGuard int

const (
	APIPermitAll = iota + 1
	APIRejectAll
	APIAuthenticated
	APIHasAuthority
	APIIsMatched
)
