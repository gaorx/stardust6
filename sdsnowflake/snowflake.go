package sdsnowflake

import (
	"github.com/gaorx/stardust6/sdlocal"
	"net"
	"sync"
	"time"

	"github.com/gaorx/stardust6/sderr"
)

const (
	nodeBits        = 10
	stepBits        = 12
	nodeMax         = -1 ^ (-1 << nodeBits)
	stepMask  int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

var Epoch int64 = 1288834974657

// Node 生成器
type Node struct {
	mux  sync.Mutex
	time int64
	node int64
	step int64
}

// New 从一个指定的值创建一个生成器，不同值创建的生成器会生成不同的ID
func New(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
		return nil, sderr.Newf("node number must be between 0 and 1023")
	}
	return &Node{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

// NewFromIP4 从IPv4地址创建一个生成器，不同机器生成的ID不会重复，但是本机的ID可能重复
func NewFromIP4(addr string) (*Node, error) {
	ip := net.ParseIP(addr)
	if ip == nil {
		return nil, sderr.With("ip", addr).Newf("sdsnowflake parse ip error")
	}
	ip4 := ip.To4()
	if len(ip4) <= 0 {
		return nil, sderr.With("ip", addr).Newf("sdsnowflake ip not IPv4")
	}
	var node int64 = 0
	h := int64([]byte(ip)[2]) & int64(0x03) // 0b00000011
	l := int64([]byte(ip)[3])
	node = (h << 1) | l
	return New(node)
}

func NewFromLocalIP4() (*Node, error) {
	localAddr := sdlocal.IPString(sdlocal.Is4(), sdlocal.IsPrivate())
	return NewFromIP4(localAddr)
}

// Generate 生成一个ID，可能重复，但尽量不重复
func (n *Node) Generate() int64 {
	n.mux.Lock()
	defer n.mux.Unlock()

	now := time.Now().UnixNano() / 1000000
	if n.time == now {
		n.step = (n.step + 1) & stepMask
		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}
	n.time = now
	return (now-Epoch)<<timeShift | (n.node << nodeShift) | (n.step)
}
