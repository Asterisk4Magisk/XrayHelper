package sharelinks

// getMuxObject get xray MuxObject
func getMuxObject(enabled bool) map[string]interface{} {
	mux := make(map[string]interface{})
	mux["enabled"] = enabled
	return mux
}

// getStreamSettingsObject get xray StreamSettingsObject
func getStreamSettingsObject(network string) map[string]interface{} {
	sockoptObject := make(map[string]interface{})
	sockoptObject["domainStrategy"] = "UseIP"

	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["network"] = network
	streamSettingsObject["sockopt"] = sockoptObject
	return streamSettingsObject
}

// getShadowsocksSettingsObject get xray Shadowsocks SettingsObject
func getShadowsocksSettingsObject(ss *Shadowsocks) map[string]interface{} {
	var serversObject []interface{}
	server := make(map[string]interface{})
	server["address"] = ss.address
	server["port"] = ss.port
	server["method"] = ss.method
	server["password"] = ss.password
	server["level"] = 0
	serversObject = append(serversObject, server)

	settingsObject := make(map[string]interface{})
	settingsObject["servers"] = serversObject
	return settingsObject
}
