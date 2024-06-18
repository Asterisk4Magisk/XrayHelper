package trojan

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getTrojanSettingsObjectXray get xray Trojan SettingsObject
func getTrojanSettingsObjectXray(trojan *Trojan) serial.OrderedMap {
	var serverArray serial.OrderedArray
	var server serial.OrderedMap
	server.Set("address", trojan.Server)
	port, _ := strconv.Atoi(trojan.Port)
	server.Set("port", port)
	server.Set("password", trojan.Password)
	server.Set("level", 0)
	serverArray = append(serverArray, server)

	var settingsObject serial.OrderedMap
	settingsObject.Set("servers", serverArray)
	return settingsObject
}
