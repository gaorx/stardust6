package sdsql

import (
	"context"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompleter(t *testing.T) {
	is := assert.New(t)
	data := testNewData()

	// Completer
	completer1 := Completer[*testPost](func(_ context.Context, post *testPost) (*testPost, error) {
		user, ok := lo.Find(data.Users, func(user *testUser) bool { return user.Id == post.Uid })
		if !ok {
			return nil, sderr.Newf("not found user")
		}
		post.User = user
		return post, nil
	})
	post1, err := ProcRow(context.Background(), &testPost{
		Id:  "postId1",
		Uid: "uid1",
	}, completer1)
	is.NoError(err)
	is.NotNil(post1)
	is.NotNil(post1.User)
	is.Equal("uid1", post1.User.Id)
	post1, err = ProcRow(context.Background(), &testPost{
		Id:  "postId1",
		Uid: "not_exists_uid",
	}, completer1)
	is.Error(err)
	is.Nil(post1)

	// InplaceCompleter
	completer2 := InplaceCompleter[*testPost](func(ctx context.Context, post *testPost) error {
		user, ok := lo.Find(data.Users, func(user *testUser) bool { return user.Id == post.Uid })
		if !ok {
			return sderr.Newf("not found user")
		}
		post.User = user
		return nil
	})
	post1, err = ProcRow(context.Background(), &testPost{
		Id:  "postId1",
		Uid: "uid1",
	}, completer2)
	is.NoError(err)
	is.NotNil(post1)
	is.NotNil(post1.User)
	is.Equal("uid1", post1.User.Id)
	post1, err = ProcRow(context.Background(), &testPost{
		Id:  "postId1",
		Uid: "not_exists_uid",
	}, completer2)
	is.Error(err)
	is.Nil(post1)
}
