package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/log"
	"XrayHelper/main/switches"
)

type SwitchCommand struct{}

func (this *SwitchCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	switcher, err := switches.NewSwitch(builds.Config.XrayHelper.CoreType)
	if err != nil {
		return err
	}
	success, err := switcher.Execute(args)
	if err != nil {
		return err
	}
	if success {
		log.HandleInfo("switch: switch success")
		// if core is running, restart it
		if len(getServicePid()) > 0 {
			log.HandleInfo("switch: detect core is running, restart it")
			stopService()
			if err := startService(); err != nil {
				log.HandleError("restart service failed, " + err.Error())
			}
		}
	} else {
		log.HandleError("switch: switch failed")
	}
	return nil
}
