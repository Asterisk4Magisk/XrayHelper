package builds

import (
	"XrayHelper/main/serial"
	"fmt"
	"runtime"
)

const (
	VersionX byte = 1
	VersionY byte = 3
	VersionZ byte = 3
	Build         = "-release"
	Intro         = "A unified helper for Android to control system proxy.\n\nTelegram channel: https://t.me/Asterisk4Magisk\nTelegram chat: https://t.me/AsteriskFactory\n\nReport issues at https://github.com/Asterisk4Magisk/XrayHelper/issues\n"
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
