package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
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
		return errors.New("service: not specify operation, available operation [start|stop|restart|status]")
	}
	if len(args) > 1 {
		return errors.New("service: too many arguments")
	}
	switch args[0] {
	case "start":
		log.HandleInfo("service: starting xray")
		//TODO
	case "stop":
		log.HandleInfo("service: stopping xray")
		//TODO
	case "restart":
		log.HandleInfo("service: restarting xray")
		//TODO
	case "status":
		log.HandleInfo("service: xray is running")
		//TODO
	default:
		return errors.New("service: unknown operation " + args[0] + ", available operation [start|stop|restart|status]")
	}
	this.Result = 0
	this.Args = args
	return nil
}
