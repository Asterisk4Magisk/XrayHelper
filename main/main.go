package main

import (
	"XrayHelper/main/builds"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(builds.VersionStatement())
		fmt.Println(builds.IntroStatement())
		return
	}
	os.Exit(LoadOption())
}
