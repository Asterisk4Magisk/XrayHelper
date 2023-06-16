package tools

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"strings"
)

func GetUid(pkgInfo string) (string, error) {
	info := strings.Split(pkgInfo, ":")
	if len(info) == 1 {
		if pkgId, ok := builds.PackageMap[info[0]]; ok {
			return pkgId, nil
		}
	} else {
		if pkgId, ok := builds.PackageMap[info[0]]; ok {
			if info[1] == "0" {
				return pkgId, nil
			}
			return info[1] + pkgId, nil
		}
	}
	return "", errors.New("cannot get uid " + info[0]).WithPrefix("tools")
}

func DisableIPV6DNS() error {
	if err := common.Ipt6.Insert("filter", "OUTPUT", 1, "-p", "udp", "--dport", "53", "-j", "REJECT"); err != nil {
		return errors.New("disable dns request on ipv6 failed, ", err).WithPrefix("tools")
	}
	return nil
}

func EnableIPV6DNS() {
	_ = common.Ipt6.Delete("filter", "OUTPUT", "-p", "udp", "--dport", "53", "-j", "REJECT")
}

func RedirectDNS(port string) error {
	if err := common.Ipt.Insert("nat", "OUTPUT", 1, "-p", "udp", "-m", "owner", "!", "--gid-owner", common.CoreGid, "--dport", "53", "-j", "DNAT", "--to-destination", "127.0.0.1:"+port); err != nil {
		return errors.New("redirect dns request failed, ", err).WithPrefix("tools")
	}
	if err := DisableIPV6DNS(); err != nil {
		return err
	}
	return nil
}

func CleanRedirectDNS(port string) {
	_ = common.Ipt.Delete("nat", "OUTPUT", "-p", "udp", "-m", "owner", "!", "--gid-owner", common.CoreGid, "--dport", "53", "-j", "DNAT", "--to-destination", "127.0.0.1:"+port)
	EnableIPV6DNS()
}
