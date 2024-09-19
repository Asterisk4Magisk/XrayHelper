package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/cgroup"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/serial"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const tagService = "service"

type ServiceCommand struct{}

func (this *ServiceCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) == 0 {
		return e.New("not specify operation, available operation [start|stop|restart|status]").WithPrefix(tagService).WithPathObj(*this)
	}
	if len(args) > 1 {
		return e.New("too many arguments").WithPrefix(tagService).WithPathObj(*this)
	}
	log.HandleInfo("service: current core type is " + builds.Config.XrayHelper.CoreType)
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
		if err := restartService(); err != nil {
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
		return e.New("unknown operation " + args[0] + ", available operation [start|stop|restart|status]").WithPrefix(tagService).WithPathObj(*this)
	}
	return nil
}

// newServices get core service
func newServices(serviceLogFile *os.File) (service common.External, err error) {
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return nil, e.New("open core config file failed, ", err).WithPrefix(tagService)
	} else {
		if confInfo.IsDir() {
			switch builds.Config.XrayHelper.CoreType {
			case "xray":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-confdir", builds.Config.XrayHelper.CoreConfig), nil
			case "v2ray":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-confdir", builds.Config.XrayHelper.CoreConfig, "-format", "jsonv5"), nil
			case "sing-box":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-C", builds.Config.XrayHelper.CoreConfig, "-D", builds.Config.XrayHelper.DataDir, "--disable-color"), nil
			case "mihomo":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "-d", builds.Config.XrayHelper.CoreConfig), nil
			case "hysteria2":
				return nil, e.New("hysteria2 CoreConfig should be a file").WithPrefix(tagService)
			default:
				return nil, e.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix(tagService)
			}
		} else {
			switch builds.Config.XrayHelper.CoreType {
			case "xray":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig), nil
			case "v2ray":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig, "-format", "jsonv5"), nil
			case "sing-box":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig, "-D", builds.Config.XrayHelper.DataDir, "--disable-color"), nil
			case "mihomo":
				return nil, e.New("mihomo CoreConfig should be a directory").WithPrefix(tagService)
			case "hysteria2":
				return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "-c", builds.Config.XrayHelper.CoreConfig), nil
			default:
				return nil, e.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix(tagService)
			}
		}
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

// startService start core service
func startService() error {
	listenFlag := false
	// check current core status
	servicePid := getServicePid()
	if len(servicePid) > 0 {
		return e.New("core is running, pid is " + servicePid).WithPrefix(tagService)
	}
	// start adguardhome before core, maybe core use adguardhome as upstream
	if builds.Config.AdgHome.Enable {
		if err := startAdgHome(); err != nil {
			stopService()
			return err
		}
	}
	// get core service log file
	serviceLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "error.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return e.New("open core log file failed, ", err).WithPrefix(tagService)
	}
	// get core service
	service, err := newServices(serviceLogFile)
	if err != nil {
		return err
	}
	// add core env variable and prepare core config
	switch builds.Config.XrayHelper.CoreType {
	case "xray", "v2ray", "sing-box":
		service.AppendEnv("XRAY_LOCATION_ASSET=" + builds.Config.XrayHelper.DataDir)
		service.AppendEnv("V2RAY_LOCATION_ASSET=" + builds.Config.XrayHelper.DataDir)
		// if enable AutoDNSStrategy
		if builds.Config.Proxy.AutoDNSStrategy {
			if err := handleRayDNS(builds.Config.Proxy.EnableIPv6); err != nil {
				return err
			}
		}
	case "mihomo":
		if err := overrideClashConfig(builds.Config.Clash.Template, path.Join(builds.Config.XrayHelper.CoreConfig, "config.yaml")); err != nil {
			return err
		}
	case "hysteria2":
		service.AppendEnv("HYSTERIA_DISABLE_UPDATE_CHECK=1")
	}
	service.SetUidGid("0", common.CoreGid)
	ignoreSignals()
	service.Start()
	if service.Err() != nil {
		return e.New("start core service failed, ", service.Err()).WithPrefix(tagService)
	}
	switch builds.Config.Proxy.Method {
	case "tproxy":
		listenFlag = common.CheckLocalPort(strconv.Itoa(service.Pid()), builds.Config.Proxy.TproxyPort, time.Duration(*builds.CoreStartTimeout)*time.Second)
	case "tun":
		// tun don't need check any local port
		listenFlag = true
	case "tun2socks":
		listenFlag = common.CheckLocalPort(strconv.Itoa(service.Pid()), builds.Config.Proxy.SocksPort, time.Duration(*builds.CoreStartTimeout)*time.Second)
	default:
		listenFlag = false
	}
	if listenFlag {
		if err := cgroup.LimitProcess(service.Pid()); err != nil {
			_ = service.Kill()
			stopService()
			return err
		}
		if err := os.WriteFile(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"), []byte(strconv.Itoa(service.Pid())), 0644); err != nil {
			_ = service.Kill()
			stopService()
			return e.New("write core pid failed, ", err).WithPrefix(tagService)
		}
	} else {
		_ = service.Kill()
		stopService()
		return e.New("core service not listen, please check error.log").WithPrefix(tagService)
	}
	return nil
}

// stopService stop core service
func stopService() {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "core.pid")); err == nil {
		pidStr := getServicePid()
		if len(pidStr) > 0 {
			pid, _ := strconv.Atoi(pidStr)
			if serviceProcess, err := os.FindProcess(pid); err == nil {
				_ = serviceProcess.Kill()
				_ = os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"))
			} else {
				log.HandleDebug(err)
			}
		}
	} else {
		log.HandleDebug(err)
	}
	stopAdgHome()
}

// restartService restart core service
func restartService() error {
	stopService()
	time.Sleep(1 * time.Second)
	return startService()
}

func startAdgHome() error {
	adgHomePath := path.Join(path.Dir(builds.Config.XrayHelper.CorePath), "adguardhome")
	adgHomeConfigPath := path.Join(builds.Config.AdgHome.WorkDir, "config.yaml")
	adgHomeConfigFile, err := os.ReadFile(adgHomeConfigPath)
	if err != nil {
		return e.New("load adgHome config failed, ", err).WithPrefix(tagService)
	}
	var adgHomeConfig serial.OrderedMap
	if err := yaml.Unmarshal(adgHomeConfigFile, &adgHomeConfig); err != nil {
		return e.New("unmarshal adgHome config failed, ", err).WithPrefix(tagService)
	}
	if http, ok := adgHomeConfig.Get("http"); ok {
		httpMap := http.Value.(serial.OrderedMap)
		// set address
		httpMap.Set("address", builds.Config.AdgHome.Address)
		adgHomeConfig.Set("http", httpMap)
	}
	if dns, ok := adgHomeConfig.Get("dns"); ok {
		dnsMap := dns.Value.(serial.OrderedMap)
		// set dnsPort
		port, _ := strconv.Atoi(builds.Config.AdgHome.DNSPort)
		dnsMap.Set("port", port)
		// set dnsStrategy
		if builds.Config.Proxy.AutoDNSStrategy {
			dnsMap.Set("aaaa_disabled", !builds.Config.Proxy.EnableIPv6)
		}
		adgHomeConfig.Set("dns", dnsMap)
	}
	// save config
	marshal, err := yaml.Marshal(adgHomeConfig)
	if err != nil {
		return e.New("marshal adgHome config failed, ", err).WithPrefix(tagService)
	}
	// write new config
	if err := os.WriteFile(adgHomeConfigPath, marshal, 0644); err != nil {
		return e.New("write adgHome config failed, ", err).WithPrefix(tagService)
	}
	// create adghome service
	service := common.NewExternal(0, nil, nil, adgHomePath,
		"--no-check-update",
		"-w", builds.Config.AdgHome.WorkDir,
		"-c", adgHomeConfigPath,
		"--pidfile", path.Join(builds.Config.XrayHelper.RunDir, "adghome.pid"),
		"-l", path.Join(builds.Config.XrayHelper.RunDir, "adghome.log"))
	service.AppendEnv("SSL_CERT_DIR=/system/etc/security/cacerts/")
	service.SetUidGid("0", common.CoreGid)
	service.Start()
	if service.Err() != nil {
		return e.New("start adgHome service failed, ", service.Err()).WithPrefix(tagService)
	}
	if err := cgroup.LimitProcess(service.Pid()); err != nil {
		return err
	}
	return nil
}

// stopAdgHome stop AdGuardHome service
func stopAdgHome() {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "adghome.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "adghome.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		pid, _ := strconv.Atoi(string(pidFile))
		if serviceProcess, err := os.FindProcess(pid); err == nil {
			_ = serviceProcess.Kill()
			_ = os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "adghome.pid"))
		} else {
			log.HandleDebug(err)
		}
	} else {
		log.HandleDebug(err)
	}
}

// ignoreSignals start a goroutine to ignore some terminal signals
func ignoreSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, // keyboard
		syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, // os termination
		syscall.SIGUSR1, syscall.SIGUSR2, // user
		syscall.SIGPIPE, syscall.SIGCHLD, syscall.SIGSEGV) // os other
	go func() {
		for {
			select {
			case sign := <-signalChan:
				log.HandleDebug("ignore signal " + sign.String() + " since core starting")
			}
		}
	}()
}

func handleRayDNS(ipv6 bool) error {
	if confInfo, err := os.Stat(builds.Config.XrayHelper.CoreConfig); err != nil {
		return e.New("open core config file failed, ", err).WithPrefix(tagService)
	} else {
		if confInfo.IsDir() {
			confDir, err := os.ReadDir(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return e.New("open config dir failed, ", err).WithPrefix(tagService)
			}
			for _, conf := range confDir {
				if !conf.IsDir() && strings.HasSuffix(conf.Name(), ".json") {
					confByte, err := os.ReadFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()))
					if err != nil {
						return e.New("read config file failed, ", err).WithPrefix(tagService)
					}
					newConfByte, err := replaceRayDNSStrategy(confByte, ipv6)
					if err != nil {
						log.HandleDebug(err)
						continue
					}
					if err := os.WriteFile(path.Join(builds.Config.XrayHelper.CoreConfig, conf.Name()), newConfByte, 0644); err != nil {
						return e.New("write new config failed, ", err).WithPrefix(tagService)
					}
				}
			}
		} else {
			confByte, err := os.ReadFile(builds.Config.XrayHelper.CoreConfig)
			if err != nil {
				return e.New("read config file failed, ", err).WithPrefix(tagService)
			}
			newConfByte, err := replaceRayDNSStrategy(confByte, ipv6)
			if err != nil {
				return err
			}
			if err := os.WriteFile(builds.Config.XrayHelper.CoreConfig, newConfByte, 0644); err != nil {
				return e.New("write new config failed, ", err).WithPrefix(tagService)
			}
		}
	}
	return nil
}

func replaceRayDNSStrategy(conf []byte, ipv6 bool) (replacedConf []byte, err error) {
	// unmarshal
	var jsonMap serial.OrderedMap
	err = json.Unmarshal(conf, &jsonMap)
	if err != nil {
		return nil, e.New("unmarshal config json failed, ", err).WithPrefix(tagService)
	}
	if dns, ok := jsonMap.Get("dns"); ok {
		dnsMap := dns.Value.(serial.OrderedMap)
		switch builds.Config.XrayHelper.CoreType {
		case "xray":
			if ipv6 {
				dnsMap.Set("queryStrategy", "UseIP")
			} else {
				dnsMap.Set("queryStrategy", "UseIPv4")
			}
		case "v2ray":
			if ipv6 {
				dnsMap.Set("queryStrategy", "USE_IP")
			} else {
				dnsMap.Set("queryStrategy", "USE_IP4")
			}
		case "sing-box":
			if ipv6 {
				dnsMap.Set("strategy", "prefer_ipv4")
			} else {
				dnsMap.Set("strategy", "ipv4_only")
			}
		default:
			return nil, e.New("unsupported core type " + builds.Config.XrayHelper.CoreType).WithPrefix(tagService)
		}
		jsonMap.Set("dns", dnsMap)
		// marshal
		marshal, err := json.MarshalIndent(jsonMap, "", "    ")
		if err != nil {
			return nil, e.New("marshal config json failed, ", err).WithPrefix(tagService)
		}
		return marshal, nil
	}
	return nil, e.New("cannot find dns object from provided conf").WithPrefix(tagService)
}

func overrideClashConfig(template string, target string) error {
	if len(template) == 0 {
		return nil
	}
	// open target config and replace with xrayhelper clash value
	targetFile, err := os.ReadFile(target)
	if err != nil {
		return e.New("load clash config failed, ", err).WithPrefix(tagService)
	}
	var targetYamlMap serial.OrderedMap
	if err := yaml.Unmarshal(targetFile, &targetYamlMap); err != nil {
		return e.New("unmarshal clash config failed, ", err).WithPrefix(tagService)
	}
	// delete origin config
	targetYamlMap.Delete("port")
	targetYamlMap.Delete("socks-port")
	targetYamlMap.Delete("redir-port")
	targetYamlMap.Delete("tproxy-port")
	targetYamlMap.Delete("mixed-port")
	targetYamlMap.Delete("authentication")
	targetYamlMap.Delete("external-controller")
	targetYamlMap.Delete("external-ui")
	targetYamlMap.Delete("secret")
	targetYamlMap.Delete("allow-lan")
	targetYamlMap.Delete("bind-address")
	targetYamlMap.Delete("tun")
	// mihomo
	targetYamlMap.Delete("ebpf")
	targetYamlMap.Delete("sniffer")
	targetYamlMap.Delete("external-controller-tls")
	targetYamlMap.Delete("tls")
	targetYamlMap.Delete("experimental")
	// open template config and replace target value with it
	templateFile, err := os.ReadFile(template)
	if err != nil {
		return e.New("load clash template config failed, ", err).WithPrefix(tagService)
	}
	var templateYamlMap serial.OrderedMap
	if err := yaml.Unmarshal(templateFile, &templateYamlMap); err != nil {
		return e.New("unmarshal clash template config failed, ", err).WithPrefix(tagService)
	} // if enable AutoDNSStrategy
	if builds.Config.Proxy.AutoDNSStrategy {
		templateYamlMap.Set("ipv6", builds.Config.Proxy.EnableIPv6)
		if dns, ok := templateYamlMap.Get("dns"); ok {
			// assert dns
			dnsMap := dns.Value.(serial.OrderedMap)
			dnsMap.Set("listen", "127.0.0.1:"+builds.Config.Clash.DNSPort)
			templateYamlMap.Set("dns", dnsMap)
		}
	}
	// save template
	marshal, err := yaml.Marshal(templateYamlMap)
	if err != nil {
		return e.New("marshal clash template config failed, ", err).WithPrefix(tagService)
	}
	// write new template config
	if err := os.WriteFile(template, marshal, 0644); err != nil {
		return e.New("write clash template config failed, ", err).WithPrefix(tagService)
	}
	// replace target
	for _, val := range templateYamlMap.Values {
		targetYamlMap.Set(val.Key, val.Value)
	}
	// save target
	marshal, err = yaml.Marshal(targetYamlMap)
	if err != nil {
		return e.New("marshal clash config failed, ", err).WithPrefix(tagService)
	}
	// write new config
	if err := os.WriteFile(target, marshal, 0644); err != nil {
		return e.New("write overridden clash config failed, ", err).WithPrefix(tagService)
	}
	return nil
}
