package tools

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"strconv"
	"strings"
)

const tagTools = "tools"

func GetUid(pkgInfo string) []string {
	var userId int
	var pkgUserId []string
	info := strings.Split(pkgInfo, ":")
	if len(info) == 2 {
		userId, _ = strconv.Atoi(info[1])
	}
	for pkgStr, pkgIdStr := range builds.PackageMap {
		if common.WildcardMatch(pkgStr, info[0]) {
			pkgId, _ := strconv.Atoi(pkgIdStr)
			pkgUserIdStr := strconv.Itoa(userId*100000 + pkgId)
			pkgUserId = append(pkgUserId, pkgUserIdStr)
		}
	}
	return pkgUserId
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

func EnableForward(device string) error {
	if err := common.Ipt.Insert("filter", "FORWARD", 1, "-i", device, "-j", "ACCEPT"); err != nil {
		return e.New("enable ipv4 forward for "+device+" incoming failed, ", err).WithPrefix(tagTools)
	}
	if err := common.Ipt.Insert("filter", "FORWARD", 1, "-o", device, "-j", "ACCEPT"); err != nil {
		return e.New("enable ipv4 forward for "+device+" outgoing failed, ", err).WithPrefix(tagTools)
	}
	if err := common.Ipt6.Insert("filter", "FORWARD", 1, "-i", device, "-j", "ACCEPT"); err != nil {
		return e.New("enable ipv6 forward for "+device+" incoming failed, ", err).WithPrefix(tagTools)
	}
	if err := common.Ipt6.Insert("filter", "FORWARD", 1, "-o", device, "-j", "ACCEPT"); err != nil {
		return e.New("enable ipv6 forward for "+device+" outgoing failed, ", err).WithPrefix(tagTools)
	}
	return nil
}

func DisableForward(device string) {
	_ = common.Ipt.Delete("filter", "FORWARD", "-i", device, "-j", "ACCEPT")
	_ = common.Ipt.Delete("filter", "FORWARD", "-o", device, "-j", "ACCEPT")
	_ = common.Ipt6.Delete("filter", "FORWARD", "-i", device, "-j", "ACCEPT")
	_ = common.Ipt6.Delete("filter", "FORWARD", "-o", device, "-j", "ACCEPT")
}
