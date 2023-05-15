package builds

import (
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"bufio"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

const PackageListFilePath = "/data/system/packages.list"

var ConfigFilePath *string
var PackageMap = make(map[string]uint32)

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
		return errors.New("load config failed, ", err).WithPrefix("config")
	}
	if err := defaults.Set(&Config); err != nil {
		return errors.New("set default config failed, ", err).WithPrefix("config")
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		return errors.New("unmarshal config failed, ", err).WithPrefix("config")
	}
	log.HandleDebug(Config.XrayHelper)
	log.HandleDebug(Config.Proxy)
	return nil
}

// LoadPackage load and parse Android package with uid list into a map
func LoadPackage() error {
	packageListFile, err := os.Open(PackageListFilePath)
	if err != nil {
		return errors.New("load package failed, ", err).WithPrefix("config")
	}
	packageScanner := bufio.NewScanner(packageListFile)
	packageScanner.Split(bufio.ScanLines)
	for packageScanner.Scan() {
		packageInfo := strings.Fields(packageScanner.Text())
		if len(packageInfo) >= 2 {
			parseUint, err := strconv.ParseUint(packageInfo[1], 10, 32)
			if err != nil {
				return errors.New("parse package failed, ", err).WithPrefix("config")
			}
			PackageMap[packageInfo[0]] = uint32(parseUint)
		}
	}
	err = packageListFile.Close()
	if err != nil {
		return errors.New("close package file failed, ", err).WithPrefix("config")
	}
	log.HandleDebug(PackageMap)
	return nil
}
