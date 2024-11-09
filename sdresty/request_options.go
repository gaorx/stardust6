package sdresty

import (
	"bytes"
	"github.com/gaorx/stardust6/sdrand"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"io"
)

// RequestOption 用于设置resty.Request的函数
type RequestOption func(req *resty.Request) *resty.Request

// QueryParam 设置一个查询参数
func QueryParam(k string, v any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		if k == "" {
			return req
		}
		return req.SetQueryParam(k, anyToStr(v))
	}
}

// QueryParams 设置多个查询参数
func QueryParams(params map[string]any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		params1 := toStringMap(params)
		if len(params1) > 0 {
			req = req.SetQueryParams(params1)
		}
		return req
	}
}

// PathParam 设置一个路径参数
func PathParam(k string, v any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		if k == "" {
			return req
		}
		return req.SetPathParam(k, anyToStr(v))
	}
}

// PathParams 设置多个路径参数
func PathParams(params map[string]any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		params1 := toStringMap(params)
		if len(params1) > 0 {
			req = req.SetPathParams(params1)
		}
		return req
	}
}

// Header 设置一个请求头
func Header(k, v string) RequestOption {
	return func(req *resty.Request) *resty.Request {
		if k == "" {
			return req
		}
		return req.SetHeader(k, v)
	}
}

// Headers 设置多个请求头
func Headers(headers map[string]string) RequestOption {
	return func(req *resty.Request) *resty.Request {
		if len(headers) > 0 {
			req = req.SetHeaders(headers)
		}
		return req
	}
}

// UserAgent 设置User-Agent请求头，从给定的UA列表中随机选择一个
func UserAgent(uas ...string) RequestOption {
	uas = lo.Filter(uas, func(ua string, _ int) bool { return ua != "" })
	return func(req *resty.Request) *resty.Request {
		switch len(uas) {
		case 0:
			return req
		case 1:
			return req.SetHeader("User-Agent", uas[0])
		default:
			return req.SetHeader("User-Agent", sdrand.Sample(uas...))
		}
	}
}

// JsonData 设置请求体为JSON数据
func JsonData(data any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		return req.SetHeader("Content-Type", "application/json").SetBody(data)
	}
}

// FormData 设置请求体为表单数据
func FormData(data map[string]any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		params1 := toStringMap(data)
		if len(params1) > 0 {
			req = req.SetFormData(params1)
		}
		return req
	}
}

// FileData 设置请求体为文件数据
func FileData(param string, filename string, file any) RequestOption {
	return func(req *resty.Request) *resty.Request {
		if param == "" {
			return req
		}
		switch f := file.(type) {
		case nil:
			return req
		case string:
			return req.SetFile(param, f)
		case []byte:
			buff := bytes.NewBuffer(f)
			return req.SetFileReader(param, filename, buff)
		case io.Reader:
			return req.SetFileReader(param, filename, f)
		default:
			return req
		}
	}
}
