package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sdwebapp"
)

// NewTestRequest 创建测试请求
func NewTestRequest(path string, req any) *sdwebapp.TestRequest {
	return sdwebapp.NewTestRequest("POST", path).SetBodyJson(req)
}
