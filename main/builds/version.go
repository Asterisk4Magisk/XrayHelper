package builds

import (
	"XrayHelper/main/utils"
	"fmt"
	"runtime"
)

var (
	VersionX byte = 0
	VersionY byte = 0
	VersionZ byte = 1
	build         = "-debug"
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", VersionX, VersionY, VersionZ)
}

func VersionStatement() string {
	return utils.Concat("XrayHelper ", Version(), build, " (", runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, ")")
}
