package sdmd5

import (
	"bytes"
	"crypto/md5"
	"github.com/gaorx/stardust6/sdbytes"
)

// Hash 计算MD5哈希
func Hash(data []byte) sdbytes.Presentable {
	sum := md5.Sum(data)
	return sum[:]
}

// Verify 验证data中的数据的签名是否与expected相等
func Verify(expected []byte, data []byte) bool {
	sum := Hash(data)
	return bytes.Equal(expected, sum)
}

// VerifyHex 验证data中的数据的签名是否与expected的十六进制字符串相等
func VerifyHex(expected string, data []byte) bool {
	hash, err := sdbytes.FromHex(expected)
	if err != nil {
		return false
	}
	return Verify(hash, data)
}
