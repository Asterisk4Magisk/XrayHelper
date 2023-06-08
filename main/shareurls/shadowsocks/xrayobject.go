package shadowsocks

import "strconv"

// getMuxObjectXray get xray MuxObject
func getMuxObjectXray(enabled bool) map[string]interface{} {
	mux := make(map[string]interface{})
	mux["enabled"] = enabled
	return mux
}

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(network string) map[string]interface{} {
	sockoptObject := make(map[string]interface{})
	sockoptObject["domainStrategy"] = "UseIP"

	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["network"] = network
	streamSettingsObject["sockopt"] = sockoptObject
	return streamSettingsObject
}

// getShadowsocksSettingsObjectXray get xray Shadowsocks SettingsObject
func getShadowsocksSettingsObjectXray(ss *Shadowsocks) map[string]interface{} {
	var serversObject []interface{}
	server := make(map[string]interface{})
	server["address"] = ss.Server
	server["port"], _ = strconv.Atoi(ss.Port)
	server["method"] = ss.Method
	server["password"] = ss.Password
	server["level"] = 0
	serversObject = append(serversObject, server)

	settingsObject := make(map[string]interface{})
	settingsObject["servers"] = serversObject
	return settingsObject
}
