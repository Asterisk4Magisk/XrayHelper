package builds

import (
	"XrayHelper/main/utils"
	"fmt"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
)

var ConfigFilePath *string

var Config struct {
	XrayHelper struct {
		Busybox       string `yaml:"busybox"`
		Xray          string `yaml:"xray"`
		XrayConfigDir string `yaml:"xrayConfigDir"`
		RunDir        string `yaml:"runDir"`
	} `yaml:"xrayHelper"`
	Proxy struct {
		Method     string   `default:"tproxy" yaml:"method"`
		EnableIPv6 bool     `default:"true" yaml:"enableIPv6"`
		Mode       string   `default:"blacklist" yaml:"mode"`
		UidList    []uint16 `yaml:"uidList"`
		ApList     []string `yaml:"apList"`
	} `yaml:"proxy"`
}

func LoadConfig() error {
	configFile, err := os.ReadFile(*ConfigFilePath)
	if err != nil {
		utils.HandleError(err)
		return err
	}
	if err := defaults.Set(&Config); err != nil {
		utils.HandleError(err)
		return err
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		utils.HandleError(err)
		return err
	}
	fmt.Printf("%+v\n", Config.XrayHelper)
	fmt.Printf("%+v\n", Config.Proxy)
	return nil
}
