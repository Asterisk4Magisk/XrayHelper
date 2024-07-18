package proxies

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/proxies/tproxy"
	"XrayHelper/main/proxies/tun"
)

const tagProxies = "proxies"

// ProxyMethod implement this interface, that program can use different proxy method
type ProxyMethod interface {
	Enable() error
	Disable()
}

func NewProxy(method string) (ProxyMethod, error) {
	switch method {
	case "tproxy":
		return new(tproxy.Tproxy), nil
	case "tun", "tun2socks":
		return new(tun.Tun), nil
	default:
		return nil, e.New("unsupported proxy method " + method).WithPrefix(tagProxies)
	}
}
