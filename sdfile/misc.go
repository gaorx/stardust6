package sdfile

import (
	"os"
	"path/filepath"
)

// IsDir 判断一个文件是否是目录
func IsDir(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

// BinDir 获取当前可执行文件所在目录
func BinDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exe)
}

// AbsByBin 以当前文件的所在目录作为起点，获取一个相对文件的绝对路径
func AbsByBin(filename string) string {
	if filename == "" {
		return BinDir()
	}
	if filepath.IsAbs(filename) {
		return filename
	}
	r, err := filepath.Abs(filepath.Join(BinDir(), filename))
	if err != nil {
		return ""
	}
	return r
}

// Exists 判断一个文件是否存在
func Exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// FirstExists 返回第一个存在的文件
func FirstExists(filenames ...string) string {
	if len(filenames) == 0 {
		return ""
	}
	for _, filename := range filenames {
		if Exists(filename) {
			return filename
		}
	}
	return ""
}

// FirstExistsAbs 返回第一个存在的文件的绝对路径
func FirstExistsAbs(filenames ...string) string {
	fn := FirstExists(filenames...)
	if fn == "" {
		return ""
	}
	abs, err := filepath.Abs(fn)
	if err != nil {
		return ""
	}
	return abs
}
