package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/proxys"
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
	if err := proxys.AddRouteTproxy(false); err != nil {
		retErr = err
	}
	if err := proxys.CreateMangleChainTproxy(false); err != nil {
		retErr = err
	}
	if err := proxys.CreateProxyChainTproxy(false); err != nil {
		retErr = err
	}
	if builds.Config.Proxy.EnableIPv6 {
		if err := proxys.AddRouteTproxy(true); err != nil {
			retErr = err
		}
		if err := proxys.CreateMangleChainTproxy(true); err != nil {
			retErr = err
		}
		if err := proxys.CreateProxyChainTproxy(true); err != nil {
			retErr = err
		}
	}
	return retErr
}

// disableTproxy disable proxy(tproxy)
func disableTproxy() {
	proxys.DeleteRouteTproxy(false)
	proxys.CleanIptablesChainTproxy(false)
	if builds.Config.Proxy.EnableIPv6 {
		proxys.DeleteRouteTproxy(true)
		proxys.CleanIptablesChainTproxy(true)
	}
}
