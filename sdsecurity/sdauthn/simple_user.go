package sdauthn

import (
	"context"
	"slices"
	"time"
)

type SimpleUser struct {
	ID          string    `json:"id"`
	Authorities []string  `json:"roles,omitempty"`
	Password    string    `json:"password,omitempty"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	AvatarUrl   string    `json:"avatarUrl,omitempty"`
	Nickname    string    `json:"nickname,omitempty"`
	Locked      bool      `json:"locked,omitempty"`
	Expiry      time.Time `json:"expiry,omitempty"`
}

func SimpleUsers(users []*SimpleUser) Loader {
	m := make(map[string]*SimpleUser)
	for _, u := range users {
		if u == nil {
			continue
		}
		if u.ID == "" && u.Username == "" {
			continue
		}
		u1 := u.clone()
		if u1.ID == "" {
			u1.ID = u1.Username
		}
		if u1.Username == "" {
			u1.Username = u1.ID
		}
		m[u1.Username] = u1
	}
	loadPrincipal := func(_ context.Context, pid PrincipalId) (*Principal, error) {
		u, ok := m[pid.ID]
		if !ok {
			return nil, ErrPrincipalNotFound
		}
		p := &Principal{
			ID:          u.ID,
			Authorities: slices.Clone(u.Authorities),
			Password:    u.Password,
			Username:    u.Username,
			Email:       u.Email,
			Phone:       u.Phone,
			AvatarUrl:   u.AvatarUrl,
			Nickname:    u.Nickname,
			Disabled:    false,
			Expiry:      u.Expiry,
			Locked:      u.Locked,
		}
		return p, nil
	}
	return LoaderFunc(loadPrincipal)
}

func (u *SimpleUser) clone() *SimpleUser {
	if u == nil {
		return nil
	}
	c := *u
	c.Authorities = slices.Clone(u.Authorities)
	return &c
}
