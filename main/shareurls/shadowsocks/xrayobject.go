package shadowsocks

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(network string) serial.OrderedMap {
	var sockoptObject serial.OrderedMap
	sockoptObject.Set("domainStrategy", "UseIP")

	var streamSettingsObject serial.OrderedMap
	streamSettingsObject.Set("network", network)
	streamSettingsObject.Set("sockopt", sockoptObject)
	return streamSettingsObject
}

// getShadowsocksSettingsObjectXray get xray Shadowsocks SettingsObject
func getShadowsocksSettingsObjectXray(ss *Shadowsocks) serial.OrderedMap {
	var serverArray serial.OrderedArray
	var server serial.OrderedMap
	server.Set("address", ss.Server)
	port, _ := strconv.Atoi(ss.Port)
	server.Set("port", port)
	server.Set("method", ss.Method)
	server.Set("password", ss.Password)
	server.Set("level", 0)
	serverArray = append(serverArray, server)

	var settingsObject serial.OrderedMap
	settingsObject.Set("servers", serverArray)
	return settingsObject
}
