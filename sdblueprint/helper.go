package sdblueprint

import (
	"github.com/samber/lo"
	"strings"
)

func checkIdUniqueness(idSources ...[]string) ([]string, bool) {
	var ids []string
	for _, idSource := range idSources {
		ids = append(ids, idSource...)
	}
	all := make(map[string]struct{})
	repetitive := make(map[string]struct{})
	for _, id := range ids {
		if _, ok := all[id]; ok {
			repetitive[id] = struct{}{}
		}
		all[id] = struct{}{}
	}
	if len(repetitive) > 0 {
		return lo.Keys(repetitive), false
	} else {
		return nil, true
	}
}

type description struct {
	doc         string
	comment     string
	hint        string
	placeholder string
}

func makePropIdsStr(schemaId string, propIds []string) string {
	refPropIds := lo.Map(propIds, func(propId string, _ int) string {
		return schemaId + "." + propId
	})
	return strings.Join(refPropIds, ",")
}

func findLangValue[V string | []string | Namer](m map[string]V, langs []string) V {
	var zero V
	if len(langs) <= 0 {
		return zero
	}
	for _, lang := range langs {
		for k, v := range m {
			if k != "" && lang != "" {
				if strings.ToLower(k) == strings.ToLower(lang) {
					return v
				}
			}
		}
	}
	return zero
}

func lastStringOf(ss []string) string {
	if len(ss) <= 0 {
		return ""
	}
	return ss[len(ss)-1]
}

func ensurePrefix(s string, prefix string) string {
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

func trimAPIPath(s string) string {
	if s == "" {
		return ""
	}
	return strings.Trim(s, "/")
}

func makeLangs(langs1 []string, langs2 []string) []string {
	langs := make([]string, 0)
	langs = append(langs, langs1...)
	langs = append(langs, langs2...)
	return lo.Uniq(langs)
}
