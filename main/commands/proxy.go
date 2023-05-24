package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/tproxy"
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

// enableTproxy enable tproxy
func enableTproxy() error {
	var retErr error
	defer func() {
		if retErr != nil {
			disableTproxy()
		}
	}()
	if err := tproxy.AddRoute(false); err != nil {
		retErr = err
	}
	if err := tproxy.CreateMangleChain(false); err != nil {
		retErr = err
	}
	if err := tproxy.CreateProxyChain(false); err != nil {
		retErr = err
	}
	if builds.Config.Proxy.EnableIPv6 {
		if err := tproxy.AddRoute(true); err != nil {
			retErr = err
		}
		if err := tproxy.CreateMangleChain(true); err != nil {
			retErr = err
		}
		if err := tproxy.CreateProxyChain(true); err != nil {
			retErr = err
		}
	}
	return retErr
}

// disableTproxy disable tproxy
func disableTproxy() {
	tproxy.DeleteRoute(false)
	tproxy.CleanIptablesChain(false)
	if builds.Config.Proxy.EnableIPv6 {
		tproxy.DeleteRoute(true)
		tproxy.CleanIptablesChain(true)
	}
}
