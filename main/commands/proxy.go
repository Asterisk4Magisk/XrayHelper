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
		switch builds.Config.Proxy.Method {
		case "tproxy":
			if err := proxys.EnableTproxy(); err != nil {
				return err
			}
		case "tun":
			if err := proxys.EnableTun(); err != nil {
				return err
			}
		default:
			return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxy").WithPathObj(*this)
		}
	case "disable":
		log.HandleInfo("proxy: disabling rules")
		switch builds.Config.Proxy.Method {
		case "tproxy":
			proxys.DisableTproxy()
		case "tun":
			proxys.DisableTun()
		default:
			return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxy").WithPathObj(*this)
		}
	case "refresh":
		log.HandleInfo("proxy: refreshing rules")
		switch builds.Config.Proxy.Method {
		case "tproxy":
			proxys.DisableTproxy()
			if err := proxys.EnableTproxy(); err != nil {
				return err
			}
		case "tun":
			proxys.DisableTun()
			if err := proxys.EnableTun(); err != nil {
				return err
			}
		default:
			return errors.New("invalid proxy method " + builds.Config.Proxy.Method).WithPrefix("proxy").WithPathObj(*this)
		}
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [enable|disable|refresh]").WithPrefix("proxy").WithPathObj(*this)
	}
	return nil
}
