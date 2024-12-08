package sdcache

// Codec 编解码器，用于将某个类型的值编解码到string或者[]byte
type Codec[T any, DST string | []byte] interface {
	// Encode 将v编码到DST，如错误返回错误
	Encode(v T) (DST, error)
	// Decode 将encoded解码到T，如出错返回错误
	Decode(encoded DST) (T, error)
}

// FuncCodec 使用内嵌函数作为实现的编解码器
type FuncCodec[T any, DST string | []byte] struct {
	EncodeFunc func(v T) (DST, error)
	DecodeFunc func(encoded DST) (T, error)
}

var _ Codec[string, string] = FuncCodec[string, string]{}

func (f FuncCodec[T, DST]) Encode(v T) (DST, error) {
	return f.EncodeFunc(v)
}

func (f FuncCodec[T, DST]) Decode(encoded DST) (T, error) {
	return f.DecodeFunc(encoded)
}
