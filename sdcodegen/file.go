package sdcodegen

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File 描述一个要写入的文件，或者一个从实际文件中读取的文件
type File struct {
	Name      string      // 文件名称，只能是相对文件名，且不能为空
	Data      []byte      // 文件中的数据
	Mode      fs.FileMode // 文件的MODE
	ModTime   time.Time   // 这个文件的编辑时间，仅在读取的文件中存在；在写入时无意义
	Discarded bool        // 为true表示这个文件在写入时是被丢弃的，实际不产生任何写入行为
}

// NewFile 创建一个新的文件
func NewFile(name string, data []byte, mode fs.FileMode) (*File, error) {
	if !filepath.IsLocal(name) {
		return nil, sderr.With("file", name).Newf("file is not local")
	}
	return &File{
		Name: name,
		Data: data,
		Mode: mode,
	}, nil
}

// ReadFile 从指定目录中读取一个文件
func ReadFile(dirname, name string) (*File, error) {
	absDirname, ok := toAbs(dirname)
	if !ok {
		return nil, sderr.With("dir", dirname).Newf("illegall dirname")
	}
	if !filepath.IsLocal(name) {
		return nil, sderr.With("file", name).Newf("file is not local")
	}
	return readFileFS(os.DirFS(absDirname), name)
}

// ReadFileFS 从指定的fs中读取一个文件
func ReadFileFS(fsys fs.FS, name string) (*File, error) {
	if fsys == nil {
		return nil, sderr.Newf("fs is nil")
	}
	if !filepath.IsLocal(name) {
		return nil, sderr.With("file", name).Newf("file is not local")
	}
	return readFileFS(fsys, name)
}

func readFileFS(fsys fs.FS, name string) (*File, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, sderr.With("file", name).Wrapf(err, "open file failed")
	}
	defer func() { _ = f.Close() }()
	fi, err := f.Stat()
	if err != nil {
		return nil, sderr.With("file", name).Wrapf(err, "stat file failed")
	}
	if fi.IsDir() {
		return nil, sderr.With("file", name).Wrapf(err, "file is a directory")
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, sderr.With("file", name).Wrapf(err, "read file failed")
	}
	return &File{
		Name:    name,
		Data:    data,
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
	}, nil
}

// Write 将文件写入指定目录
func (f *File) Write(dirname string, logger Logger) error {
	absDirname, ok := toAbs(dirname)
	if !ok {
		return sderr.With("dir", dirname).Newf("illegall dirname")
	}

	if !filepath.IsLocal(f.Name) {
		return sderr.With("file", f.Name).Newf("file is not local")
	}
	if f.Discarded {
		if logger != nil {
			logger.LogDiscard(f.Name, f.Mode)
		}
		return nil
	}
	fullName := filepath.Join(absDirname, f.Name)
	fullDir := filepath.Dir(fullName)
	if os.MkdirAll(fullDir, 0755) != nil {
		return sderr.With("dir", fullDir).Newf("make dir failed")
	}
	err := os.WriteFile(fullName, f.Data, f.Mode)
	if err != nil {
		if logger != nil {
			logger.LogWrite(f.Name, f.Mode, err)
		}
		return sderr.With("file", fullName).Wrapf(err, "write file failed")
	}
	if logger != nil {
		logger.LogWrite(f.Name, f.Mode, nil)
	}
	return nil
}

// String 以字符串形式返回文件信息，但不包括文件内容
func (f *File) String() string {
	if f.IsZero() {
		return ""
	}
	modeStr := f.Mode.String()
	sizeStr := fmt.Sprintf("%9dB", len(f.Data))
	timeStr := "-------------------"
	if !f.ModTime.IsZero() {
		timeStr = f.ModTime.Format(time.DateTime)
	}
	name := f.Name
	if name != "" && f.Discarded {
		name = name + " [DISCARD]"
	}
	return fmt.Sprintf("%s  %s  %s  %s", modeStr, sizeStr, timeStr, name)
}

// StringText 以字符串形式返回文件信息，包括文件内容
func (f *File) StringText() string {
	if f.IsZero() {
		return ""
	}
	var b strings.Builder
	_, _ = fmt.Fprintln(&b, f.Name)
	_, _ = fmt.Fprintln(&b, strings.Repeat("-", 60))
	_, _ = fmt.Fprintln(&b, string(f.Data))
	return b.String()
}

// IsZero 判断文件是否为空
func (f *File) IsZero() bool {
	return f == nil || (f.Name == "" && len(f.Data) <= 0 && f.Mode == 0 && f.ModTime.IsZero() && f.Discarded == false)
}

// Text 以字符串形式返回文件内容
func (f *File) Text() string {
	return string(f.Data)
}

// SetName 设置文件名称
func (f *File) SetName(name string) *File {
	f.Name = name
	return f
}

// SetDiscarded 设置文件是否被丢弃
func (f *File) SetDiscarded(b bool) *File {
	f.Discarded = b
	return f
}

// SetMode 设置文件MODE
func (f *File) SetMode(mode fs.FileMode) *File {
	f.Mode = mode
	return f
}

// SetText 以文本的方式设置文件内容
func (f *File) SetText(text string) *File {
	f.Data = []byte(text)
	return f
}
