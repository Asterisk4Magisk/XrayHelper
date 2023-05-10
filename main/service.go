package main

import (
	"XrayHelper/main/utils"
	"errors"
)

type ServiceCommand struct {
	Args   []string
	Result int64
}

func (this *ServiceCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("not specify service operation")
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
		return errors.New("unknown service operation " + args[0])
	}
	this.Result = 0
	this.Args = args
	return nil
}
