package sdblueprint

var (
	String  = primitiveType(KString)
	Bool    = primitiveType(KBool)
	Bytes   = primitiveType(KBytes)
	Int     = primitiveType(KInt)
	Int64   = primitiveType(KInt64)
	Uint    = primitiveType(KUint)
	Uint64  = primitiveType(KUint64)
	Float64 = primitiveType(KFloat64)
	Time    = primitiveType(KTime)
)

type primitiveType TypeKind

var _ Type = primitiveType(0)

func (p primitiveType) Kind() TypeKind {
	return TypeKind(p)
}

func (p primitiveType) Schema() Schema {
	return nil
}

func (p primitiveType) Elem() Type {
	return nil
}

func (p primitiveType) MakeArray() Type {
	return arrayType{p}
}
