package vmess

import (
	"strconv"
	"strings"
)

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
	vnext["address"] = vmess.Server
	vnext["port"], _ = strconv.Atoi(string(vmess.Port))

	var usersObject []interface{}
	user := make(map[string]interface{})
	user["id"] = vmess.Id
	user["alterId"], _ = strconv.Atoi(string(vmess.AlterId))
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
	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["network"] = vmess.Network
	switch vmess.Network {
	case "tcp":
		tcpSettingsObject := make(map[string]interface{})
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
		tcpSettingsObject["header"] = headerObject
		streamSettingsObject["tcpSettings"] = tcpSettingsObject
	case "kcp":
		kcpSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = vmess.Type
		kcpSettingsObject["congestion"] = false
		kcpSettingsObject["downlinkCapacity"] = 100
		kcpSettingsObject["header"] = headerObject
		kcpSettingsObject["mtu"] = 1350
		kcpSettingsObject["readBufferSize"] = 1
		kcpSettingsObject["seed"] = vmess.Path
		kcpSettingsObject["tti"] = 50
		kcpSettingsObject["uplinkCapacity"] = 12
		kcpSettingsObject["writeBufferSize"] = 1
		streamSettingsObject["kcpSettings"] = kcpSettingsObject
	case "ws":
		wsSettingsObject := make(map[string]interface{})
		headersObject := make(map[string]interface{})
		headersObject["Host"] = vmess.Host
		wsSettingsObject["headers"] = headersObject
		wsSettingsObject["path"] = vmess.Path
		streamSettingsObject["wsSettings"] = wsSettingsObject
	case "h2":
		httpSettingsObject := make(map[string]interface{})
		var host []interface{}
		host = append(host, vmess.Host)
		httpSettingsObject["host"] = host
		httpSettingsObject["path"] = vmess.Path
		streamSettingsObject["httpSettings"] = httpSettingsObject
	case "quic":
		quicSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = vmess.Type
		quicSettingsObject["header"] = headerObject
		quicSettingsObject["key"] = vmess.Path
		quicSettingsObject["security"] = vmess.Host
		streamSettingsObject["quicSettings"] = quicSettingsObject
	case "grpc":
		grpcSettingsObject := make(map[string]interface{})
		if vmess.Type == "multi" {
			grpcSettingsObject["multiMode"] = true
		} else {
			grpcSettingsObject["multiMode"] = false
		}
		grpcSettingsObject["serviceName"] = vmess.Path
		streamSettingsObject["grpcSettings"] = grpcSettingsObject
	}
	streamSettingsObject["security"] = vmess.Tls
	if len(vmess.Tls) > 0 {
		tlsSettingsObject := make(map[string]interface{})
		var alpn []interface{}
		alpnSlice := strings.Split(string(vmess.Alpn), ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsSettingsObject["alpn"] = alpn
			}
		}
		tlsSettingsObject["allowInsecure"] = false
		tlsSettingsObject["fingerprint"] = vmess.FingerPrint
		tlsSettingsObject["serverName"] = vmess.Sni
		streamSettingsObject["tlsSettings"] = tlsSettingsObject
	}
	sockoptObject := make(map[string]interface{})
	sockoptObject["domainStrategy"] = "UseIP"
	streamSettingsObject["sockopt"] = sockoptObject
	return streamSettingsObject
}
