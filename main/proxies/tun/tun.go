package tun

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/proxies/tools"
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type Tun struct{}

func (this *Tun) Enable() error {
	if builds.Config.Proxy.Method == "tun2socks" {
		if err := startTun2socks(); err != nil {
			this.Disable()
			return err
		}
		if err := addRoute(false); err != nil {
			this.Disable()
			return err
		}
		if err := createMangleChain(false); err != nil {
			this.Disable()
			return err
		}
		if err := createProxyChain(false); err != nil {
			this.Disable()
			return err
		}
		if builds.Config.Proxy.EnableIPv6 {
			if err := addRoute(true); err != nil {
				this.Disable()
				return err
			}
			if err := createMangleChain(true); err != nil {
				this.Disable()
				return err
			}
			if err := createProxyChain(true); err != nil {
				this.Disable()
				return err
			}
		}
		// handleDns, some core not support sniffing(eg: clash), need redirect dns request to local dns port
		switch builds.Config.XrayHelper.CoreType {
		case "clash.meta", "mihomo":
			if err := tools.RedirectDNS(builds.Config.Clash.DNSPort); err != nil {
				this.Disable()
				return err
			}
		default:
			if !builds.Config.Proxy.EnableIPv6 {
				if err := tools.DisableIPV6DNS(); err != nil {
					this.Disable()
					return err
				}
			}
		}
	} else {
		if !tunDeviceReady(builds.Config.Proxy.TunDevice) {
			return e.New("cannot find your tun device " + builds.Config.Proxy.TunDevice + " did you configure core correctly?").WithPrefix("tun").WithPathObj(*this)
		}
	}
	return nil
}

func (this *Tun) Disable() {
	if builds.Config.Proxy.Method == "tun2socks" {
		deleteRoute(false)
		cleanIptablesChain(false)
		//always clean ipv6 rules
		deleteRoute(true)
		cleanIptablesChain(true)
		stopTun2socks()
		//always clean dns rules
		tools.EnableIPV6DNS()
		tools.CleanRedirectDNS(builds.Config.Clash.DNSPort)
	}
}

func tunDeviceReady(checkDev string) bool {
	for i := 0; i < *builds.CoreStartTimeout; i++ {
		time.Sleep(1 * time.Second)
		if common.CheckLocalDevice(checkDev) {
			return true
		}
	}
	return false
}

func startTun2socks() error {
	tun2socksPath := path.Join(path.Dir(builds.Config.XrayHelper.CorePath), "tun2socks")
	tun2socksConfigPath := path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.yml")
	var tunConfig struct {
		Tunnel struct {
			Name       string `yaml:"name"`
			Mtu        int    `yaml:"mtu"`
			MultiQueue bool   `yaml:"multi-queue"`
			IPv4       string `yaml:"ipv4"`
			IPv6       string `yaml:"ipv6"`
		} `yaml:"tunnel"`
		Socks5 struct {
			Port    int    `yaml:"port"`
			Address string `yaml:"address"`
			Udp     string `yaml:"udp"`
		} `yaml:"socks5"`
	}
	tunConfig.Tunnel.Name = builds.Config.Proxy.TunDevice
	tunConfig.Tunnel.Mtu = common.Tun2socksMTU
	tunConfig.Tunnel.MultiQueue = common.Tun2socksMultiQueue
	tunConfig.Tunnel.IPv4 = common.Tun2socksIPv4
	tunConfig.Tunnel.IPv6 = common.Tun2socksIPv6
	tunConfig.Socks5.Port, _ = strconv.Atoi(builds.Config.Proxy.SocksPort)
	tunConfig.Socks5.Address = "127.0.0.1"
	tunConfig.Socks5.Udp = common.Tun2socksUdpMode
	configByte, err := yaml.Marshal(&tunConfig)
	if err != nil {
		return e.New("generate tun2socks config failed, ", err).WithPrefix("tun")
	}
	if err := os.WriteFile(tun2socksConfigPath, configByte, 0644); err != nil {
		return e.New("write tun2socks config failed, ", err).WithPrefix("tun")
	}
	tun2socksLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return e.New("open tun2socks log file failed, ", err).WithPrefix("tun")
	}
	service := common.NewExternal(0, tun2socksLogFile, tun2socksLogFile, tun2socksPath, tun2socksConfigPath)
	service.Start()
	if service.Err() != nil {
		return e.New("start tun2socks failed, ", service.Err()).WithPrefix("tun")
	}
	if tunDeviceReady(builds.Config.Proxy.TunDevice) {
		if err := os.WriteFile(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.pid"), []byte(strconv.Itoa(service.Pid())), 0644); err != nil {
			_ = service.Kill()
			return e.New("write tun2socks pid failed, ", err).WithPrefix("tun")
		}
	} else {
		_ = service.Kill()
		return e.New("start tun2socks failed, please check tun2socks.log").WithPrefix("tun")
	}
	return nil
}

func stopTun2socks() {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		pid, _ := strconv.Atoi(string(pidFile))
		if serviceProcess, err := os.FindProcess(pid); err == nil {
			_ = serviceProcess.Kill()
			_ = os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.pid"))
		} else {
			log.HandleDebug(err)
		}
	} else {
		log.HandleDebug(err)
	}
	err := os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.yml"))
	if err != nil {
		log.HandleDebug(err)
	}
}

// addRoute Add ip route to proxy
func addRoute(ipv6 bool) error {
	var errMsg bytes.Buffer
	if !ipv6 {
		common.NewExternal(0, nil, &errMsg, "ip", "rule", "add", "fwmark", common.TunMarkId, "lookup", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return e.New("add ip rule failed, ", errMsg.String()).WithPrefix("tun")
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "route", "add", "default", "dev", builds.Config.Proxy.TunDevice, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return e.New("add ip route failed, ", errMsg.String()).WithPrefix("tun")
		}
	} else {
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "add", "fwmark", common.TunMarkId, "lookup", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return e.New("add ip rule failed, ", errMsg.String()).WithPrefix("tun")
		}
		errMsg.Reset()
		// when device do not have ipv6 address, route all ipv6 traffic to tun
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "add", "from", "all", "lookup", common.TunTableId, "prio", "31999").Run()
		if errMsg.Len() > 0 {
			return e.New("add ip rule failed, ", errMsg.String()).WithPrefix("tun")
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "add", "default", "dev", builds.Config.Proxy.TunDevice, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return e.New("add ip route failed, ", errMsg.String()).WithPrefix("tun")
		}
	}
	return nil
}

// deleteRoute Delete ip route to proxy
func deleteRoute(ipv6 bool) {
	var errMsg bytes.Buffer
	if !ipv6 {
		common.NewExternal(0, nil, &errMsg, "ip", "rule", "del", "fwmark", common.TunMarkId, "lookup", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip rule: " + errMsg.String())
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "route", "flush", "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip route: " + errMsg.String())
		}
	} else {
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "del", "fwmark", common.TunMarkId, "lookup", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip rule: " + errMsg.String())
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "del", "from", "all", "lookup", common.TunTableId, "prio", "31999").Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip rule: " + errMsg.String())
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "flush", "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip route: " + errMsg.String())
		}
	}
}

// createProxyChain Create XT chain for local applications
func createProxyChain(ipv6 bool) error {
	var currentProto string
	currentIpt := common.Ipt
	currentProto = "ipv4"
	if ipv6 {
		currentIpt = common.Ipt6
		currentProto = "ipv6"
	}
	if currentIpt == nil {
		return e.New("get iptables failed").WithPrefix("tun")
	}
	if err := currentIpt.NewChain("mangle", "XT"); err != nil {
		return e.New("create "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
	}
	// bypass tun2socks
	if err := currentIpt.Append("mangle", "XT", "-o", builds.Config.Proxy.TunDevice, "-j", "RETURN"); err != nil {
		return e.New("ignore tun2socks interface "+builds.Config.Proxy.TunDevice+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
	}
	// bypass ignore list
	for _, ignore := range builds.Config.Proxy.IgnoreList {
		if err := currentIpt.Append("mangle", "XT", "-o", ignore, "-j", "RETURN"); err != nil {
			return e.New("apply ignore interface "+ignore+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
		}
	}
	// bypass intraNet list
	if currentProto == "ipv4" {
		for _, intraIp := range common.IntraNet {
			if err := currentIpt.Append("mangle", "XT", "-d", intraIp, "-j", "RETURN"); err != nil {
				return e.New("bypass intraNet "+intraIp+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
	} else {
		for _, intraIp6 := range common.IntraNet6 {
			if err := currentIpt.Append("mangle", "XT", "-d", intraIp6, "-j", "RETURN"); err != nil {
				return e.New("bypass intraNet "+intraIp6+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
	}
	// bypass Core itself
	if err := currentIpt.Append("mangle", "XT", "-m", "owner", "--gid-owner", common.CoreGid, "-j", "RETURN"); err != nil {
		return e.New("bypass core gid on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
	}
	// start processing proxy rules
	// if PkgList has no package, should proxy everything
	if len(builds.Config.Proxy.PkgList) == 0 {
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create local applications proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create local applications proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	} else if builds.Config.Proxy.Mode == "blacklist" {
		// bypass PkgList
		for _, pkg := range builds.Config.Proxy.PkgList {
			uid, err := tools.GetUid(pkg)
			if err != nil {
				log.HandleDebug(err)
				continue
			}
			if err := currentIpt.Insert("mangle", "XT", 1, "-m", "owner", "--uid-owner", uid, "-j", "RETURN"); err != nil {
				return e.New("bypass package "+pkg+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
		// allow others
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create local applications proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create local applications proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	} else if builds.Config.Proxy.Mode == "whitelist" {
		// allow PkgList
		for _, pkg := range builds.Config.Proxy.PkgList {
			uid, err := tools.GetUid(pkg)
			if err != nil {
				log.HandleDebug(err)
				continue
			}
			if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-m", "owner", "--uid-owner", uid, "-j", "TUN2SOCKS"); err != nil {
				return e.New("create package "+pkg+" proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
			}
			if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-m", "owner", "--uid-owner", uid, "-j", "TUN2SOCKS"); err != nil {
				return e.New("create package "+pkg+" proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
		// allow root user(eg: magisk, ksud, netd...)
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-m", "owner", "--uid-owner", "0", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create root user proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-m", "owner", "--uid-owner", "0", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create root user proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		// allow dns_tether user(eg: dnsmasq...)
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-m", "owner", "--uid-owner", "1052", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create dns_tether user proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-m", "owner", "--uid-owner", "1052", "-j", "TUN2SOCKS"); err != nil {
			return e.New("create dns_tether user proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	} else {
		return e.New("invalid proxy mode " + builds.Config.Proxy.Mode).WithPrefix("tun")
	}
	// allow IntraList
	for _, intra := range builds.Config.Proxy.IntraList {
		if (currentProto == "ipv4" && !common.IsIPv6(intra)) || (currentProto == "ipv6" && common.IsIPv6(intra)) {
			if err := currentIpt.Insert("mangle", "XT", 1, "-p", "tcp", "-d", intra, "-j", "TUN2SOCKS"); err != nil {
				return e.New("allow intra "+intra+" on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
			}
			if err := currentIpt.Insert("mangle", "XT", 1, "-p", "udp", "-d", intra, "-j", "TUN2SOCKS"); err != nil {
				return e.New("allow intra "+intra+" on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
	}
	// mark all dns request
	if !strings.Contains(builds.Config.XrayHelper.CoreType, "clash") {
		if err := currentIpt.Insert("mangle", "XT", 1, "-p", "udp", "-m", "owner", "!", "--gid-owner", common.CoreGid, "--dport", "53", "-j", "TUN2SOCKS"); err != nil {
			return e.New("mark all dns request on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	}
	// apply rules to OUTPUT
	if err := currentIpt.Append("mangle", "OUTPUT", "-j", "XT"); err != nil {
		return e.New("apply mangle chain XT to OUTPUT failed, ", err).WithPrefix("tun")
	}
	return nil
}

// createMangleChain Create TUN2SOCKS chain for AP interface, there will be problem on some device
func createMangleChain(ipv6 bool) error {
	var currentProto string
	currentIpt := common.Ipt
	currentProto = "ipv4"
	if ipv6 {
		currentIpt = common.Ipt6
		currentProto = "ipv6"
	}
	if currentIpt == nil {
		return e.New("get iptables failed").WithPrefix("tun")
	}
	if err := currentIpt.NewChain("mangle", "TUN2SOCKS"); err != nil {
		return e.New("create "+currentProto+" mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
	}
	// bypass intraNet list
	if currentProto == "ipv4" {
		for _, intraIp := range common.IntraNet {
			if err := currentIpt.Append("mangle", "TUN2SOCKS", "-d", intraIp, "-j", "RETURN"); err != nil {
				return e.New("bypass intraNet "+intraIp+" on "+currentProto+" mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
		}
	} else {
		for _, intraIp6 := range common.IntraNet6 {
			if err := currentIpt.Append("mangle", "TUN2SOCKS", "-d", intraIp6, "-j", "RETURN"); err != nil {
				return e.New("bypass intraNet "+intraIp6+" on "+currentProto+" mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
		}
	}
	// allow IntraList
	for _, intra := range builds.Config.Proxy.IntraList {
		if (currentProto == "ipv4" && !common.IsIPv6(intra)) || (currentProto == "ipv6" && common.IsIPv6(intra)) {
			if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "tcp", "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
				return e.New("allow intra "+intra+" on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
			if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "udp", "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
				return e.New("allow intra "+intra+" on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
		}
	}
	// mark all traffic
	if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "tcp", "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
		return e.New("create all traffic proxy on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
	}
	if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "udp", "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
		return e.New("create all traffic proxy on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
	}
	// trans ApList to chain XRAY
	for _, ap := range builds.Config.Proxy.ApList {
		// allow ApList to IntraList
		for _, intra := range builds.Config.Proxy.IntraList {
			if (currentProto == "ipv4" && !common.IsIPv6(intra)) || (currentProto == "ipv6" && common.IsIPv6(intra)) {
				if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "tcp", "-i", ap, "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
					return e.New("allow intra "+intra+" on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
				}
				if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "udp", "-i", ap, "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
					return e.New("allow intra "+intra+" on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
				}
			}
		}
		if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "tcp", "-i", ap, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
			return e.New("create ap interface "+ap+" proxy on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "udp", "-i", ap, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
			return e.New("create ap interface "+ap+" proxy on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
		}
	}
	// mark all dns request
	if !strings.Contains(builds.Config.XrayHelper.CoreType, "clash") {
		if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "udp", "--dport", "53", "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
			return e.New("mark all dns request on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
		}
	}
	// apply rules to PREROUTING
	if err := currentIpt.Append("mangle", "PREROUTING", "-j", "TUN2SOCKS"); err != nil {
		return e.New("apply mangle chain TUN2SOCKS to PREROUTING failed, ", err).WithPrefix("tun")
	}
	return nil
}

// cleanIptablesChain Clean all changed iptables rules by XrayHelper
func cleanIptablesChain(ipv6 bool) {
	currentIpt := common.Ipt
	if ipv6 {
		currentIpt = common.Ipt6
	}
	if currentIpt == nil {
		return
	}
	_ = currentIpt.Delete("mangle", "OUTPUT", "-j", "XT")
	_ = currentIpt.Delete("mangle", "PREROUTING", "-j", "TUN2SOCKS")
	_ = currentIpt.ClearAndDeleteChain("mangle", "XT")
	_ = currentIpt.ClearAndDeleteChain("mangle", "TUN2SOCKS")
}
