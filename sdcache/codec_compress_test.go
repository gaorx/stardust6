package sdcache

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcompress/sdlz4"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestLz4(t *testing.T) {
	is := assert.New(t)

	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 将一个JSON编解码器转换为LZ4编解码器，将JSON编码后的字节数组压缩
	codec := Lz4(
		AnyToJsonBytes[[]user](),
		sdlz4.Lz4Fast,
	)

	var data0 []user
	for i := 0; i < 100; i++ {
		data0 = append(data0, user{
			Name: fmt.Sprintf("Tom% d", i),
			Age:  rand.Intn(100),
		})
	}
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)

	// 测试null
	encoded, err = codec.Encode(nil)
	is.NoError(err)
	data1, err = codec.Decode(encoded)
	is.NoError(err)
	is.Equal(([]user)(nil), data1)
}
