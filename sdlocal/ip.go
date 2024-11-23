package sdlocal

import (
	"net"

	"github.com/gaorx/stardust6/sderr"
)

// IPPredicate 用于断言IP是否符合要求
type IPPredicate func(net.Interface, net.IP) bool

// Not 取反
func (p IPPredicate) Not() IPPredicate {
	return func(iface net.Interface, ip net.IP) bool {
		return !p(iface, ip)
	}
}

// NetInterfaceNames 获取所有IP
func NetInterfaceNames() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, sderr.Wrapf(err, "get net interfaces error")
	}
	var ifaceNames []string
	for _, iface := range ifaces {
		ifaceNames = append(ifaceNames, iface.Name)
	}
	return ifaceNames, nil
}

// IPs 获取所有通过断言的IP
func IPs(predicates ...IPPredicate) ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, sderr.Wrapf(err, "get net interfaces error")
	}
	var ips []net.IP
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := extractIP(iface, addr)
			if len(ip) > 0 && predicateIP(iface, ip, predicates) {
				ips = append(ips, ip)
			}
		}
	}
	return ips, nil
}

// IP 获取第一个通过断言的IP
func IP(predicates ...IPPredicate) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, sderr.Wrapf(err, "get net interfaces error")
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := extractIP(iface, addr)
			if len(ip) > 0 && predicateIP(iface, ip, predicates) {
				return ip, nil
			}
		}
	}
	return nil, sderr.Newf("not found ip")
}

// IPString 获取第一个通过断言的IP的字符串形式
func IPString(predicates ...IPPredicate) string {
	ip, err := IP(predicates...)
	if err != nil || len(ip) <= 0 {
		return ""
	}
	return ip.String()
}

// PrivateIP4String 获取第一个内网IP
func PrivateIP4String(ifaceNames ...string) string {
	if len(ifaceNames) > 0 {
		return IPString(Is4(), IsPrivate(), NameIs(ifaceNames[0], ifaceNames[1:]...))
	} else {
		return IPString(Is4(), IsPrivate())
	}
}

// Is4 断言IP地址是IPv4
func Is4() IPPredicate {
	return func(_ net.Interface, ip net.IP) bool {
		ip4 := ip.To4()
		return len(ip4) > 0
	}
}

// NameIs 断言IFace的名称在指定的列表中
func NameIs(ifaceName string, others ...string) IPPredicate {
	return func(iface net.Interface, ip net.IP) bool {
		if iface.Name == ifaceName {
			return true
		}
		for _, other := range others {
			if iface.Name == other {
				return true
			}
		}
		return false
	}
}

// IsLoopback 断言IP地址是环回地址
func IsLoopback() IPPredicate {
	return func(_ net.Interface, ip net.IP) bool {
		return ip.IsLoopback()
	}
}

// IsPrivate 断言IP地址是私有地址
func IsPrivate() IPPredicate {
	return func(_ net.Interface, ip net.IP) bool {
		return ip.IsPrivate()
	}
}

// IsMulticast 断言IP地址是多播地址
func IsMulticast() IPPredicate {
	return func(_ net.Interface, ip net.IP) bool {
		return ip.IsMulticast()
	}
}

func predicateIP(iface net.Interface, ip net.IP, predicates []IPPredicate) bool {
	ok := true
	for _, pred := range predicates {
		if pred != nil && !pred(iface, ip) {
			ok = false
			break
		}
	}
	return ok
}

func extractIP(_ net.Interface, addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		return nil
	}
}
