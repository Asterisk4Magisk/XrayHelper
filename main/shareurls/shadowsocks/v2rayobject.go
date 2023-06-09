package shadowsocks

import "strconv"

// getMuxObjectV2ray get v2ray MuxObject
func getMuxObjectV2ray(enabled bool) map[string]interface{} {
	mux := make(map[string]interface{})
	mux["enabled"] = enabled
	return mux
}

// getStreamSettingsObjectV2ray get v2ray StreamSettingsObject
func getStreamSettingsObjectV2ray(transport string) map[string]interface{} {
	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["transport"] = transport
	return streamSettingsObject
}

// getShadowsocksSettingsObjectV2ray get v2ray Shadowsocks SettingsObject
func getShadowsocksSettingsObjectV2ray(ss *Shadowsocks) map[string]interface{} {
	settingsObject := make(map[string]interface{})
	settingsObject["address"] = ss.Server
	settingsObject["port"], _ = strconv.Atoi(ss.Port)
	settingsObject["method"] = ss.Method
	settingsObject["password"] = ss.Password
	return settingsObject
}
