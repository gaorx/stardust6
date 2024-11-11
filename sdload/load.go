package sdload

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"net/url"

	"github.com/BurntSushi/toml"
	"github.com/gaorx/stardust6/sderr"
)

// Bytes 加载字节数据
func Bytes(loc string) ([]byte, error) {
	var scheme string
	u, err := url.Parse(loc)
	if err != nil {
		scheme = ""
	} else {
		scheme = u.Scheme
	}
	l, ok := loaders[scheme]
	if !ok {
		return nil, sderr.With("loc", loc).Newf("unknown scheme")
	}
	data, err := l.LoadBytes(loc)
	if err != nil {
		return nil, sderr.With("loc", loc).Wrapf(err, "load error")
	}
	return data, nil
}

// Text 加载文本数据
func Text(loc string) (string, error) {
	data, err := Bytes(loc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSON 加载JSON格式的数据
func JSON[T any](loc string) (T, error) {
	var empty, r T
	data, err := Bytes(loc)
	if err != nil {
		return empty, err
	}
	err = json.Unmarshal(data, &r)
	if err != nil {
		return empty, sderr.Wrapf(err, "parse json error")
	}
	return r, nil
}

// TOML 加载TOML格式的数据
func TOML[T any](loc string) (T, error) {
	var empty, r T
	data, err := Bytes(loc)
	if err != nil {
		return empty, err
	}
	err = toml.Unmarshal(data, &r)
	if err != nil {
		return empty, sderr.Wrapf(err, "parse TOML error")
	}
	return r, nil
}

// YAML 加载YAML格式的数据
func YAML[T any](loc string) (T, error) {
	var empty, r T
	data, err := Bytes(loc)
	if err != nil {
		return empty, err
	}
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return empty, sderr.Wrapf(err, "parse YAML error")
	}
	return r, nil
}
