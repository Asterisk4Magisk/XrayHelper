package vmess

import "strconv"

// getMuxObjectXray get xray MuxObject
func getMuxObjectXray(enabled bool) map[string]interface{} {
	mux := make(map[string]interface{})
	mux["enabled"] = enabled
	return mux
}

// getVmessSettingsObjectXray get xray Vmess SettingsObject
func getVmessSettingsObjectXray(vmess *Vmess) map[string]interface{} {
	var vnextsObject []interface{}
	vnext := make(map[string]interface{})
	vnext["address"] = vmess.Address
	vnext["port"], _ = strconv.Atoi(vmess.Port)

	var usersObject []interface{}
	user := make(map[string]interface{})
	user["id"] = vmess.Id
	user["alterId"], _ = strconv.Atoi(vmess.AlterId)
	user["security"] = vmess.Security
	user["level"] = 0
	usersObject = append(usersObject, user)

	vnext["users"] = usersObject
	vnextsObject = append(vnextsObject, vnext)
	settingsObject := make(map[string]interface{})
	settingsObject["vnext"] = vnextsObject
	return settingsObject
}

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(vmess *Vmess) map[string]interface{} {
	// TODO
	return nil
	//sockoptObject := make(map[string]interface{})
	//sockoptObject["domainStrategy"] = "UseIP"
	//
	//streamSettingsObject := make(map[string]interface{})
	//streamSettingsObject["network"] = vmess.Network
	//streamSettingsObject["security"] = vmess.Tls
	//
	//streamSettingsObject["sockopt"] = sockoptObject
	//return streamSettingsObject
}
