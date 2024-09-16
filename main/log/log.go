package log

import (
	"XrayHelper/main/serial"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
	"strings"
	"time"
)

var Verbose *bool

func init() {
	out, err := exec.Command("/system/bin/getprop", "persist.sys.timezone").Output()
	if err != nil {
		return
	}
	z, err := time.LoadLocation(strings.TrimSpace(string(out)))
	if err != nil {
		return
	}
	time.Local = z
}

// HandleError record error log
func HandleError(v any) {
	if str := serial.ToString(v); str != "" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.RedString("ERROR"), ":", str)
	}
}

// HandleInfo record info log
func HandleInfo(v any) {
	if str := serial.ToString(v); str != "" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.GreenString("INFO"), ":", str)
	}
}

// HandleDebug record debug log
func HandleDebug(v any) {
	if *Verbose {
		if str := serial.ToString(v); str != "" {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.BlueString("DEBUG"), ":", str)
		}
	}
}
