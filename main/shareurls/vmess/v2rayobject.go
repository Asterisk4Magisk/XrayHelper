package vmess

import (
	"strconv"
	"strings"
)

// getMuxObjectV2ray get v2ray MuxObject
func getMuxObjectV2ray(enabled bool) map[string]interface{} {
	mux := make(map[string]interface{})
	mux["enabled"] = enabled
	return mux
}

// getVmessSettingsObjectV2ray get v2ray Vmess SettingsObject
func getVmessSettingsObjectV2ray(vmess *Vmess) map[string]interface{} {
	settingsObject := make(map[string]interface{})
	settingsObject["address"] = vmess.Server
	settingsObject["port"], _ = strconv.Atoi(vmess.Port)
	settingsObject["uuid"] = vmess.Id
	return settingsObject
}

// getStreamSettingsObjectV2ray get v2ray StreamSettingsObject
func getStreamSettingsObjectV2ray(vmess *Vmess) map[string]interface{} {
	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["transport"] = vmess.Network
	switch vmess.Network {
	case "tcp":
		transportSettingsObject := make(map[string]interface{})
		transportSettingsObject["acceptProxyProtocol"] = false
		headerObject := make(map[string]interface{})
		switch vmess.Type {
		case "http":
			requestObject := make(map[string]interface{})
			headers := make(map[string]interface{})
			var connection []interface{}
			connection = append(connection, "keep-alive")
			var host []interface{}
			host = append(host, vmess.Host)
			var acceptEncoding []interface{}
			acceptEncoding = append(acceptEncoding, "gzip, deflate")
			var userAgent []interface{}
			userAgent = append(userAgent, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
				"Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1")
			headers["Connection"] = connection
			headers["Host"] = host
			headers["Pragma"] = "no-cache"
			headers["Accept-Encoding"] = acceptEncoding
			headers["User-Agent"] = userAgent
			requestObject["headers"] = headers
			headerObject["type"] = vmess.Type
			headerObject["request"] = requestObject
		default:
			headerObject["type"] = "none"
		}
		transportSettingsObject["header"] = headerObject
		streamSettingsObject["transportSettings"] = transportSettingsObject
	case "kcp":
		transportSettingsObject := make(map[string]interface{})
		transportSettingsObject["congestion"] = false
		transportSettingsObject["downlinkCapacity"] = 100
		transportSettingsObject["mtu"] = 1350
		transportSettingsObject["readBufferSize"] = 1
		transportSettingsObject["seed"] = vmess.Path
		transportSettingsObject["tti"] = 50
		transportSettingsObject["uplinkCapacity"] = 12
		transportSettingsObject["writeBufferSize"] = 1
		streamSettingsObject["transportSettings"] = transportSettingsObject
	case "ws":
		transportSettingsObject := make(map[string]interface{})
		headersObject := make(map[string]interface{})
		headersObject["Host"] = vmess.Host
		transportSettingsObject["headers"] = headersObject
		transportSettingsObject["path"] = vmess.Path
		streamSettingsObject["transportSettings"] = transportSettingsObject
	case "meek":
		transportSettingsObject := make(map[string]interface{})
		transportSettingsObject["url"] = vmess.Path
		streamSettingsObject["transportSettings"] = transportSettingsObject
	case "quic":
		transportSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = vmess.Type
		transportSettingsObject["header"] = headerObject
		transportSettingsObject["key"] = vmess.Path
		transportSettingsObject["security"] = vmess.Host
		streamSettingsObject["transportSettings"] = transportSettingsObject
	case "grpc":
		transportSettingsObject := make(map[string]interface{})
		if vmess.Type == "multi" {
			transportSettingsObject["multiMode"] = true
		} else {
			transportSettingsObject["multiMode"] = false
		}
		transportSettingsObject["serviceName"] = vmess.Path
		streamSettingsObject["transportSettings"] = transportSettingsObject
	}
	streamSettingsObject["security"] = vmess.Security
	switch vmess.Security {
	case "tls":
		securitySettings := make(map[string]interface{})
		var alpn []interface{}
		alpnSlice := strings.Split(vmess.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				securitySettings["nextProtocol"] = alpn
			}
		}
		securitySettings["disableSystemRoot"] = false
		securitySettings["serverName"] = vmess.Sni
		streamSettingsObject["securitySettings"] = securitySettings
	case "utls":
		securitySettings := make(map[string]interface{})
		tlsConfig := make(map[string]interface{})
		var alpn []interface{}
		alpnSlice := strings.Split(vmess.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsConfig["nextProtocol"] = alpn
			}
		}
		tlsConfig["disableSystemRoot"] = false
		tlsConfig["serverName"] = vmess.Sni
		securitySettings["tlsConfig"] = tlsConfig
		securitySettings["imitate"] = vmess.FingerPrint
		streamSettingsObject["securitySettings"] = securitySettings
	}
	return streamSettingsObject
}
