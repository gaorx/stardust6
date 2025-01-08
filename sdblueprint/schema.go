package sdblueprint

type Schema interface {
	Id() string
	Impl() SchemaImpl
	Names() Names
	Refs() Refs
	Categories() []string
	Properties() []*Property
	Property(id string) *Property
	AsType() Type
}

type SchemaImpl int

const (
	SchemaTable SchemaImpl = iota + 1
	SchemaObject
)
