package sdsimpleapi

import (
	"github.com/gaorx/stardust6/sderr"
	"maps"
)

const (
	CodeOK      = "OK"
	CodeUnknown = "UNKNOWN"
)

type ResultInterface interface {
	SimpleAPIResultCode() string
	SimpleAPIResultData() any
	SimpleAPIResultMessage() string
	SimpleAPIResultMeta() map[string]any
}

type Result struct {
	Code    string         `json:"code"`
	Data    any            `json:"data,omitempty"`
	Message string         `json:"message,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

func NewResult(code string, data any, message string, meta map[string]any) *Result {
	return &Result{Code: code, Data: data, Message: message, Meta: meta}
}

func OK(data any) *Result {
	return NewResult(CodeOK, data, "", nil)
}

func Err(code string, err any) *Result {
	return NewResult(code, nil, sderr.PublicMsg(sderr.Ensure(err)), nil)
}

func (r *Result) Error() string {
	return r.Message
}

func (r *Result) SetCode(code string) *Result {
	r.Code = code
	return r
}

func (r *Result) SetData(data any) *Result {
	r.Data = data
	return r
}

func (r *Result) SetMessage(message string) *Result {
	r.Message = message
	return r
}

func (r *Result) SetMetas(meta map[string]any) *Result {
	r.Meta = maps.Clone(meta)
	return r
}

func (r *Result) WithMetas(meta map[string]any) *Result {
	if len(meta) <= 0 {
		return r
	}
	r.ensureMeta()
	maps.Copy(r.Meta, meta)
	return r
}

func (r *Result) WithMeta(key string, value any) *Result {
	r.ensureMeta()
	r.Meta[key] = value
	return r
}

func (r *Result) SimpleAPIResultCode() string {
	return r.Code
}

func (r *Result) SimpleAPIResultData() any {
	return r.Data
}

func (r *Result) SimpleAPIResultMessage() string {
	return r.Message
}

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
