package tproxy

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"bytes"
)

const tagDummy = "dummy"

func createDummyDevice() error {
	var errMsg bytes.Buffer
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "add", common.DummyDevice, "type", "dummy").Run()
	if errMsg.Len() > 0 {
		return e.New("add dummy device failed, ", errMsg.String()).WithPrefix(tagDummy)
	}
	errMsg.Reset()
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "addr", "add", common.DummyIp, "dev", common.DummyDevice).Run()
	if errMsg.Len() > 0 {
		return e.New("add dummy ip failed, ", errMsg.String()).WithPrefix(tagDummy)
	}
	errMsg.Reset()
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "set", common.DummyDevice, "up").Run()
	if errMsg.Len() > 0 {
		return e.New("set dummy up failed, ", errMsg.String()).WithPrefix(tagDummy)
	}
	return nil
}

func removeDummyDevice() {
	var errMsg bytes.Buffer
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "set", common.DummyDevice, "down").Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("set dummy up down: " + errMsg.String())
	}
	errMsg.Reset()
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "del", common.DummyDevice, "type", "dummy").Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("delete dummy device: " + errMsg.String())
	}
}

func addDummyRoute() error {
	var errMsg bytes.Buffer
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "add", "not", "from", "all", "fwmark", common.DummyMarkId, "table", common.DummyTableId).Run()
	if errMsg.Len() > 0 {
		return e.New("add dummy rule failed, ", errMsg.String()).WithPrefix(tagDummy)
	}
	errMsg.Reset()
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "add", "local", "default", "dev", common.DummyDevice, "table", common.DummyTableId).Run()
	if errMsg.Len() > 0 {
		return e.New("add dummy route failed, ", errMsg.String()).WithPrefix(tagDummy)
	}
	return nil
}

func deleteDummyRoute() {
	var errMsg bytes.Buffer
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "del", "not", "from", "all", "fwmark", common.DummyMarkId, "table", common.DummyTableId).Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("delete dummy rule: " + errMsg.String())
	}
	errMsg.Reset()
	common.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "del", "local", "default", "dev", common.DummyDevice, "table", common.DummyTableId).Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("delete dummy route: " + errMsg.String())
	}
}

func createDummyOutputChain() error {
	if err := common.Ipt6.NewChain("mangle", "DUMMY"); err != nil {
		return e.New("create ipv6 mangle chain DUMMY failed, ", err).WithPrefix(tagDummy)
	}
	if err := common.Ipt6.Append("mangle", "DUMMY", "-p", "tcp", "-j", "MARK", "--set-mark", common.DummyMarkId); err != nil {
		return e.New("set mark on tcp mangle chain DUMMY failed, ", err).WithPrefix(tagDummy)
	}
	if err := common.Ipt6.Append("mangle", "DUMMY", "-p", "udp", "-j", "MARK", "--set-mark", common.DummyMarkId); err != nil {
		return e.New("set mark on udp mangle chain DUMMY failed, ", err).WithPrefix(tagDummy)
	}
	if err := common.Ipt6.Append("mangle", "OUTPUT", "-j", "DUMMY"); err != nil {
		return e.New("apply ipv6 mangle chain DUMMY on OUTPUT failed, ", err).WithPrefix(tagDummy)
	}
	return nil
}

func createDummyPreroutingChain() error {
	if err := common.Ipt6.NewChain("mangle", "XD"); err != nil {
		return e.New("create ipv6 mangle chain XD failed, ", err).WithPrefix(tagDummy)
	}
	if err := common.Ipt6.Append("mangle", "XD", "-i", common.DummyDevice, "-p", "tcp", "-j", "TPROXY", "--on-ip", "::", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", common.DummyMarkId); err != nil {
		return e.New("set mark on tcp mangle chain XD failed, ", err).WithPrefix(tagDummy)
	}
	if err := common.Ipt6.Append("mangle", "XD", "-i", common.DummyDevice, "-p", "udp", "-j", "TPROXY", "--on-ip", "::", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", common.DummyMarkId); err != nil {
		return e.New("set mark on udp mangle chain XD failed, ", err).WithPrefix(tagDummy)
	}
	if err := common.Ipt6.Append("mangle", "PREROUTING", "-j", "XD"); err != nil {
		return e.New("apply ipv6 mangle chain XD on PREROUTING failed, ", err).WithPrefix(tagDummy)
	}
	return nil
}

func cleanDummyChain() {
	_ = common.Ipt6.Delete("mangle", "OUTPUT", "-j", "DUMMY")
	_ = common.Ipt6.Delete("mangle", "PREROUTING", "-j", "XD")
	_ = common.Ipt6.ClearAndDeleteChain("mangle", "DUMMY")
	_ = common.Ipt6.ClearAndDeleteChain("mangle", "XD")
}

func enableDummy() error {
	if err := createDummyDevice(); err != nil {
		return err
	}
	if err := addDummyRoute(); err != nil {
		return err
	}
	if err := createDummyPreroutingChain(); err != nil {
		return err
	}
	if err := createDummyOutputChain(); err != nil {
		return err
	}
	return nil
}

func disableDummy() {
	cleanDummyChain()
	deleteDummyRoute()
	removeDummyDevice()
}
