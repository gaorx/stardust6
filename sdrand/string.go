package sdrand

import (
	"github.com/samber/lo"
)

var (
	LowerCaseLettersCharset      = lo.LowerCaseLettersCharset
	UpperCaseLettersCharset      = lo.UpperCaseLettersCharset
	LettersCharset               = lo.LettersCharset
	NumbersCharset               = lo.NumbersCharset
	AlphanumericCharset          = lo.AlphanumericCharset
	LowerCaseAlphanumericCharset = append(LowerCaseLettersCharset, NumbersCharset...)
	UpperCaseAlphanumericCharset = append(UpperCaseLettersCharset, NumbersCharset...)
	SpecialCharset               = lo.SpecialCharset
	AllCharset                   = lo.AllCharset
)

// String 生成指定长度的随机字符串
func String(n int, set []rune) string {
	return lo.RandomString(n, set)
}
