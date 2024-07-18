package addon

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/serial"
	"strings"
)

// GetTransportObjectSingbox get transport Object sing-box
func GetTransportObjectSingbox(addon *Addon, network string) serial.OrderedMap {
	var transportObject serial.OrderedMap
	switch network {
	case "http", "h2":
		transportObject.Set("type", "http")
		if len(addon.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, addon.Host)
			transportObject.Set("host", host)
		}
		if len(addon.Path) > 0 {
			transportObject.Set("path", addon.Path)
		}
	case "ws":
		transportObject.Set("type", "ws")
		if len(addon.Path) > 0 {
			transportObject.Set("path", addon.Path)
		}
		if len(addon.Host) > 0 {
			var headersObject serial.OrderedMap
			headersObject.Set("Host", addon.Host)
			transportObject.Set("headers", headersObject)
		}
		transportObject.Set("early_data_header_name", "Sec-WebSocket-Protocol")
	case "quic":
		transportObject.Set("type", "quic")
	case "grpc":
		transportObject.Set("type", "grpc")
		if len(addon.Path) > 0 {
			transportObject.Set("service_name", addon.Path)
		}
	case "httpupgrade":
		transportObject.Set("type", "httpupgrade")
		if len(addon.Host) > 0 {
			transportObject.Set("host", addon.Host)
		}
		if len(addon.Path) > 0 {
			transportObject.Set("path", addon.Path)
		}
	case "splithttp":
		transportObject.Set("type", "splithttp")
		if len(addon.Host) > 0 {
			transportObject.Set("host", addon.Host)
		}
		if len(addon.Path) > 0 {
			transportObject.Set("path", addon.Path)
		}
	}
	return transportObject
}

// GetTlsObjectSingbox get tls Object sing-box
func GetTlsObjectSingbox(addon *Addon, security string) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	if len(security) > 0 {
		tlsObject.Set("enabled", true)
		if len(addon.Sni) > 0 {
			tlsObject.Set("server_name", addon.Sni)
		}
		var alpn serial.OrderedArray
		alpnSlice := strings.Split(addon.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject.Set("alpn", alpn)
			}
		}
		var utlsObject serial.OrderedMap
		if len(addon.FingerPrint) > 0 {
			utlsObject.Set("enabled", true)
			utlsObject.Set("fingerprint", addon.FingerPrint)
			tlsObject.Set("utls", utlsObject)
		}
		if security == "reality" {
			var realityObject serial.OrderedMap
			realityObject.Set("enabled", true)
			realityObject.Set("public_key", addon.PublicKey)
			realityObject.Set("short_id", addon.ShortId)
			tlsObject.Set("reality", realityObject)
		}
	} else {
		tlsObject.Set("enabled", false)
	}
	if builds.Config.XrayHelper.AllowInsecure {
		tlsObject.Set("insecure", true)
	}
	return tlsObject
}
