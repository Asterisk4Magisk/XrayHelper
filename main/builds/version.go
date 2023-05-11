package builds

import (
	"XrayHelper/main/utils"
	"fmt"
	"runtime"
)

var (
	Version_x byte = 0
	Version_y byte = 0
	Version_z byte = 1
	build          = "-debug"
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", Version_x, Version_y, Version_z)
}

func VersionStatement() string {
	return utils.Concat("XrayHelper ", Version(), build, " (", runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, ")")
}
