package sdwebapp

import (
	"github.com/gaorx/stardust6/sderr"
)

var (
	_ Component = (*Route)(nil)
	_ Component = (Routables)(nil)
)

func (r *Route) Apply(app *App) error {
	if r == nil {
		return nil
	}
	r1, err := r.build()
	if err != nil {
		return sderr.Wrap(err)
	}
	if r1 == nil {
		return nil
	}
	r2 := app.Add(r1.method, r1.path, r1.handler, r1.middlewares...)
	if r1.name != "" {
		r2.Name = r1.name
	}
	return nil
}

func (rs Routes) Apply(app *App) error {
	for _, route := range rs {
		err := route.Apply(app)
		if err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

func (rs Routables) Apply(app *App) error {
	for _, routable := range rs {
		if routable == nil {
			continue
		}
		for _, route := range routable.ToRoutes(app) {
			err := route.Apply(app)
			if err != nil {
				return sderr.Wrap(err)
			}
		}
	}
	return nil
}
