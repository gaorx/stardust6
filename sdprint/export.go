package sdprint

// Print 在stdout打印一组数据，中间不插入空格
func Print(a ...any) Printer {
	return Stdout().Print(a...)
}

// Println 在stdout打印一组数据并换行，中间插入空格
func Println(a ...any) Printer {
	return Stdout().Println(a...)
}

// Printf 在stdout格式化打印
func Printf(format string, a ...any) Printer {
	return Stdout().Printf(format, a...)
}

// Ln 在stdout打印换行
func Ln() Printer {
	return Stdout().Ln()
}

// Json 在stdout打印json
func Json(v any) Printer {
	return Stdout().Json(v)
}

// Jsonln 在stdout打印json并换行
func Jsonln(v any) Printer {
	return Stdout().Jsonln(v)
}

// PlainJson 在stdout打印json，不带格式
func PlainJson(v any) Printer {
	return Stdout().PlainJson(v)
}

// PlainJsonln 在stdout打印json并换行，不带格式
func PlainJsonln(v any) Printer {
	return Stdout().PlainJsonln(v)
}

// Pretty 在stdout打印pretty格式的数据
func Pretty(v any) Printer {
	return Stdout().Pretty(v)
}

// Prettyln 在stdout打印pretty格式的数据并换行
func Prettyln(v any) Printer {
	return Stdout().Prettyln(v)
}
