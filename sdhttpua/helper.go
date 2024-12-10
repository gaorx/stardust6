package sdhttpua

import (
	"github.com/samber/lo"
	"sort"
)

func uniqAndSort(ss []string) []string {
	ss = lo.Uniq(ss)
	sort.Strings(ss)
	return ss
}
