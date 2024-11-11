package sdsha512

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"github.com/gaorx/stardust6/sdbytes"
)

// Hash 计算SHA512哈希
func Hash(data []byte) sdbytes.Presentable {
	sum := sha512.Sum512(data)
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

// HashHmac 计算HMAC哈希
func HashHmac(key []byte, data []byte) sdbytes.Presentable {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

// VerifyHmac 验证data中的数据的签名是否与expected相等
func VerifyHmac(expected []byte, key []byte, data []byte) bool {
	actual := HashHmac(key, data)
	return bytes.Equal(actual, expected)
}

// VerifyHmacHex 验证data中的数据的签名是否与expected的十六进制字符串相等
func VerifyHmacHex(expected string, key []byte, data []byte) bool {
	hash, err := sdbytes.FromHex(expected)
	if err != nil {
		return false
	}
	return VerifyHmac(hash, key, data)
}
