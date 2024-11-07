package sdfile

import (
	"github.com/samber/lo"
	"os"

	"github.com/gaorx/stardust6/sderr"
)

// UseTempFile 创建一个临时文件，可以对其进行操作，操作完成后会自动删除此文件
func UseTempFile(dir, pattern string, action func(*os.File)) error {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return sderr.With("dir", dir).
			With("pattern", pattern).
			Wrapf(err, "create temp file error")
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	action(f)
	return nil
}

// UseTempDir 创建一个临时目录，可以在其中进行操作，操作完成后会自动删除此目录
func UseTempDir(dir, pattern string, action func(string)) error {
	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return sderr.With("dir", dir).
			With("pattern", pattern).
			Wrapf(err, "create temp dir error")
	}
	defer func() {
		_ = os.RemoveAll(name)
	}()
	action(name)
	return nil
}

// UseTempFileFor 创建一个临时文件，可以对其进行操作，操作完成后会自动删除此文件，之后可以返回一个值
func UseTempFileFor[R any](dir, pattern string, action func(*os.File) (R, error)) (R, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return lo.Empty[R](), sderr.With("dir", dir).
			With("pattern", pattern).
			Wrapf(err, "create temp file error")
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	if r, err := action(f); err != nil {
		return lo.Empty[R](), sderr.With("dir", dir).
			With("pattern", pattern).
			Wrapf(err, "call file action error")
	} else {
		return r, nil
	}
}

// UseTempDirFor 创建一个临时目录，可以在其中进行操作，操作完成后会自动删除此目录，之后可以返回一个值
func UseTempDirFor[R any](dir, pattern string, action func(string) (R, error)) (R, error) {
	var empty R
	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return empty, sderr.With("dir", dir).
			With("pattern", pattern).
			Wrapf(err, "create temp dir error")
	}
	defer func() {
		_ = os.RemoveAll(name)
	}()

	if r, err := action(name); err != nil {
		return empty, sderr.With("dir", dir).
			With("pattern", pattern).
			Wrapf(err, "call dir action error")
	} else {
		return r, nil
	}
}
