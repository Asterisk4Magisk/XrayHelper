package utils

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

var Verbose *bool

func HandleError(err error) {
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.RedString("ERROR"), ":", err.Error())
	}
}

func HandleInfo(str string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.GreenString("INFO"), ":", str)
}

func HandleDebug(str string) {
	if *Verbose {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"), color.BlueString("DEBUG"), ":", str)
	}
}
