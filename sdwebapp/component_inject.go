package sdwebapp

import (
	"github.com/labstack/echo/v4"
)

func Inject(injectables ...Injectable) ComponentFunc {
	m := map[string]any{}
	for _, injectable := range injectables {
		if injectable != nil {
			for _, attr := range injectable.Injections() {
				if attr.K != "" && attr.V != nil {
					m[attr.K] = attr.V
				}
			}
		}
	}
	return func(app *App) error {
		app.Pre(func(next HandlerFunc) HandlerFunc {
			return func(c echo.Context) error {
				for k, v := range m {
					c.Set(k, v)
				}
				return next(c)
			}
		})
		return nil
	}
}
