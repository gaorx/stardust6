package sdcache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToString(t *testing.T) {
	is := assert.New(t)
	codec := StringToString()
	data0 := "hello"
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	is.Equal("hello", encoded)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}

func TestBytesToHex(t *testing.T) {
	is := assert.New(t)
	codec := BytesToHex()
	data0 := []byte{0x00, 0x01, 0x02, 0x03}
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	is.Equal("00010203", encoded)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}

func TestInt64ToString(t *testing.T) {
	is := assert.New(t)
	codec := Int64ToString()
	data0 := int64(-123)
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	is.Equal("-123", encoded)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}

func TestUint64ToString(t *testing.T) {
	is := assert.New(t)
	codec := Uint64ToString()
	data0 := uint64(123)
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	is.Equal("123", encoded)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}

func TestBytesToString(t *testing.T) {
	is := assert.New(t)
	codec := BytesToString()
	data0 := []byte("hello")
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	is.Equal("hello", encoded)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}

func TestStringToBytes(t *testing.T) {
	is := assert.New(t)
	codec := StringToBytes()
	data0 := "hello"
	encoded, err := codec.Encode(data0)
	is.NoError(err)
	is.Equal([]byte("hello"), encoded)
	data1, err := codec.Decode(encoded)
	is.NoError(err)
	is.Equal(data0, data1)
}
