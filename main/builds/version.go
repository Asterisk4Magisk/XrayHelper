package builds

import (
	"XrayHelper/main/utils"
	"fmt"
	"runtime"
)

const (
	VersionX byte = 0
	VersionY byte = 0
	VersionZ byte = 1
	Build         = "-debug"
	Intro         = "An xray helper for Android to control system proxy."
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", VersionX, VersionY, VersionZ)
}

func VersionStatement() string {
	return utils.Concat("XrayHelper ", Version(), Build, " (", runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, ")")
}

func IntroStatement() string {
	return Intro
}
