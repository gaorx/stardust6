package sdsqlgen

import (
	"fmt"
	"github.com/gaorx/stardust6/sdbytes"
	"strings"
)

var Mysql Dialect = mysqlDialect{}

type mysqlDialect struct{}

func (d mysqlDialect) QuoteId(id string) string {
	return fmt.Sprintf("`%s`", id)
}

func (d mysqlDialect) EscapeRune(r []rune, c rune) []rune {
	switch c {
	case 0:
		r = append(r, '\\', '0')
	case '\'':
		r = append(r, '\\', '\'')
	case '"':
		r = append(r, '\\', '"')
	case '\b':
		r = append(r, '\\', 'b')
	case '\n':
		r = append(r, '\\', 'n')
	case '\r':
		r = append(r, '\\', 'r')
	case '\t':
		r = append(r, '\\', 't')
	case 26:
		r = append(r, '\\', 'Z')
	case '\\':
		r = append(r, '\\', '\\')
	case '%':
		r = append(r, '\\', '%')
	case '_':
		r = append(r, '\\', '_')
	default:
		r = append(r, c)
	}
	return r
}

func (d mysqlDialect) MakeBlob(data []byte) string {
	return fmt.Sprintf("UNHEX('%s')", sdbytes.ToHexU(data))
}

func (d mysqlDialect) AllowDefaultValue(typ string) bool {
	typ = strings.ToUpper(typ)
	return !strings.Contains(typ, "BLOB") &&
		!strings.Contains(typ, "TEXT") &&
		!strings.Contains(typ, "GEOMETRY") &&
		!strings.Contains(typ, "JSON")
}
