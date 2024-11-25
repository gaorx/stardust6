package sdauthn

import (
	"context"
	"github.com/gaorx/stardust6/sdtime"
	"time"
)

// Request 发起认证的请求
type Request interface {
	// PrincipalId 返回请求对应的PrincipalId
	PrincipalId() PrincipalId
	// ValidateRequest 验证请求，如果验证通过返回nil，否则返回错误
	ValidateRequest(ctx context.Context, p *Principal, requestExpiration Expiration, now time.Time) error
}

const (
	// PrincipalUid 使用UID作为ID的请求类型
	PrincipalUid = "uid"
	// PrincipalUsername 使用用户名作为ID的请求类型
	PrincipalUsername = "username"
	// PrincipalEmail 使用邮箱作为ID的请求类型
	PrincipalEmail = "email"
	// PrincipalPhone 使用手机号作为ID的请求类型
	PrincipalPhone = "phone"
)

// PrincipalId 表示要验证的主体的ID
type PrincipalId struct {
	// Type ID类型
	Type string `json:"type"`
	// Domain 请求所在的领域
	Domain string `json:"domain,omitempty"`
	// ID 请求的ID，根据type不同也可以是username或者email等
	ID string `json:"id"`
}

// UserToken 表示用户令牌的请求
type UserToken struct {
	// Domain 请求所在的领域
	Domain string `json:"d,omitempty"`
	// ID 令牌中的用户ID
	ID string `json:"id"`
	// At 令牌的生成时间
	At int64 `json:"at,omitempty"`
}

// UsernameAndPassword 表示用户名和密码的请求
type UsernameAndPassword struct {
	// Domain 请求所在的领域
	Domain string `json:"domain,omitempty"`
	// Username 用户名
	Username string `json:"username"`
	// Password 密码
	Password string `json:"password"`
}

var (
	_ Request = (*PrincipalId)(nil)
	_ Request = (*UserToken)(nil)
	_ Request = (*UsernameAndPassword)(nil)
)

// NewPrincipalId 创建一个PrincipalId
func NewPrincipalId(typ, domain, id string) PrincipalId {
	return PrincipalId{Type: typ, Domain: domain, ID: id}
}

// NewUserId 创建一个使用UID作为ID的PrincipalId
func NewUserId(uid string) PrincipalId {
	return NewPrincipalId(PrincipalUid, "", uid)
}

// PrincipalId 返回请求对应的PrincipalId
func (pid PrincipalId) PrincipalId() PrincipalId {
	return pid
}

// ValidateRequest 验证请求，如果验证通过返回nil，否则返回错误
func (pid PrincipalId) ValidateRequest(_ context.Context, p *Principal, _ Expiration, now time.Time) error {
	if p == nil || p.ID != pid.ID {
		return ErrPrincipalNotFound
	}
	return p.ValidateSelf(now)
}

// IsZero 是否是空值
func (pid PrincipalId) IsZero() bool {
	return pid.Type == "" && pid.Domain == "" && pid.ID == ""
}

// In 返回在某个领域下的PrincipalId
func (pid PrincipalId) In(domain string) PrincipalId {
	pid1 := pid
	pid1.Domain = domain
	return pid1
}

// String 返回字符串表示
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

// NewUserToken 创建一个用户令牌
func NewUserToken(id string, at time.Time) UserToken {
	return UserToken{ID: id, At: sdtime.ToUnixS(at)}
}

// PrincipalId 返回请求对应的PrincipalId
func (t UserToken) PrincipalId() PrincipalId {
	return PrincipalId{PrincipalUid, t.Domain, t.ID}
}

// ValidateRequest 验证请求，如果验证通过返回nil，否则返回错误
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

// In 返回在某个领域下的用户令牌
func (t UserToken) In(domain string) UserToken {
	t1 := t
	t1.Domain = domain
	return t1
}

// NewUsernameAndPassword 创建一个用户名和密码的请求
func NewUsernameAndPassword(username, password string) UsernameAndPassword {
	return UsernameAndPassword{Username: username, Password: password}
}

// PrincipalId 返回请求对应的PrincipalId
func (u UsernameAndPassword) PrincipalId() PrincipalId {
	return PrincipalId{PrincipalUsername, u.Domain, u.Username}
}

// ValidateRequest 验证请求，如果验证通过返回nil，否则返回错误
func (u UsernameAndPassword) ValidateRequest(_ context.Context, p *Principal, _ Expiration, now time.Time) error {
	if p == nil {
		return ErrPrincipalNotFound
	}
	if p.Password != u.Password {
		return ErrCredentialInvalid
	}
	return p.ValidateSelf(now)
}

// In 返回在某个领域下的用户名和密码的请求
func (u UsernameAndPassword) In(domain string) UsernameAndPassword {
	u1 := u
	u1.Domain = domain
	return u1
}
