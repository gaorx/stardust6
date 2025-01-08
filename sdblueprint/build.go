package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
)

func Build(builder func(*ProjectBuilder)) (*Project, error) {
	c := newBuildContext()
	c.project.ensureDefault()

	// build
	builder(&c.project)

	// setup generators
	for _, bg := range c.project.generators {
		bg.Setup(&c.project)
	}

	// add default post hooks
	c.project.SchemaRefs(ByAll(), defaultSchemaRefs)

	// prepare
	if err := c.project.prepare(c); err != nil {
		return nil, sderr.Wrap(err)
	}

	// build
	p, err := c.project.build(c)
	if err != nil {
		return nil, sderr.Wrap(err)
	}

	// post build
	if err := p.postBuild(p); err != nil {
		return nil, sderr.Wrap(err)
	}

	// post do
	if len(p.postHooks) > 0 {
		for _, hook := range p.postHooks {
			if err := hook(p); err != nil {
				return nil, sderr.Wrap(err)
			}
		}
	}
	return p, nil
}

func MustBuild(builder func(*ProjectBuilder)) *Project {
	p, err := Build(builder)
	if err != nil {
		panic(err)
	}
	return p
}

type buildable interface {
	postBuild(p *Project) error
}

type builder[T buildable] interface {
	prepare(c *buildContext) error
	build(c *buildContext) (T, error)
}

type buildContext struct {
	project               ProjectBuilder
	buildingTable         *TableBuilder
	buildingObject        *ObjectBuilder
	buildingProperty      *PropertyBuilder
	buildingPropertyDepth int
}

func newBuildContext() *buildContext {
	return &buildContext{
		project: ProjectBuilder{},
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
