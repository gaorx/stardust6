package sdcache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnyToJsonBytes(t *testing.T) {
	is := assert.New(t)

	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 测试结构指针类型
	codec1 := AnyToJsonBytes[*user]()
	data0 := &user{
		Name: "Tom",
		Age:  18,
	}
	encoded, err := codec1.Encode(data0)
	is.NoError(err)
	data1, err := codec1.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)

	// 测试null
	encoded, err = codec1.Encode(nil)
	is.NoError(err)
	is.Equal("null", string(encoded))
	data1, err = codec1.Decode(encoded)
	is.NoError(err)
	is.Equal((*user)(nil), data1)

	// 测试结构体
	codec2 := AnyToJsonBytes[user]()
	encoded, err = codec2.Encode(*data0)
	is.NoError(err)
	data1, err = codec1.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}

func TestAnyToJsonStr(t *testing.T) {
	is := assert.New(t)

	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 测试结构指针类型
	codec1 := AnyToJsonStr[*user]()
	data0 := &user{
		Name: "Tom",
		Age:  18,
	}
	encoded, err := codec1.Encode(data0)
	is.NoError(err)
	data1, err := codec1.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)

	// 测试null
	encoded, err = codec1.Encode(nil)
	is.NoError(err)
	is.Equal("null", string(encoded))
	data1, err = codec1.Decode(encoded)
	is.NoError(err)
	is.Equal((*user)(nil), data1)

	// 测试结构体
	codec2 := AnyToJsonStr[user]()
	encoded, err = codec2.Encode(*data0)
	is.NoError(err)
	data1, err = codec1.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}
