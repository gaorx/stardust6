package sdstrings

import (
	"github.com/fatih/camelcase"
	"github.com/iancoleman/strcase"
)

// ToSnakeL 转换字符串到小写的Snake形式
func ToSnakeL(s string) string {
	return strcase.ToSnake(s)
}

// ToSnakeU 转换字符串到大写的Snake形式
func ToSnakeU(s string) string {
	return strcase.ToScreamingSnake(s)
}

// ToKebabL 转换字符串到小写的Kebab形式
func ToKebabL(s string) string {
	return strcase.ToKebab(s)
}

// ToKebabU 转换字符串到大写的Kebab形式
func ToKebabU(s string) string {
	return strcase.ToScreamingKebab(s)
}

// ToDelimitedL 转换字符串到小写的指定分隔符形式
func ToDelimitedL(s string, delimiter uint8) string {
	return strcase.ToDelimited(s, delimiter)
}

// ToDelimitedU 转换字符串到大写的指定分隔符形式
func ToDelimitedU(s string, delimiter uint8) string {
	return strcase.ToScreamingDelimited(s, delimiter, "", true)
}

// ToCamelL 转换字符串到首词小写的Camel形式
func ToCamelL(s string) string {
	return strcase.ToLowerCamel(s)
}

// ToCamelU 转换字符串到首词大写的Camel形式
func ToCamelU(s string) string {
	return strcase.ToCamel(s)
}

// SplitCamel 拆分Camel形式的字符串到单词切片
func SplitCamel(s string) []string {
	return camelcase.Split(s)
}
