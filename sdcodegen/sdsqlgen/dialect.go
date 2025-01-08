package sdsqlgen

type Dialect interface {
	QuoteId(id string) string
	EscapeRune(r []rune, c rune) []rune
	MakeBlob(data []byte) string
	AllowDefaultValue(typ string) bool
}
