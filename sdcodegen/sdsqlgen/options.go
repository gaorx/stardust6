package sdsqlgen

type CreateTableOptions struct {
	Dialect      Dialect
	Scope        string
	IfNotExists  bool
	PostModifier string
	Comment      string
}

type FieldOptions struct {
	Dialect      Dialect
	Comma        bool
	PrimaryKey   bool
	AutoIncr     bool
	Nullable     bool
	DefaultValue any
	Comment      string
	Other        string
}
