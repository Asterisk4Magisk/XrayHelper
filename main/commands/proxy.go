package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/proxy"
)

type ProxyCommand struct{}

func (this *ProxyCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if err := builds.LoadPackage(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("not specify operation, available operation [enable|disable|refresh]").WithPrefix("proxy").WithPathObj(*this)
	}
	if len(args) > 1 {
		return errors.New("too many arguments").WithPrefix("proxy").WithPathObj(*this)
	}
	log.HandleInfo("proxy: current method is " + builds.Config.Proxy.Method)
	switch args[0] {
	case "enable":
		log.HandleInfo("proxy: enabling rules")
		if builds.Config.Proxy.Method == "tproxy" {
			if err := enableTproxy(); err != nil {
				return err
			}
		} else {
			return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxy").WithPathObj(*this)
		}
	case "disable":
		log.HandleInfo("proxy: disabling rules")
		if builds.Config.Proxy.Method == "tproxy" {
			disableTproxy()
		} else {
			return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxy").WithPathObj(*this)
		}
	case "refresh":
		log.HandleInfo("proxy: refreshing rules")
		if builds.Config.Proxy.Method == "tproxy" {
			disableTproxy()
			if err := enableTproxy(); err != nil {
				return err
			}
		} else {
			return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxy").WithPathObj(*this)
		}
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [enable|disable|refresh]").WithPrefix("proxy").WithPathObj(*this)
	}
	return nil
}

// enableTproxy enable proxy(tproxy)
func enableTproxy() error {
	var retErr error
	defer func() {
		if retErr != nil {
			disableTproxy()
		}
	}()
	if err := proxy.AddRoute(false); err != nil {
		retErr = err
	}
	if err := proxy.CreateMangleChain(false); err != nil {
		retErr = err
	}
	if err := proxy.CreateProxyChain(false); err != nil {
		retErr = err
	}
	if builds.Config.Proxy.EnableIPv6 {
		if err := proxy.AddRoute(true); err != nil {
			retErr = err
		}
		if err := proxy.CreateMangleChain(true); err != nil {
			retErr = err
		}
		if err := proxy.CreateProxyChain(true); err != nil {
			retErr = err
		}
	}
	return retErr
}

// disableTproxy disable proxy(tproxy)
func disableTproxy() {
	proxy.DeleteRoute(false)
	proxy.CleanIptablesChain(false)
	if builds.Config.Proxy.EnableIPv6 {
		proxy.DeleteRoute(true)
		proxy.CleanIptablesChain(true)
	}
}
