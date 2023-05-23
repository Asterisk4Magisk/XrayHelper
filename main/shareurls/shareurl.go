package shareurls

import (
	"XrayHelper/main/errors"
	"strings"
)

const (
	socksPrefix  = "socks://"
	ssPrefix     = "ss://"
	vmessPrefix  = "vmess://"
	vlessPrefix  = "vless://"
	trojanPrefix = "trojan://"
)

// ShareUrl implement this interface, that node can be converted to xray OutoundObject
type ShareUrl interface {
	GetNodeInfo() string
	ToOutoundWithTag(coreType string, tag string) (interface{}, error)
}

// ParseShareUrl parse the url, return a ShareUrl
func ParseShareUrl(link string) (ShareUrl, error) {
	if strings.HasPrefix(link, socksPrefix) {
		return parseSocksShareUrl(link)
	}
	if strings.HasPrefix(link, ssPrefix) {
		return parseShadowsocksShareUrl(link)
	}
	if strings.HasPrefix(link, vmessPrefix) {
		return parseVmessShareUrl(strings.TrimPrefix(link, vmessPrefix))
	}
	if strings.HasPrefix(link, vlessPrefix) {
		return parseVLESSShareUrl(link)
	}
	if strings.HasPrefix(link, trojanPrefix) {
		return parseTrojanShareUrl(link)
	}
	return nil, errors.New("not a supported share link").WithPrefix("shareurls")
}
