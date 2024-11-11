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
	BypassSelf       bool   `short:"b" long:"bypass-self" description:"bypass xrayhelper self network traffic (tproxy/tun2socks only)"`
	ConfigFilePath   string `short:"c" long:"config" default:"/data/adb/xray/xrayhelper.yml" description:"specify configuration file"`
	CoreStartTimeout int    `short:"t" long:"core-start-timeout" default:"15" description:"core listen check timeout (second)"`
	VerboseFlag      bool   `short:"v" long:"verbose" description:"show verbose debug information"`
	VersionFlag      bool   `short:"V" long:"version" description:"show current version"`

	Service commands.ServiceCommand `command:"service" description:"control core service"`
	Proxy   commands.ProxyCommand   `command:"proxy" description:"control system proxy"`
	Update  commands.UpdateCommand  `command:"update" description:"update core, adghome, tun2socks, geodata, yacd-meta, metacubexd or subscribe"`
	Switch  commands.SwitchCommand  `command:"switch" description:"switch proxy node or clash config"`
	Api     commands.ApiCommand     `command:"api" description:"xrayhelper api for webui"`
}

// LoadOption load Option, the program entry
func LoadOption() (exitCode int) {
	// if no args, show Intro
	if len(os.Args) == 1 {
		fmt.Println(builds.VersionStatement())
		fmt.Println(builds.IntroStatement())
		return
	}
	log.Verbose = &Option.VerboseFlag
	builds.ConfigFilePath = &Option.ConfigFilePath
	builds.CoreStartTimeout = &Option.CoreStartTimeout
	builds.BypassSelf = &Option.BypassSelf
	parser := flags.NewParser(&Option, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		var flagsError *flags.Error
		if errors.As(err, &flagsError) {
			if errors.Is((*flagsError).Type, flags.ErrCommandRequired) {
				if Option.VersionFlag {
					fmt.Println(builds.Version())
					err = nil
				} else {
					exitCode = 127
				}
			} else if errors.Is((*flagsError).Type, flags.ErrHelp) {
				fmt.Println(builds.VersionStatement())
				fmt.Println(err.Error())
				err = nil
			} else {
				exitCode = 126
			}
			log.HandleError(err)
		} else {
			log.HandleError(err)
			exitCode = 1
		}
	}
	return
}
