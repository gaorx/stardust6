package sdblueprint

import "github.com/samber/lo"

type Schemas []Schema

func (ss Schemas) ForEach(f func(s Schema)) Schemas {
	if f == nil {
		return ss
	}
	for _, s := range ss {
		f(s)
	}
	return ss
}

func (ss Schemas) Ids() []string {
	ids := lo.Map(ss, func(s Schema, _ int) string { return s.Id() })
	return lo.Uniq(ids)
}

func (ss Schemas) Find(pred SchemaPredicate) Schemas {
	var filtered []Schema
	for _, s := range ss {
		if pred.apply(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (ss Schemas) Get(id string) Schema {
	if id == "" {
		return nil
	}
	for _, s := range ss {
		if s.Id() == id {
			return s
		}
	}
	return nil
}
