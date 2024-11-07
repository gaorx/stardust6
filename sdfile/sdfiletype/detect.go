package sdfiletype

import (
	"github.com/h2non/filetype"
	"io"
)

// DetectBytes 检测一个字节数组的内容可能是哪种文件
func DetectBytes(data []byte) *Type {
	t, err := filetype.Match(data)
	if err != nil {
		return nil
	}
	return &Type{t}
}

// DetectReader 检测一个Reader的内容可能是哪种文件
func DetectReader(r io.Reader) *Type {
	t, err := filetype.MatchReader(r)
	if err != nil {
		return nil
	}
	return &Type{t}
}

// DetectFile 检测一个文件的内容可能是哪种文件
func DetectFile(path string) *Type {
	t, err := filetype.MatchFile(path)
	if err != nil {
		return nil
	}
	return &Type{t}
}
