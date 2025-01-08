package sdblueprint

import (
	"slices"
)

type SchemaPredicate func(Schema) bool

func (f SchemaPredicate) apply(schema Schema) bool {
	if f == nil {
		return true
	}
	return f(schema)
}

func ByAnd(preds ...SchemaPredicate) SchemaPredicate {
	return func(schema Schema) bool {
		for _, f := range preds {
			if !f.apply(schema) {
				return false
			}
		}
		return true
	}
}

func ByOr(preds ...SchemaPredicate) SchemaPredicate {
	return func(schema Schema) bool {
		for _, f := range preds {
			if f.apply(schema) {
				return true
			}
		}
		return false
	}
}

func ByAll() SchemaPredicate {
	return func(_ Schema) bool {
		return true
	}
}

func ByNone() SchemaPredicate {
	return func(_ Schema) bool {
		return false
	}
}

func ByImpl(impl SchemaImpl) SchemaPredicate {
	return func(schema Schema) bool {
		return schema.Impl() == impl
	}
}

func ByTable() SchemaPredicate {
	return ByImpl(SchemaTable)
}

func ByObject() SchemaPredicate {
	return ByImpl(SchemaObject)
}

func ByCategory(category string) SchemaPredicate {
	return func(schema Schema) bool {
		if category == "" {
			return false
		}
		return slices.Contains(schema.Categories(), category)
	}
}

func ByInclude(ids ...string) SchemaPredicate {
	return func(schema Schema) bool {
		if len(ids) <= 0 {
			return false
		}
		return slices.Contains(ids, schema.Id())
	}
}

func ByExclude(ids ...string) SchemaPredicate {
	return func(schema Schema) bool {
		if len(ids) <= 0 {
			return true
		}
		return !slices.Contains(ids, schema.Id())
	}
}
