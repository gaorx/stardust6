package sdblueprint

type TypeKind int

const (
	KInvalid TypeKind = iota
	KString
	KBool
	KBytes
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
	Schema() Schema
	Elem() Type
	MakeArray() Type
}

type Schema interface {
	Id() string
	Names() Names
	Properties() []*Property
	Property(id string) *Property
	AsType() Type
}
