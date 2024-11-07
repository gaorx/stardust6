package sdtextenc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChinese(t *testing.T) {
	is := assert.New(t)
	const s = "你好，世界！"
	is.Equal(s, GB2312.MustDecode(GB2312.MustEncode(s)))
	is.Equal(s, GBK.MustDecode(GBK.MustEncode(s)))
	is.Equal(s, GB18030.MustDecode(GB18030.MustEncode(s)))
}
