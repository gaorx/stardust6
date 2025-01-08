package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust6/sderr"
	"slices"
)

type (
	IndexKind  int
	IndexOrder int
)

type Index struct {
	kind     IndexKind
	name     string
	comment  string
	cols     []string
	refTable string
	refCols  []string
	order    IndexOrder
}

type IndexBuilder Index

var _ builder[*Index] = (*IndexBuilder)(nil)

const (
	IndexSimple IndexKind = iota + 1
	IndexUnique
	IndexPK
	IndexFK
	IndexFullText
)

const (
	IndexAsc  IndexOrder = 0
	IndexDesc            = 1
)

func (idx *Index) Name() string {
	return idx.name
}

func (idx *Index) Comment() string {
	return idx.comment
}

func (idx *Index) Kind() IndexKind {
	return idx.kind
}

func (idx *Index) Columns() []string {
	return idx.cols
}

func (idx *Index) ReferencedTable() string {
	return idx.refTable
}

func (idx *Index) ReferencedColumns() []string {
	return idx.refCols
}

func (idx *Index) Order() IndexOrder {
	return idx.order
}

func (idx *Index) ColumnsSql(p *Project, tableId string) []string {
	t := p.TableById(tableId)
	if t == nil {
		panic(sderr.Newf("table not found %s", tableId))
	}
	var sqlCols []string
	for _, colId := range idx.cols {
		col := t.Column(colId)
		if col == nil {
			panic(sderr.Newf("column not found %s", colId))
		}
		sqlCols = append(sqlCols, col.Names().Sql())
	}
	return sqlCols
}

func (idx *Index) ReferencedTableSql(p *Project) string {
	if idx.refTable == "" {
		return ""
	}
	t := p.TableById(idx.refTable)
	if t == nil {
		panic(sderr.Newf("table not found %s", idx.refTable))
	}
	return t.Names().Sql()
}

func (idx *Index) ReferencedColumnsSql(p *Project) []string {
	if idx.refTable == "" {
		return nil
	}
	t := p.TableById(idx.refTable)
	if t == nil {
		panic(sderr.Newf("table not found %s", idx.refTable))
	}
	var sqlCols []string
	for _, colId := range idx.refCols {
		col := t.Column(colId)
		if col == nil {
			panic(sderr.Newf("column not found %s", colId))
		}
		sqlCols = append(sqlCols, col.Names().Sql())
	}
	return sqlCols
}

func (idx *Index) asBuilder() *IndexBuilder {
	return (*IndexBuilder)(idx)
}

func (idx *Index) postBuild(_ *Project) error {
	return nil
}

func (b *IndexBuilder) Name(name string) *IndexBuilder {
	b.name = name
	return b
}

func (b *IndexBuilder) Comment(comment string) *IndexBuilder {
	b.comment = comment
	return b
}

func (b *IndexBuilder) Asc() *IndexBuilder {
	b.order = IndexAsc
	return b
}

func (b *IndexBuilder) Desc() *IndexBuilder {
	b.order = IndexDesc
	return b
}

func (b *IndexBuilder) checkColumns() error {
	if len(b.cols) <= 0 {
		return fmt.Errorf("no index columns")
	}
	if b.kind == IndexFK {
		if b.refTable == "" {
			return fmt.Errorf("no reference table")
		}
		if len(b.cols) != len(b.refCols) {
			return fmt.Errorf("reference columns not match")
		}
	}
	return nil
}

func (b *IndexBuilder) prepare(_ *buildContext) error {
	return nil
}

func (b *IndexBuilder) build(_ *buildContext) (*Index, error) {
	if err := b.checkColumns(); err != nil {
		return nil, sderr.Wrap(err)
	}

	return &Index{
		kind:     b.kind,
		name:     b.name,
		comment:  b.comment,
		cols:     slices.Clone(b.cols),
		refTable: b.refTable,
		refCols:  slices.Clone(b.refCols),
		order:    b.order,
	}, nil
}
