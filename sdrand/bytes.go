package sdrand

import (
	"math/rand"
)

// Bytes 返回长度为 n 的随机字节切片
func Bytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(16))
	}
	return b
}
