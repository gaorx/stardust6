package sdcodegen

import (
	"io/fs"
)

// DirFS 从fs中读取所有文件，然后将文件添加到生成器中
func DirFS(fsys fs.FS) Module {
	if fsys == nil {
		return filesMod{}
	}
	files, err := ReadDirFS(fsys)
	if err != nil {
		return filesMod{}
	}
	return filesMod{files: files}
}

// Dir 从目录中读取所有文件，然后将文件添加到生成器中
func Dir(dirname string) Module {
	if dirname == "" {
		return filesMod{}
	}
	files, err := ReadDir(dirname)
	if err != nil {
		return filesMod{}
	}
	return filesMod{files: files}
}

type filesMod struct {
	files []*File
}

func (m filesMod) Apply(i Interface) {
	i.AddFiles(m.files)
}
