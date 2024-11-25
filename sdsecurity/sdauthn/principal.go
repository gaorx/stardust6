package sdauthn

import (
	"slices"
	"time"
)

// Principal principal，表示要验证的主体
type Principal struct {
	ID          string    `json:"id"`
	Authorities []string  `json:"authorities,omitempty"`
	Password    string    `json:"password,omitempty"`
	SmsVC       []string  `json:"smsVC,omitempty"`
	EmailVC     []string  `json:"emailVCm,omitempty"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	AvatarUrl   string    `json:"avatarUrl,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	Disabled    bool      `json:"disabled,omitempty"`
	Locked      bool      `json:"locked,omitempty"`
	Expiry      time.Time `json:"expiry,omitempty"`
}

// PrincipalIdForDomain 返回这个Principal在某个domain下的PrincipalId
func (p *Principal) PrincipalIdForDomain(domain string) PrincipalId {
	if p == nil {
		return PrincipalId{}
	}
	return PrincipalId{Domain: domain, Type: PrincipalUid, ID: p.ID}
}

// ValidateSelf 验证自身是否有效
func (p *Principal) ValidateSelf(now time.Time) error {
	if p == nil {
		return ErrPrincipalNotFound
	}
	if p.ID == "" {
		return ErrPrincipalNotFound
	}
	if p.Disabled {
		return ErrPrincipalDisabled
	}
	if !p.Expiry.IsZero() {
		if now.After(p.Expiry) {
			return ErrPrincipalExpired
		}
	}
	if p.Locked {
		return ErrPrincipalLocked
	}
	return nil
}

// HasAuthority 是否有某个权限
func (p *Principal) HasAuthority(authority string) bool {
	if p == nil {
		return false
	}
	return slices.Contains(p.Authorities, authority)
}

// HasAllAuthorities 是否有所有权限
func (p *Principal) HasAllAuthorities(authorities ...string) bool {
	if p == nil {
		return false
	}
	if len(p.Authorities) <= 0 || len(authorities) <= 0 {
		return false
	}
	for _, a := range authorities {
		if !p.HasAuthority(a) {
			return false
		}
	}
	return true
}

// HasAnyAuthority 是否有几个权限中的任一个
func (p *Principal) HasAnyAuthority(authorities ...string) bool {
	if p == nil {
		return false
	}
	if len(p.Authorities) <= 0 || len(authorities) <= 0 {
		return false
	}
	for _, a := range authorities {
		if p.HasAuthority(a) {
			return true
		}
	}
	return false
}
