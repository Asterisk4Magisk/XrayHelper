package sharelinks

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

// ShareLink implement this interface, that node can be converted to xray OutoundObject
type ShareLink interface {
	GetNodeInfo() string
	ToOutoundWithTag(tag string) interface{}
}

// NewShareLink parse the url, return a ShareLink
func NewShareLink(link string) (ShareLink, error) {
	if strings.HasPrefix(link, socksPrefix) {
		return newSocksShareLink(strings.TrimPrefix(link, socksPrefix))
	}
	if strings.HasPrefix(link, ssPrefix) {
		return newShadowsocksShareLink(strings.TrimPrefix(link, ssPrefix))
	}
	if strings.HasPrefix(link, vmessPrefix) {
		return newVmessShareLink(strings.TrimPrefix(link, vmessPrefix))
	}
	if strings.HasPrefix(link, vlessPrefix) {
		return newVLESSShareLink(strings.TrimPrefix(link, vlessPrefix))
	}
	if strings.HasPrefix(link, trojanPrefix) {
		return newTrojanShareLink(strings.TrimPrefix(link, trojanPrefix))
	}
	return nil, errors.New(link + " is not a supported share link").WithPrefix("sharelink")
}
