package sdsemver

// ToInt 将"x.y.z"字符串形式的语义化版本转换到int64整数形式
func ToInt(s string) (int64, error) {
	v, err := Parse(s)
	if err != nil {
		return 0, err
	}
	return v.ToInt(), nil
}

// ToString 将int64整数形式的语义化版本转换到"x.y.z"字符串形式
func ToString(i int64) (string, error) {
	v, err := FromInt(i)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
