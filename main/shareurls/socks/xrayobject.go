package socks

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

// getSocksSettingsObjectXray get xray Socks SettingsObject
func getSocksSettingsObjectXray(socks *Socks) map[string]interface{} {
	var serversObject []interface{}
	server := make(map[string]interface{})
	// v2rayNg share "null" user for no auth socks server
	if socks.User != "null" {
		var usersObject []interface{}
		user := make(map[string]interface{})
		user["user"] = socks.User
		user["pass"] = socks.Password
		user["level"] = 0
		usersObject = append(usersObject, user)
		server["users"] = usersObject
	}
	server["address"] = socks.Address
	server["port"], _ = strconv.Atoi(socks.Port)
	serversObject = append(serversObject, server)
	settingsObject := make(map[string]interface{})
	settingsObject["servers"] = serversObject
	return settingsObject
}
