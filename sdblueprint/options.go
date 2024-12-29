package sdblueprint

type Options struct {
	TableNamers  map[string]Namer
	ColumnNamers map[string]Namer
}

func (opts *Options) ensureDefault() {
	// table namers
	if opts.TableNamers == nil {
		opts.TableNamers = map[string]Namer{}
	}
	opts.TableNamers["go"] = ToPascal
	opts.TableNamers["sql"] = ToPascal
	opts.TableNamers["ts"] = ToPascal

	// column namers
	if opts.ColumnNamers == nil {
		opts.ColumnNamers = map[string]Namer{}
	}
	opts.ColumnNamers["go"] = ToPascal
	opts.ColumnNamers["sql"] = ToCamel
	opts.ColumnNamers["js"] = ToCamel
	opts.ColumnNamers["json"] = ToCamel
}
