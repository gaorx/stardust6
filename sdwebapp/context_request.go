package sdwebapp

import (
	"encoding/json"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/gaorx/stardust6/sdparse"
	"github.com/samber/lo"
	"io"
	"net/http"
	"slices"
)

type Argument = sdparse.Presentable

func (c Context) RequestBodyBytes() ([]byte, error) {
	reader := c.Request().Body
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, sderr.Wrapf(err, "sdecho read request body error")
	}
	return r, nil
}

func (c Context) RequestBodyString() (string, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (c Context) RequestBodyJsonValue() (sdjson.Value, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return sdjson.Value{}, err
	}
	v, err := sdjson.UnmarshalValueBytes(b)
	if err != nil {
		return sdjson.Value{}, err
	}
	return v, nil
}

func (c Context) RequestBodyJsonObject() (sdjson.Object, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c Context) RequestBodyJsonArray() (sdjson.Array, error) {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return nil, err
	}
	var a []any
	err = json.Unmarshal(b, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c Context) RequestBodyJson(v any) error {
	b, err := c.RequestBodyBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (c Context) PathParams() map[string]string {
	names := c.ParamNames()
	m := map[string]string{}
	for _, n := range names {
		m[n] = c.Param(n)
	}
	return m
}

func (c Context) QueryArg(name string, candidates ...string) Argument {
	names := append([]string{name}, candidates...)
	for _, n := range names {
		if v := c.QueryParam(n); v != "" {
			return Argument(v)
		}
	}
	return ""
}

func (c Context) PathArg(name string, candidates ...string) Argument {
	names := append([]string{name}, candidates...)
	for _, n := range names {
		if v := c.Param(n); v != "" {
			return Argument(v)
		}
	}
	return ""
}

func (c Context) FormArg(name string, candidates ...string) Argument {
	names := append([]string{name}, candidates...)
	for _, n := range names {
		if v := c.FormValue(n); v != "" {
			return Argument(v)
		}
	}
	return ""
}

var formSupportMethods = []string{
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodOptions,
	http.MethodConnect,
}

func (c Context) VariousArg(name string, candidates ...string) Argument {
	names := append([]string{name}, candidates...)
	for _, n := range names {
		if v := c.QueryParam(n); v != "" {
			return Argument(v)
		}
		if v := c.Param(n); v != "" {
			return Argument(v)
		}
		if slices.Contains(formSupportMethods, c.Request().Method) {
			if v := c.FormValue(n); v != "" {
				return Argument(v)
			}
		}
	}
	return ""
}

func (c Context) HandleMultipartFormFile(files []string, h func(file FileHeader) error) error {
	form, err := c.MultipartForm()
	if err != nil {
		return sderr.Wrapf(err, "parse multiple form error")
	}
	if len(files) <= 0 {
		files = lo.Keys(form.File)
	}
	for _, fileField := range files {
		fileParts := form.File[fileField]
		for _, filePart := range fileParts {
			if h != nil {
				if err := h(FileHeader{
					File:       fileField,
					FileHeader: filePart,
				}); err != nil {
					return sderr.With("file", fileField).Wrapf(err, "handle upload file error")
				}
			}
		}
	}
	return nil
}
