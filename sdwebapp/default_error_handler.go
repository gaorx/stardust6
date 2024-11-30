package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"
)

func DefaultHttpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	he, ok := sderr.As[*echo.HTTPError](err)
	if ok {
		if he.Internal != nil {
			if herr, ok := sderr.As[*echo.HTTPError](he.Internal); ok {
				he = herr
			}
		}
	} else {
		errMsg := http.StatusText(http.StatusInternalServerError)
		if c.QueryParam("_show_error") == "1" {
			errMsg = errMsg + strings.Repeat("\r\n", 2) + err.Error()
		}
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: errMsg,
		}
	}

	var errMsg string
	if m, ok := he.Message.(string); ok {
		errMsg = m
	} else {
		errMsg = "Unknown error"
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.String(he.Code, errMsg)
	}
	if err != nil {
		slog.With("error", err.Error()).Error("http error handler error")
	}
}
