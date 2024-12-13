package sdwebapp

import (
	"github.com/gaorx/stardust6/sdjson"
	"github.com/samber/lo"
	"net/http/httptest"
)

type TestResponse struct {
	*httptest.ResponseRecorder
}

func (res *TestResponse) Let(f func(response *TestResponse)) {
	if f != nil {
		f(res)
	}
}

func (res *TestResponse) HeaderMap() map[string]string {
	m := map[string]string{}
	for k, vals := range res.Header() {
		if len(vals) > 0 {
			m[k] = vals[0]
		}
	}
	return m
}

func (res *TestResponse) BodyBytes() []byte {
	return res.Body.Bytes()
}

func (res *TestResponse) BodyText() string {
	return res.Body.String()
}

func (res *TestResponse) BodyJson(p any) {
	lo.Must0(sdjson.UnmarshalBytes(res.Body.Bytes(), p))
}

func (res *TestResponse) BodyJsonValue() sdjson.Value {
	return lo.Must(sdjson.UnmarshalValueBytes(res.Body.Bytes()))
}

func (res *TestResponse) BodyJsonObject() sdjson.Object {
	jv := res.BodyJsonValue()
	jo, ok := jv.ToObject(false)
	if !ok {
		panic("body is not a json object")
	}
	return jo
}

func TestResponseBodyJson[T any](res *TestResponse) T {
	var body T
	res.BodyJson(&body)
	return body
}
