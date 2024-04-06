package trojan

import (
	"XrayHelper/main/serial"
	"strings"
)

// getTrojanTlsObjectSingbox get sing-box Trojan tls Object
func getTrojanTlsObjectSingbox(trojan *Trojan) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	if len(trojan.Security) > 0 {
		tlsObject.Set("enabled", true)
		if len(trojan.Sni) > 0 {
			tlsObject.Set("server_name", trojan.Sni)
		}
		var alpn serial.OrderedArray
		alpnSlice := strings.Split(trojan.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject.Set("alpn", alpn)
			}
		}
		if trojan.Security == "reality" {
			var realityObject serial.OrderedMap
			realityObject.Set("enabled", true)
			realityObject.Set("public_key", trojan.PublicKey)
			realityObject.Set("short_id", trojan.ShortId)
			tlsObject.Set("reality", realityObject)
		}
	} else {
		tlsObject.Set("enabled", false)
	}
	return tlsObject
}

// getTrojanTransportObjectSingbox get sing-box Trojan transport Object
func getTrojanTransportObjectSingbox(trojan *Trojan) serial.OrderedMap {
	var transportObject serial.OrderedMap
	switch trojan.Network {
	case "tcp", "http", "h2":
		transportObject.Set("type", "http")
		if len(trojan.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, trojan.Host)
			transportObject.Set("host", host)
		}
		if len(trojan.Path) > 0 {
			transportObject.Set("path", trojan.Path)
		}
	case "ws":
		transportObject.Set("type", "ws")
		if len(trojan.Path) > 0 {
			transportObject.Set("path", trojan.Path)
		}
		if len(trojan.Host) > 0 {
			var headersObject serial.OrderedMap
			headersObject.Set("Host", trojan.Host)
			transportObject.Set("headers", headersObject)
		}
		transportObject.Set("early_data_header_name", "Sec-WebSocket-Protocol")
	case "quic":
		transportObject.Set("type", "quic")
	case "grpc":
		transportObject.Set("type", "grpc")
		if len(trojan.Path) > 0 {
			transportObject.Set("service_name", trojan.Path)
		}
	case "httpupgrade":
		transportObject.Set("type", "httpupgrade")
		if len(trojan.Host) > 0 {
			transportObject.Set("host", trojan.Host)
		}
		if len(trojan.Path) > 0 {
			transportObject.Set("path", trojan.Path)
		}
	}
	return transportObject
}
