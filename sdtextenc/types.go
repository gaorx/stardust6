package sdtextenc

// Encoding 编解码接口
type Encoding interface {
	// Encode 将UTF8字符串编码为字节数组
	Encode(s string) ([]byte, error)
	// Decode 将字节数组解码为UTF8字符串
	Decode(encoded []byte) (string, error)
	// MustEncode 将UTF8字符串编码为字节数组，如果失败则panic
	MustEncode(s string) []byte
	// MustDecode 将字节数组解码为UTF8字符串，如果失败则panic
	MustDecode(encoded []byte) string
	// DecodeDef 将字节数组解码为UTF8字符串，如果失败则返回默认值
	DecodeDef(encoded []byte, def string) string
}

type core interface {
	Encode(s string) ([]byte, error)
	Decode(encoded []byte) (string, error)
}

type coreExtension struct{ core }

func newEncoding(e core) Encoding {
	return coreExtension{e}
}

func (e coreExtension) MustEncode(s string) []byte {
	b, err := e.core.Encode(s)
	if err != nil {
		panic(err)
	}
	return b
}

func (e coreExtension) MustDecode(encoded []byte) string {
	s, err := e.core.Decode(encoded)
	if err != nil {
		panic(err)
	}
	return s
}

func (e coreExtension) DecodeDef(encoded []byte, def string) string {
	s, err := e.core.Decode(encoded)
	if err != nil {
		return def
	}
	return s
}
