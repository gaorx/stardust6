package sdsnowflake

import (
	"github.com/samber/lo"
)

var (
	defaultNode *Node
	zeroNode    *Node
)

func init() {
	defaultNode = lo.Must(NewFromLocalIP4())
	zeroNode = lo.Must(New(0))
}

// Generate 随机生成一个int64类型的ID，使用本机IP作为因素，不同机器生成的ID不会重复，但是本机的ID可能重复
func Generate() int64 {
	return defaultNode.Generate()
}

// GenerateZero 随机生成一个int64类型的ID，生成的ID可能重复
func GenerateZero() int64 {
	return zeroNode.Generate()
}
