package tproxy

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/utils"
	"bytes"
	"github.com/coreos/go-iptables/iptables"
)

const (
	tableId = "233"
	xrayGid = "3003"
	markId  = "1111"
)

var (
	ipt, _   = iptables.NewWithProtocol(iptables.ProtocolIPv4)
	ipt6, _  = iptables.NewWithProtocol(iptables.ProtocolIPv6)
	intraNet = []string{"0.0.0.0/8", "10.0.0.0/8", "100.64.0.0/10", "127.0.0.0/8", "169.254.0.0/16",
		"172.16.0.0/12", "192.0.0.0/24", "192.0.2.0/24", "192.88.99.0/24", "192.168.0.0/16", "198.51.100.0/24",
		"203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4", "255.255.255.255/32"}
	intraNet6 = []string{"::/128", "::1/128", "::ffff:0:0/96", "100::/64", "64:ff9b::/96", "2001::/32",
		"2001:10::/28", "2001:20::/28", "2001:db8::/32", "2002::/16", "fc00::/7", "fe80::/10", "ff00::/8"}
)

// AddRoute Add ip route to proxy
func AddRoute(ipv6 bool) error {
	var outMsg bytes.Buffer
	if !ipv6 {
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "rule", "add", "fwmark", markId, "table", tableId).Run()
		if outMsg.Len() > 0 {
			return errors.New("add ip rule failed, ", outMsg.String()).WithPrefix("tproxy")
		}
		outMsg.Reset()
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "route", "add", "local", "default", "dev", "lo", "table", tableId).Run()
		if outMsg.Len() > 0 {
			return errors.New("add ip route failed, ", outMsg.String()).WithPrefix("tproxy")
		}
	} else {
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "-6", "rule", "add", "fwmark", markId, "table", tableId).Run()
		if outMsg.Len() > 0 {
			return errors.New("add ip rule failed, ", outMsg.String()).WithPrefix("tproxy")
		}
		outMsg.Reset()
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "-6", "route", "add", "local", "default", "dev", "lo", "table", tableId).Run()
		if outMsg.Len() > 0 {
			return errors.New("add ip route failed, ", outMsg.String()).WithPrefix("tproxy")
		}
	}
	return nil
}

// DeleteRoute Delete ip route to proxy
func DeleteRoute(ipv6 bool) {
	var outMsg bytes.Buffer
	if !ipv6 {
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "rule", "del", "fwmark", markId, "table", tableId).Run()
		log.HandleDebug("delete ip rule: " + outMsg.String())
		outMsg.Reset()
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "route", "flush", "table", tableId).Run()
		log.HandleDebug("delete ip route: " + outMsg.String())
	} else {
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "-6", "rule", "del", "fwmark", markId, "table", tableId).Run()
		log.HandleDebug("delete ip rule: " + outMsg.String())
		outMsg.Reset()
		utils.NewExternal(0, &outMsg, &outMsg, "ip", "-6", "route", "flush", "table", tableId).Run()
		log.HandleDebug("delete ip route: " + outMsg.String())
	}
}

// CreateProxyChain Create PROXY chain for local applications
func CreateProxyChain(ipv6 bool) error {
	var currentProto string
	currentIpt := ipt
	if ipv6 {
		currentIpt = ipt6
	}
	if currentIpt == nil {
		return errors.New("get iptables failed").WithPrefix("tproxy")
	}
	if currentIpt.Proto() == iptables.ProtocolIPv4 {
		currentProto = "ipv4"
	} else {
		currentProto = "ipv6"
	}
	if err := currentIpt.NewChain("mangle", "PROXY"); err != nil {
		currentIpt.Proto()
		return errors.New("create "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
	}
	// bypass ignore list
	for _, ignore := range builds.Config.Proxy.IgnoreList {
		if err := currentIpt.Append("mangle", "PROXY", "-o", ignore, "-j", "RETURN"); err != nil {
			return errors.New("apply ignore interface "+ignore+" on "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
	}
	// bypass intraNet list
	if currentProto == "ipv4" {
		for _, intraIp := range intraNet {
			if err := currentIpt.Append("mangle", "PROXY", "-d", intraIp, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp+" on "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
			}
		}
	} else {
		for _, intraIp6 := range intraNet6 {
			if err := currentIpt.Append("mangle", "PROXY", "-d", intraIp6, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp6+" on "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
			}
		}
		if externalIPv6, err := utils.GetIPv6Addr(); err != nil {
			if len(externalIPv6) > 0 {
				for _, external := range externalIPv6 {
					if err := currentIpt.Append("mangle", "PROXY", "-d", external, "-j", "RETURN"); err != nil {
						return errors.New("bypass externalIPv6 "+external+" on "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
					}
				}
			}
		} else {
			return err
		}
	}
	// bypass Xray itself
	if err := currentIpt.Append("mangle", "PROXY", "-m", "owner", "--gid-owner", xrayGid, "-j", "RETURN"); err != nil {
		return errors.New("bypass xray gid on "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
	}
	// start processing proxy rules
	// if PkgList has no package, should proxy everything
	if len(builds.Config.Proxy.PkgList) == 0 {
		if err := currentIpt.Append("mangle", "PROXY", "-p", "tcp", "-j", "MARK", "--set-mark", markId); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" tcp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
		if err := currentIpt.Append("mangle", "PROXY", "-p", "udp", "-j", "MARK", "--set-mark", markId); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" udp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
	} else if builds.Config.Proxy.Mode == "blacklist" {
		// bypass PkgList
		for _, pkg := range builds.Config.Proxy.PkgList {
			if uid, ok := builds.PackageMap[pkg]; ok {
				if err := currentIpt.Insert("mangle", "PROXY", 0, "-m", "owner", "--uid-owner", uid, "-j", "RETURN"); err != nil {
					return errors.New("bypass package "+pkg+" on "+currentProto+" mangle chain PROXY failed, ", err).WithPrefix("tproxy")
				}
			}
		}
		// allow others
		if err := currentIpt.Append("mangle", "PROXY", "-p", "tcp", "-j", "MARK", "--set-mark", markId); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" tcp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
		if err := currentIpt.Append("mangle", "PROXY", "-p", "udp", "-j", "MARK", "--set-mark", markId); err != nil {
			return errors.New("create local applications proxy on "+currentProto+" udp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
	} else if builds.Config.Proxy.Mode == "whitelist" {
		// allow PkgList
		for _, pkg := range builds.Config.Proxy.PkgList {
			if uid, ok := builds.PackageMap[pkg]; ok {
				if err := currentIpt.Append("mangle", "PROXY", "-p", "tcp", "-m", "owner", "--uid-owner", uid, "-j", "MARK", "--set-mark", markId); err != nil {
					return errors.New("create package "+pkg+" proxy on "+currentProto+" tcp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
				}
				if err := currentIpt.Append("mangle", "PROXY", "-p", "udp", "-m", "owner", "--uid-owner", uid, "-j", "MARK", "--set-mark", markId); err != nil {
					return errors.New("create package "+pkg+" proxy on "+currentProto+" udp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
				}
			}
		}
		// allow root user(eg: magisk, netd, dnsmasq...)
		if err := currentIpt.Append("mangle", "PROXY", "-p", "tcp", "-m", "owner", "--uid-owner", "0", "-j", "MARK", "--set-mark", markId); err != nil {
			return errors.New("create root user proxy on "+currentProto+" tcp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
		if err := currentIpt.Append("mangle", "PROXY", "-p", "udp", "-m", "owner", "--uid-owner", "0", "-j", "MARK", "--set-mark", markId); err != nil {
			return errors.New("create root user proxy on "+currentProto+" udp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
	} else {
		return errors.New("invalid proxy mode " + builds.Config.Proxy.Mode).WithPrefix("tproxy")
	}
	// apply rules to OUTPUT
	if err := currentIpt.Append("mangle", "OUTPUT", "-j", "PROXY"); err != nil {
		return errors.New("apply mangle chain PROXY to OUTPUT failed, ", err).WithPrefix("tproxy")
	}
	return nil
}

// CreateMangleChain Create XRAY chain for AP interface
func CreateMangleChain(ipv6 bool) error {
	var currentProto string
	currentIpt := ipt
	if ipv6 {
		currentIpt = ipt6
	}
	if currentIpt == nil {
		return errors.New("get iptables failed").WithPrefix("tproxy")
	}
	if currentIpt.Proto() == iptables.ProtocolIPv4 {
		currentProto = "ipv4"
	} else {
		currentProto = "ipv6"
	}
	if err := currentIpt.NewChain("mangle", "XRAY"); err != nil {
		currentIpt.Proto()
		return errors.New("create "+currentProto+" mangle chain XRAY failed, ", err).WithPrefix("tproxy")
	}
	// bypass intraNet list
	if currentProto == "ipv4" {
		for _, intraIp := range intraNet {
			if err := currentIpt.Append("mangle", "XRAY", "-d", intraIp, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp+" on "+currentProto+" mangle chain XRAY failed, ", err).WithPrefix("tproxy")
			}
		}
	} else {
		for _, intraIp6 := range intraNet6 {
			if err := currentIpt.Append("mangle", "XRAY", "-d", intraIp6, "-j", "RETURN"); err != nil {
				return errors.New("bypass intraNet "+intraIp6+" on "+currentProto+" mangle chain XRAY failed, ", err).WithPrefix("tproxy")
			}
		}
		if externalIPv6, err := utils.GetIPv6Addr(); err != nil {
			if len(externalIPv6) > 0 {
				for _, external := range externalIPv6 {
					if err := currentIpt.Append("mangle", "XRAY", "-d", external, "-j", "RETURN"); err != nil {
						return errors.New("bypass externalIPv6 "+external+" on "+currentProto+" mangle chain XRAY failed, ", err).WithPrefix("tproxy")
					}
				}
			}
		} else {
			return err
		}
	}
	// mark all traffic
	if err := currentIpt.Append("mangle", "XRAY", "-p", "tcp", "-m", "mark", "--mark", markId, "-j", "TPROXY", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", markId); err != nil {
		return errors.New("create all traffic proxy on "+currentProto+" tcp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
	}
	if err := currentIpt.Append("mangle", "XRAY", "-p", "udp", "-m", "mark", "--mark", markId, "-j", "TPROXY", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", markId); err != nil {
		return errors.New("create all traffic proxy on "+currentProto+" udp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
	}
	// trans ApList to chain XRAY
	for _, ap := range builds.Config.Proxy.ApList {
		if err := currentIpt.Append("mangle", "XRAY", "-p", "tcp", "-i", ap, "-j", "TPROXY", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", markId); err != nil {
			return errors.New("create ap interface "+ap+" proxy on "+currentProto+" tcp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
		if err := currentIpt.Append("mangle", "XRAY", "-p", "udp", "-i", ap, "-j", "TPROXY", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", markId); err != nil {
			return errors.New("create ap interface "+ap+" proxy on "+currentProto+" udp mangle chain PROXY failed, ", err).WithPrefix("tproxy")
		}
	}
	// apply rules to PREROUTING
	if err := currentIpt.Append("mangle", "PREROUTING", "-j", "XRAY"); err != nil {
		return errors.New("apply mangle chain XRAY to PREROUTING failed, ", err).WithPrefix("tproxy")
	}
	return nil
}

// CleanIptablesChain Clean all changed iptables rules by XrayHelper
func CleanIptablesChain(ipv6 bool) {
	currentIpt := ipt
	if ipv6 {
		currentIpt = ipt6
	}
	if currentIpt == nil {
		return
	}
	_ = currentIpt.Delete("mangle", "OUTPUT", "-j", "PROXY")
	_ = currentIpt.Delete("mangle", "PREROUTING", "-j", "XRAY")
	_ = currentIpt.ClearAndDeleteChain("mangle", "PROXY")
	_ = currentIpt.ClearAndDeleteChain("mangle", "XRAY")
}
