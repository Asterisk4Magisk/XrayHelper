package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/utils"
	"errors"
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
		return errors.New("not specify proxy operation, available operation [enable|disable|refresh]")
	}
	if len(args) > 1 {
		return errors.New("too many proxy arguments")
	}
	utils.HandleInfo("current proxy method is " + builds.Config.Proxy.Method)
	switch args[0] {
	case "enable":
		utils.HandleInfo("enabling proxy rule")
		//TODO
	case "disable":
		utils.HandleInfo("disabling proxy rule")
		//TODO
	case "refresh":
		utils.HandleInfo("refreshing proxy rule")
		//TODO
	default:
		return errors.New("unknown proxy operation " + args[0] + ", available operation [enable|disable|refresh]")
	}
	this.Result = 0
	this.Args = args
	return nil
}
