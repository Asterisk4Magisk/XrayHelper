package builds

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"bufio"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const packageListPath = "/data/system/packages.list"
const tagConfig = "config"

var ConfigFilePath *string
var CoreStartTimeout *int
var BypassSelf *bool
var PackageMap = make(map[string]string)

// Config the program configuration, yml
var Config struct {
	XrayHelper struct {
		CoreType      string   `default:"xray" yaml:"coreType"`
		CorePath      string   `yaml:"corePath"`
		CoreConfig    string   `yaml:"coreConfig"`
		DataDir       string   `yaml:"dataDir"`
		RunDir        string   `yaml:"runDir"`
		CPULimit      string   `default:"1.0" yaml:"cpuLimit"`
		MemLimit      string   `default:"-1" yaml:"memLimit"`
		ProxyTag      string   `default:"proxy" yaml:"proxyTag"`
		AllowInsecure bool     `default:"false" yaml:"allowInsecure"`
		SubList       []string `yaml:"subList"`
		UserAgent     string   `yaml:"userAgent"`
	} `yaml:"xrayHelper"`
	Clash struct {
		DNSPort  string `default:"65533" yaml:"dnsPort"`
		Template string `yaml:"template"`
	} `yaml:"clash"`
	AdgHome struct {
		Enable  bool   `default:"false" yaml:"enable"`
		Address string `default:"127.0.0.1:65530" yaml:"address"`
		WorkDir string `yaml:"workDir"`
		DNSPort string `default:"65531" yaml:"dnsPort"`
	} `yaml:"adgHome"`
	Proxy struct {
		Method          string   `default:"tproxy" yaml:"method"`
		TproxyPort      string   `default:"65535" yaml:"tproxyPort"`
		SocksPort       string   `default:"65534" yaml:"socksPort"`
		TunDevice       string   `default:"xtun" yaml:"tunDevice"`
		EnableIPv6      bool     `default:"false" yaml:"enableIPv6"`
		AutoDNSStrategy bool     `default:"true" yaml:"autoDNSStrategy"`
		Mode            string   `default:"blacklist" yaml:"mode"`
		PkgList         []string `yaml:"pkgList"`
		ApList          []string `yaml:"apList"`
		IgnoreList      []string `yaml:"ignoreList"`
		IntraList       []string `yaml:"intraList"`
	} `yaml:"proxy"`
}

// LoadConfig load program configuration file, should be called before any command Execute
func LoadConfig() error {
	configFile, err := os.ReadFile(*ConfigFilePath)
	if err != nil {
		return e.New("load config failed, ", err).WithPrefix(tagConfig)
	}
	if err := defaults.Set(&Config); err != nil {
		return e.New("set default config failed, ", err).WithPrefix(tagConfig)
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		return e.New("unmarshal config failed, ", err).WithPrefix(tagConfig)
	}
	log.HandleDebug(Config.XrayHelper)
	log.HandleDebug(Config.Clash)
	log.HandleDebug(Config.AdgHome)
	log.HandleDebug(Config.Proxy)
	return nil
}

// LoadPackage load and parse Android package with uid list into a map
func LoadPackage() error {
	packageListFile, err := os.Open(packageListPath)
	if err != nil {
		return e.New("load package failed, ", err).WithPrefix(tagConfig)
	}
	packageScanner := bufio.NewScanner(packageListFile)
	packageScanner.Split(bufio.ScanLines)
	for packageScanner.Scan() {
		packageInfo := strings.Fields(packageScanner.Text())
		if len(packageInfo) >= 2 {
			PackageMap[packageInfo[0]] = packageInfo[1]
		}
	}
	if err := packageListFile.Close(); err != nil {
		return e.New("close package file failed, ", err).WithPrefix(tagConfig)
	}
	log.HandleDebug(PackageMap)
	return nil
}
