package sdhttpfile

import (
	"io"
	"net/http"

	"github.com/gaorx/stardust6/sderr"
)

// ReadBytes 从http.FileSystem中读取文件内容
func ReadBytes(hfs http.FileSystem, name string) ([]byte, error) {
	if hfs == nil {
		return nil, sderr.Newf("nil hfs")
	}
	f, err := hfs.Open(name)
	if err != nil {
		return nil, sderr.With("name", name).Wrapf(err, "open error")
	}
	defer func() { _ = f.Close() }()
	r, err := io.ReadAll(f)
	if err != nil {
		return nil, sderr.With("name", name).Wrapf(err, "read error")
	}
	return r, nil
}

// ReadText 从http.FileSystem中读取文件的文本内容
func ReadText(hfs http.FileSystem, name string) (string, error) {
	b, err := ReadBytes(hfs, name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadTextDef 从http.FileSystem中读取文件的文本内容，如果失败则返回默认值
func ReadTextDef(hfs http.FileSystem, name, def string) string {
	s, err := ReadText(hfs, name)
	if err != nil {
		return def
	}
	return s
}
