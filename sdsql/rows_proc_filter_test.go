package sdsql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	is := assert.New(t)

	data := testNewData()
	filtered, err := ProcRows(
		context.Background(),
		data.Posts,
		Filter[*testPost](func(_ context.Context, post *testPost) bool {
			return post.Uid == "uid1"
		}),
	)
	is.NoError(err)
	is.Len(filtered, 2)
	is.Equal("uid1", filtered[0].Uid)
	is.Equal("uid1", filtered[1].Uid)

	filtered, err = ProcRows(
		context.Background(),
		data.Posts,
		Filter[*testPost](func(ctx context.Context, post *testPost) bool {
			return post.Uid == "not_exists"
		}),
	)
	is.NoError(err)
	is.Len(filtered, 0)
}
