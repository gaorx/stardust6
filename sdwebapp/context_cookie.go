package sdwebapp

import (
	"encoding/base64"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/gaorx/stardust6/sdparse"
	"net/http"
	"strconv"
)

func (c Context) DeleteCookie(name string, path string) {
	c.SetCookieString(name, "", path, -1)
}

func (c Context) CookieArg(name string) Argument {
	return Argument(c.CookieString(name))
}

func (c Context) CookieString(name string) string {
	v, err := c.Cookie(name)
	if err != nil {
		return ""
	}
	return v.Value
}

func (c Context) SetCookieString(name, val string, path string, maxAge int) {
	c.SetCookie(&http.Cookie{
		Name:   name,
		Value:  val,
		Path:   path,
		MaxAge: maxAge,
	})
}

func (c Context) CookieInt(name string) int {
	return sdparse.Int(c.CookieString(name))
}

func (c Context) SetCookieInt(name string, val int, path string, maxAge int) {
	c.SetCookieString(name, strconv.Itoa(val), path, maxAge)
}

func (c Context) CookieInt64(name string) int64 {
	return sdparse.Int64(c.CookieString(name))
}

func (c Context) SetCookieInt64(name string, val int64, path string, maxAge int) {
	c.SetCookieString(name, strconv.FormatInt(val, 10), path, maxAge)
}

func (c Context) CookieJson(name string, v any) error {
	base64Str := c.CookieString(name)
	jsonBytes, err := base64.URLEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}
	err = sdjson.UnmarshalBytes(jsonBytes, v)
	if err != nil {
		return sderr.Wrapf(err, "sdecho unmarshal cookie json error")
	}
	return nil
}

func (c Context) SetCookieJson(name string, v any, path string, maxAge int) error {
	jsonBytes, err := sdjson.MarshalBytes(v)
	if err != nil {
		return sderr.Wrapf(err, "sdecho marshal json cookie error")
	}
	c.SetCookieString(name, base64.URLEncoding.EncodeToString(jsonBytes), path, maxAge)
	return nil
}

func (c Context) CookieJsonObject(name string) sdjson.Object {
	base64Str := c.CookieString(name)
	jsonBytes, err := base64.URLEncoding.DecodeString(base64Str)
	if err != nil {
		return nil
	}
	v, err := sdjson.UnmarshalValueBytes(jsonBytes)
	if err != nil {
		return nil
	}
	return v.AsObject()
}
