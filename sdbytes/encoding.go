package sdbytes

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/gaorx/stardust6/sderr"
	"strings"
)

// ToHex 将一个字节数组转换为十六进制字符串
func ToHex(d []byte, upper bool) string {
	if len(d) <= 0 {
		return ""
	}
	s := hex.EncodeToString(d)
	if upper {
		s = strings.ToUpper(s)
	}
	return s
}

func ToHexL(d []byte) string {
	return ToHex(d, false)
}

func ToHexU(d []byte) string {
	return ToHex(d, true)
}

// FromHex 将一个十六进制字符串转换为字节数组
func FromHex(encoded string) ([]byte, error) {
	d, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return d, nil
}

// ToBase64Std 将一个字节数组转换为标准 Base64 字符串
func ToBase64Std(d []byte) string {
	if len(d) <= 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(d)
}

// FromBase64Std 将一个标准 Base64 字符串转换为字节数组
func FromBase64Std(encoded string) ([]byte, error) {
	d, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return d, nil
}

// ToBase64Url 将一个字节数组转换为 URL Base64 字符串
func ToBase64Url(d []byte) string {
	if len(d) <= 0 {
		return ""
	}
	return base64.URLEncoding.EncodeToString(d)
}

// FromBase64Url 将一个 URL Base64 字符串转换为字节数组
func FromBase64Url(encoded string) ([]byte, error) {
	d, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return d, nil
}
