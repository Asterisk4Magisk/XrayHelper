package tun

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strconv"
	"time"
)

func StartTun() error {
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
	tunConfig.Tunnel.Name = common.TunDevice
	tunConfig.Tunnel.Mtu = common.TunMTU
	tunConfig.Tunnel.MultiQueue = common.TunMultiQueue
	tunConfig.Tunnel.IPv4 = common.TunIPv4
	tunConfig.Tunnel.IPv6 = common.TunIPv6
	tunConfig.Socks5.Port, _ = strconv.Atoi(builds.Config.Proxy.SocksPort)
	tunConfig.Socks5.Address = "127.0.0.1"
	tunConfig.Socks5.Udp = common.TunUdpMode
	configByte, err := yaml.Marshal(&tunConfig)
	if err != nil {
		return errors.New("generate tun2socks config failed, ", err).WithPrefix("tun")
	}
	if err := os.WriteFile(tun2socksConfigPath, configByte, 0644); err != nil {
		return errors.New("write tun2socks config failed, ", err).WithPrefix("tun")
	}
	service := common.NewExternal(0, nil, nil, tun2socksPath, tun2socksConfigPath)
	service.Start()
	deviceReady := false
	for i := 0; i < 15; i++ {
		time.Sleep(1 * time.Second)
		if common.CheckLocalIP(common.TunIPv4) {
			deviceReady = true
			break
		}
	}
	if deviceReady {
		if err := os.WriteFile(path.Join(builds.Config.XrayHelper.RunDir, "tun2socks.pid"), []byte(strconv.Itoa(service.Pid())), 0644); err != nil {
			_ = service.Kill()
			return errors.New("write tun2socks pid failed, ", err).WithPrefix("tun")
		}
	} else {
		_ = service.Kill()
		return errors.New("start tun2socks service failed, ", service.Err()).WithPrefix("tun")
	}
	return nil
}

func StopTun() {
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

// AddRoute Add ip route to proxy
func AddRoute(ipv6 bool) error {
	var errMsg bytes.Buffer
	if !ipv6 {
		common.NewExternal(0, nil, &errMsg, "ip", "rule", "add", "fwmark", common.TunMarkId, "lookup", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return errors.New("add ip rule failed, ", errMsg.String()).WithPrefix("tun")
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "route", "add", "default", "dev", common.TunDevice, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return errors.New("add ip route failed, ", errMsg.String()).WithPrefix("tun")
		}
	} else {
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "add", "not", "from", "all", "fwmark", common.TunMarkId, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return errors.New("add ip rule failed, ", errMsg.String()).WithPrefix("tun")
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "add", "local", "default", "dev", common.TunDevice, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			return errors.New("add ip route failed, ", errMsg.String()).WithPrefix("tun")
		}
	}
	return nil
}

// DeleteRoute Delete ip route to proxy
func DeleteRoute(ipv6 bool) {
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
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "del", "not", "from", "all", "fwmark", common.TunMarkId, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip rule: " + errMsg.String())
		}
		errMsg.Reset()
		common.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "del", "local", "default", "dev", common.TunDevice, "table", common.TunTableId).Run()
		if errMsg.Len() > 0 {
			log.HandleDebug("delete ip route: " + errMsg.String())
		}
	}
}

// CreateProxyChain Create PROXY chain for local applications
func CreateProxyChain(ipv6 bool) error {
	var currentProto string
	currentIpt := common.Ipt
	currentProto = "ipv4"
	if ipv6 {
		currentIpt = common.Ipt6
		currentProto = "ipv6"
	}
	if currentIpt == nil {
		return errors.New("get iptables failed").WithPrefix("tun")
	}
	if err := currentIpt.NewChain("mangle", "XT"); err != nil {
		return errors.New("create "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
	}
	// bypass ignore list
	for _, ignore := range builds.Config.Proxy.IgnoreList {
		if err := currentIpt.Append("mangle", "XT", "-o", ignore, "-j", "RETURN"); err != nil {
			return errors.New("apply ignore interface "+ignore+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
		}
	}
	// bypass intraNet list
	if currentProto == "ipv4" {
		for _, intraIp := range common.IntraNet {
			if err := currentIpt.Append("mangle", "XT", "-d", intraIp, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
	} else {
		for _, intraIp6 := range common.IntraNet6 {
			if err := currentIpt.Append("mangle", "XT", "-d", intraIp6, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp6+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
	}
	// bypass Core itself
	if err := currentIpt.Append("mangle", "XT", "-m", "owner", "--gid-owner", common.CoreGid, "-j", "RETURN"); err != nil {
		return errors.New("bypass core gid on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
	}
	// start processing proxy rules
	// if PkgList has no package, should proxy everything
	if len(builds.Config.Proxy.PkgList) == 0 {
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-j", "TUN2SOCKS"); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-j", "TUN2SOCKS"); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	} else if builds.Config.Proxy.Mode == "blacklist" {
		// bypass PkgList
		for _, pkg := range builds.Config.Proxy.PkgList {
			if uid, ok := builds.PackageMap[pkg]; ok {
				if err := currentIpt.Insert("mangle", "XT", 1, "-m", "owner", "--uid-owner", uid, "-j", "RETURN"); err != nil {
					return errors.New("bypass package "+pkg+" on "+currentProto+" mangle chain XT failed, ", err).WithPrefix("tun")
				}
			}
		}
		// allow others
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-j", "TUN2SOCKS"); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-j", "TUN2SOCKS"); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	} else if builds.Config.Proxy.Mode == "whitelist" {
		// allow PkgList
		for _, pkg := range builds.Config.Proxy.PkgList {
			if uid, ok := builds.PackageMap[pkg]; ok {
				if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-m", "owner", "--uid-owner", uid, "-j", "TUN2SOCKS"); err != nil {
					return errors.New("create package "+pkg+" proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
				}
				if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-m", "owner", "--uid-owner", uid, "-j", "TUN2SOCKS"); err != nil {
					return errors.New("create package "+pkg+" proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
				}
			}
		}
		// allow root user(eg: magisk, netd, dnsmasq...)
		if err := currentIpt.Append("mangle", "XT", "-p", "tcp", "-m", "owner", "--uid-owner", "0", "-j", "TUN2SOCKS"); err != nil {
			return errors.New("create root user proxy on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "XT", "-p", "udp", "-m", "owner", "--uid-owner", "0", "-j", "TUN2SOCKS"); err != nil {
			return errors.New("create root user proxy on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
		}
	} else {
		return errors.New("invalid proxy mode " + builds.Config.Proxy.Mode).WithPrefix("tun")
	}
	// allow IntraList
	for _, intra := range builds.Config.Proxy.IntraList {
		if (currentProto == "ipv4" && !common.IsIPv6(intra)) || (currentProto == "ipv6" && common.IsIPv6(intra)) {
			if err := currentIpt.Insert("mangle", "XT", 1, "-p", "tcp", "-d", intra, "-j", "TUN2SOCKS"); err != nil {
				return errors.New("allow intra "+intra+" on "+currentProto+" tcp mangle chain XT failed, ", err).WithPrefix("tun")
			}
			if err := currentIpt.Insert("mangle", "XT", 1, "-p", "udp", "-d", intra, "-j", "TUN2SOCKS"); err != nil {
				return errors.New("allow intra "+intra+" on "+currentProto+" udp mangle chain XT failed, ", err).WithPrefix("tun")
			}
		}
	}
	// apply rules to OUTPUT
	if err := currentIpt.Append("mangle", "OUTPUT", "-j", "XT"); err != nil {
		return errors.New("apply mangle chain XT to OUTPUT failed, ", err).WithPrefix("tun")
	}
	return nil
}

// CreateMangleChain Create XRAY chain for AP interface
func CreateMangleChain(ipv6 bool) error {
	var currentProto string
	currentIpt := common.Ipt
	currentProto = "ipv4"
	if ipv6 {
		currentIpt = common.Ipt6
		currentProto = "ipv6"
	}
	if currentIpt == nil {
		return errors.New("get iptables failed").WithPrefix("tun")
	}
	if err := currentIpt.NewChain("mangle", "TUN2SOCKS"); err != nil {
		return errors.New("create "+currentProto+" mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
	}
	// bypass intraNet list
	if currentProto == "ipv4" {
		for _, intraIp := range common.IntraNet {
			if err := currentIpt.Append("mangle", "TUN2SOCKS", "-d", intraIp, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp+" on "+currentProto+" mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
		}
	} else {
		for _, intraIp6 := range common.IntraNet6 {
			if err := currentIpt.Append("mangle", "TUN2SOCKS", "-d", intraIp6, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp6+" on "+currentProto+" mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
		}
	}
	// allow IntraList
	for _, intra := range builds.Config.Proxy.IntraList {
		if (currentProto == "ipv4" && !common.IsIPv6(intra)) || (currentProto == "ipv6" && common.IsIPv6(intra)) {
			if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "tcp", "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
				return errors.New("allow intra "+intra+" on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
			if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "udp", "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
				return errors.New("allow intra "+intra+" on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
			}
		}
	}
	// mark all traffic
	if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "tcp", "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
		return errors.New("create all traffic proxy on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
	}
	if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "udp", "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
		return errors.New("create all traffic proxy on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
	}
	// trans ApList to chain XRAY
	for _, ap := range builds.Config.Proxy.ApList {
		// allow ApList to IntraList
		for _, intra := range builds.Config.Proxy.IntraList {
			if (currentProto == "ipv4" && !common.IsIPv6(intra)) || (currentProto == "ipv6" && common.IsIPv6(intra)) {
				if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "tcp", "-i", ap, "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
					return errors.New("allow intra "+intra+" on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
				}
				if err := currentIpt.Insert("mangle", "TUN2SOCKS", 1, "-p", "udp", "-i", ap, "-d", intra, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
					return errors.New("allow intra "+intra+" on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
				}
			}
		}
		if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "tcp", "-i", ap, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
			return errors.New("create ap interface "+ap+" proxy on "+currentProto+" tcp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
		}
		if err := currentIpt.Append("mangle", "TUN2SOCKS", "-p", "udp", "-i", ap, "-j", "MARK", "--set-xmark", common.TunMarkId); err != nil {
			return errors.New("create ap interface "+ap+" proxy on "+currentProto+" udp mangle chain TUN2SOCKS failed, ", err).WithPrefix("tun")
		}
	}
	// apply rules to PREROUTING
	if err := currentIpt.Append("mangle", "PREROUTING", "-j", "TUN2SOCKS"); err != nil {
		return errors.New("apply mangle chain TUN2SOCKS to PREROUTING failed, ", err).WithPrefix("tun")
	}
	return nil
}

// CleanIptablesChain Clean all changed iptables rules by XrayHelper
func CleanIptablesChain(ipv6 bool) {
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
