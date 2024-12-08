package sdsql

import (
	"context"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregator(t *testing.T) {
	is := assert.New(t)
	data := testNewData()
	aggr := &Aggregator[*testUser, string, *testPost]{
		Collect: func(_ context.Context, user *testUser) []string { return []string{user.Id} },
		Fetch: func(_ context.Context, uids []string) (map[string]*testPost, error) {
			return lo.SliceToMap(data.Posts, func(post *testPost) (string, *testPost) {
				return post.Id, post
			}), nil
		},
		CompleteInplace: func(_ context.Context, user *testUser, posts map[string]*testPost) error {
			postsMap := lo.PickBy(posts, func(postId string, post *testPost) bool {
				return post.Uid == user.Id
			})
			user.Posts = lo.Values(postsMap)
			return nil
		},
	}
	users, err := ProcRows(context.Background(), []*testUser{
		{Id: "uid1"},
		{Id: "uid2"},
	}, aggr)
	is.NoError(err)
	is.Len(users, 2)

	// user1
	is.Len(users[0].Posts, 2)
	is.Equal("uid1", users[0].Id)
	is.Equal("postId1", users[0].Posts[0].Id)
	is.Equal("postId2", users[0].Posts[1].Id)

	// user2
	is.Len(users[1].Posts, 1)
	is.Equal("uid2", users[1].Id)
	is.Equal("postId3", users[1].Posts[0].Id)
}
