package shadowsocks

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
	server["Address"] = ss.Address
	server["Port"] = ss.Port
	server["Method"] = ss.Method
	server["Password"] = ss.Password
	server["level"] = 0
	serversObject = append(serversObject, server)

	settingsObject := make(map[string]interface{})
	settingsObject["servers"] = serversObject
	return settingsObject
}
