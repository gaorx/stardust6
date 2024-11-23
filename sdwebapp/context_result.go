package sdwebapp

import (
	"github.com/labstack/echo/v4"
)

func (c Context) ResultRenderer() ResultRenderer {
	v := c.Get(akResultRender)
	switch v1 := v.(type) {
	case nil:
		return nil
	case ResultRenderer:
		if v1 == nil {
			return nil
		}
		return v1
	case func(echo.Context, *Result) error:
		if v1 == nil {
			return nil
		}
		return v1
	default:
		return nil
	}
}

func (c Context) SetResultRenderer(r ResultRenderer) {
	if r != nil {
		c.Set(akResultRender, r)
	}
}
