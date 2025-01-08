package sdblueprint

import (
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdcodegen/sdsqlgen"
	"github.com/gaorx/stardust6/sderr"
	"path/filepath"
	"strings"
)

type DDLMysqlGenerator struct {
	dir      string
	filename string
	pred     SchemaPredicate
}

var _ Generator = (*DDLMysqlGenerator)(nil)

func DDLMysql() *DDLMysqlGenerator {
	return &DDLMysqlGenerator{}
}

func (bg *DDLMysqlGenerator) To(d string) *DDLMysqlGenerator {
	bg.dir = d
	return bg
}

func (bg *DDLMysqlGenerator) File(fn string) *DDLMysqlGenerator {
	bg.filename = fn
	return bg
}

func (bg *DDLMysqlGenerator) Predicate(pred SchemaPredicate) *DDLMysqlGenerator {
	bg.pred = pred
	return bg
}

func (bg *DDLMysqlGenerator) ensureDefault() {
	if bg.dir == "" {
		bg.dir = "sqls"
	}
	if bg.filename == "" {
		bg.filename = "ddl.sql"
	}
}

func (bg *DDLMysqlGenerator) Setup(p *ProjectBuilder) {
	bg.ensureDefault()
}

func (bg *DDLMysqlGenerator) Generate(p *Project, cg *sdcodegen.Generator) {
	h := func(c *sdcodegen.Context) {
		bg.genTables(sdsqlgen.C(c), p)
	}
	cg.Add(filepath.Join(bg.dir, bg.filename), h)
}

func (bg *DDLMysqlGenerator) genTables(c *sdsqlgen.Context, p *Project) {
	tables := p.Tables(bg.pred)
	c.Newl()
	c.PrintWarning(3).Newl()
	for _, t := range tables {
		bg.genTable(c, p, t)
	}
}

func (bg *DDLMysqlGenerator) genTable(c *sdsqlgen.Context, p *Project, t *Table) {
	tableName := t.Names().Sql()
	c.Commentf("Tables %s", t.Id())
	c.CreateTable(bg.quoteId(tableName), func() {
		cols := t.Columns()
		indexes := t.Indexes()
		for i, col := range cols {
			colName := col.Names().Sql()
			if colName == "-" {
				continue
			}
			colRef := bg.genFieldRef(p, col.Type())
			c.Tab().Field(bg.quoteId(colName), colRef, &sdsqlgen.FieldOptions{
				Dialect:      sdsqlgen.Mysql,
				Comma:        i < len(cols)+len(indexes)-1,
				PrimaryKey:   false,
				AutoIncr:     col.AutoIncr(),
				Nullable:     col.Nullable(),
				DefaultValue: col.DefaultValue(),
				Comment:      col.Comment(),
			})
		}
		for i, idx := range indexes {
			c.Tab().FieldLine(bg.genIndex(p, idx, t.Id()), &sdsqlgen.FieldOptions{
				Comma: i < len(indexes)-1,
			})
		}
	}, &sdsqlgen.CreateTableOptions{
		Comment: t.Comment(),
	}).Newl()
}

func (bg *DDLMysqlGenerator) genFieldRef(p *Project, typ Type) string {
	sqlRef := typ.Refs().Sql()
	if sqlRef != "" {
		return sqlRef
	}
	switch typ.Kind() {
	case KString:
		return p.SqlTypeString()
	case KBool:
		return "TINYINT"
	case KBytes:
		return p.SqlTypeBytes()
	case KEnum:
		return "TINYINT"
	case KInt:
		return "INT"
	case KInt64:
		return "BIGINT"
	case KUint:
		return "INT UNSIGNED"
	case KUint64:
		return "BIGINT UNSIGNED"
	case KFloat64:
		return "DOUBLE"
	case KTime:
		return p.SqlTypeTime()
	case KSchema:
		return p.SqlTypeSchema()
	case KArray:
		return p.SqlTypeArray()
	default:
		panic(sderr.Newf("unknown type kind: %d", typ.Kind()))
	}
}

func (bg *DDLMysqlGenerator) genIndex(p *Project, idx *Index, tableId string) string {
	var b strings.Builder
	if idx.Name() != "" {
		b.WriteString("CONSTRAINT ")
		b.WriteString(bg.quoteId(idx.Name()))
	}
	switch idx.Kind() {
	case IndexPK:
		b.WriteString(" PRIMARY KEY ")
		b.WriteString(sqlJoinIds(idx.ColumnsSql(p, tableId), bg.quoteId, true))
	case IndexSimple:
		b.WriteString(" INDEX ")
		b.WriteString(sqlJoinIds(idx.ColumnsSql(p, tableId), bg.quoteId, true))
	case IndexUnique:
		b.WriteString(" UNIQUE INDEX ")
		b.WriteString(sqlJoinIds(idx.ColumnsSql(p, tableId), bg.quoteId, true))
	case IndexFK:
		b.WriteString(" FOREIGN KEY ")
		b.WriteString(sqlJoinIds(idx.ColumnsSql(p, tableId), bg.quoteId, true))
		b.WriteString(" REFERENCES ")
		b.WriteString(bg.quoteId(idx.ReferencedTableSql(p)))
		b.WriteString(" ")
		b.WriteString(sqlJoinIds(idx.ReferencedColumnsSql(p), bg.quoteId, true))
	case IndexFullText:
		b.WriteString(" FULLTEXT INDEX ")
		b.WriteString(sqlJoinIds(idx.ColumnsSql(p, tableId), bg.quoteId, true))
	}
	return b.String()
}

func (bg *DDLMysqlGenerator) quoteId(id string) string {
	return sdsqlgen.QuoteId(id, DialectMysql)
}
