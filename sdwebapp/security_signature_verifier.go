package sdwebapp

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"slices"
	"strings"
)

type SignatureVerifier struct {
	Verifier     func(echo.Context) bool
	PathPrefixes []string
}

type SignatureVerifiers []SignatureVerifier

func VerifySign(verifier func(echo.Context) bool) SignatureVerifier {
	return SignatureVerifier{
		Verifier:     verifier,
		PathPrefixes: nil,
	}
}

func (v SignatureVerifier) For(pathPrefixes ...string) SignatureVerifier {
	v.PathPrefixes = slices.Clone(pathPrefixes)
	return v
}

func (v SignatureVerifier) ToMiddleware() MiddlewareFunc {
	return SignatureVerifiers{v}.ToMiddleware()
}

func (vs SignatureVerifiers) ToMiddleware() MiddlewareFunc {
	type item struct {
		pathPrefix string
		verifier   func(echo.Context) bool
	}
	var items []item
	for _, v := range vs {
		for _, pathPrefix := range v.PathPrefixes {
			if pathPrefix != "" && v.Verifier != nil {
				items = append(items, item{pathPrefix: pathPrefix, verifier: v.Verifier})
			}
		}
	}
	if len(items) <= 0 {
		return nil
	}
	return func(next HandlerFunc) HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			matchedItem, matched := lo.Find(items, func(item item) bool {
				return strings.HasPrefix(path, item.pathPrefix)
			})
			if matched {
				if ok := matchedItem.verifier(c); !ok {
					return echo.ErrBadRequest
				}
			}
			return next(c)
		}
	}
}
