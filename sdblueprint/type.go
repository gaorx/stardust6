package sdblueprint

type TypeKind int

const (
	KInvalid TypeKind = iota
	KString
	KBool
	KBytes
	KEnum
	KInt
	KInt64
	KUint
	KUint64
	KFloat64
	KTime
	KSchema
	KArray
)

type Type interface {
	Kind() TypeKind
	Refs() Refs
	WithRefs(langRefs map[string]string) Type
	WithRef(lang, ref string) Type
	Schema() Schema
	Elem() Type
	MakeArray() Type
}
