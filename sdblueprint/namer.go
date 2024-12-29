package sdblueprint

import (
	"github.com/gaorx/stardust6/sdstrings"
)

type Namer func(id string) string

var (
	ToSnakeL Namer = sdstrings.ToSnakeL
	ToSnakeU Namer = sdstrings.ToSnakeU
	ToKebabL Namer = sdstrings.ToKebabL
	ToKebabU Namer = sdstrings.ToKebabU
	ToCamel  Namer = sdstrings.ToCamelL
	ToPascal Namer = sdstrings.ToCamelU
)

func (namer Namer) Prefix(prefix string) Namer {
	return func(id string) string {
		return prefix + namer(id)
	}
}

func (namer Namer) Suffix(suffix string) Namer {
	return func(id string) string {
		return namer(id) + suffix
	}
}

func (namer Namer) PrefixSuffix(prefix, suffix string) Namer {
	return func(id string) string {
		return prefix + namer(id) + suffix
	}
}
