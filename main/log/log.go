package log

import (
	"XrayHelper/main/serial"
	"fmt"
	"github.com/fatih/color"
	"time"
)

var Verbose *bool

// HandleError record error log
func HandleError(v interface{}) {
	if str := serial.ToString(v); str != "" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.RedString("ERROR"), ":", str)
	}
}

// HandleInfo record info log
func HandleInfo(v interface{}) {
	if str := serial.ToString(v); str != "" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.GreenString("INFO"), ":", str)
	}
}

// HandleDebug record debug log
func HandleDebug(v interface{}) {
	if *Verbose {
		if str := serial.ToString(v); str != "" {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.BlueString("DEBUG"), ":", str)
		}
	}
}
