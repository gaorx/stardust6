package sdobjectstore

// Interface 对象存储接口
type Interface interface {
	Store(src Source, objectName string) (*Target, error)
}
