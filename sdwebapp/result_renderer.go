package sdwebapp

import (
	"fmt"
	"github.com/gaorx/stardust6/sdfile/sdfiletype"
	"github.com/gaorx/stardust6/sdjson"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type ResultRenderer func(c echo.Context, r *Result) error

const akResultRender = "sdwebapp.result.renderer"

func renderResultDefault(c echo.Context, r *Result) error {
	if r == nil {
		return nil
	}
	r1 := r.Normalize(c)
	if r1.Err != nil {
		return c.String(r1.StatusCode, r1.Err.Error())
	}
	if r1.Data == nil {
		return c.NoContent(r1.StatusCode)
	}
	if r1.View != "" {
		return c.Render(r1.StatusCode, r1.View, r1.Data)
	}
	switch d1 := r1.Data.(type) {
	case string:
		return c.String(r1.StatusCode, d1)
	case []byte:
		contentType := r1.ContentType
		if contentType == "" {
			contentType = sdfiletype.DetectBytes(d1).MimeTypeOr(echo.MIMEOctetStream)
		}
		return c.Blob(r1.StatusCode, contentType, d1)
	case error:
		if d1 == nil {
			return c.NoContent(r1.StatusCode)
		}
		return c.String(r1.StatusCode, d1.Error())
	case io.Reader:
		contentType := r1.ContentType
		if contentType == "" {
			contentType = echo.MIMEOctetStream
		}
		return c.Stream(r1.StatusCode, contentType, d1)
	default:
		j, err := sdjson.MarshalBytes(r1.Data)
		if err != nil {
			return c.String(http.StatusInternalServerError, "unsupported result data type")
		}
		return c.Blob(r1.StatusCode, echo.MIMEApplicationJSON, j)
	}
}

func AsText(c echo.Context, r *Result) error {
	if r == nil {
		return nil
	}
	r1 := r.Normalize(c)
	if r1.Err != nil {
		return c.String(r1.StatusCode, r1.Err.Error())
	}
	switch d1 := r1.Data.(type) {
	case nil:
		return c.NoContent(r1.StatusCode)
	case string:
		return c.String(r1.StatusCode, d1)
	case fmt.Stringer:
		return c.String(r1.StatusCode, d1.String())
	default:
		j := sdjson.MarshalPretty(d1)
		return c.String(r1.StatusCode, j)
	}
}

func AsJson(c echo.Context, r *Result) error {
	if r == nil {
		return nil
	}
	r1 := r.Normalize(c)
	if r1.Err != nil {
		return c.String(r1.StatusCode, r1.Err.Error())
	}
	return c.JSON(r1.StatusCode, r1.Data)
}

func AsPrettyJson(c echo.Context, r *Result) error {
	if r == nil {
		return nil
	}
	r1 := r.Normalize(c)
	if r1.Err != nil {
		return c.String(r1.StatusCode, r1.Err.Error())
	}
	return c.JSONPretty(r1.StatusCode, r1.Data, "  ")
}
