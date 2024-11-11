package sdbytes

// Presentable 便于战士用的的字节切片，可转换成hex或者base64
type Presentable []byte

// Summarize 返回字节切片的摘要
func (p Presentable) String() string {
	return Summarize(p)
}

// Hex 返回字节切片的十六进制字符串，upper为true是大写
func (p Presentable) Hex(upper bool) string {
	return ToHex(p, upper)
}

// HexL 返回字节切片的小写十六进制字符串
func (p Presentable) HexL() string {
	return ToHex(p, false)
}

// HexU 返回字节切片的大写十六进制字符串
func (p Presentable) HexU() string {
	return ToHex(p, true)
}

// Base64Std 返回字节切片的标准Base64字符串
func (p Presentable) Base64Std() string {
	return ToBase64Std(p)
}

// Base64Url 返回字节切片的URL Base64字符串
func (p Presentable) Base64Url() string {
	return ToBase64Url(p)
}
