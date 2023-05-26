package proxys

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/proxys/tproxy"
	"XrayHelper/main/proxys/tun"
)

// EnableTproxy enable tproxy
func EnableTproxy() error {
	if err := tproxy.AddRoute(false); err != nil {
		DisableTproxy()
		return err
	}
	if err := tproxy.CreateMangleChain(false); err != nil {
		DisableTproxy()
		return err
	}
	if err := tproxy.CreateProxyChain(false); err != nil {
		DisableTproxy()
		return err
	}
	if builds.Config.Proxy.EnableIPv6 {
		if err := tproxy.AddRoute(true); err != nil {
			DisableTproxy()
			return err
		}
		if err := tproxy.CreateMangleChain(true); err != nil {
			DisableTproxy()
			return err
		}
		if err := tproxy.CreateProxyChain(true); err != nil {
			DisableTproxy()
			return err
		}
	}
	return nil
}

// DisableTproxy disable tproxy
func DisableTproxy() {
	tproxy.DeleteRoute(false)
	tproxy.CleanIptablesChain(false)
	if builds.Config.Proxy.EnableIPv6 {
		tproxy.DeleteRoute(true)
		tproxy.CleanIptablesChain(true)
	}
}

// EnableTun enable tun
func EnableTun() error {
	if err := tun.StartTun(); err != nil {
		DisableTun()
		return err
	}
	if err := tun.AddRoute(false); err != nil {
		DisableTun()
		return err
	}
	if err := tun.CreateMangleChain(false); err != nil {
		DisableTun()
		return err
	}
	if err := tun.CreateProxyChain(false); err != nil {
		DisableTun()
		return err
	}
	if builds.Config.Proxy.EnableIPv6 {
		if err := tun.AddRoute(true); err != nil {
			DisableTun()
			return err
		}
		if err := tun.CreateMangleChain(true); err != nil {
			DisableTun()
			return err
		}
		if err := tun.CreateProxyChain(true); err != nil {
			DisableTun()
			return err
		}
	}
	return nil
}

// DisableTun disable tun
func DisableTun() {
	tun.DeleteRoute(false)
	tun.CleanIptablesChain(false)
	if builds.Config.Proxy.EnableIPv6 {
		tun.DeleteRoute(true)
		tun.CleanIptablesChain(true)
	}
	tun.StopTun()
}
