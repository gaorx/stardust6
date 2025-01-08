package sdblueprint

import (
	"github.com/samber/lo"
)

type Tables []*Table

func (ts Tables) ForEach(f func(t *Table)) Tables {
	if f == nil {
		return ts
	}
	for _, t := range ts {
		f(t)
	}
	return ts
}

func (ts Tables) Ids() []string {
	ids := lo.Map(ts, func(t *Table, _ int) string { return t.Id() })
	return lo.Uniq(ids)
}

func (ts Tables) Find(pred SchemaPredicate) Tables {
	var filtered []*Table
	for _, t := range ts {
		if pred.apply(t) {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func (ts Tables) Get(id string) *Table {
	if id == "" {
		return nil
	}
	for _, t := range ts {
		if t.id == id {
			return t
		}
	}
	return nil
}
