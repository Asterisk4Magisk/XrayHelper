package clash

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"fmt"
	"os"
	"path"
	"strconv"
)

type ClashSwitch struct{}

func (this *ClashSwitch) Execute(args []string) (bool, error) {
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return false, errors.New("open core config file failed, ", err).WithPrefix("clashswitch").WithPathObj(*this)
	} else {
		if !confInfo.IsDir() {
			return false, errors.New("clash CoreConfig should be a directory").WithPrefix("clashswitch").WithPathObj(*this)
		}
	}
	clashConfig := path.Join(builds.Config.XrayHelper.CoreConfig, "config.yaml")
	if len(args) > 1 {
		return false, errors.New("too many arguments").WithPrefix("clashswitch").WithPathObj(*this)
	}
	if len(args) == 1 && args[0] == "custom" {
		_ = os.Remove(clashConfig)
		if _, err := common.CopyFile(path.Join(builds.Config.XrayHelper.DataDir, "clashCustom.yaml"), clashConfig); err != nil {
			return false, err
		}
	} else {
		for index, clashSubUrl := range builds.Config.XrayHelper.SubList {
			fmt.Printf("[%d] %s\n", index, clashSubUrl)
		}
		fmt.Print("Please choose a clash subscribe: ")
		index := 0
		_, err := fmt.Scanln(&index)
		if err != nil {
			return false, errors.New("invalid input, ", err).WithPrefix("clashswitch").WithPathObj(*this)
		}
		if index < 0 || index >= len(builds.Config.XrayHelper.SubList) {
			return false, errors.New("invalid node number").WithPrefix("clashswitch").WithPathObj(*this)
		}
		_ = os.Remove(clashConfig)
		if _, err := common.CopyFile(path.Join(builds.Config.XrayHelper.DataDir, "clashSub"+strconv.Itoa(index)+".yaml"), clashConfig); err != nil {
			return false, err
		}
	}
	return true, nil
}
