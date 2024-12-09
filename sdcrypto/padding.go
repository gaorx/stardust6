package sdcrypto

import (
	"bytes"
	"github.com/gaorx/stardust6/sderr"
)

type (
	// Padding 填充函数
	Padding func(data []byte, blockSize int) ([]byte, error)
	// Unpadding 反填充函数
	Unpadding func(data []byte, blockSize int) ([]byte, error)
)

// Pkcs5 PKCS5填充
func Pkcs5(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, sderr.With("block_size", blockSize).Newf("illegal block size")
	}
	padding := blockSize - len(data)%blockSize
	padded := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(bytes.Clone(data), padded...), nil
}

// UnPkcs5 PKCS5反填充
func UnPkcs5(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, sderr.With("block_size", blockSize).Newf("illegal block size")
	}
	if len(data) < blockSize {
		return nil, sderr.Newf("data too short")
	}
	lastByte := int(data[len(data)-1])
	if lastByte <= 0 || lastByte > blockSize {
		return nil, sderr.Newf("illegal padding size")
	}
	return bytes.Clone(data[:len(data)-lastByte]), nil
}

// Zeros 补充0填充
func Zeros(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, sderr.With("block_size", blockSize).Newf("illegal block size")
	}
	padding := blockSize - len(data)%blockSize
	padded := bytes.Repeat([]byte{0}, padding)
	result := make([]byte, 0, len(data)+padding)
	result = append(result, data...)
	return append(result, padded...), nil
}

// UnZeros 去除0填充
func UnZeros(data []byte, blockSize int) ([]byte, error) {
	return bytes.TrimRightFunc(data,
		func(r rune) bool {
			return r == 0
		}), nil
}
