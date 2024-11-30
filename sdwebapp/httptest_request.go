package sdwebapp

import (
	"bytes"
	"github.com/gaorx/stardust6/sdbytes"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/samber/lo"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type TestRequest struct {
	*http.Request
}

func NewTestRequest(method, path string) *TestRequest {
	return &TestRequest{httptest.NewRequest(method, path, nil)}
}

func (req *TestRequest) Call(app *App) *TestResponse {
	return app.TestCall(req)
}

func (req *TestRequest) SetBodyBytes(contentType string, body []byte) *TestRequest {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	req.SetHeader("Content-Type", contentType)
	req.setBody(bytes.NewReader(body))
	return req
}

func (req *TestRequest) SetBodyText(contentType string, body string) *TestRequest {
	return req.SetBodyBytes(contentType, []byte(body))
}

func (req *TestRequest) SetBodyJson(v any) *TestRequest {
	j := lo.Must(sdjson.MarshalBytes(v))
	return req.SetBodyBytes("application/json", j)
}

func (req *TestRequest) SetHeader(k, v string) *TestRequest {
	if k != "" {
		req.Header.Set(k, v)
	}
	return req
}

func (req *TestRequest) SetHeaders(headers map[string]string) *TestRequest {
	for k, v := range headers {
		if k != "" {
			req.Header.Set(k, v)
		}
	}
	return req
}

func (req *TestRequest) SetBasicAuth(username, password string) *TestRequest {
	return req.SetHeader("Authorization", "Basic "+sdbytes.P([]byte(username+":"+password)).Base64Std())
}

func (req *TestRequest) SetBearerAuth(token string) *TestRequest {
	return req.SetHeader("Authorization", "Bearer "+token)
}

func (req *TestRequest) setBody(body io.Reader) {
	switch v := body.(type) {
	case *bytes.Buffer:
		req.ContentLength = int64(v.Len())
	case *bytes.Reader:
		req.ContentLength = int64(v.Len())
	case *strings.Reader:
		req.ContentLength = int64(v.Len())
	default:
		req.ContentLength = -1
	}
	if rc, ok := body.(io.ReadCloser); ok {
		req.Body = rc
	} else {
		req.Body = io.NopCloser(body)
	}
}
