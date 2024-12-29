package sdblueprint

import (
	"github.com/gaorx/stardust6/sdstrings"
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

func makeLangMap(langAndValues []string) map[string]string {
	langMap := make(map[string]string)
	if len(langAndValues) <= 0 {
		return langMap
	}
	for _, langAndValue := range langAndValues {
		lang, value := sdstrings.Split2s(langAndValue, ":")
		if lang != "" && value != "" {
			langMap[lang] = value
		}
	}
	return langMap
}
