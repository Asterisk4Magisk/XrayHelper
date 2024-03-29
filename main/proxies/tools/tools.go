package tools

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"strconv"
	"strings"
)

const tagTools = "tools"

func GetUid(pkgInfo string) (string, error) {
	var appId, userId int
	info := strings.Split(pkgInfo, ":")
	if pkgId, ok := builds.PackageMap[info[0]]; ok {
		appId, _ = strconv.Atoi(pkgId)
	} else {
		return "", e.New("cannot get uid from " + pkgInfo).WithPrefix(tagTools)
	}
	if len(info) == 2 {
		userId, _ = strconv.Atoi(info[1])
	}
	return strconv.Itoa(userId*100000 + appId), nil
}

func DisableIPV6DNS() error {
	if err := common.Ipt6.Insert("filter", "OUTPUT", 1, "-p", "udp", "--dport", "53", "-j", "REJECT"); err != nil {
		return e.New("disable dns request on ipv6 failed, ", err).WithPrefix(tagTools)
	}
	return nil
}

func EnableIPV6DNS() {
	_ = common.Ipt6.Delete("filter", "OUTPUT", "-p", "udp", "--dport", "53", "-j", "REJECT")
}

func RedirectDNS(port string) error {
	if err := common.Ipt.Insert("nat", "OUTPUT", 1, "-p", "udp", "-m", "owner", "!", "--gid-owner", common.CoreGid, "--dport", "53", "-j", "DNAT", "--to-destination", "127.0.0.1:"+port); err != nil {
		return e.New("redirect dns request failed, ", err).WithPrefix(tagTools)
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
