package sdblueprint

type schemaType struct {
	s Schema
}

var _ Type = schemaType{}

func (s schemaType) Kind() TypeKind {
	return KSchema
}

func (s schemaType) Schema() Schema {
	return s.s
}

func (s schemaType) Elem() Type {
	return nil
}

func (s schemaType) MakeArray() Type {
	return arrayType{s}
}
