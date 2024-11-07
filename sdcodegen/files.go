package sdcodegen

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"io"
	"io/fs"
	"os"
	"strings"
)

// Files 描述一组文件，它可以从一个目录中读取文件，也可以写入到一个目录中
type Files []*File

// ReadDir 从指定目录中读取文件列表
func ReadDir(dirname string, nameFilters ...StringFilter) (Files, error) {
	absDirname, ok := toAbs(dirname)
	if !ok {
		return nil, sderr.With("dir", dirname).Newf("illegall dirname")
	}
	return readDirFS(os.DirFS(absDirname), nameFilters)
}

// ReadDirFS 从指定的fs中读取文件列表
func ReadDirFS(fsys fs.FS, nameFilters ...StringFilter) (Files, error) {
	if fsys == nil {
		return nil, sderr.Newf("fsys is nil")
	}
	return readDirFS(fsys, nameFilters)
}

func readDirFS(fsys fs.FS, nameFilters []StringFilter) (Files, error) {
	dirStat, err := fs.Stat(fsys, ".")
	if err != nil {
		if IsNotExistErr(err) {
			return nil, nil
		}
		return nil, sderr.Wrapf(err, "stat dir failed")
	}
	if !dirStat.IsDir() {
		return nil, sderr.Newf("not a dir")
	}

	nameFilters = lo.Filter(nameFilters, func(sf StringFilter, _ int) bool {
		return sf != nil
	})

	isFiltered := func(name string) bool {
		if len(nameFilters) <= 0 {
			return true
		}
		for _, filter := range nameFilters {
			if filter != nil && filter(name) {
				return true
			}
		}
		return false
	}

	var files []*File
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err0 error) error {
		if err0 != nil {
			return err0
		}
		if d.IsDir() {
			return nil
		}
		name := path
		if !isFiltered(name) {
			return nil
		}
		f, err := readFileFS(fsys, name)
		if err != nil {
			return sderr.Wrap(err)
		}
		files = append(files, f)
		return nil
	})
	if err != nil {
		return nil, sderr.Wrapf(err, "walk dir failed")
	}
	return files, nil
}

// Write 将这组文件写入到指定目录中
func (files Files) Write(dirname string, logger Logger) error {
	absDirname, ok := toAbs(dirname)
	if !ok {
		return sderr.With("dir", dirname).Newf("illegall dirname")
	}
	for _, f := range files {
		err := f.Write(absDirname, logger)
		if err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

// Len 返回文件数量
func (files Files) Len() int {
	return len(files)
}

// String 返回描述
func (files Files) String() string {
	return fmt.Sprintf("%d files", len(files))
}

// StringText 返回详细描述，保罗每个文件的内容
func (files Files) StringText() string {
	sep := strings.Repeat("=", 60)
	var b strings.Builder
	for i := 0; i < len(files); i++ {
		b.WriteString(sep + "\n")
		f := files[i]
		s := f.StringText()
		b.WriteString(s)
		if !strings.HasSuffix(s, "\n") {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// Get 通过name查找这个列表中的文件，如果不存在则返回<nil>
func (files Files) Get(name string) *File {
	if len(files) <= 0 {
		return nil
	}
	for _, f := range files {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// StringList 返回文件列表的描述，每行描述一个文件
func (files Files) StringList(excludeDiscarded bool) []string {
	var listFiles []*File
	for _, f := range files {
		if f.IsZero() {
			continue
		}
		if excludeDiscarded && f.Discarded {
			continue
		}
		listFiles = append(listFiles, f)
	}
	return lo.Map(listFiles, func(f *File, _ int) string {
		return f.String()
	})
}

// LL 获取像ls -l 类似格式的输出的文件列表
func (files Files) LL(w io.Writer, excludeDiscarded bool) {
	var b strings.Builder
	_, _ = fmt.Fprintln(&b, fmt.Sprintf("total %d", len(files)))
	for _, s := range files.StringList(excludeDiscarded) {
		_, _ = fmt.Fprintln(&b, s)
	}
	_, _ = fmt.Fprintln(w, b.String())
}
