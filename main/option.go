package main

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/commands"
	"XrayHelper/main/log"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var Option struct {
	BypassSelf       bool   `short:"b" long:"bypass-self" description:"bypass xrayhelper self network traffic (tproxy only)"`
	ConfigFilePath   string `short:"c" long:"config" default:"/data/adb/xray/xrayhelper.yml" description:"specify configuration file"`
	CoreStartTimeout int    `short:"t" long:"core-start-timeout" default:"15" description:"core listen check timeout (second)"`
	VerboseFlag      bool   `short:"v" long:"verbose" description:"show verbose debug information"`
	VersionFlag      bool   `short:"V" long:"version" description:"show current version"`

	Service commands.ServiceCommand `command:"service" description:"control core service"`
	Proxy   commands.ProxyCommand   `command:"proxy" description:"control system proxy"`
	Update  commands.UpdateCommand  `command:"update" description:"update core, tun2socks, geodata, yacd-meta or subscribe"`
	Switch  commands.SwitchCommand  `command:"switch" description:"switch proxy node or clash config"`
}

// LoadOption load Option, the program entry
func LoadOption() int {
	// if no args, show Intro
	if len(os.Args) == 1 {
		fmt.Println(builds.VersionStatement())
		fmt.Println(builds.IntroStatement())
		return 0
	}
	log.Verbose = &Option.VerboseFlag
	builds.ConfigFilePath = &Option.ConfigFilePath
	builds.CoreStartTimeout = &Option.CoreStartTimeout
	builds.BypassSelf = &Option.BypassSelf
	rCode := 0
	parser := flags.NewParser(&Option, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		var flagsError *flags.Error
		if errors.As(err, &flagsError) {
			if errors.Is((*flagsError).Type, flags.ErrCommandRequired) {
				if Option.VersionFlag {
					fmt.Println(builds.Version())
					err = nil
				} else {
					rCode = 127
				}
			} else if errors.Is((*flagsError).Type, flags.ErrHelp) {
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
