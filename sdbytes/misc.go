package sdbytes

import (
	"fmt"
)

// Summarize 返回一个字节数组的摘要展示，用于观察字节数组的前几个字节和后几个字节的内容
func Summarize(d []byte) string {
	const hextable = "0123456789abcdef"
	hb := func(b byte) string {
		buff := [2]byte{}
		buff[0] = hextable[b>>4]
		buff[1] = hextable[b&0x0f]
		return string(buff[:])
	}

	if d == nil {
		return "bytes(0) nil"
	}
	l := len(d)
	switch len(d) {
	case 0:
		return "bytes(0) []"
	case 1:
		return fmt.Sprintf("bytes(1) [%s]", hb(d[0]))
	case 2:
		return fmt.Sprintf("bytes(2) [%s %s]", hb(d[0]), hb(d[1]))
	case 3:
		return fmt.Sprintf("bytes(3) [%s %s %s]", hb(d[0]), hb(d[1]), hb(d[2]))
	case 4:
		return fmt.Sprintf("bytes(4) [%s %s %s %s]", hb(d[0]), hb(d[1]), hb(d[2]), hb(d[3]))
	case 5:
		return fmt.Sprintf("bytes(5) [%s %s %s %s %s]", hb(d[0]), hb(d[1]), hb(d[2]), hb(d[3]), hb(d[4]))
	case 6:
		return fmt.Sprintf("bytes(6) [%s %s %s %s %s %s]", hb(d[0]), hb(d[1]), hb(d[2]), hb(d[3]), hb(d[4]), hb(d[5]))
	default:
		return fmt.Sprintf("bytes(%d) [%s %s %s %s .. %s %s]", l, hb(d[0]), hb(d[1]), hb(d[2]), hb(d[3]), hb(d[l-2]), hb(d[l-1]))
	}
}
