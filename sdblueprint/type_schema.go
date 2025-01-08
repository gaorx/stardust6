package sdblueprint

type schemaType struct {
	s Schema
}

var _ Type = schemaType{}

func newSchemaType(s Schema) schemaType {
	return schemaType{s: s}
}

func (s schemaType) Kind() TypeKind {
	return KSchema
}

func (s schemaType) Schema() Schema {
	return s.s
}

func (s schemaType) Refs() Refs {
	return s.s.Refs()
}

func (s schemaType) WithRefs(langRefs map[string]string) Type {
	return newRefsType(s, langRefs)
}

func (s schemaType) WithRef(lang, ref string) Type {
	return newRefsType(s, map[string]string{lang: ref})
}

func (s schemaType) Elem() Type {
	return nil
}

func (s schemaType) MakeArray() Type {
	return arrayType{s}
}
