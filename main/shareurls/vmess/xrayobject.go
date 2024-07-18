package vmess

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getVmessSettingsObjectXray get xray Vmess SettingsObject
func getVmessSettingsObjectXray(vmess *Vmess) serial.OrderedMap {
	var vnextArray serial.OrderedArray
	var vnext serial.OrderedMap
	vnext.Set("address", vmess.Server)
	port, _ := strconv.Atoi(string(vmess.Port))
	vnext.Set("port", port)

	var userArray serial.OrderedArray
	var user serial.OrderedMap
	user.Set("id", vmess.Id)
	alterId, _ := strconv.Atoi(string(vmess.AlterId))
	user.Set("alterId", alterId)
	user.Set("security", vmess.Security)
	user.Set("level", 0)
	userArray = append(userArray, user)

	vnext.Set("users", userArray)
	vnextArray = append(vnextArray, vnext)
	var settingsObject serial.OrderedMap
	settingsObject.Set("vnext", vnextArray)
	return settingsObject
}
