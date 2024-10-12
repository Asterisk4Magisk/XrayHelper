package addon

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/serial"
	"strings"
)

// GetMuxObjectXray get xray MuxObject
func GetMuxObjectXray(enabled bool) serial.OrderedMap {
	var mux serial.OrderedMap
	mux.Set("enabled", enabled)
	return mux
}

// GetStreamSettingsObjectXray get addon StreamSettingsObject Xray
func GetStreamSettingsObjectXray(addon *Addon, network string, security string) serial.OrderedMap {
	var streamSettingsObject serial.OrderedMap
	streamSettingsObject.Set("network", network)
	switch network {
	case "tcp":
		var tcpSettingsObject serial.OrderedMap
		var headerObject serial.OrderedMap
		switch addon.Type {
		case "http":
			headerObject.Set("type", addon.Type)
			if len(addon.Host) > 0 {
				var requestObject serial.OrderedMap
				var headers serial.OrderedMap
				var host serial.OrderedArray
				host = append(host, addon.Host)
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
		if len(addon.Type) > 0 {
			var headerObject serial.OrderedMap
			headerObject.Set("type", addon.Type)
			kcpSettingsObject.Set("header", headerObject)
		}
		kcpSettingsObject.Set("congestion", false)
		kcpSettingsObject.Set("downlinkCapacity", 100)
		kcpSettingsObject.Set("mtu", 1350)
		kcpSettingsObject.Set("readBufferSize", 1)
		if len(addon.Path) > 0 {
			kcpSettingsObject.Set("seed", addon.Path)
		}
		kcpSettingsObject.Set("tti", 50)
		kcpSettingsObject.Set("uplinkCapacity", 12)
		kcpSettingsObject.Set("writeBufferSize", 1)
		streamSettingsObject.Set("kcpSettings", kcpSettingsObject)
	case "ws":
		var wsSettingsObject serial.OrderedMap
		if len(addon.Host) > 0 {
			var headersObject serial.OrderedMap
			headersObject.Set("Host", addon.Host)
			wsSettingsObject.Set("headers", headersObject)
		}
		if len(addon.Path) > 0 {
			wsSettingsObject.Set("path", addon.Path)
		}
		streamSettingsObject.Set("wsSettings", wsSettingsObject)
	case "http", "h2", "h3":
		var httpSettingsObject serial.OrderedMap
		if len(addon.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, addon.Host)
			httpSettingsObject.Set("host", host)
		}
		if len(addon.Path) > 0 {
			httpSettingsObject.Set("path", addon.Path)
		}
		streamSettingsObject.Set("httpSettings", httpSettingsObject)
	case "httpupgrade":
		var httpupgradeSettingsObject serial.OrderedMap
		if len(addon.Host) > 0 {
			httpupgradeSettingsObject.Set("host", addon.Host)
		}
		if len(addon.Path) > 0 {
			httpupgradeSettingsObject.Set("path", addon.Path)
		}
		streamSettingsObject.Set("httpupgradeSettings", httpupgradeSettingsObject)
	case "splithttp":
		var splithttpSettingsObject serial.OrderedMap
		if len(addon.Host) > 0 {
			splithttpSettingsObject.Set("host", addon.Host)
		}
		if len(addon.Path) > 0 {
			splithttpSettingsObject.Set("path", addon.Path)
		}
		streamSettingsObject.Set("splithttpSettings", splithttpSettingsObject)
	case "quic":
		var quicSettingsObject serial.OrderedMap
		if len(addon.Type) > 0 {
			var headerObject serial.OrderedMap
			headerObject.Set("type", addon.Type)
			quicSettingsObject.Set("header", headerObject)
		}
		if len(addon.Path) > 0 {
			quicSettingsObject.Set("key", addon.Path)
		}
		if len(addon.Host) > 0 {
			quicSettingsObject.Set("security", addon.Host)
		}
		streamSettingsObject.Set("quicSettings", quicSettingsObject)
	case "grpc":
		var grpcSettingsObject serial.OrderedMap
		if addon.Type == "multi" {
			grpcSettingsObject.Set("multiMode", true)
		} else {
			grpcSettingsObject.Set("multiMode", false)
		}
		if len(addon.Host) > 0 {
			grpcSettingsObject.Set("authority", addon.Host)
		}
		if len(addon.Path) > 0 {
			grpcSettingsObject.Set("serviceName", addon.Path)
		}
		streamSettingsObject.Set("grpcSettings", grpcSettingsObject)
	}
	switch security {
	case "tls":
		streamSettingsObject.Set("security", security)
		var tlsSettingsObject serial.OrderedMap
		var alpn serial.OrderedArray
		alpnSlice := strings.Split(addon.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsSettingsObject.Set("alpn", alpn)
			}
		}
		if len(addon.FingerPrint) > 0 {
			tlsSettingsObject.Set("fingerprint", addon.FingerPrint)
		}
		if len(addon.Sni) > 0 {
			tlsSettingsObject.Set("serverName", addon.Sni)
		}
		if builds.Config.XrayHelper.AllowInsecure {
			tlsSettingsObject.Set("allowInsecure", true)
		} else {
			tlsSettingsObject.Set("allowInsecure", false)
		}
		streamSettingsObject.Set("tlsSettings", tlsSettingsObject)
	case "reality":
		streamSettingsObject.Set("security", security)
		var realitySettingsObject serial.OrderedMap
		if len(addon.FingerPrint) > 0 {
			realitySettingsObject.Set("fingerprint", addon.FingerPrint)
		}
		if len(addon.Sni) > 0 {
			realitySettingsObject.Set("serverName", addon.Sni)
		}
		realitySettingsObject.Set("publicKey", addon.PublicKey)
		realitySettingsObject.Set("shortId", addon.ShortId)
		realitySettingsObject.Set("spiderX", addon.SpiderX)
		if builds.Config.XrayHelper.AllowInsecure {
			realitySettingsObject.Set("allowInsecure", true)
		} else {
			realitySettingsObject.Set("allowInsecure", false)
		}
		streamSettingsObject.Set("realitySettings", realitySettingsObject)
	}
	var sockoptObject serial.OrderedMap
	sockoptObject.Set("domainStrategy", "UseIP")
	streamSettingsObject.Set("sockopt", sockoptObject)
	return streamSettingsObject
}
