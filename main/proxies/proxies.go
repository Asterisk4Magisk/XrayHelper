package proxies

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"XrayHelper/main/proxies/tproxy"
	"XrayHelper/main/proxies/tun"
)

func Enable() error {
	switch builds.Config.Proxy.Method {
	case "tproxy":
		if err := tproxy.AddRoute(false); err != nil {
			Disable()
			return err
		}
		if err := tproxy.CreateMangleChain(false); err != nil {
			Disable()
			return err
		}
		if err := tproxy.CreateProxyChain(false); err != nil {
			Disable()
			return err
		}
		if builds.Config.Proxy.EnableIPv6 {
			if err := tproxy.AddRoute(true); err != nil {
				Disable()
				return err
			}
			if err := tproxy.CreateMangleChain(true); err != nil {
				Disable()
				return err
			}
			if err := tproxy.CreateProxyChain(true); err != nil {
				Disable()
				return err
			}
		}
	case "tun":
		if err := tun.StartTun(); err != nil {
			Disable()
			return err
		}
		if err := tun.AddRoute(false); err != nil {
			Disable()
			return err
		}
		if err := tun.CreateMangleChain(false); err != nil {
			Disable()
			return err
		}
		if err := tun.CreateProxyChain(false); err != nil {
			Disable()
			return err
		}
		if builds.Config.Proxy.EnableIPv6 {
			if err := tun.AddRoute(true); err != nil {
				Disable()
				return err
			}
			if err := tun.CreateMangleChain(true); err != nil {
				Disable()
				return err
			}
			if err := tun.CreateProxyChain(true); err != nil {
				Disable()
				return err
			}
		}
	default:
		return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxies")
	}
	if err := handleDns(); err != nil {
		Disable()
		return err
	}
	return nil
}

func Disable() {
	switch builds.Config.Proxy.Method {
	case "tproxy":
		tproxy.DeleteRoute(false)
		tproxy.CleanIptablesChain(false)
		if builds.Config.Proxy.EnableIPv6 {
			tproxy.DeleteRoute(true)
			tproxy.CleanIptablesChain(true)
		}
	case "tun":
		tun.DeleteRoute(false)
		tun.CleanIptablesChain(false)
		if builds.Config.Proxy.EnableIPv6 {
			tun.DeleteRoute(true)
			tun.CleanIptablesChain(true)
		}
		tun.StopTun()
	}
	dehandleDns()
}

func handleDns() error {
	if !builds.Config.Proxy.EnableIPv6 {
		if err := common.Ipt6.Insert("filter", "OUTPUT", 1, "-p", "udp", "--dport", "53", "-j", "REJECT"); err != nil {
			return errors.New("disable dns request on ipv6 failed, ", err).WithPrefix("proxies")
		}
	}
	return nil
}

func dehandleDns() {
	if !builds.Config.Proxy.EnableIPv6 {
		_ = common.Ipt6.Delete("filter", "OUTPUT", "-p", "udp", "--dport", "53", "-j", "REJECT")
	}
}
