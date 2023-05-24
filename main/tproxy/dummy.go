package tproxy

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/errors"
	"XrayHelper/main/log"
	"XrayHelper/main/utils"
	"bytes"
)

const (
	dummyDevice  = "xdummy"
	dummyIp      = "fd01:5ca1:ab1e:8d97:497f:8b48:b9aa:85cd/128"
	dummyMarkId  = "164"
	dummyTableId = "164"
)

func createDummyDevice() error {
	var errMsg bytes.Buffer
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "add", dummyDevice, "type", "dummy").Run()
	if errMsg.Len() > 0 {
		return errors.New("add dummy device failed, ", errMsg.String()).WithPrefix("dummy")
	}
	errMsg.Reset()
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "addr", "add", dummyIp, "dev", dummyDevice).Run()
	if errMsg.Len() > 0 {
		return errors.New("add dummy ip failed, ", errMsg.String()).WithPrefix("dummy")
	}
	errMsg.Reset()
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "set", dummyDevice, "up").Run()
	if errMsg.Len() > 0 {
		return errors.New("set dummy up failed, ", errMsg.String()).WithPrefix("dummy")
	}
	return nil
}

func removeDummyDevice() {
	var errMsg bytes.Buffer
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "set", dummyDevice, "down").Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("set dummy up down: " + errMsg.String())
	}
	errMsg.Reset()
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "link", "del", dummyDevice, "type", "dummy").Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("delete dummy device: " + errMsg.String())
	}
}

func addDummyRoute() error {
	var errMsg bytes.Buffer
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "add", "not", "from", "all", "fwmark", dummyMarkId, "table", dummyTableId).Run()
	if errMsg.Len() > 0 {
		return errors.New("add dummy rule failed, ", errMsg.String()).WithPrefix("dummy")
	}
	errMsg.Reset()
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "add", "default", "dev", dummyDevice, "table", dummyTableId).Run()
	if errMsg.Len() > 0 {
		return errors.New("add dummy route failed, ", errMsg.String()).WithPrefix("dummy")
	}
	return nil
}

func deleteDummyRoute() {
	var errMsg bytes.Buffer
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "rule", "del", "not", "from", "all", "fwmark", dummyMarkId, "table", dummyTableId).Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("delete dummy rule: " + errMsg.String())
	}
	errMsg.Reset()
	utils.NewExternal(0, nil, &errMsg, "ip", "-6", "route", "del", "default", "dev", dummyDevice, "table", dummyTableId).Run()
	if errMsg.Len() > 0 {
		log.HandleDebug("delete dummy route: " + errMsg.String())
	}
}

func createDummyOutputChain() error {
	if err := ipt6.NewChain("mangle", "DUMMY"); err != nil {
		return errors.New("create ipv6 mangle chain DUMMY failed, ", err).WithPrefix("dummy")
	}
	if err := ipt6.Append("mangle", "DUMMY", "-p", "tcp", "-j", "MARK", "--set-mark", dummyMarkId); err != nil {
		return errors.New("set mark on tcp mangle chain DUMMY failed, ", err).WithPrefix("dummy")
	}
	if err := ipt6.Append("mangle", "DUMMY", "-p", "udp", "-j", "MARK", "--set-mark", dummyMarkId); err != nil {
		return errors.New("set mark on udp mangle chain DUMMY failed, ", err).WithPrefix("dummy")
	}
	if err := ipt6.Append("mangle", "OUTPUT", "-j", "DUMMY"); err != nil {
		return errors.New("apply ipv6 mangle chain DUMMY on OUTPUT failed, ", err).WithPrefix("dummy")
	}
	return nil
}

func createDummyPreroutingChain() error {
	if err := ipt6.NewChain("mangle", "XD"); err != nil {
		return errors.New("create ipv6 mangle chain XD failed, ", err).WithPrefix("dummy")
	}
	if err := ipt6.Append("mangle", "XD", "-i", dummyDevice, "-p", "tcp", "-j", "TPROXY", "--on-ip", "::1", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", dummyMarkId); err != nil {
		return errors.New("set mark on tcp mangle chain XD failed, ", err).WithPrefix("dummy")
	}
	if err := ipt6.Append("mangle", "XD", "-i", dummyDevice, "-p", "udp", "-j", "TPROXY", "--on-ip", "::1", "--on-port", builds.Config.Proxy.TproxyPort, "--tproxy-mark", dummyMarkId); err != nil {
		return errors.New("set mark on udp mangle chain XD failed, ", err).WithPrefix("dummy")
	}
	if err := ipt6.Append("mangle", "PREROUTING", "-j", "XD"); err != nil {
		return errors.New("apply ipv6 mangle chain XD on PREROUTING failed, ", err).WithPrefix("dummy")
	}
	return nil
}

func cleanDummyChain() {
	_ = ipt6.Delete("mangle", "OUTPUT", "-j", "DUMMY")
	_ = ipt6.Delete("mangle", "PREROUTING", "-j", "XD")
	_ = ipt6.ClearAndDeleteChain("mangle", "DUMMY")
	_ = ipt6.ClearAndDeleteChain("mangle", "XD")
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
