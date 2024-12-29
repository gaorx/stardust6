package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"maps"
	"slices"
	"strings"
)

type Object struct {
	id         string
	categories []string
	names      names
	anns       map[string][]string
	description
	props []*Property
}

type ObjectBuilder Object

var (
	_ Schema           = (*Object)(nil)
	_ builder[*Object] = (*ObjectBuilder)(nil)
)

func (o *Object) Id() string {
	return o.id
}

func (o *Object) Categories() []string {
	return o.categories
}

func (o *Object) Names() Names {
	return &o.names
}

func (o *Object) Anns() Anns {
	return o.anns
}

func (o *Object) Properties() []*Property {
	return o.props
}

func (o *Object) Property(id string) *Property {
	if id == "" {
		return nil
	}
	for _, prop := range o.props {
		if prop.id == id {
			return prop
		}
	}
	return nil
}

func (o *Object) AsType() Type {
	return schemaType{o}
}

func (o *Object) asBuilder() *ObjectBuilder {
	return (*ObjectBuilder)(o)
}

func (b *ObjectBuilder) Id(id string) *ObjectBuilder {
	b.id = id
	return b
}

func (b *ObjectBuilder) Categories(categories ...string) *ObjectBuilder {
	b.categories = categories
	return b
}

func (b *ObjectBuilder) Names(langAndNames ...string) *ObjectBuilder {
	b.names.add(langAndNames)
	return b
}

func (b *ObjectBuilder) Ann(lang string, ann ...string) *ObjectBuilder {
	if lang == "" || len(ann) <= 0 {
		return b
	}
	if b.anns == nil {
		b.anns = map[string][]string{}
	}
	b.anns[strings.ToLower(lang)] = ann
	return b
}

func (b *ObjectBuilder) Doc(doc string) *ObjectBuilder {
	b.doc = doc
	return b
}

func (b *ObjectBuilder) Comment(comment string) *ObjectBuilder {
	b.comment = comment
	return b
}

func (b *ObjectBuilder) Prop(id string, typ Type) *PropertyBuilder {
	prop := &Property{
		id:  id,
		typ: typ,
	}
	b.props = append(b.props, prop)
	return prop.asBuilder()
}

func (b *ObjectBuilder) RefProp(refProp string) *PropertyBuilder {
	prop := &Property{
		refId: refProp,
	}
	b.props = append(b.props, prop)
	return prop.asBuilder()
}

func (b *ObjectBuilder) RefProps(refProps ...string) {
	for _, refProp := range refProps {
		b.RefProp(refProp)
	}
}

func (b *ObjectBuilder) RefPropsIn(schemaId string, propIds ...string) {
	for _, propId := range propIds {
		if schemaId != "" && propId != "" {
			b.RefProp(schemaId + "." + propId)
		}
	}
}

func (b *ObjectBuilder) asObject() *Object {
	return (*Object)(b)
}

func (b *ObjectBuilder) prepare(c *buildContext) error {
	c.setBuildingObject(b)
	defer c.unsetBuildingObject()

	for _, prop := range b.props {
		if err := prop.asBuilder().prepare(c); err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

func (b *ObjectBuilder) build(c *buildContext) (*Object, error) {
	c.setBuildingObject(b)
	defer c.unsetBuildingObject()

	// 检测Id是否为空
	for _, prop := range b.props {
		if prop.id == "" {
			return nil, sderr.Newf("property id is empty")
		}
	}

	// 检测属性Id不能重复
	if repetitive, ok := checkIdUniqueness(
		lo.Map(b.props, func(col *Property, _ int) string { return col.id }),
	); !ok {
		return nil, sderr.Newf("property id repetitive (%s)", makePropIdsStr(c.buildingSchemaId(), repetitive))
	}

	// 创建新的属性
	var newProps []*Property
	for _, prop := range b.props {
		prop1, err := prop.asBuilder().build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newProps = append(newProps, prop1)
	}

	return &Object{
		id:         b.id,
		categories: slices.Clone(b.categories),
		names: names{
			id:       b.id,
			names:    maps.Clone(b.names.names),
			defaults: c.options.TableNamers,
		},
		anns:        maps.Clone(b.anns),
		description: b.description,
		props:       newProps,
	}, nil
}
