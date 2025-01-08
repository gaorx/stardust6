package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen/sdsqlgen"
)

var DialectMysql Dialect = dialectMysql{}

type dialectMysql struct{}

func (d dialectMysql) QuoteId(id string) string {
	return sdsqlgen.Mysql.QuoteId(id)
}

func (d dialectMysql) EscapeRune(r []rune, c rune) []rune {
	return sdsqlgen.Mysql.EscapeRune(r, c)
}

func (d dialectMysql) MakeBlob(data []byte) string {
	return sdsqlgen.Mysql.MakeBlob(data)
}

func (d dialectMysql) AllowDefaultValue(typ string) bool {
	return sdsqlgen.Mysql.AllowDefaultValue(typ)
}

func (d dialectMysql) DefaultStringType() string {
	return "VARCHAR(255)"
}

func (d dialectMysql) DefaultBytesType() string {
	return "BLOB"
}

func (d dialectMysql) DefaultTimeType() string {
	return "DATETIME"
}

func (d dialectMysql) DefaultSchemaType() string {
	return "TEXT"
}

func (d dialectMysql) DefaultArrayType() string {
	return "TEXT"
}
