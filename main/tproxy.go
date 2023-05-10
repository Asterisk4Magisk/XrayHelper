package main

import (
	"XrayHelper/main/utils"
	"errors"
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
)

type TproxyCommand struct {
	ConfigPath string `short:"c" long:"config" default:"./config.yml" description:"specify configuration file"`
	Args       []string
	Result     int64
}

func (this *TproxyCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("not specify tproxy operation")
	}
	if len(args) > 1 {
		return errors.New("too many tproxy arguments")
	}
	configFile, err := os.ReadFile(this.ConfigPath)
	if err != nil {
		return err
	}
	if err := defaults.Set(&Config); err != nil {
		return err
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		return err
	}
	fmt.Printf("%+v\n", Config.XrayHelper)
	fmt.Printf("%+v\n", Config.Proxy)
	switch args[0] {
	case "start":
		utils.HandleInfo("starting tproxy")
		//TODO
	case "stop":
		utils.HandleInfo("stopping tproxy")
		//TODO
	case "restart":
		utils.HandleInfo("restarting tproxy")
		//TODO
	case "status":
		utils.HandleInfo("tproxy tproxy")
		//TODO
	default:
		return errors.New("unknown tproxy operation " + args[0])
	}
	this.Result = 0
	this.Args = args
	return nil
}
