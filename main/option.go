package main

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/commands"
	"XrayHelper/main/utils"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var Option struct {
	ConfigFilePath string                  `short:"c" long:"config" default:"./config.yml" description:"specify configuration file"`
	VerboseFlag    bool                    `short:"v" long:"verbose" description:"show verbose debug information"`
	VersionFlag    bool                    `short:"V" long:"version" description:"show current version"`
	Service        commands.ServiceCommand `command:"service" description:"control xray service"`
	Proxy          commands.ProxyCommand   `command:"proxy" description:"control system proxy"`
}

func LoadOption() {
	if len(os.Args) == 1 {
		fmt.Println(builds.VersionStatement())
		fmt.Println(builds.IntroStatement())
		return
	}
	utils.Verbose = &Option.VerboseFlag
	builds.ConfigFilePath = &Option.ConfigFilePath
	parser := flags.NewParser(&Option, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			typ := err.(*flags.Error).Type
			if typ == flags.ErrCommandRequired {
				if Option.VersionFlag {
					fmt.Println(builds.Version())
				}
				err = nil
			}
			if typ == flags.ErrHelp {
				fmt.Println(builds.VersionStatement())
				fmt.Println(err.Error())
				err = nil
			}
			utils.HandleError(err)
		} else {
			utils.HandleError(err)
		}
	}
}
