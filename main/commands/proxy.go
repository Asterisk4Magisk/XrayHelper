package commands

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/proxies"
)

const tagProxy = "proxy"

type ProxyCommand struct{}

func (this *ProxyCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if err := builds.LoadPackage(); err != nil {
		return err
	}
	if len(args) == 0 {
		return e.New("not specify operation, available operation [enable|disable|refresh]").WithPrefix(tagProxy).WithPathObj(*this)
	}
	if len(args) > 1 {
		return e.New("too many arguments").WithPrefix(tagProxy).WithPathObj(*this)
	}
	log.HandleInfo("proxy: current proxy method is " + builds.Config.Proxy.Method)
	proxy, err := proxies.NewProxy(builds.Config.Proxy.Method)
	if err != nil {
		return err
	}
	switch args[0] {
	case "enable":
		log.HandleInfo("proxy: enabling rules")
		if len(getServicePid()) > 0 {
			if err := proxy.Enable(); err != nil {
				return err
			}
		} else {
			log.HandleInfo("proxy: service not running, please check it")
		}
	case "disable":
		log.HandleInfo("proxy: disabling rules")
		proxy.Disable()
	case "refresh":
		log.HandleInfo("proxy: refreshing rules")
		proxy.Disable()
		if len(getServicePid()) > 0 {
			if err := proxy.Enable(); err != nil {
				return err
			}
		} else {
			log.HandleInfo("proxy: service not running, please check it")
		}
	default:
		return e.New("unknown operation " + args[0] + ", available operation [enable|disable|refresh]").WithPrefix(tagProxy).WithPathObj(*this)
	}
	return nil
}
