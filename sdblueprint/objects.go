package sdblueprint

import "github.com/samber/lo"

type Objects []*Object

func (os Objects) ForEach(f func(o *Object)) Objects {
	if f == nil {
		return os
	}
	for _, o := range os {
		f(o)
	}
	return os
}

func (os Objects) Ids() []string {
	ids := lo.Map(os, func(o *Object, _ int) string { return o.Id() })
	return lo.Uniq(ids)
}

func (os Objects) Find(pred SchemaPredicate) Objects {
	var filtered []*Object
	for _, o := range os {
		if pred.apply(o) {
			filtered = append(filtered, o)
		}
	}
	return filtered
}

func (os Objects) Get(id string) *Object {
	if id == "" {
		return nil
	}
	for _, o := range os {
		if o.id == id {
			return o
		}
	}
	return nil
}
