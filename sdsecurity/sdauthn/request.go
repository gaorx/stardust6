package sdauthn

import (
	"context"
	"github.com/gaorx/stardust6/sdtime"
	"time"
)

type Request interface {
	PrincipalId() PrincipalId
	ValidateRequest(ctx context.Context, p *Principal, requestExpiration Expiration, now time.Time) error
}

const (
	PrincipalUid      = "uid"
	PrincipalUsername = "username"
	PrincipalEmail    = "email"
	PrincipalPhone    = "phone"
)

type PrincipalId struct {
	Type   string `json:"type"`
	Domain string `json:"domain,omitempty"`
	ID     string `json:"id"`
}

type UserToken struct {
	Domain string `json:"d,omitempty"`
	ID     string `json:"id"`
	At     int64  `json:"at,omitempty"`
}

type UsernameAndPassword struct {
	Domain   string `json:"domain,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	_ Request = (*PrincipalId)(nil)
	_ Request = (*UserToken)(nil)
	_ Request = (*UsernameAndPassword)(nil)
)

func NewPrincipalId(typ, domain, id string) PrincipalId {
	return PrincipalId{Type: typ, Domain: domain, ID: id}
}

func NewUserId(uid string) PrincipalId {
	return NewPrincipalId(PrincipalUid, "", uid)
}

func (pid PrincipalId) PrincipalId() PrincipalId {
	return pid
}

func (pid PrincipalId) ValidateRequest(_ context.Context, p *Principal, _ Expiration, now time.Time) error {
	if p == nil || p.ID != pid.ID {
		return ErrPrincipalNotFound
	}
	return p.ValidateSelf(now)
}

func (pid PrincipalId) IsZero() bool {
	return pid.Type == "" && pid.Domain == "" && pid.ID == ""
}

func (pid PrincipalId) In(domain string) PrincipalId {
	pid1 := pid
	pid1.Domain = domain
	return pid1
}

func (pid PrincipalId) String() string {
	if pid.IsZero() {
		return ""
	}
	typ := pid.Type
	if typ == "" {
		typ = "unknown"
	}
	id := pid.ID
	if id == "" {
		id = "-"
	}
	s := typ + ":" + id
	if pid.Domain != "" {
		s += "@" + pid.Domain
	}
	return s
}

func NewUserToken(id string, at time.Time) UserToken {
	return UserToken{ID: id, At: sdtime.ToUnixS(at)}
}

func (t UserToken) PrincipalId() PrincipalId {
	return PrincipalId{PrincipalUid, t.Domain, t.ID}
}

func (t UserToken) ValidateRequest(_ context.Context, p *Principal, requestExpiration Expiration, now time.Time) error {
	if p == nil || p.ID != t.ID {
		return ErrPrincipalNotFound
	}
	if requestExpiration != nil && t.At > 0 {
		start := sdtime.FromUnixS(t.At)
		if requestExpiration.IsExpired(start, now) {
			return ErrCredentialExpired
		}
	}
	return p.ValidateSelf(now)
}

func (t UserToken) In(domain string) UserToken {
	t1 := t
	t1.Domain = domain
	return t1
}

func NewUsernameAndPassword(username, password string) UsernameAndPassword {
	return UsernameAndPassword{Username: username, Password: password}
}

func (u UsernameAndPassword) PrincipalId() PrincipalId {
	return PrincipalId{PrincipalUsername, u.Domain, u.Username}
}

func (u UsernameAndPassword) ValidateRequest(_ context.Context, p *Principal, _ Expiration, now time.Time) error {
	if p == nil {
		return ErrPrincipalNotFound
	}
	if p.Password != u.Password {
		return ErrCredentialInvalid
	}
	return p.ValidateSelf(now)
}

func (u UsernameAndPassword) In(domain string) UsernameAndPassword {
	u1 := u
	u1.Domain = domain
	return u1
}
