package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjson"
	"maps"
)

const (
	CodeOK        = "OK"
	CodeUnknown   = "UNKNOWN"
	CodeDataError = "DATA_ERROR"
	CodeNotFound  = "NOT_FOUND"
)

// ResultInterface 结果接口，一个Handler返回此接口的实例是可以被当作Result处理的
type ResultInterface interface {
	SimpleAPIResultCode() string
	SimpleAPIResultData() any
	SimpleAPIResultMessage() string
	SimpleAPIResultMeta() map[string]any
}

// Result API结果，使用 OK 和 Err 函数创建
type Result struct {
	Code    string        `json:"code"`
	Data    any           `json:"data,omitempty"`
	Message string        `json:"message,omitempty"`
	Meta    sdjson.Object `json:"meta,omitempty"`
}

// ResultT API结果，用于测试作为json结构来解析的
type ResultT[T any] struct {
	Code    string        `json:"code"`
	Data    T             `json:"data,omitempty"`
	Message string        `json:"message,omitempty"`
	Meta    sdjson.Object `json:"meta,omitempty"`
}

// NewResult 创建一个结果
func NewResult(code string, data any, message string, meta map[string]any) *Result {
	return &Result{Code: code, Data: data, Message: message, Meta: meta}
}

// OK 创建一个成功结果
func OK(data any) *Result {
	return NewResult(CodeOK, data, "", nil)
}

// Err 创建一个错误结果
func Err(code string, err any) *Result {
	return NewResult(code, nil, sderr.PublicMsg(sderr.Ensure(err)), nil)
}

// Error 实现error接口
func (r *Result) Error() string {
	return r.Message
}

// SetCode 设置结果错误码(业务逻辑错误码，非HTTP状态码)
func (r *Result) SetCode(code string) *Result {
	r.Code = code
	return r
}

// SetData 设置结果数据
func (r *Result) SetData(data any) *Result {
	r.Data = data
	return r
}

// SetMessage 设置结果消息
func (r *Result) SetMessage(message string) *Result {
	r.Message = message
	return r
}

// SetMetas 设置结果元数据，通常分页数据在其中
func (r *Result) SetMetas(meta map[string]any) *Result {
	r.Meta = maps.Clone(meta)
	return r
}

// WithMetas 附加结果元数据，通常分页数据在其中
func (r *Result) WithMetas(meta map[string]any) *Result {
	if len(meta) <= 0 {
		return r
	}
	r.ensureMeta()
	maps.Copy(r.Meta, meta)
	return r
}

// WithMeta 附加结果元数据
func (r *Result) WithMeta(key string, value any) *Result {
	r.ensureMeta()
	r.Meta[key] = value
	return r
}

// SimpleAPIResultCode 实现 ResultInterface 接口
func (r *Result) SimpleAPIResultCode() string {
	return r.Code
}

// SimpleAPIResultData 实现 ResultInterface 接口
func (r *Result) SimpleAPIResultData() any {
	return r.Data
}

// SimpleAPIResultMessage 实现 ResultInterface 接口
func (r *Result) SimpleAPIResultMessage() string {
	return r.Message
}

// SimpleAPIResultMeta 实现 ResultInterface 接口
func (r *Result) SimpleAPIResultMeta() map[string]any {
	return r.Meta
}

func (r *Result) ensureMeta() {
	if r.Meta == nil {
		r.Meta = map[string]any{}
	}
}

func (r *Result) assignFrom(i ResultInterface) {
	if i == nil {
		return
	}
	r.Code = i.SimpleAPIResultCode()
	data0 := i.SimpleAPIResultData()
	if r.Data == nil && data0 != nil {
		r.Data = data0
	}
	r.Message = i.SimpleAPIResultMessage()
	r.Meta = i.SimpleAPIResultMeta()
}

func (r *Result) trimMeta() {
	if len(r.Meta) <= 0 {
		return
	}
	meta1 := map[string]any{}
	for k, v := range r.Meta {
		if k != "" && v != nil {
			meta1[k] = v
		}
	}
	r.Meta = meta1
}
