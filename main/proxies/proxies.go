package proxies

import (
	"XrayHelper/main/errors"
	"XrayHelper/main/proxies/tproxy"
	"XrayHelper/main/proxies/tun"
)

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
		return nil, errors.New("unsupported proxy method " + method).WithPrefix("proxies")
	}
}
