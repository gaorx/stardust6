package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"github.com/samber/lo"
	"maps"
)

type Table struct {
	id         string
	categories []string
	names      names
	refs       references
	anns       annotations
	description
	cols    []*Property
	indexes []*Index
}

type TableBuilder Table

var (
	_ Schema          = (*Table)(nil)
	_ builder[*Table] = (*TableBuilder)(nil)
)

func (t *Table) Id() string {
	return t.id
}

func (t *Table) Impl() SchemaImpl {
	return SchemaTable
}

func (t *Table) Categories() []string {
	return t.categories
}

func (t *Table) Names() Names {
	return &t.names
}

func (t *Table) Comment() string {
	return t.comment
}

func (t *Table) Doc() string {
	return t.doc
}

func (t *Table) Refs() Refs {
	return &t.refs
}

func (t *Table) Anns() Anns {
	return &t.anns
}

func (t *Table) Columns() []*Property {
	return t.cols
}

func (t *Table) Column(id string) *Property {
	if id == "" {
		return nil
	}
	for _, col := range t.cols {
		if col.id == id {
			return col
		}
	}
	return nil
}

func (t *Table) Properties() []*Property {
	return t.Columns()
}

func (t *Table) Property(id string) *Property {
	if id == "" {
		return nil
	}
	return t.Column(id)
}

func (t *Table) Indexes() []*Index {
	return t.indexes
}

func (t *Table) asBuilder() *TableBuilder {
	return (*TableBuilder)(t)
}

func (t *Table) AsType() Type {
	return schemaType{t}
}

func (t *Table) postBuild(p *Project) error {
	for _, col := range t.cols {
		if err := col.postBuild(p); err != nil {
			return sderr.Wrap(err)
		}
	}
	for _, idx := range t.indexes {
		if err := idx.postBuild(p); err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

func (b *TableBuilder) Id(id string) *TableBuilder {
	b.id = id
	return b
}

func (b *TableBuilder) Categories(categories ...string) *TableBuilder {
	b.categories = categories
	return b
}

func (b *TableBuilder) Names(langNames map[string]string) *TableBuilder {
	b.names.add(langNames)
	return b
}

func (b *TableBuilder) Name(lang, name string) *TableBuilder {
	return b.Names(map[string]string{lang: name})
}

func (b *TableBuilder) Refs(langRefs map[string]string) *TableBuilder {
	b.refs.add(langRefs)
	return b
}

func (b *TableBuilder) Ref(lang, ref string) *TableBuilder {
	return b.Refs(map[string]string{lang: ref})
}

func (b *TableBuilder) Ann(lang string, ann ...string) *TableBuilder {
	b.anns.add(lang, ann)
	return b
}

func (b *TableBuilder) Doc(doc string) *TableBuilder {
	b.doc = doc
	return b
}

func (b *TableBuilder) Comment(comment string) *TableBuilder {
	b.comment = comment
	return b
}

func (b *TableBuilder) Col(id string, typ Type) *PropertyBuilder {
	col := &Property{
		id:  id,
		typ: typ,
	}
	b.cols = append(b.cols, col)
	return col.asBuilder()
}

func (b *TableBuilder) RefCol(refProp string) *PropertyBuilder {
	col := &Property{
		refId: refProp,
	}
	b.cols = append(b.cols, col)
	return col.asBuilder()
}

func (b *TableBuilder) RefCols(refProps ...string) {
	for _, refProp := range refProps {
		b.RefCol(refProp)
	}
}

func (b *TableBuilder) RefColsIn(schemaId string, propIds ...string) {
	for _, propId := range propIds {
		if schemaId != "" && propId != "" {
			b.RefCol(schemaId + "." + propId)
		}
	}
}

func (b *TableBuilder) Index(cols ...string) *IndexBuilder {
	cols1 := lo.Filter(cols, func(col string, _ int) bool {
		return col != ""
	})
	if len(cols1) <= 0 {
		panic("no index columns")
	}
	idx := &Index{
		kind: IndexSimple,
		cols: cols1,
	}
	b.indexes = append(b.indexes, idx)
	return idx.asBuilder()
}

func (b *TableBuilder) Unique(cols ...string) *IndexBuilder {
	idx := &Index{
		kind: IndexUnique,
		cols: cols,
	}
	b.indexes = append(b.indexes, idx)
	return idx.asBuilder()
}

func (b *TableBuilder) PrimaryKey(cols ...string) *IndexBuilder {
	idx := &Index{
		kind: IndexPK,
		cols: cols,
	}
	b.indexes = append(b.indexes, idx)
	return idx.asBuilder()
}

func (b *TableBuilder) ForeignKey(cols []string, refTable string, refCols []string) *IndexBuilder {
	idx := &Index{
		kind:     IndexFK,
		cols:     cols,
		refTable: refTable,
		refCols:  refCols,
	}
	b.indexes = append(b.indexes, idx)
	return idx.asBuilder()
}

func (b *TableBuilder) FullText(col string) *IndexBuilder {
	idx := &Index{
		kind: IndexFullText,
		cols: []string{col},
	}
	b.indexes = append(b.indexes, idx)
	return idx.asBuilder()
}

func (b *TableBuilder) asTable() *Table {
	return (*Table)(b)
}

func (b *TableBuilder) prepare(c *buildContext) error {
	c.setBuildingTable(b)
	defer c.unsetBuildingTable()

	for _, col := range b.cols {
		if err := col.asBuilder().prepare(c); err != nil {
			return sderr.Wrap(err)
		}
	}
	return nil
}

func (b *TableBuilder) build(c *buildContext) (*Table, error) {
	c.setBuildingTable(b)
	defer c.unsetBuildingTable()

	// 检测Id是否为空
	if b.id == "" {
		return nil, sderr.Newf("table id is empty")
	}

	// 检测列是否为空
	if len(b.cols) <= 0 {
		return nil, sderr.Newf("no columns")
	}

	// 检测列名是否为空
	for _, col := range b.cols {
		if col.id == "" {
			return nil, sderr.Newf("column id is empty")
		}
	}

	// 检测列名不能重复
	if repetitive, ok := checkIdUniqueness(
		lo.Map(b.cols, func(col *Property, _ int) string { return col.id }),
	); !ok {
		return nil, sderr.Newf("column id repetitive (%s)", makePropIdsStr(c.buildingSchemaId(), repetitive))
	}

	// 构建新的索引，将列中的关于索引的选项加过来
	var indexBuilders []*IndexBuilder
	for _, col := range b.cols {
		if col.indexKind != 0 {
			indexBuilders = append(indexBuilders, &IndexBuilder{
				kind: col.indexKind,
				cols: []string{col.id},
			})
			// 清空col中关于索引的信息
			col.indexKind = 0
		}
	}
	// 再将原来的索引加过来
	for _, idx := range b.indexes {
		indexBuilders = append(indexBuilders, idx.asBuilder())
	}

	// 检测所有索引中的PrimaryKey只能有一项，且不能缺失
	pkCount := lo.CountBy(indexBuilders, func(idx *IndexBuilder) bool {
		return idx.kind == IndexPK
	})
	if pkCount <= 0 {
		return nil, sderr.Newf("no primary key")
	}
	if pkCount > 1 {
		return nil, sderr.Newf("multiple primary key")
	}

	// 创建新的列
	var newCols []*Property
	for _, col := range b.cols {
		col1, err := col.asBuilder().build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newCols = append(newCols, col1)
	}
	// 创建新的索引
	var newIndexes []*Index
	for _, idx := range indexBuilders {
		idx1, err := idx.build(c)
		if err != nil {
			return nil, sderr.Wrap(err)
		}
		newIndexes = append(newIndexes, idx1)
	}
	return &Table{
		id:         b.id,
		categories: b.categories,
		names: names{
			id:       b.id,
			names:    maps.Clone(b.names.names),
			defaults: c.project.namersSchema,
		},
		refs: references{
			base: b.refs.base,
			refs: maps.Clone(b.refs.refs),
		},
		anns:        annotations{anns: maps.Clone(b.anns.anns)},
		description: b.description,
		cols:        newCols,
		indexes:     newIndexes,
	}, nil
}
