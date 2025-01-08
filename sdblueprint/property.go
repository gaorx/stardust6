package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"maps"
)

type Property struct {
	id    string
	typ   Type
	names names
	anns  annotations
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

func (p *Property) Id() string {
	return p.id
}

func (p *Property) Type() Type {
	return p.typ
}

func (p *Property) Names() Names {
	return &p.names
}

func (p *Property) Anns() Anns {
	return &p.anns
}

func (p *Property) Doc() string {
	return p.doc
}

func (p *Property) Comment() string {
	return p.comment
}

func (p *Property) Hint() string {
	return p.hint
}

func (p *Property) Placeholder() string {
	return p.placeholder
}

func (p *Property) AutoIncr() bool {
	return p.autoIncr
}

func (p *Property) Nullable() bool {
	return p.nullable
}

func (p *Property) DefaultValue() any {
	return p.defaultValue
}

func (p *Property) asBuilder() *PropertyBuilder {
	return (*PropertyBuilder)(p)
}

func (p *Property) postBuild(proj *Project) error {
	if err := resolveType(proj, p.typ); err != nil {
		return sderr.Wrap(err)
	}
	return nil
}

func (b *PropertyBuilder) Id(id string) *PropertyBuilder {
	b.id = id
	return b
}

func (b *PropertyBuilder) Names(langNames map[string]string) *PropertyBuilder {
	b.names.add(langNames)
	return b
}

func (b *PropertyBuilder) Name(lang, name string) *PropertyBuilder {
	return b.Names(map[string]string{lang: name})
}

func (b *PropertyBuilder) Ann(lang string, ann ...string) *PropertyBuilder {
	b.anns.add(lang, ann)
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

func (b *PropertyBuilder) FullText() *PropertyBuilder {
	b.indexKind = IndexFullText
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
	b.names.mergeOther(&refProp.names)
	b.anns.mergeOther(&refProp.anns)
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
			defaults: c.project.namersProp,
		},
		anns:        annotations{anns: maps.Clone(b.anns.anns)},
		description: b.description,
		dbAttrs:     b.dbAttrs,
		refId:       b.refId,
	}, nil
}
