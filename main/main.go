package main

import (
	"XrayHelper/main/utils"
	"fmt"
	"github.com/jessevdk/go-flags"
)

var Option struct {
	VerboseFlag bool           `long:"verbose" description:"show verbose debug information"`
	VersionFlag bool           `short:"v" long:"version" description:"show current version"`
	Service     ServiceCommand `command:"service"`
	Tproxy      TproxyCommand  `command:"tproxy"`
}

func main() {
	utils.Verbose = &Option.VerboseFlag
	parser := flags.NewParser(&Option, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.Parse()
	if err != nil {
		if _, ok := err.(*flags.Error); ok {
			typ := err.(*flags.Error).Type
			if typ == flags.ErrCommandRequired {
				if Option.VersionFlag {
					fmt.Println(Version())
				}
				err = nil
			}
			if typ == flags.ErrHelp {
				fmt.Println(VersionStatement())
				fmt.Println(err.Error())
				err = nil
			}
			utils.HandleError(err)
		} else {
			utils.HandleError(err)
		}
	}
}
