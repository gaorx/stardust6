package sdauthn

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	is := assert.New(t)

	loader := SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321"},
	})

	// 通过UID查询
	p, err := Load(context.Background(), NewUserId("uid2"), loader)
	is.NoError(err)
	is.Equal("uid2", p.ID)
	_, err = Load(context.Background(), NewUserId("uid3"), loader)
	is.ErrorIs(err, ErrPrincipalNotFound)

	p, err = Load(context.Background(), NewUsernameAndPassword("user1", ""), loader)
	is.NoError(err)
	is.Equal("uid1", p.ID)
	_, err = Load(context.Background(), NewUsernameAndPassword("user3", ""), loader)
	is.ErrorIs(err, ErrPrincipalNotFound)
	_, err = Load(context.Background(), NewUsernameAndPassword("uid2", ""), loader)
	is.ErrorIs(err, ErrPrincipalNotFound)
}

func TestAuthenticate(t *testing.T) {
	is := assert.New(t)

	loader := SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321"},
	})

	// username和password正确的情况
	p, err := Authenticate(
		context.Background(),
		NewUsernameAndPassword("user2", "654321"),
		loader,
		NoExpiration(),
		time.Now(),
	)
	is.NoError(err)
	is.Equal("uid2", p.ID)

	// username正确，但是password错误的情况
	p, err = Authenticate(
		context.Background(),
		NewUsernameAndPassword("user2", "invalid_password"),
		loader,
		NoExpiration(),
		time.Now(),
	)
	is.ErrorIs(err, ErrCredentialInvalid)
	is.Equal("uid2", p.ID)

	// username不正确
	p, err = Authenticate(
		context.Background(),
		NewUsernameAndPassword("invalid_username", "654321"),
		loader,
		NoExpiration(),
		time.Now(),
	)
	is.ErrorIs(err, ErrPrincipalNotFound)
	is.Nil(p)

	// username和password正确，但是用户被disable
	loader = SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321", Disabled: true},
	})
	p, err = Authenticate(
		context.Background(),
		NewUsernameAndPassword("user2", "654321"),
		loader,
		NoExpiration(),
		time.Now(),
	)
	is.ErrorIs(err, ErrPrincipalDisabled)
	is.Equal("uid2", p.ID)

	// username和password正确，但是用户被lock
	loader = SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321", Locked: true},
	})
	p, err = Authenticate(
		context.Background(),
		NewUsernameAndPassword("user2", "654321"),
		loader,
		NoExpiration(),
		time.Now(),
	)
	is.ErrorIs(err, ErrPrincipalLocked)
	is.Equal("uid2", p.ID)

	// username和password正确，但是用户超期了
	loader = SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321", Expiry: time.Now().Add(-1 * time.Second)},
	})
	p, err = Authenticate(
		context.Background(),
		NewUsernameAndPassword("user2", "654321"),
		loader,
		NoExpiration(),
		time.Now(),
	)
	is.ErrorIs(err, ErrPrincipalExpired)
	is.Equal("uid2", p.ID)

	// 使用user token登陆
	loader = SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321"},
	})
	p, err = Authenticate(
		context.Background(),
		NewUserToken("uid2", time.Now().Add(-2*time.Second)),
		loader,
		ExpireIn(3*time.Second),
		time.Now(),
	)
	is.NoError(err)
	is.Equal("uid2", p.ID)

	// 使用user token登陆，但是token已经过期
	loader = SimpleUsers([]*SimpleUser{
		{ID: "uid1", Username: "user1", Password: "123456"},
		{ID: "uid2", Username: "user2", Password: "654321"},
	})
	p, err = Authenticate(
		context.Background(),
		NewUserToken("uid2", time.Now().Add(-2*time.Second)),
		loader,
		ExpireIn(1*time.Second),
		time.Now(),
	)
	is.ErrorIs(err, ErrCredentialExpired)
	is.Equal("uid2", p.ID)
}
