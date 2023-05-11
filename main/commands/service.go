package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/utils"
	"errors"
)

type ServiceCommand struct {
	Args   []string
	Result int64
}

func (this *ServiceCommand) Execute(args []string) error {
	err := builds.LoadConfig()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("not specify service operation, available operation [start|stop|restart|status]")
	}
	if len(args) > 1 {
		return errors.New("too many service arguments")
	}
	switch args[0] {
	case "start":
		utils.HandleInfo("starting xray service")
		//TODO
	case "stop":
		utils.HandleInfo("stopping xray service")
		//TODO
	case "restart":
		utils.HandleInfo("restarting xray service")
		//TODO
	case "status":
		utils.HandleInfo("xray is running")
		//TODO
	default:
		return errors.New("unknown service operation " + args[0] + ", available operation [start|stop|restart|status]")
	}
	this.Result = 0
	this.Args = args
	return nil
}
