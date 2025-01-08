package sdblueprint

import (
	"github.com/samber/lo"
	"strings"
)

func sqlJoinIds(ids []string, quoter func(string) string, bracket bool) string {
	quotedIds := lo.Map(ids, func(id string, _ int) string {
		return quoter(id)
	})
	joined := strings.Join(quotedIds, ", ")
	if bracket {
		joined = "(" + joined + ")"
	}
	return joined
}
