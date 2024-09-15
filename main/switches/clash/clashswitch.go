package clash

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path"
	"strconv"
	"strings"
)

const tagClashswitch = "clashswitch"

var clashUrl []string

type ClashSwitch struct{}

func (this *ClashSwitch) Execute(args []string) (bool, error) {
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return false, e.New("open core config file failed, ", err).WithPrefix(tagClashswitch).WithPathObj(*this)
	} else {
		if !confInfo.IsDir() {
			return false, e.New("clash CoreConfig should be a directory").WithPrefix(tagClashswitch).WithPathObj(*this)
		}
	}
	clashConfig := path.Join(builds.Config.XrayHelper.CoreConfig, "config.yaml")
	if len(args) > 1 {
		return false, e.New("too many arguments").WithPrefix(tagClashswitch).WithPathObj(*this)
	}
	if len(args) == 1 {
		_ = os.Remove(clashConfig)
		if _, err := common.CopyFile(path.Join(builds.Config.XrayHelper.CoreConfig, args[0]), clashConfig); err != nil {
			return false, err
		}
	} else {
		loadClashUrl()
		if len(clashUrl) > 0 {
			for index, clashSubUrl := range clashUrl {
				fmt.Printf(color.GreenString("[%d]")+" %s\n", index, clashSubUrl)
			}
			fmt.Print("Please choose a clash subscribe: ")
			index := 0
			if _, err := fmt.Scanln(&index); err != nil {
				return false, e.New("invalid input, ", err).WithPrefix(tagClashswitch).WithPathObj(*this)
			}
			if err := change(index); err != nil {
				return false, err
			}
		} else {
			return false, e.New("do not have any clash subscribe url in subList").WithPrefix(tagClashswitch).WithPathObj(*this)
		}
	}
	return true, nil
}

func (this *ClashSwitch) Get(custom bool) serial.OrderedArray {
	var result serial.OrderedArray
	loadClashUrl()
	if len(clashUrl) > 0 {
		for _, url := range clashUrl {
			result = append(result, url)
		}
	}
	return result
}

func (this *ClashSwitch) Set(custom bool, index int) error {
	loadClashUrl()
	return change(index)
}

func (this *ClashSwitch) Choose(custom bool, index int) any {
	loadClashUrl()
	if index >= 0 && index < len(clashUrl) {
		return clashUrl[index]
	}
	return nil
}

func loadClashUrl() {
	if len(clashUrl) > 0 {
		return
	}
	for _, subUrl := range builds.Config.XrayHelper.SubList {
		if strings.HasPrefix(subUrl, "clash+") {
			clashUrl = append(clashUrl, strings.TrimPrefix(subUrl, "clash+"))
		}
	}
}

func change(index int) error {
	clashConfig := path.Join(builds.Config.XrayHelper.CoreConfig, "config.yaml")
	if index < 0 || index >= len(builds.Config.XrayHelper.SubList) {
		return e.New("invalid number").WithPrefix(tagClashswitch)
	}
	_ = os.Remove(clashConfig)
	if _, err := common.CopyFile(path.Join(builds.Config.XrayHelper.DataDir, "clashSub"+strconv.Itoa(index)+".yaml"), clashConfig); err != nil {
		return err
	}
	return nil
}
