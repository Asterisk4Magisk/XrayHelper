package shareurls

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"strings"
)

const (
	tagShareurl     = "shareurl"
	socksPrefix     = "socks://"
	ssPrefix        = "ss://"
	vmessPrefix     = "vmess://"
	vlessPrefix     = "vless://"
	trojanPrefix    = "trojan://"
	hysteriaPrefix  = "hysteria://"
	hysteria2Prefix = "hysteria2://"
	hy2Prefix       = "hy2://"
)

// ShareUrl implement this interface, that node can be converted to core OutoundObject
type ShareUrl interface {
	GetNodeInfo() string
	ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error)
}

// Parse return a ShareUrl
func Parse(link string) (ShareUrl, error) {
	if strings.HasPrefix(link, socksPrefix) {
		return parseSocks(link)
	}
	if strings.HasPrefix(link, ssPrefix) {
		return parseShadowsocks(link)
	}
	if strings.HasPrefix(link, vmessPrefix) {
		return parseVmess(link)
	}
	if strings.HasPrefix(link, vlessPrefix) {
		return parseVLESS(link)
	}
	if strings.HasPrefix(link, trojanPrefix) {
		return parseTrojan(link)
	}
	if strings.HasPrefix(link, hysteriaPrefix) {
		return parseHysteria(link)
	}
	if strings.HasPrefix(link, hysteria2Prefix) || strings.HasPrefix(link, hy2Prefix) {
		return parseHysteria2(link)
	}
	return nil, e.New("not a supported share link").WithPrefix(tagShareurl)
}
