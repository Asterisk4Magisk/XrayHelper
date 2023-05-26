package common

import "github.com/coreos/go-iptables/iptables"

const (
	CoreGid       = "3005"
	TproxyTableId = "233"
	TproxyMarkId  = "1111"
	DummyDevice   = "xdummy"
	DummyIp       = "fd01:5ca1:ab1e:8d97:497f:8b48:b9aa:85cd/128"
	DummyMarkId   = "164"
	DummyTableId  = "164"
	TunDevice     = "xtun"
	TunMTU        = 8500
	TunMultiQueue = true
	TunIPv4       = "10.10.12.1"
	TunIPv6       = "fd02:5ca1:ab1e:8d97:497f:8b48:b9aa:85cd"
	TunUdpMode    = "udp"
	TunTableId    = "168"
	TunMarkId     = "168"
)

var (
	Ipt, _   = iptables.NewWithProtocol(iptables.ProtocolIPv4)
	Ipt6, _  = iptables.NewWithProtocol(iptables.ProtocolIPv6)
	IntraNet = []string{"0.0.0.0/8", "10.0.0.0/8", "100.64.0.0/10", "127.0.0.0/8", "169.254.0.0/16",
		"172.16.0.0/12", "192.0.0.0/24", "192.0.2.0/24", "192.88.99.0/24", "192.168.0.0/16", "198.51.100.0/24",
		"203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4", "255.255.255.255/32"}
	IntraNet6 = []string{"::/128", "::1/128", "::ffff:0:0/96", "100::/64", "64:ff9b::/96", "2001::/32",
		"2001:10::/28", "2001:20::/28", "2001:db8::/32", "2002::/16", "fc00::/7", "fe80::/10", "ff00::/8"}
	ExternalIPv6 []string
)

func init() {
	ExternalIPv6, _ = getExternalIPv6Addr()
}
