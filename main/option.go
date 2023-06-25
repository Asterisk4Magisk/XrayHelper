package main

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/commands"
	"XrayHelper/main/log"
	"fmt"
	"github.com/jessevdk/go-flags"
)

var Option struct {
	ConfigFilePath string                  `short:"c" long:"config" default:"/data/adb/xray/xrayhelper.yml" description:"specify configuration file"`
	VerboseFlag    bool                    `short:"v" long:"verbose" description:"show verbose debug information"`
	VersionFlag    bool                    `short:"V" long:"version" description:"show current version"`
	Service        commands.ServiceCommand `command:"service" description:"control core service"`
	Proxy          commands.ProxyCommand   `command:"proxy" description:"control system proxy"`
	Update         commands.UpdateCommand  `command:"update" description:"update core, tun2socks, geodata, yacd or subscribe"`
	Switch         commands.SwitchCommand  `command:"switch" description:"switch proxy node or clash config"`
}

// LoadOption load Option, the program entry
func LoadOption() int {
	rCode := 0
	log.Verbose = &Option.VerboseFlag
	builds.ConfigFilePath = &Option.ConfigFilePath
	parser := flags.NewParser(&Option, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			typ := err.(*flags.Error).Type
			if typ == flags.ErrCommandRequired {
				if Option.VersionFlag {
					fmt.Println(builds.Version())
					err = nil
				} else {
					rCode = 127
				}
			} else if typ == flags.ErrHelp {
				fmt.Println(builds.VersionStatement())
				fmt.Println(err.Error())
				err = nil
			} else {
				rCode = 126
			}
			log.HandleError(err)
		} else {
			log.HandleError(err)
			rCode = 1
		}
	}
	return rCode
}
