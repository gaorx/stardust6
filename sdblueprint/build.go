package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
)

func Build(builder func(*ProjectBuilder), opts *Options) (*Project, error) {
	c := newBuildContext(opts)
	builder(&c.project)
	if err := c.project.prepare(c); err != nil {
		return nil, sderr.Wrap(err)
	}
	p, err := c.project.build(c)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	return p, nil
}

func MustBuild(builder func(*ProjectBuilder), opts *Options) *Project {
	p, err := Build(builder, opts)
	if err != nil {
		panic(err)
	}
	return p
}

func NewObject(builder func(*ObjectBuilder), opts *Options) (*Object, error) {
	p, err := Build(func(b *ProjectBuilder) {
		b.Object(builder)
	}, opts)
	if err != nil {
		return nil, sderr.Wrap(err)
	}
	if len(p.objects) != 1 {
		return nil, sderr.Newf("expect one object")
	}
	return p.objects[0], nil
}

func MustObject(builder func(*ObjectBuilder), opts *Options) *Object {
	o, err := NewObject(builder, opts)
	if err != nil {
		panic(err)
	}
	return o
}

type builder[T any] interface {
	prepare(c *buildContext) error
	build(c *buildContext) (T, error)
}

type buildContext struct {
	project               ProjectBuilder
	options               Options
	buildingTable         *TableBuilder
	buildingObject        *ObjectBuilder
	buildingProperty      *PropertyBuilder
	buildingPropertyDepth int
}

func newBuildContext(opts *Options) *buildContext {
	opts1 := lo.FromPtr(opts)
	opts1.ensureDefault()
	return &buildContext{
		project: ProjectBuilder{},
		options: opts1,
	}
}

func (c *buildContext) buildingSchemaId() string {
	if c.buildingTable != nil {
		return c.buildingTable.id
	}
	if c.buildingObject != nil {
		return c.buildingObject.id
	}
	return ""
}

func (c *buildContext) setBuildingTable(t *TableBuilder) {
	c.buildingTable = t
}

func (c *buildContext) unsetBuildingTable() {
	c.buildingTable = nil
}

func (c *buildContext) setBuildingObject(o *ObjectBuilder) {
	c.buildingObject = o
}

func (c *buildContext) unsetBuildingObject() {
	c.buildingObject = nil
}

func (c *buildContext) setBuildingProperty(p *PropertyBuilder) {
	c.buildingProperty = p
}

func (c *buildContext) unsetBuildingProperty() {
	c.buildingProperty = nil
}

func (c *buildContext) increaseBuildingPropertyDepth() {
	c.buildingPropertyDepth += 1
}

func (c *buildContext) decreaseBuildingPropertyDepth() {
	c.buildingPropertyDepth -= 1
}
