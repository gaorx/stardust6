package sdsql

// Page 表示数据查询中的分页信息
type Page struct {
	num   int64
	size  int64
	base0 bool
}

// Page0 创建一个从0开始的分页信息
func Page0(num, size int64) Page {
	p := Page{num: num, size: size, base0: true}
	ensurePageNum(&p)
	ensurePageSize(&p)
	return p
}

// Page1 创建一个从1开始的分页信息
func Page1(num, size int64) Page {
	p := Page{num: num, size: size, base0: false}
	ensurePageNum(&p)
	ensurePageSize(&p)
	return p
}

// Num 返回页码
func (p Page) Num() int64 {
	return p.num
}

// Size 返回每页大小
func (p Page) Size() int64 {
	return p.size
}

// LimitAndOffset 返回数据库查询中这个分页在SQL中需要的limit和offset值
func (p Page) LimitAndOffset() (int64, int64) {
	limit, num := p.size, p.num
	if p.base0 {
		return limit, num * limit
	} else {
		return limit, (num - 1) * limit
	}
}

func ensurePageNum(p *Page) {
	if p.base0 {
		if p.num <= 0 {
			p.num = 0
		}
	} else {
		if p.num <= 1 {
			p.num = 1
		}
	}
}

func ensurePageSize(p *Page) {
	const (
		defaultSize = 1
		maxSize     = 100000
	)
	if p.size <= 0 {
		p.size = defaultSize
	} else if p.size > maxSize {
		p.size = maxSize
	}
}
