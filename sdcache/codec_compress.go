package sdcache

import (
	"github.com/gaorx/stardust6/sdcompress/sdlz4"
)

// Lz4 可以另外一个到[]byte的编解码器进行LZ4压缩
func Lz4[V any](codec Codec[V, []byte], level sdlz4.Level) Codec[V, []byte] {
	return FuncCodec[V, []byte]{
		EncodeFunc: func(v V) ([]byte, error) {
			data, err := codec.Encode(v)
			if err != nil {
				return nil, err
			}
			return sdlz4.Compress(data, level)
		},
		DecodeFunc: func(encoded []byte) (V, error) {
			uncompressed, err := sdlz4.Uncompress(encoded)
			if err != nil {
				var zero V
				return zero, err
			}
			return codec.Decode(uncompressed)
		},
	}
}
