package sdauthn

import (
	"context"
	"github.com/samber/lo"
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
	Disabled    bool      `json:"disabled,omitempty"`
	Locked      bool      `json:"locked,omitempty"`
	Expiry      time.Time `json:"expiry,omitempty"`
}

func SimpleUsers(users []*SimpleUser) Loader {
	users = lo.Filter(users, func(u *SimpleUser, _ int) bool {
		return u != nil
	})
	users = lo.Map(users, func(u *SimpleUser, _ int) *SimpleUser {
		return u.clone()
	})
	lo.ForEach(users, func(u *SimpleUser, _ int) {
		if u.ID != "" {
			if u.Username == "" {
				u.Username = u.ID
			}
		} else {
			if u.Username != "" {
				u.ID = u.Username
			}
		}
	})

	userToPrincipal := func(u *SimpleUser) *Principal {
		return &Principal{
			ID:          u.ID,
			Authorities: slices.Clone(u.Authorities),
			Password:    u.Password,
			Username:    u.Username,
			Email:       u.Email,
			Phone:       u.Phone,
			AvatarUrl:   u.AvatarUrl,
			Nickname:    u.Nickname,
			Disabled:    u.Disabled,
			Expiry:      u.Expiry,
			Locked:      u.Locked,
		}
	}

	loadPrincipal := func(_ context.Context, pid PrincipalId) (*Principal, error) {
		for _, u := range users {
			switch pid.Type {
			case PrincipalUid:
				if u.ID != "" && u.ID == pid.ID {
					return userToPrincipal(u), nil
				}
			case PrincipalUsername:
				if u.Username != "" && u.Username == pid.ID {
					return userToPrincipal(u), nil
				}
			case PrincipalEmail:
				if u.Email != "" && u.Email == pid.ID {
					return userToPrincipal(u), nil
				}
			case PrincipalPhone:
				if u.Phone != "" && u.Phone == pid.ID {
					return userToPrincipal(u), nil
				}
			}
		}
		return nil, ErrPrincipalNotFound
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
