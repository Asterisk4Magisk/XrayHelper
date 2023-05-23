package builds

import (
	"XrayHelper/main/serial"
	"fmt"
	"runtime"
)

const (
	VersionX byte = 0
	VersionY byte = 0
	VersionZ byte = 1
	Build         = "-release"
	Intro         = "An xray helper for Android to control system proxy."
)

func Version() string {
	return fmt.Sprintf("%v.%v.%v", VersionX, VersionY, VersionZ)
}

func VersionStatement() string {
	return serial.Concat("XrayHelper ", Version(), Build, " (", runtime.Version(), " ", runtime.GOOS, "/", runtime.GOARCH, ")")
}

func IntroStatement() string {
	return Intro
}
