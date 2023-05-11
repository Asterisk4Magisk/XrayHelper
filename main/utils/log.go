package utils

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

var Verbose *bool

func HandleError(v interface{}) {
	str := ToString(v)
	if str != "" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.RedString("ERROR"), ":", str)
	}
}

func HandleInfo(v interface{}) {
	str := ToString(v)
	if str != "" {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.GreenString("INFO"), ":", str)
	}
}

func HandleDebug(v interface{}) {
	if *Verbose {
		str := ToString(v)
		if str != "" {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.BlueString("DEBUG"), ":", str)
		}
	}
}
