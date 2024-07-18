package vless

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getVLESSSettingsObjectXray get xray VLESS SettingsObject
func getVLESSSettingsObjectXray(vless *VLESS) serial.OrderedMap {
	var vnextArray serial.OrderedArray
	var vnext serial.OrderedMap
	vnext.Set("address", vless.Server)
	port, _ := strconv.Atoi(vless.Port)
	vnext.Set("port", port)

	var userArray serial.OrderedArray
	var user serial.OrderedMap
	user.Set("id", vless.Id)
	user.Set("flow", vless.Flow)
	user.Set("encryption", vless.Encryption)
	user.Set("level", 0)
	userArray = append(userArray, user)

	vnext.Set("users", userArray)
	vnextArray = append(vnextArray, vnext)
	var settingsObject serial.OrderedMap
	settingsObject.Set("vnext", vnextArray)
	return settingsObject
}
