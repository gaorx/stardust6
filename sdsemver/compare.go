package sdsemver

// Equal 判断两个版本是否相等
func Equal(a, b V) bool {
	return a == b
}

// Compare 比较两个语义化版本，在排序中使用
func Compare(a, b V) int {
	if a.Major < b.Major {
		return -1
	} else if a.Major > b.Major {
		return 1
	}
	if a.Minor < b.Minor {
		return -1
	} else if a.Minor > b.Minor {
		return 1
	}
	if a.Patch < b.Patch {
		return -1
	} else if a.Patch > b.Patch {
		return 1
	}
	return 0
}
