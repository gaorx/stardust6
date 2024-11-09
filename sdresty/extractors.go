package sdresty

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/go-resty/resty/v2"
	"os"
	"path/filepath"
)

// Nothing 表示无意义的返回结果
type Nothing struct{}

// ForResponse 直接返回resty.Response
func ForResponse() ExtractorFunc[*resty.Response] {
	return func(resp *resty.Response, err error) (*resty.Response, error) {
		if err != nil {
			return nil, sderr.Wrapf(err, "request failed")
		}
		return resp, nil
	}
}

// ForBytesBody 从Response提取body的字节数据
func ForBytesBody() ExtractorFunc[[]byte] {
	return func(resp *resty.Response, err error) ([]byte, error) {
		if err != nil {
			return nil, sderr.Wrapf(err, "request failed")
		}
		data := resp.Body()
		return data, nil
	}
}

// ForTextBody 从Response提取body的中文本
func ForTextBody() ExtractorFunc[string] {
	return func(resp *resty.Response, err error) (string, error) {
		if err != nil {
			return "", sderr.Wrapf(err, "request failed")
		}
		data := resp.Body()
		return string(data), nil
	}
}

// ForJsonBody 从Response提取body的json数据
func ForJsonBody[T any]() ExtractorFunc[T] {
	return func(resp *resty.Response, err error) (T, error) {
		if err != nil {
			var zero T
			return zero, sderr.Wrapf(err, "request failed")
		}
		data := resp.Body()
		v, err := sdjson.UnmarshalBytesT[T](data)
		if err != nil {
			var zero T
			return zero, sderr.Wrapf(err, "unmarshal response json body failed")
		}
		return v, nil
	}
}

// ForJsonValue 从Response提取body的json数据，并转换成sdjson.Value
func ForJsonValue() ExtractorFunc[sdjson.Value] {
	return func(resp *resty.Response, err error) (sdjson.Value, error) {
		if err != nil {
			return sdjson.Value{}, sderr.Wrapf(err, "request failed")
		}
		data := resp.Body()
		v, err := sdjson.UnmarshalValueBytes(data)
		if err != nil {
			return sdjson.Value{}, sderr.Wrapf(err, "unmarshal response json body failed")
		}
		return v, nil
	}
}

// SaveToFile 保存Response的body到文件，mkdirp表示是否自动文件的创建目录
func SaveToFile(path string, mkdirp bool) ExtractorFunc[Nothing] {
	return func(resp *resty.Response, err error) (Nothing, error) {
		if err != nil {
			return Nothing{}, sderr.Wrapf(err, "request failed")
		}
		if mkdirp {
			err = os.MkdirAll(filepath.Dir(path), 0600)
			if err != nil {
				return Nothing{}, sderr.Wrapf(err, "create directory failed")
			}
		}
		w, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return Nothing{}, sderr.Wrapf(err, "open file failed")
		}
		defer func() { _ = w.Close() }()

		_, err = w.Write(resp.Body())
		if err != nil {
			return Nothing{}, sderr.Wrapf(err, "save response body to file failed")
		}
		return Nothing{}, nil
	}
}
