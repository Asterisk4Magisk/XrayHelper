package vmessaead

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getVmessSettingsObjectXray get xray Vmess SettingsObject
func getVmessSettingsObjectXray(vmess *VmessAEAD) serial.OrderedMap {
	var vnextArray serial.OrderedArray
	var vnext serial.OrderedMap
	vnext.Set("address", vmess.Server)
	port, _ := strconv.Atoi(vmess.Port)
	vnext.Set("port", port)

	var userArray serial.OrderedArray
	var user serial.OrderedMap
	user.Set("id", vmess.Id)
	user.Set("alterId", 0)
	user.Set("security", vmess.Encryption)
	user.Set("level", 0)
	userArray = append(userArray, user)

	vnext.Set("users", userArray)
	vnextArray = append(vnextArray, vnext)
	var settingsObject serial.OrderedMap
	settingsObject.Set("vnext", vnextArray)
	return settingsObject
}
