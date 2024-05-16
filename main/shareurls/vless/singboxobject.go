package vless

import (
	"XrayHelper/main/serial"
	"strings"
)

// getVLESSTlsObjectSingbox get sing-box VLESS tls Object
func getVLESSTlsObjectSingbox(vless *VLESS) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	if len(vless.Security) > 0 {
		tlsObject.Set("enabled", true)
		if len(vless.Sni) > 0 {
			tlsObject.Set("server_name", vless.Sni)
		}
		var alpn serial.OrderedArray
		alpnSlice := strings.Split(vless.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject.Set("alpn", alpn)
			}
		}
		var utlsObject serial.OrderedMap
		if len(vless.FingerPrint) > 0 {
			utlsObject.Set("enabled", true)
			utlsObject.Set("fingerprint", vless.FingerPrint)
			tlsObject.Set("utls", utlsObject)
		}
		if vless.Security == "reality" {
			var realityObject serial.OrderedMap
			realityObject.Set("enabled", true)
			realityObject.Set("public_key", vless.PublicKey)
			realityObject.Set("short_id", vless.ShortId)
			tlsObject.Set("reality", realityObject)
		}
	} else {
		tlsObject.Set("enabled", false)
	}
	return tlsObject
}

// getVLESSTransportObjectSingbox get sing-box VLESS transport Object
func getVLESSTransportObjectSingbox(vless *VLESS) serial.OrderedMap {
	var transportObject serial.OrderedMap
	switch vless.Network {
	case "http", "h2":
		transportObject.Set("type", "http")
		if len(vless.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, vless.Host)
			transportObject.Set("host", host)
		}
		if len(vless.Path) > 0 {
			transportObject.Set("path", vless.Path)
		}
	case "ws":
		transportObject.Set("type", "ws")
		if len(vless.Path) > 0 {
			transportObject.Set("path", vless.Path)
		}
		if len(vless.Host) > 0 {
			var headersObject serial.OrderedMap
			headersObject.Set("Host", vless.Host)
			transportObject.Set("headers", headersObject)
		}
		transportObject.Set("early_data_header_name", "Sec-WebSocket-Protocol")
	case "quic":
		transportObject.Set("type", "quic")
	case "grpc":
		transportObject.Set("type", "grpc")
		if len(vless.Path) > 0 {
			transportObject.Set("service_name", vless.Path)
		}
	case "httpupgrade":
		transportObject.Set("type", "httpupgrade")
		if len(vless.Host) > 0 {
			transportObject.Set("host", vless.Host)
		}
		if len(vless.Path) > 0 {
			transportObject.Set("path", vless.Path)
		}
	}
	return transportObject
}
