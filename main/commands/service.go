package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"encoding/json"
	"github.com/tailscale/hujson"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strconv"
	"time"
)

var service common.External

type ServiceCommand struct{}

func (this *ServiceCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.New("not specify operation, available operation [start|stop|restart|status]").WithPrefix("service").WithPathObj(*this)
	}
	if len(args) > 1 {
		return errors.New("too many arguments").WithPrefix("service").WithPathObj(*this)
	}
	switch args[0] {
	case "start":
		log.HandleInfo("service: starting core")
		if err := startService(); err != nil {
			return err
		}
		log.HandleInfo("service: core is running, pid is " + getServicePid())
	case "stop":
		log.HandleInfo("service: stopping core")
		stopService()
		log.HandleInfo("service: core is stopped")
	case "restart":
		log.HandleInfo("service: restarting core")
		stopService()
		if err := startService(); err != nil {
			return err
		}
		log.HandleInfo("service: core is running, pid is " + getServicePid())
	case "status":
		pidStr := getServicePid()
		if len(pidStr) > 0 {
			log.HandleInfo("service: core is running, pid is " + pidStr)
		} else {
			log.HandleInfo("service: core is stopped")
		}
	default:
		return errors.New("unknown operation " + args[0] + ", available operation [start|stop|restart|status]").WithPrefix("service").WithPathObj(*this)
	}
	return nil
}

// startService start core service
func startService() error {
	listenFlag := false
	servicePid := getServicePid()
	if len(servicePid) > 0 {
		return errors.New("core is running, pid is " + servicePid).WithPrefix("service")
	}
	serviceLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "error.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("open core log file failed, ", err).WithPrefix("service")
	}
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return errors.New("open core config file failed, ", err).WithPrefix("service")
	} else {
		if confInfo.IsDir() {
			switch builds.Config.XrayHelper.CoreType {
			case "xray":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-confdir", builds.Config.XrayHelper.CoreConfig)
			case "v2ray":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-confdir", builds.Config.XrayHelper.CoreConfig, "-format", "jsonv5")
			case "sing-box":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-C", builds.Config.XrayHelper.CoreConfig, "-D", builds.Config.XrayHelper.DataDir, "--disable-color")
			case "clash":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "-d", builds.Config.XrayHelper.CoreConfig)
			default:
				return errors.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix("service")
			}
		} else {
			switch builds.Config.XrayHelper.CoreType {
			case "xray":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig)
			case "v2ray":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig, "-format", "jsonv5")
			case "sing-box":
				service = common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig, "-D", builds.Config.XrayHelper.DataDir, "--disable-color")
			case "clash":
				return errors.New("clash CoreConfig should be a directory").WithPrefix("service")
			default:
				return errors.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix("service")
			}
		}
	}
	switch builds.Config.XrayHelper.CoreType {
	case "xray", "v2ray", "sing-box":
		service.AppendEnv("XRAY_LOCATION_ASSET=" + builds.Config.XrayHelper.DataDir)
		service.AppendEnv("V2RAY_LOCATION_ASSET=" + builds.Config.XrayHelper.DataDir)
		if err := handleRayDNS(builds.Config.Proxy.EnableIPv6); err != nil {
			return err
		}
	case "clash":
		if err := overrideClashConfig(builds.Config.Clash.Template, path.Join(builds.Config.XrayHelper.CoreConfig, "config.yaml")); err != nil {
			return err
		}
	}
	if err := service.SetUidGid("0", common.CoreGid); err != nil {
		return err
	}
	service.Start()
	for i := 0; i < 15; i++ {
		time.Sleep(1 * time.Second)
		if builds.Config.Proxy.Method == "tproxy" {
			if common.CheckPort("tcp", "127.0.0.1", builds.Config.Proxy.TproxyPort) {
				listenFlag = true
				break
			}
		} else if builds.Config.Proxy.Method == "tun" {
			if common.CheckPort("tcp", "127.0.0.1", builds.Config.Proxy.SocksPort) {
				listenFlag = true
				break
			}
		} else {
			listenFlag = false
			break
		}
	}
	if listenFlag {
		if err := os.WriteFile(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"), []byte(strconv.Itoa(service.Pid())), 0644); err != nil {
			_ = service.Kill()
			return errors.New("write core pid failed, ", err).WithPrefix("service")
		}
	} else {
		_ = service.Kill()
		return errors.New("start core service failed, ", service.Err()).WithPrefix("service")
	}
	return nil
}

// stopService stop core service
func stopService() {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "core.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		pid, _ := strconv.Atoi(string(pidFile))
		if serviceProcess, err := os.FindProcess(pid); err == nil {
			_ = serviceProcess.Kill()
			_ = os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"))
		} else {
			log.HandleDebug(err)
		}
	} else {
		log.HandleDebug(err)
	}
}

// getServicePid get core pid from pid file
func getServicePid() string {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "core.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		return string(pidFile)
	} else {
		log.HandleDebug(err)
	}
	return ""
}

func handleRayDNS(ipv6 bool) error {
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return errors.New("open core config file failed, ", err).WithPrefix("service")
	} else {
		if confInfo.IsDir() {
			confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return errors.New("open config dir failed, ", err).WithPrefix("service")
			}
			for _, conf := range confDir {
				if !conf.IsDir() {
					confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()))
					if err != nil {
						return errors.New("read config file failed, ", err).WithPrefix("service")
					}
					newConfByte, err := replaceRayDNSStrategy(confByte, ipv6)
					if err != nil {
						log.HandleDebug(err)
						continue
					}
					if err := os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), newConfByte, 0644); err != nil {
						return errors.New("write new config failed, ", err).WithPrefix("service")
					}
				}
			}
		} else {
			confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return errors.New("read config file failed, ", err).WithPrefix("service")
			}
			newConfByte, err := replaceRayDNSStrategy(confByte, ipv6)
			if err != nil {
				return err
			}
			if err := os.WriteFile(builds.Config.XrayHelper.CoreConfig, newConfByte, 0644); err != nil {
				return errors.New("write new config failed, ", err).WithPrefix("service")
			}
		}
	}
	return nil
}

func replaceRayDNSStrategy(conf []byte, ipv6 bool) (replacedConf []byte, err error) {
	// standardize origin json (remove comment)
	standardize, err := hujson.Standardize(conf)
	if err != nil {
		return nil, errors.New("standardize config json failed, ", err).WithPrefix("service")
	}
	// unmarshal
	var jsonValue interface{}
	err = json.Unmarshal(standardize, &jsonValue)
	if err != nil {
		return nil, errors.New("unmarshal config json failed, ", err).WithPrefix("service")
	}
	// assert json to map
	jsonMap, ok := jsonValue.(map[string]interface{})
	if !ok {
		return nil, errors.New("assert config json to map failed").WithPrefix("service")
	}
	dns, ok := jsonMap["dns"]
	if !ok {
		return nil, errors.New("cannot find dns").WithPrefix("service")
	}
	// assert dns
	dnsMap, ok := dns.(map[string]interface{})
	if !ok {
		return nil, errors.New("assert dns to map failed").WithPrefix("service")
	}
	switch builds.Config.XrayHelper.CoreType {
	case "xray":
		if ipv6 {
			dnsMap["queryStrategy"] = "UseIP"
		} else {
			dnsMap["queryStrategy"] = "UseIPv4"
		}
	case "v2ray":
		if ipv6 {
			dnsMap["queryStrategy"] = "USE_IP"
		} else {
			dnsMap["queryStrategy"] = "USE_IP4"
		}
	case "sing-box":
		if ipv6 {
			dnsMap["strategy"] = "prefer_ipv4"
		} else {
			dnsMap["strategy"] = "ipv4_only"
		}
	default:
		return nil, errors.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix("service")
	}
	// replace
	jsonMap["dns"] = dnsMap
	// marshal
	marshal, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		return nil, errors.New("marshal config json failed, ", err).WithPrefix("service")
	}
	return marshal, nil
}

func overrideClashConfig(template string, target string) error {
	// open target config and replace with xrayhelper clash value
	targetFile, err := os.ReadFile(target)
	if err != nil {
		return errors.New("load clash config failed, ", err).WithPrefix("service")
	}
	var targetYamlValue interface{}
	if err := yaml.Unmarshal(targetFile, &targetYamlValue); err != nil {
		return errors.New("unmarshal clash config failed, ", err).WithPrefix("service")
	}
	targetYamlMap, ok := targetYamlValue.(map[string]interface{})
	if !ok {
		return errors.New("assert clash config to map failed").WithPrefix("service")
	}
	// delete useless config
	delete(targetYamlMap, "port")
	delete(targetYamlMap, "redir-port")
	delete(targetYamlMap, "mixed-port")
	delete(targetYamlMap, "authentication")
	delete(targetYamlMap, "external-controller")
	delete(targetYamlMap, "external-ui")
	delete(targetYamlMap, "secret")
	// set xrayhelper necessary parameters value
	targetYamlMap["tproxy-port"] = builds.Config.Proxy.TproxyPort
	targetYamlMap["socks-port"] = builds.Config.Proxy.SocksPort
	// open template config and replace target value with it
	templateFile, err := os.ReadFile(template)
	if err != nil {
		return errors.New("load clash template config failed, ", err).WithPrefix("service")
	}
	var templateYamlValue interface{}
	if err := yaml.Unmarshal(templateFile, &templateYamlValue); err != nil {
		return errors.New("unmarshal clash template config failed, ", err).WithPrefix("service")
	}
	templateYamlMap, ok := templateYamlValue.(map[string]interface{})
	if !ok {
		return errors.New("assert clash template config to map failed").WithPrefix("service")
	}
	// replace
	for key, value := range templateYamlMap {
		targetYamlMap[key] = value
	}
	// marshal
	marshal, err := yaml.Marshal(targetYamlMap)
	if err != nil {
		return errors.New("marshal clash config failed, ", err).WithPrefix("service")
	}
	// write new config
	if err := os.WriteFile(target, marshal, 0644); err != nil {
		return errors.New("write overridden clash config failed, ", err).WithPrefix("service")
	}
	return nil
}
