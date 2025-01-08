package sdblueprint

var (
	String  = primitiveType(KString)
	Bool    = primitiveType(KBool)
	Bytes   = primitiveType(KBytes)
	Enum    = primitiveType(KEnum)
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

func (p primitiveType) Refs() Refs {
	return &references{}
}

func (p primitiveType) WithRefs(langRefs map[string]string) Type {
	return newRefsType(p, langRefs)
}

func (p primitiveType) WithRef(lang, ref string) Type {
	return newRefsType(p, map[string]string{lang: ref})
}

func (p primitiveType) Schema() Schema {
	return nil
}

func (p primitiveType) Elem() Type {
	return nil
}

func (p primitiveType) MakeArray() Type {
	return arrayType{elem: p}
}
