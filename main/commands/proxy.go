package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
)

type ProxyCommand struct {
	Args   []string
	Result int64
}

func (this *ProxyCommand) Execute(args []string) error {
	err := builds.LoadConfig()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("proxy: not specify operation, available operation [enable|disable|refresh]")
	}
	if len(args) > 1 {
		return errors.New("proxy: too many arguments")
	}
	log.HandleInfo("proxy: current method is " + builds.Config.Proxy.Method)
	switch args[0] {
	case "enable":
		log.HandleInfo("proxy: enabling rules")
		//TODO
	case "disable":
		log.HandleInfo("proxy: disabling rules")
		//TODO
	case "refresh":
		log.HandleInfo("proxy: refreshing rules")
		//TODO
	default:
		return errors.New("proxy: unknown operation " + args[0] + ", available operation [enable|disable|refresh]")
	}
	this.Result = 0
	this.Args = args
	return nil
}
