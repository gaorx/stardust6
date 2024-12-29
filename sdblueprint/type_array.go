package sdblueprint

type arrayType struct {
	elem Type
}

var _ Type = arrayType{}

func (a arrayType) Kind() TypeKind {
	return KArray
}

func (a arrayType) Schema() Schema {
	return nil
}

func (a arrayType) Elem() Type {
	return a.elem
}

func (a arrayType) MakeArray() Type {
	return arrayType{a}
}
