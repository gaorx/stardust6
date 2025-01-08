package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdrand"
	"github.com/gaorx/stardust6/sdstrings"
	"github.com/samber/lo"
	"maps"
	"slices"
	"strings"
)

type Project struct {
	tables          []*Table
	objects         []*Object
	namelessObjects []*Object
	apis            []*APIDecl
	generators      []Generator

	// options
	root         string
	modName      string
	postHooks    []func(*Project) error
	namersSchema map[string]Namer
	namersProp   map[string]Namer
	sqlDialect   Dialect
	sqlTypString string
	sqlTypBytes  string
	sqlTypTime   string
	sqlTypSchema string
	sqlTypArray  string
}

type ProjectBuilder Project

var (
	_ builder[*Project] = (*ProjectBuilder)(nil)
)

func (p *Project) Tables(pred SchemaPredicate) Tables {
	return Tables(p.tables).Find(pred)
}

func (p *Project) Objects(pred SchemaPredicate) Objects {
	return Objects(p.objects).Find(pred)
}

func (p *Project) Schemas(pred SchemaPredicate) Schemas {
	return p.allSchemas().Find(pred)
}

func (p *Project) TableById(id string) *Table {
	return Tables(p.tables).Get(id)
}

func (p *Project) ObjectById(id string) *Object {
	return Objects(p.objects).Get(id)
}

func (p *Project) SchemaById(id string) Schema {
	return p.allSchemas().Get(id)
}

func (p *Project) schemaByIdWithNameless(id string) Schema {
	schema := p.SchemaById(id)
	if schema != nil {
		return schema
	}
	o := Objects(p.namelessObjects).Get(id)
	if o != nil {
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

func (p *Project) Root() string {
	return p.root
}

func (p *Project) ModName() string {
	return p.modName
}

func (p *Project) SqlTypeString() string {
	return p.sqlTypString
}

func (p *Project) SqlTypeBytes() string {
	return p.sqlTypBytes
}

func (p *Project) SqlTypeTime() string {
	return p.sqlTypTime
}

func (p *Project) SqlTypeSchema() string {
	return p.sqlTypSchema
}

func (p *Project) SqlTypeArray() string {
	return p.sqlTypArray
}

func (p *Project) APIs() APIDecls {
	return p.apis
}

func (p *Project) allSchemas() Schemas {
	var schemas Schemas
	for _, t := range p.tables {
		schemas = append(schemas, t)
	}
	for _, o := range p.objects {
		schemas = append(schemas, o)
	}
	return schemas
}

func (p *Project) postBuild(proj *Project) error {
	for _, t := range p.tables {
		if err := t.postBuild(proj); err != nil {
			return sderr.Wrap(err)
		}
	}
	for _, o := range p.objects {
		if err := o.postBuild(proj); err != nil {
			return sderr.Wrap(err)
		}
	}
	for _, o := range p.namelessObjects {
		if err := o.postBuild(proj); err != nil {
			return sderr.Wrap(err)
		}
	}
	for _, api := range p.apis {
		if err := api.postBuild(proj); err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

func (b *ProjectBuilder) Tables(builders ...func(*TableBuilder)) *ProjectBuilder {
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

func (b *ProjectBuilder) Objects(builders ...func(*ObjectBuilder)) *ProjectBuilder {
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

func (b *ProjectBuilder) MakeObject(builder func(*ObjectBuilder)) Type {
	if builder == nil {
		panic("object builder is nil")
	}
	id := "Nameless" + sdrand.String(8, sdrand.UpperCaseAlphanumericCharset)
	o := &ObjectBuilder{}
	builder(o)
	o.id = id
	b.namelessObjects = append(b.namelessObjects, o.asObject())
	return SchemaType(id)
}

func (b *ProjectBuilder) APIs(apis ...*APIDeclBuilder) *ProjectBuilder {
	for _, api := range apis {
		if api == nil {
			continue
		}
		b.apis = append(b.apis, api.asAPI())
	}
	return b
}

func (b *ProjectBuilder) Generate(generators ...Generator) *ProjectBuilder {
	for _, g := range generators {
		if g != nil {
			b.generators = append(b.generators, g)
		}
	}
	return b
}

func (b *ProjectBuilder) Root(root string) *ProjectBuilder {
	b.root = root
	return b
}

func (b *ProjectBuilder) ModName(modName string) *ProjectBuilder {
	b.modName = modName
	return b
}

func (b *ProjectBuilder) OnPost(hook func(*Project) error) *ProjectBuilder {
	if hook == nil {
		return b
	}
	b.postHooks = append(b.postHooks, hook)
	return b
}

func (b *ProjectBuilder) DefaultNamersForSchema(namers map[string]Namer) *ProjectBuilder {
	for lang, namer := range namers {
		if lang != "" && namer != nil {
			if b.namersSchema == nil {
				b.namersSchema = make(map[string]Namer)
			}
			b.namersSchema[lang] = namer
		}
	}
	return b
}

func (b *ProjectBuilder) DefaultNamersForProperty(namers map[string]Namer) *ProjectBuilder {
	for lang, namer := range namers {
		if lang != "" && namer != nil {
			if b.namersProp == nil {
				b.namersProp = make(map[string]Namer)
			}
			b.namersProp[lang] = namer
		}
	}
	return b
}

func (b *ProjectBuilder) SchemaRefs(pred SchemaPredicate, refs func(*Project, Schema) map[string]string) *ProjectBuilder {
	if refs == nil {
		return b
	}
	b.OnPost(func(p *Project) error {
		p.Schemas(pred).ForEach(func(schema Schema) {
			switch x := schema.(type) {
			case *Table:
				x.refs.mergeRefs(refs(p, schema))
			case *Object:
				x.refs.mergeRefs(refs(p, schema))
			}
		})
		return nil
	})
	return b
}

func (b *ProjectBuilder) SqlDialect(d Dialect) *ProjectBuilder {
	b.sqlDialect = d
	return b
}

func (b *ProjectBuilder) DefaultSqlTypeForString(sqlTyp string) *ProjectBuilder {
	b.sqlTypString = sqlTyp
	return b
}

func (b *ProjectBuilder) DefaultSqlTypeForBytes(sqlTyp string) *ProjectBuilder {
	b.sqlTypBytes = sqlTyp
	return b
}

func (b *ProjectBuilder) DefaultSqlTypeForTime(sqlTyp string) *ProjectBuilder {
	b.sqlTypTime = sqlTyp
	return b
}

func (b *ProjectBuilder) DefaultSqlTypeForSchema(sqlTyp string) *ProjectBuilder {
	b.sqlTypSchema = sqlTyp
	return b
}

func (b *ProjectBuilder) DefaultSqlTypeForArray(sqlTyp string) *ProjectBuilder {
	b.sqlTypArray = sqlTyp
	return b
}

func (b *ProjectBuilder) ensureDefault() {
	// root
	if b.root == "" {
		b.root = getGoModDir()
	}
	if b.modName == "" {
		b.modName = getGoModName()
	}
	if b.modName == "" {
		b.modName = "example.com/app"
	}

	// table namers
	if b.namersSchema == nil {
		b.namersSchema = map[string]Namer{}
	}
	mergeNamers(b.namersSchema, map[string]Namer{
		"go":  ToPascal,
		"sql": ToPascal,
		"js":  ToPascal,
	})

	// column namers
	if b.namersProp == nil {
		b.namersProp = map[string]Namer{}
	}
	mergeNamers(b.namersProp, map[string]Namer{
		"go":   ToPascal,
		"sql":  ToCamel,
		"js":   ToCamel,
		"json": ToCamel,
	})

	// sql dialect
	if b.sqlDialect == nil {
		b.sqlDialect = DialectMysql
	}

	// string sql type
	if b.sqlTypString == "" {
		b.sqlTypString = b.sqlDialect.DefaultStringType()
	}
	if b.sqlTypBytes == "" {
		b.sqlTypBytes = b.sqlDialect.DefaultBytesType()
	}
	if b.sqlTypTime == "" {
		b.sqlTypTime = b.sqlDialect.DefaultTimeType()
	}
	if b.sqlTypSchema == "" {
		b.sqlTypSchema = b.sqlDialect.DefaultSchemaType()
	}
	if b.sqlTypArray == "" {
		b.sqlTypArray = b.sqlDialect.DefaultArrayType()
	}
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
	for _, o := range b.namelessObjects {
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
	for _, api := range b.apis {
		if api.id == "" {
			return nil, sderr.Newf("API id is empty")
		}
	}

	// 检测ID是否重复
	if repetitive, ok := checkIdUniqueness(
		lo.Map(b.tables, func(t *Table, _ int) string { return t.id }),
		lo.Map(b.objects, func(o *Object, _ int) string { return o.id }),
		lo.Map(b.apis, func(api *APIDecl, _ int) string { return api.id }),
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
	var newNamelessObjects []*Object
	for _, o := range b.namelessObjects {
		o1, err := o.asBuilder().build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newNamelessObjects = append(newNamelessObjects, o1)
	}

	// build APIs
	var newAPIs []*APIDecl
	for _, api := range b.apis {
		api1, err := api.asBuilder().build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newAPIs = append(newAPIs, api1)
	}

	return &Project{
		tables:          newTables,
		objects:         newObjects,
		namelessObjects: newNamelessObjects,
		apis:            newAPIs,
		generators:      slices.Clone(b.generators),

		// options
		root:         b.root,
		modName:      b.modName,
		postHooks:    slices.Clone(b.postHooks),
		namersSchema: maps.Clone(b.namersSchema),
		namersProp:   maps.Clone(b.namersProp),
		sqlTypString: b.sqlTypString,
		sqlTypBytes:  b.sqlTypBytes,
		sqlTypTime:   b.sqlTypTime,
		sqlTypSchema: b.sqlTypSchema,
		sqlTypArray:  b.sqlTypArray,
	}, nil
}
