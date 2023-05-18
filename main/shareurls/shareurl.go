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
	ToOutoundWithTag(tag string) interface{}
}

// NewShareUrl parse the url, return a ShareUrl
func NewShareUrl(link string) (ShareUrl, error) {
	if strings.HasPrefix(link, socksPrefix) {
		return newSocksShareUrl(strings.TrimPrefix(link, socksPrefix))
	}
	if strings.HasPrefix(link, ssPrefix) {
		return newShadowsocksShareUrl(strings.TrimPrefix(link, ssPrefix))
	}
	if strings.HasPrefix(link, vmessPrefix) {
		return newVmessShareUrl(strings.TrimPrefix(link, vmessPrefix))
	}
	if strings.HasPrefix(link, vlessPrefix) {
		return newVLESSShareUrl(strings.TrimPrefix(link, vlessPrefix))
	}
	if strings.HasPrefix(link, trojanPrefix) {
		return newTrojanShareUrl(strings.TrimPrefix(link, trojanPrefix))
	}
	return nil, errors.New("not a supported share link").WithPrefix("shareurls")
}
