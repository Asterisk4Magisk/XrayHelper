package socks

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getMuxObjectXray get xray MuxObject
func getMuxObjectXray(enabled bool) serial.OrderedMap {
	var mux serial.OrderedMap
	mux.Set("enabled", enabled)
	return mux
}

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(network string) serial.OrderedMap {
	var sockoptObject serial.OrderedMap
	sockoptObject.Set("domainStrategy", "UseIP")

	var streamSettingsObject serial.OrderedMap
	streamSettingsObject.Set("network", network)
	streamSettingsObject.Set("sockopt", sockoptObject)
	return streamSettingsObject
}

// getSocksSettingsObjectXray get xray Socks SettingsObject
func getSocksSettingsObjectXray(socks *Socks) serial.OrderedMap {
	var serverArray serial.OrderedArray
	var server serial.OrderedMap
	// v2rayNg share "null" user for no auth socks server
	if socks.User != "null" {
		var userArray serial.OrderedArray
		var user serial.OrderedMap
		user.Set("user", socks.User)
		user.Set("pass", socks.Password)
		user.Set("level", 0)
		userArray = append(userArray, user)
		server.Set("users", userArray)
	}
	server.Set("address", socks.Server)
	port, _ := strconv.Atoi(socks.Port)
	server.Set("port", port)
	serverArray = append(serverArray, server)
	var settingsObject serial.OrderedMap
	settingsObject.Set("servers", serverArray)
	return settingsObject
}
