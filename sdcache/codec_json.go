package sdcache

import (
	"github.com/gaorx/stardust6/sdjson"
)

// AnyToJsonBytes 转换任意类型到JSON形式的字节数组
func AnyToJsonBytes[V any]() Codec[V, []byte] {
	return FuncCodec[V, []byte]{
		EncodeFunc: func(v V) ([]byte, error) {
			return sdjson.MarshalBytes(v)
		},
		DecodeFunc: func(encoded []byte) (V, error) {
			return sdjson.UnmarshalBytesT[V](encoded)
		},
	}
}

// AnyToJsonStr 转换任意类型到JSON形式的字符串
func AnyToJsonStr[V any]() Codec[V, string] {
	return FuncCodec[V, string]{
		EncodeFunc: func(v V) (string, error) {
			return sdjson.MarshalString(v)
		},
		DecodeFunc: func(encoded string) (V, error) {
			return sdjson.UnmarshalStringT[V](encoded)
		},
	}
}
