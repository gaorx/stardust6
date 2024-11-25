package sdauthn

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPrincipal(t *testing.T) {
	is := assert.New(t)

	// PrincipalId
	p := &Principal{
		ID: "uid1",
	}
	is.Equal(
		NewPrincipalId(PrincipalUid, "domain1", "uid1"),
		p.PrincipalIdForDomain("domain1"),
	)

	// ValidateSelf
	p = &Principal{ID: ""}
	is.Equal(ErrPrincipalNotFound, p.ValidateSelf(time.Now()))
	p = &Principal{ID: "uid", Disabled: true}
	is.Equal(ErrPrincipalDisabled, p.ValidateSelf(time.Now()))
	p = &Principal{ID: "uid", Locked: true}
	is.Equal(ErrPrincipalLocked, p.ValidateSelf(time.Now()))
	p = &Principal{ID: "uid", Expiry: time.Now().Add(-1 * time.Second)}
	is.Equal(ErrPrincipalExpired, p.ValidateSelf(time.Now()))

	// HasAuthority
	p = &Principal{ID: "uid", Authorities: nil}
	is.False(p.HasAuthority("a"))
	p = &Principal{ID: "uid", Authorities: []string{"b", "a"}}
	is.True(p.HasAuthority("a"))

	// HasAllAuthorities
	p = &Principal{ID: "uid", Authorities: nil}
	is.False(p.HasAllAuthorities())
	p = &Principal{ID: "uid", Authorities: []string{"a", "b", "c"}}
	is.True(p.HasAllAuthorities("a", "c"))
	p = &Principal{ID: "uid", Authorities: []string{"a", "b", "c"}}
	is.False(p.HasAllAuthorities("a", "d"))

	// HasAnyAuthority
	p = &Principal{ID: "uid", Authorities: nil}
	is.False(p.HasAnyAuthority())
	p = &Principal{ID: "uid", Authorities: []string{"a", "b", "c"}}
	is.True(p.HasAnyAuthority("a", "d"))
	p = &Principal{ID: "uid", Authorities: []string{"a", "b", "c"}}
	is.False(p.HasAnyAuthority("d", "e"))
}
