package sdobjectstore

import (
	"github.com/gaorx/stardust6/sderr"
	"io/fs"
)

// Store 对象存储，扩展了Interface功能，加入了一些更方便使用的方法
type Store struct {
	Interface
}

// IsNil 其包含的是否为nil
func (s Store) IsNil() bool {
	return s.Interface == nil
}

// StoreFile 存储一个文件到对象存储中
func (s Store) StoreFile(filename, objectName string) (*Target, error) {
	return s.Store(File(filename, ""), objectName)
}

// StoreData 存储数据到对象存储中
func (s Store) StoreData(data []byte, objectName string) (*Target, error) {
	return s.Store(Bytes(data, ""), objectName)
}

// StoreFileFS 存储一个在FS中的文件到对象存储中
func (s Store) StoreFileFS(fsys fs.FS, fn string, objectName string) (*Target, error) {
	data, err := fs.ReadFile(fsys, fn)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return s.Store(Bytes(data, ""), objectName)
}

// Dir 记录一个目录创建一个对象存储，将数据存储在其中
func Dir(root string) Store {
	return Store{dir{root}}
}
