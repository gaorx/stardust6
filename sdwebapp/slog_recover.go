package sdwebapp

import (
	"context"
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdlo"
	"github.com/gaorx/stardust6/sdparse"
	"github.com/gaorx/stardust6/sdslog"
	"github.com/gaorx/stardust6/sdtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/lo"
	"log/slog"
	"time"
)

func SlogSkipAll() middleware.Skipper {
	return func(ec echo.Context) bool {
		return true
	}
}

type SlogOptions struct {
	Logger  *slog.Logger
	Level   slog.Level
	Skipper middleware.Skipper
}

func SlogRecover(opts *SlogOptions) echo.MiddlewareFunc {
	opts1 := lo.FromPtr(opts)
	if opts1.Logger == nil {
		opts1.Logger = sdslog.DiscardLogger
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			if opts1.Skipper != nil && opts1.Skipper(ec) {
				return next(ec)
			}
			req := ec.Request()
			res := ec.Response()
			startAt := time.Now()
			var nextErr, panicErr, finalErr error
			panicErr = sdlo.TryWithPanicError(func() {
				nextErr = next(ec)
			})

			if panicErr != nil {
				finalErr = panicErr
			} else {
				finalErr = nextErr
			}
			if finalErr != nil {
				ec.Error(finalErr)
			}
			elapsedHuman := time.Since(startAt)
			elapsedMs := sdtime.ToMillisF(elapsedHuman)
			statusCode := res.Status
			method := req.Method
			path := req.URL.Path
			if path == "" {
				path = "/"
			}

			bytesIn, err := sdparse.Int64E(req.Header.Get(echo.HeaderContentLength))
			if err != nil {
				bytesIn = 0
			}

			logAttrs := []any{
				slog.Float64("latency", elapsedMs),
				slog.Duration("latency_h", elapsedHuman),
				slog.String("remote_ip", ec.RealIP()),
				slog.Int64("bytes_in", bytesIn),
				slog.Int64("bytes_out", res.Size),
			}
			if finalErr == nil {
				opts1.Logger.With(logAttrs...).Log(
					context.Background(),
					opts1.Level,
					fmt.Sprintf("%d %s %s", statusCode, method, path),
				)
			} else {
				logAttrs = append(logAttrs, sdslog.E(err))
				opts1.Logger.With(logAttrs...).Error(fmt.Sprintf("%d %s %s", statusCode, method, path))
			}
			return sderr.Wrapf(finalErr, "logging recover middleware error")
		}
	}
}
