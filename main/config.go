package main

type Config struct {
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
