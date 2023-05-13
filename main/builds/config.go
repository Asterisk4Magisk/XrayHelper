package builds

import (
	"XrayHelper/main/log"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
)

var ConfigFilePath *string

// Config the program configuration, yml
var Config struct {
	XrayHelper struct {
		Xray          string `yaml:"xray"`
		XrayConfigDir string `yaml:"xrayConfigDir"`
		RunDir        string `yaml:"runDir"`
	} `yaml:"xrayHelper"`
	Proxy struct {
		Method     string   `default:"tproxy" yaml:"method"`
		EnableIPv6 bool     `default:"true" yaml:"enableIPv6"`
		Mode       string   `default:"blacklist" yaml:"mode"`
		PkgList    []string `yaml:"pkgList"`
		ApList     []string `yaml:"apList"`
	} `yaml:"proxy"`
}

// LoadConfig load program configuration file, should be called before any command Execute
func LoadConfig() error {
	configFile, err := os.ReadFile(*ConfigFilePath)
	if err != nil {
		return err
	}
	if err := defaults.Set(&Config); err != nil {
		return err
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		return err
	}
	log.HandleDebug(Config.XrayHelper)
	log.HandleDebug(Config.Proxy)
	return nil
}
