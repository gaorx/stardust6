package sdauthn

import (
	"slices"
	"time"
)

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

func (p *Principal) PrincipalIdForDomain(domain string) PrincipalId {
	if p == nil {
		return PrincipalId{}
	}
	return PrincipalId{Domain: domain, Type: PrincipalUid, ID: p.ID}
}

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

func (p *Principal) HasAuthority(authority string) bool {
	if p == nil {
		return false
	}
	return slices.Contains(p.Authorities, authority)
}

func (p *Principal) HasAllAuthorities(authorities ...string) bool {
	if p == nil {
		return false
	}
	for _, a := range authorities {
		if !p.HasAuthority(a) {
			return false
		}
	}
	return true
}

func (p *Principal) HasAnyAuthority(authorities ...string) bool {
	if p == nil {
		return false
	}
	for _, a := range authorities {
		if p.HasAuthority(a) {
			return true
		}
	}
	return false
}
