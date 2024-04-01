package trojan

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

// getTrojanSettingsObjectXray get xray Trojan SettingsObject
func getTrojanSettingsObjectXray(trojan *Trojan) map[string]interface{} {
	var serversObject []interface{}
	server := make(map[string]interface{})
	server["address"] = trojan.Server
	server["port"], _ = strconv.Atoi(trojan.Port)
	server["password"] = trojan.Password
	server["level"] = 0
	serversObject = append(serversObject, server)

	settingsObject := make(map[string]interface{})
	settingsObject["servers"] = serversObject
	return settingsObject
}

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(trojan *Trojan) map[string]interface{} {
	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["network"] = trojan.Network
	switch trojan.Network {
	case "tcp":
		tcpSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		switch trojan.Type {
		case "http":
			requestObject := make(map[string]interface{})
			headers := make(map[string]interface{})
			var connection []interface{}
			connection = append(connection, "keep-alive")
			var host []interface{}
			host = append(host, trojan.Host)
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
			headerObject["type"] = trojan.Type
			headerObject["request"] = requestObject
		default:
			headerObject["type"] = "none"
		}
		tcpSettingsObject["header"] = headerObject
		streamSettingsObject["tcpSettings"] = tcpSettingsObject
	case "kcp":
		kcpSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = trojan.Type
		kcpSettingsObject["congestion"] = false
		kcpSettingsObject["downlinkCapacity"] = 100
		kcpSettingsObject["header"] = headerObject
		kcpSettingsObject["mtu"] = 1350
		kcpSettingsObject["readBufferSize"] = 1
		kcpSettingsObject["seed"] = trojan.Path
		kcpSettingsObject["tti"] = 50
		kcpSettingsObject["uplinkCapacity"] = 12
		kcpSettingsObject["writeBufferSize"] = 1
		streamSettingsObject["kcpSettings"] = kcpSettingsObject
	case "ws":
		wsSettingsObject := make(map[string]interface{})
		headersObject := make(map[string]interface{})
		headersObject["Host"] = trojan.Host
		wsSettingsObject["headers"] = headersObject
		wsSettingsObject["path"] = trojan.Path
		streamSettingsObject["wsSettings"] = wsSettingsObject
	case "http", "h2":
		httpSettingsObject := make(map[string]interface{})
		var host []interface{}
		host = append(host, trojan.Host)
		httpSettingsObject["host"] = host
		httpSettingsObject["path"] = trojan.Path
		streamSettingsObject["httpSettings"] = httpSettingsObject
	case "httpupgrade":
		httpupgradeSettingsObject := make(map[string]interface{})
		var host []interface{}
		host = append(host, trojan.Host)
		httpupgradeSettingsObject["host"] = host
		httpupgradeSettingsObject["path"] = trojan.Path
		streamSettingsObject["httpupgrade"] = httpupgradeSettingsObject
	case "quic":
		quicSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = trojan.Type
		quicSettingsObject["header"] = headerObject
		quicSettingsObject["key"] = trojan.Path
		quicSettingsObject["security"] = trojan.Host
		streamSettingsObject["quicSettings"] = quicSettingsObject
	case "grpc":
		grpcSettingsObject := make(map[string]interface{})
		if trojan.Type == "multi" {
			grpcSettingsObject["multiMode"] = true
		} else {
			grpcSettingsObject["multiMode"] = false
		}
		grpcSettingsObject["authority"] = trojan.Host
		grpcSettingsObject["serviceName"] = trojan.Path
		streamSettingsObject["grpcSettings"] = grpcSettingsObject
	}
	streamSettingsObject["security"] = trojan.Security
	switch trojan.Security {
	case "tls":
		tlsSettingsObject := make(map[string]interface{})
		var alpn []interface{}
		alpnSlice := strings.Split(trojan.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsSettingsObject["alpn"] = alpn
			}
		}
		tlsSettingsObject["allowInsecure"] = false
		tlsSettingsObject["fingerprint"] = trojan.FingerPrint
		tlsSettingsObject["serverName"] = trojan.Sni
		streamSettingsObject["tlsSettings"] = tlsSettingsObject
	case "reality":
		realitySettingsObject := make(map[string]interface{})
		realitySettingsObject["allowInsecure"] = false
		realitySettingsObject["fingerprint"] = trojan.FingerPrint
		realitySettingsObject["serverName"] = trojan.Sni
		realitySettingsObject["publicKey"] = trojan.PublicKey
		realitySettingsObject["shortId"] = trojan.ShortId
		realitySettingsObject["spiderX"] = trojan.SpiderX
		streamSettingsObject["realitySettings"] = realitySettingsObject
	}
	sockoptObject := make(map[string]interface{})
	sockoptObject["domainStrategy"] = "UseIP"
	streamSettingsObject["sockopt"] = sockoptObject
	return streamSettingsObject
}
