package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdstrings"
	"github.com/samber/lo"
	"slices"
	"strings"
)

type Project struct {
	tables  []*Table
	objects []*Object
}

type ProjectBuilder Project

var (
	_ builder[*Project] = (*ProjectBuilder)(nil)
)

func (p *Project) Tables() []*Table {
	return p.tables
}

func (p *Project) TablesByCategory(category string) []*Table {
	return lo.Filter(p.tables, func(t *Table, _ int) bool {
		return slices.Contains(t.categories, category)
	})
}

func (p *Project) Objects() []*Object {
	return p.objects
}

func (p *Project) ObjectsByCategory(category string) []*Object {
	return lo.Filter(p.objects, func(o *Object, _ int) bool {
		return slices.Contains(o.categories, category)
	})
}

func (p *Project) Schemas() []Schema {
	var schemas []Schema
	for _, t := range p.tables {
		schemas = append(schemas, t)
	}
	for _, o := range p.objects {
		schemas = append(schemas, o)
	}
	return schemas
}

func (p *Project) TableById(id string) *Table {
	if id == "" {
		return nil
	}
	for _, t := range p.tables {
		if t.id == id {
			return t
		}
	}
	return nil
}

func (p *Project) ObjectById(id string) *Object {
	if id == "" {
		return nil
	}
	for _, o := range p.objects {
		if o.id == id {
			return o
		}
	}
	return nil
}

func (p *Project) SchemaById(id string) Schema {
	if id == "" {
		return nil
	}
	if t := p.TableById(id); t != nil {
		return t
	}
	if o := p.ObjectById(id); o != nil {
		return o
	}
	return nil
}

func (p *Project) PropertyById(schemaId, propId string) *Property {
	schema := p.SchemaById(schemaId)
	if schema == nil {
		return nil
	}
	return schema.Property(propId)
}

func (p *Project) PropertyByRef(refPropId string) *Property {
	schemaId, propId := sdstrings.Split2s(refPropId, ".")
	return p.PropertyById(schemaId, propId)
}

func (b *ProjectBuilder) Table(builders ...func(*TableBuilder)) *ProjectBuilder {
	for _, b0 := range builders {
		if b0 == nil {
			continue
		}
		t := &TableBuilder{}
		b0(t)
		b.tables = append(b.tables, t.asTable())
	}
	return b
}

func (b *ProjectBuilder) Object(builders ...func(*ObjectBuilder)) *ProjectBuilder {
	for _, b0 := range builders {
		if b0 == nil {
			continue
		}
		o := &ObjectBuilder{}
		b0(o)
		b.objects = append(b.objects, o.asObject())
	}
	return b
}

func (b *ProjectBuilder) asProject() *Project {
	return (*Project)(b)
}

func (b *ProjectBuilder) prepare(c *buildContext) error {
	// prepare tables
	for _, t := range b.tables {
		if err := t.asBuilder().prepare(c); err != nil {
			return sderr.Wrap(err)
		}
	}
	// prepares objects
	for _, o := range b.objects {
		if err := o.asBuilder().prepare(c); err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

func (b *ProjectBuilder) build(c *buildContext) (*Project, error) {
	// 检测Id是否为空
	for _, t := range b.tables {
		if t.id == "" {
			return nil, sderr.Newf("table id is empty")
		}
	}
	for _, o := range b.objects {
		if o.id == "" {
			return nil, sderr.Newf("object id is empty")
		}
	}

	// 检测ID是否重复
	if repetitive, ok := checkIdUniqueness(
		lo.Map(b.tables, func(t *Table, _ int) string { return t.id }),
		lo.Map(b.objects, func(o *Object, _ int) string { return o.id }),
	); !ok {
		return nil, sderr.Newf("schema id repetitive (%s)", strings.Join(repetitive, ","))
	}

	// build tables
	var newTables []*Table
	for _, t := range b.tables {
		t1, err := t.asBuilder().build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newTables = append(newTables, t1)
	}

	// build objects
	var newObjects []*Object
	for _, o := range b.objects {
		o1, err := o.asBuilder().build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newObjects = append(newObjects, o1)
	}

	return &Project{
		tables:  newTables,
		objects: newObjects,
	}, nil
}
