package sduuid

import (
	"github.com/gaorx/stardust6/sdbytes"
	"github.com/gofrs/uuid"
)

// UUID UUID类型
type UUID = uuid.UUID

// Encode 将UUID编码为字节数组
func Encode(id UUID) sdbytes.Presentable {
	return id.Bytes()
}

// NewV1 生成V1版本的UUID
func NewV1() sdbytes.Presentable {
	v, err := uuid.NewV1()
	if err != nil {
		return nil
	}
	return v.Bytes()
}

// NewV4 生成V4版本的UUID
func NewV4() sdbytes.Presentable {
	v, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	return v.Bytes()
}
