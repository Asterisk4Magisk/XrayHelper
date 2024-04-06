package vless

import (
	"XrayHelper/main/serial"
	"strconv"
	"strings"
)

// getMuxObjectXray get xray MuxObject
func getMuxObjectXray(enabled bool) serial.OrderedMap {
	var mux serial.OrderedMap
	mux.Set("enabled", enabled)
	return mux
}

// getVLESSSettingsObjectXray get xray VLESS SettingsObject
func getVLESSSettingsObjectXray(vless *VLESS) serial.OrderedMap {
	var vnextArray serial.OrderedArray
	var vnext serial.OrderedMap
	vnext.Set("address", vless.Server)
	port, _ := strconv.Atoi(vless.Port)
	vnext.Set("port", port)

	var userArray serial.OrderedArray
	var user serial.OrderedMap
	user.Set("id", vless.Id)
	user.Set("flow", vless.Flow)
	user.Set("encryption", vless.Encryption)
	user.Set("level", 0)
	userArray = append(userArray, user)

	vnext.Set("users", userArray)
	vnextArray = append(vnextArray, vnext)
	var settingsObject serial.OrderedMap
	settingsObject.Set("vnext", vnextArray)
	return settingsObject
}

// getStreamSettingsObjectXray get xray StreamSettingsObject
func getStreamSettingsObjectXray(vless *VLESS) serial.OrderedMap {
	var streamSettingsObject serial.OrderedMap
	streamSettingsObject.Set("network", vless.Network)
	switch vless.Network {
	case "tcp":
		var tcpSettingsObject serial.OrderedMap
		var headerObject serial.OrderedMap
		switch vless.Type {
		case "http":
			headerObject.Set("type", vless.Type)
			if len(vless.Host) > 0 {
				var requestObject serial.OrderedMap
				var headers serial.OrderedMap
				var host serial.OrderedArray
				host = append(host, vless.Host)
				var connection serial.OrderedArray
				connection = append(connection, "keep-alive")
				var acceptEncoding serial.OrderedArray
				acceptEncoding = append(acceptEncoding, "gzip, deflate")
				var userAgent serial.OrderedArray
				userAgent = append(userAgent, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
					"Mozilla/5.0 (iPhone; CPU iPhone OS 16_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Mobile/15E148 Safari/604.1")
				headers.Set("Host", host)
				headers.Set("Connection", connection)
				headers.Set("Pragma", "no-cache")
				headers.Set("Accept-Encoding", acceptEncoding)
				headers.Set("User-Agent", userAgent)
				requestObject.Set("headers", headers)
				headerObject.Set("request", requestObject)
			}
		default:
			headerObject.Set("type", "none")
		}
		tcpSettingsObject.Set("header", headerObject)
		streamSettingsObject.Set("tcpSettings", tcpSettingsObject)
	case "kcp":
		var kcpSettingsObject serial.OrderedMap
		if len(vless.Type) > 0 {
			var headerObject serial.OrderedMap
			headerObject.Set("type", vless.Type)
			kcpSettingsObject.Set("header", headerObject)
		}
		kcpSettingsObject.Set("congestion", false)
		kcpSettingsObject.Set("downlinkCapacity", 100)
		kcpSettingsObject.Set("mtu", 1350)
		kcpSettingsObject.Set("readBufferSize", 1)
		if len(vless.Path) > 0 {
			kcpSettingsObject.Set("seed", vless.Path)
		}
		kcpSettingsObject.Set("tti", 50)
		kcpSettingsObject.Set("uplinkCapacity", 12)
		kcpSettingsObject.Set("writeBufferSize", 1)
		streamSettingsObject.Set("kcpSettings", kcpSettingsObject)
	case "ws":
		var wsSettingsObject serial.OrderedMap
		if len(vless.Host) > 0 {
			var headersObject serial.OrderedMap
			headersObject.Set("Host", vless.Host)
			wsSettingsObject.Set("headers", headersObject)
		}
		if len(vless.Path) > 0 {
			wsSettingsObject.Set("path", vless.Path)
		}
		streamSettingsObject.Set("wsSettings", wsSettingsObject)
	case "http", "h2":
		var httpSettingsObject serial.OrderedMap
		if len(vless.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, vless.Host)
			httpSettingsObject.Set("host", host)
		}
		if len(vless.Path) > 0 {
			httpSettingsObject.Set("path", vless.Path)
		}
		streamSettingsObject.Set("httpSettings", httpSettingsObject)
	case "httpupgrade":
		var httpupgradeSettingsObject serial.OrderedMap
		if len(vless.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, vless.Host)
			httpupgradeSettingsObject.Set("host", host)
		}
		if len(vless.Path) > 0 {
			httpupgradeSettingsObject.Set("path", vless.Path)
		}
		streamSettingsObject.Set("httpupgrade", httpupgradeSettingsObject)
	case "quic":
		var quicSettingsObject serial.OrderedMap
		if len(vless.Type) > 0 {
			var headerObject serial.OrderedMap
			headerObject.Set("type", vless.Type)
			quicSettingsObject.Set("header", headerObject)
		}
		if len(vless.Path) > 0 {
			quicSettingsObject.Set("key", vless.Path)
		}
		if len(vless.Host) > 0 {
			quicSettingsObject.Set("security", vless.Host)
		}
		streamSettingsObject.Set("quicSettings", quicSettingsObject)
	case "grpc":
		var grpcSettingsObject serial.OrderedMap
		if vless.Type == "multi" {
			grpcSettingsObject.Set("multiMode", true)
		} else {
			grpcSettingsObject.Set("multiMode", false)
		}
		if len(vless.Host) > 0 {
			grpcSettingsObject.Set("authority", vless.Host)
		}
		if len(vless.Path) > 0 {
			grpcSettingsObject.Set("serviceName", vless.Path)
		}
		streamSettingsObject.Set("grpcSettings", grpcSettingsObject)
	}
	streamSettingsObject.Set("security", vless.Security)
	switch vless.Security {
	case "tls":
		var tlsSettingsObject serial.OrderedMap
		var alpn serial.OrderedArray
		alpnSlice := strings.Split(vless.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsSettingsObject.Set("alpn", alpn)
			}
		}
		tlsSettingsObject.Set("allowInsecure", false)
		if len(vless.FingerPrint) > 0 {
			tlsSettingsObject.Set("fingerprint", vless.FingerPrint)
		}
		if len(vless.Sni) > 0 {
			tlsSettingsObject.Set("serverName", vless.Sni)
		}
		streamSettingsObject.Set("tlsSettings", tlsSettingsObject)
	case "reality":
		var realitySettingsObject serial.OrderedMap
		realitySettingsObject.Set("allowInsecure", false)
		if len(vless.FingerPrint) > 0 {
			realitySettingsObject.Set("fingerprint", vless.FingerPrint)
		}
		if len(vless.Sni) > 0 {
			realitySettingsObject.Set("serverName", vless.Sni)
		}
		realitySettingsObject.Set("publicKey", vless.PublicKey)
		realitySettingsObject.Set("shortId", vless.ShortId)
		realitySettingsObject.Set("spiderX", vless.SpiderX)
		streamSettingsObject.Set("realitySettings", realitySettingsObject)
	}
	var sockoptObject serial.OrderedMap
	sockoptObject.Set("domainStrategy", "UseIP")
	streamSettingsObject.Set("sockopt", sockoptObject)
	return streamSettingsObject
}
