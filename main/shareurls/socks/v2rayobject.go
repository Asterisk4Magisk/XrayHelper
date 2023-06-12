package socks

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

// getSocksSettingsObjectV2ray get v2ray Socks SettingsObject
func getSocksSettingsObjectV2ray(socks *Socks) map[string]interface{} {
	// v2ray v5 not support auth
	settingsObject := make(map[string]interface{})
	settingsObject["address"] = socks.Server
	settingsObject["port"], _ = strconv.Atoi(socks.Port)
	return settingsObject
}
