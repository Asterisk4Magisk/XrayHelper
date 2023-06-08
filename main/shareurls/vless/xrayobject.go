package vless

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

// getVLESSSettingsObjectXray get xray VLESS SettingsObject
func getVLESSSettingsObjectXray(vless *VLESS) map[string]interface{} {
	var vnextsObject []interface{}
	vnext := make(map[string]interface{})
	vnext["address"] = vless.Server
	vnext["port"], _ = strconv.Atoi(vless.Port)

	var usersObject []interface{}
	user := make(map[string]interface{})
	user["id"] = vless.Id
	user["flow"] = vless.Flow
	user["encryption"] = vless.Encryption
	user["level"] = 0
	usersObject = append(usersObject, user)

	vnext["users"] = usersObject
	vnextsObject = append(vnextsObject, vnext)
	settingsObject := make(map[string]interface{})
	settingsObject["vnext"] = vnextsObject
	return settingsObject
}

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(vless *VLESS) map[string]interface{} {
	streamSettingsObject := make(map[string]interface{})
	streamSettingsObject["network"] = vless.Network
	switch vless.Network {
	case "tcp":
		tcpSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		switch vless.Type {
		case "http":
			requestObject := make(map[string]interface{})
			headers := make(map[string]interface{})
			var connection []interface{}
			connection = append(connection, "keep-alive")
			var host []interface{}
			host = append(host, vless.Host)
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
			headerObject["type"] = vless.Type
			headerObject["request"] = requestObject
		default:
			headerObject["type"] = "none"
		}
		tcpSettingsObject["header"] = headerObject
		streamSettingsObject["tcpSettings"] = tcpSettingsObject
	case "kcp":
		kcpSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = vless.Type
		kcpSettingsObject["congestion"] = false
		kcpSettingsObject["downlinkCapacity"] = 100
		kcpSettingsObject["header"] = headerObject
		kcpSettingsObject["mtu"] = 1350
		kcpSettingsObject["readBufferSize"] = 1
		kcpSettingsObject["seed"] = vless.Path
		kcpSettingsObject["tti"] = 50
		kcpSettingsObject["uplinkCapacity"] = 12
		kcpSettingsObject["writeBufferSize"] = 1
		streamSettingsObject["kcpSettings"] = kcpSettingsObject
	case "ws":
		wsSettingsObject := make(map[string]interface{})
		headersObject := make(map[string]interface{})
		headersObject["Host"] = vless.Host
		wsSettingsObject["headers"] = headersObject
		wsSettingsObject["path"] = vless.Path
		streamSettingsObject["wsSettings"] = wsSettingsObject
	case "h2":
		httpSettingsObject := make(map[string]interface{})
		var host []interface{}
		host = append(host, vless.Host)
		httpSettingsObject["host"] = host
		httpSettingsObject["path"] = vless.Path
		streamSettingsObject["httpSettings"] = httpSettingsObject
	case "quic":
		quicSettingsObject := make(map[string]interface{})
		headerObject := make(map[string]interface{})
		headerObject["type"] = vless.Type
		quicSettingsObject["header"] = headerObject
		quicSettingsObject["key"] = vless.Path
		quicSettingsObject["security"] = vless.Host
		streamSettingsObject["quicSettings"] = quicSettingsObject
	case "grpc":
		grpcSettingsObject := make(map[string]interface{})
		if vless.Type == "multi" {
			grpcSettingsObject["multiMode"] = true
		} else {
			grpcSettingsObject["multiMode"] = false
		}
		grpcSettingsObject["serviceName"] = vless.Path
		streamSettingsObject["grpcSettings"] = grpcSettingsObject
	}
	streamSettingsObject["security"] = vless.Security
	switch vless.Security {
	case "tls":
		tlsSettingsObject := make(map[string]interface{})
		var alpn []interface{}
		alpnSlice := strings.Split(vless.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsSettingsObject["alpn"] = alpn
			}
		}
		tlsSettingsObject["allowInsecure"] = false
		tlsSettingsObject["fingerprint"] = vless.FingerPrint
		tlsSettingsObject["serverName"] = vless.Sni
		streamSettingsObject["tlsSettings"] = tlsSettingsObject
	case "reality":
		realitySettingsObject := make(map[string]interface{})
		realitySettingsObject["allowInsecure"] = false
		realitySettingsObject["fingerprint"] = vless.FingerPrint
		realitySettingsObject["serverName"] = vless.Sni
		realitySettingsObject["publicKey"] = vless.PublicKey
		realitySettingsObject["shortId"] = vless.ShortId
		realitySettingsObject["spiderX"] = vless.SpiderX
		streamSettingsObject["realitySettings"] = realitySettingsObject
	}
	sockoptObject := make(map[string]interface{})
	sockoptObject["domainStrategy"] = "UseIP"
	streamSettingsObject["sockopt"] = sockoptObject
	return streamSettingsObject
}
