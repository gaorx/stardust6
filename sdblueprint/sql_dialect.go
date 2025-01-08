package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen/sdsqlgen"
)

type Dialect interface {
	sdsqlgen.Dialect
	DefaultStringType() string
	DefaultBytesType() string
	DefaultTimeType() string
	DefaultSchemaType() string
	DefaultArrayType() string
}
