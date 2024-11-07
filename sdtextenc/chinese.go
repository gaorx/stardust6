package sdtextenc

import (
	"github.com/gaorx/stardust6/sderr"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	// GB2312 编解码
	GB2312 = newEncoding(textEncoding{simplifiedchinese.HZGB2312})

	// GBK 编解码
	GBK = newEncoding(textEncoding{simplifiedchinese.GBK})

	// GB18030 编解码
	GB18030 = newEncoding(textEncoding{simplifiedchinese.GB18030})
)

type textEncoding struct {
	encoding encoding.Encoding
}

func (e textEncoding) Encode(s string) ([]byte, error) {
	b, err := e.encoding.NewEncoder().Bytes([]byte(s))
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return b, nil
}

func (e textEncoding) Decode(encoded []byte) (string, error) {
	if encoded == nil {
		return "", sderr.Newf("encoded is nil")
	}
	b, err := e.encoding.NewDecoder().Bytes(encoded)
	if err != nil {
		return "", sderr.Wrap(err)
	}
	return string(b), nil
}
