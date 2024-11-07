package sdfile

import (
	"io"
	"os"

	"github.com/gaorx/stardust6/sderr"
)

// ReadBytes 从一个文件中读取所有数据到一个字节数组
func ReadBytes(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, sderr.With("file", filename).Wrapf(err, "read file error")
	}
	return data, nil
}

// ReadBytesDef 从一个文件中读取所有数据到一个字节数组，如果出错则返回默认值
func ReadBytesDef(filename string, def []byte) []byte {
	data, err := ReadBytes(filename)
	if err != nil {
		return def
	}
	return data
}

// WriteBytes 将一个字节数组写入到一个文件，可以设定文件权限
func WriteBytes(filename string, data []byte, perm os.FileMode) error {
	err := os.WriteFile(filename, data, perm)
	if err != nil {
		return sderr.With("file", filename).Wrapf(err, "write file error")
	}
	return nil
}

// AppendBytes 将一个字节数组追加到一个文件，可以设定文件权限
func AppendBytes(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return sderr.With("file", filename).Wrapf(err, "open append error")
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	_ = f.Close()
	return sderr.Wrapf(err, "write for append error")
}

// ReadText 从一个文件中读取所有数据到一个字符串
func ReadText(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", sderr.With("file", filename).Wrapf(err, "read text error")
	}
	return string(data), nil
}

// ReadTextDef 从一个文件中读取所有数据到一个字符串，如果出错则返回默认值
func ReadTextDef(filename, def string) string {
	data, err := ReadText(filename)
	if err != nil {
		return def
	}
	return data
}

// WriteText 将一个字符串写入到一个文件，可以设定文件权限
func WriteText(filename string, text string, perm os.FileMode) error {
	return WriteBytes(filename, []byte(text), perm)
}

// AppendText 将一个字符串追加到一个文件，可以设定文件权限
func AppendText(filename string, text string, perm os.FileMode) error {
	return AppendBytes(filename, []byte(text), perm)
}
