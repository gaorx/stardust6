package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"maps"
	"strings"
)

type Property struct {
	id    string
	typ   Type
	names names
	sigs  map[string]Sig
	anns  map[string][]string
	description
	dbAttrs
	refId string
}

type dbAttrs struct {
	indexKind    IndexKind
	autoIncr     bool
	nullable     bool
	defaultValue any
}

type PropertyBuilder Property

var (
	_ builder[*Property] = (*PropertyBuilder)(nil)
)

func (p *Property) Names() Names {
	return &p.names
}

func (p *Property) Types() Sigs {
	return p.sigs
}

func (p *Property) Anns() Anns {
	return p.anns
}

func (p *Property) asBuilder() *PropertyBuilder {
	return (*PropertyBuilder)(p)
}

func (b *PropertyBuilder) Id(id string) *PropertyBuilder {
	b.id = id
	return b
}

func (b *PropertyBuilder) Names(langAndNames ...string) *PropertyBuilder {
	b.names.add(langAndNames)
	return b
}

func (b *PropertyBuilder) Sigs(langAndSigs ...string) *PropertyBuilder {
	if len(langAndSigs) <= 0 {
		return b
	}
	if b.sigs == nil {
		b.sigs = map[string]Sig{}
	}
	for lang, v := range makeLangMap(langAndSigs) {
		if lang != "" && v != "" {
			b.sigs[strings.ToLower(lang)] = SigOf(v)
		}
	}
	return b
}

func (b *PropertyBuilder) Ann(lang string, anns ...string) *PropertyBuilder {
	if lang == "" || len(anns) <= 0 {
		return b
	}
	if b.anns == nil {
		b.anns = map[string][]string{}
	}
	b.anns[strings.ToLower(lang)] = anns
	return b
}

func (b *PropertyBuilder) Doc(doc string) *PropertyBuilder {
	b.doc = doc
	return b
}

func (b *PropertyBuilder) Comment(comment string) *PropertyBuilder {
	b.comment = comment
	return b
}

func (b *PropertyBuilder) Hint(hint string) *PropertyBuilder {
	b.hint = hint
	return b
}

func (b *PropertyBuilder) Placeholder(placeholder string) *PropertyBuilder {
	b.placeholder = placeholder
	return b
}

func (b *PropertyBuilder) RefId(id string) *PropertyBuilder {
	b.refId = id
	return b
}

func (b *PropertyBuilder) PrimaryKey() *PropertyBuilder {
	b.indexKind = IndexPK
	return b
}

func (b *PropertyBuilder) Index() *PropertyBuilder {
	b.indexKind = IndexSimple
	return b
}

func (b *PropertyBuilder) Unique() *PropertyBuilder {
	b.indexKind = IndexUnique
	return b
}

func (b *PropertyBuilder) AutoIncr() *PropertyBuilder {
	b.autoIncr = true
	return b
}

func (b *PropertyBuilder) Nullable() *PropertyBuilder {
	b.nullable = true
	return b
}

func (b *PropertyBuilder) DefaultValue(value any) *PropertyBuilder {
	b.defaultValue = value
	return b
}

func (b *PropertyBuilder) isRef() bool {
	return b.refId != ""
}

func (b *PropertyBuilder) deref(c *buildContext) error {
	if !b.isRef() {
		return nil
	}
	if c.buildingPropertyDepth > 10 {
		return sderr.Newf("property reference too deep")
	}
	c.increaseBuildingPropertyDepth()
	defer c.decreaseBuildingPropertyDepth()

	refProp := c.project.asProject().PropertyByRef(b.refId)
	if refProp == nil {
		return sderr.Newf("property reference not found (%s)", b.refId)
	}
	if err := refProp.asBuilder().deref(c); err != nil {
		return sderr.Wrap(err)
	}

	// refProp的属性和当前属性合并
	if b.id == "" {
		b.id = refProp.id
	}
	if b.typ == nil {
		b.typ = refProp.typ
	}
	if len(refProp.names.names) > 0 {
		if b.names.names == nil {
			b.names.names = map[string]string{}
		}
		for lang, name := range refProp.names.names {
			if _, ok := b.names.names[lang]; !ok {
				b.names.names[lang] = name
			}
		}
	}
	if len(refProp.sigs) > 0 {
		if b.sigs == nil {
			b.sigs = map[string]Sig{}
		}
		for lang, typ := range refProp.sigs {
			if _, ok := b.sigs[lang]; !ok {
				b.sigs[lang] = typ
			}
		}
	}
	if len(refProp.anns) > 0 {
		if b.anns == nil {
			b.anns = map[string][]string{}
		}
		for lang, anns := range refProp.anns {
			if _, ok := b.anns[lang]; !ok {
				b.anns[lang] = anns
			}
		}
	}
	if b.doc == "" {
		b.doc = refProp.doc
	}
	if b.comment == "" {
		b.comment = refProp.comment
	}
	if b.hint == "" {
		b.hint = refProp.hint
	}
	if b.placeholder == "" {
		b.placeholder = refProp.placeholder
	}
	b.refId = "" // deref DONE
	return nil
}

func (b *PropertyBuilder) prepare(c *buildContext) error {
	c.setBuildingProperty(b)
	defer c.unsetBuildingProperty()

	if err := b.deref(c); err != nil {
		return sderr.Wrap(err)
	}
	return nil
}

func (b *PropertyBuilder) build(c *buildContext) (*Property, error) {
	c.setBuildingProperty(b)
	defer c.unsetBuildingProperty()

	return &Property{
		id:  b.id,
		typ: b.typ,
		names: names{
			id:       b.id,
			names:    maps.Clone(b.names.names),
			defaults: c.options.ColumnNamers,
		},
		sigs:        maps.Clone(b.sigs),
		anns:        maps.Clone(b.anns),
		description: b.description,
		dbAttrs:     b.dbAttrs,
		refId:       b.refId,
	}, nil
}
