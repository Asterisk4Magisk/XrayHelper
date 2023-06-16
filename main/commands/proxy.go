package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/proxies"
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
	proxy, err := proxies.NewProxy(builds.Config.Proxy.Method)
	if err != nil {
		return err
	}
	switch args[0] {
	case "enable":
		log.HandleInfo("proxy: enabling rules")
		if err := proxy.Enable(); err != nil {
			return err
		}
	case "disable":
		log.HandleInfo("proxy: disabling rules")
		proxy.Disable()
	case "refresh":
		log.HandleInfo("proxy: refreshing rules")
		proxy.Disable()
		if err := proxy.Enable(); err != nil {
			return err
		}
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [enable|disable|refresh]").WithPrefix("proxy").WithPathObj(*this)
	}
	return nil
}
