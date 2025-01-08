package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
)

type lazySchemaType struct {
	schemaId string
	schema   Schema
}

var _ Type = (*lazySchemaType)(nil)

func SchemaType(id string) Type {
	return &lazySchemaType{schemaId: id}
}

func (s *lazySchemaType) Kind() TypeKind {
	return KSchema
}

func (s *lazySchemaType) Schema() Schema {
	if s.schema == nil {
		panic("schema not loaded")
	}
	return s.schema
}

func (s *lazySchemaType) Refs() Refs {
	if s.schema == nil {
		panic("schema not loaded")
	}
	return s.schema.Refs()
}

func (s *lazySchemaType) WithRefs(langRefs map[string]string) Type {
	return newRefsType(s, langRefs)
}

func (s *lazySchemaType) WithRef(lang, ref string) Type {
	return newRefsType(s, map[string]string{lang: ref})
}

func (s *lazySchemaType) Elem() Type {
	return nil
}

func (s *lazySchemaType) MakeArray() Type {
	return arrayType{s}
}

func resolveType(p *Project, typ Type) error {
	if typ == nil {
		return nil
	}
	if refTyp, ok := typ.(refsType); ok {
		return resolveType(p, refTyp.base)
	}
	switch typ.Kind() {
	case KSchema:
		if lazy, ok := typ.(*lazySchemaType); ok {
			return resolveSchema(p, lazy)
		}
		return nil
	case KArray:
		return resolveType(p, typ.Elem())
	default:
		return nil
	}
}

func resolveSchema(p *Project, lazy *lazySchemaType) error {
	if lazy.schema != nil {
		return nil
	}
	schema := p.schemaByIdWithNameless(lazy.schemaId)
	if schema == nil {
		return sderr.Newf("schema not found: %s", lazy.schemaId)
	}
	lazy.schema = schema
	return nil
}
