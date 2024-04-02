package trojan

import (
	"strings"
)

// getTrojanTlsObjectSingbox get sing-box Trojan tls Object
func getTrojanTlsObjectSingbox(trojan *Trojan) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	if len(trojan.Security) > 0 {
		tlsObject["enabled"] = true
		if len(trojan.Sni) > 0 {
			tlsObject["server_name"] = trojan.Sni
		}
		var alpn []interface{}
		alpnSlice := strings.Split(trojan.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject["alpn"] = alpn
			}
		}
		if trojan.Security == "reality" {
			realityObject := make(map[string]interface{})
			realityObject["enabled"] = true
			realityObject["public_key"] = trojan.PublicKey
			realityObject["short_id"] = trojan.ShortId
			tlsObject["reality"] = realityObject
		}
	} else {
		tlsObject["enabled"] = false
	}
	return tlsObject
}

// getTrojanTransportObjectSingbox get sing-box Trojan transport Object
func getTrojanTransportObjectSingbox(trojan *Trojan) map[string]interface{} {
	transportObject := make(map[string]interface{})
	switch trojan.Network {
	case "tcp", "http", "h2":
		transportObject["type"] = "http"
		if len(trojan.Host) > 0 {
			var host []interface{}
			host = append(host, trojan.Host)
			transportObject["host"] = host
		}
		if len(trojan.Path) > 0 {
			transportObject["path"] = trojan.Path
		}
	case "ws":
		transportObject["type"] = "ws"
		if len(trojan.Path) > 0 {
			transportObject["path"] = trojan.Path
		}
		if len(trojan.Host) > 0 {
			headersObject := make(map[string]interface{})
			headersObject["Host"] = trojan.Host
			transportObject["headers"] = headersObject
		}
		transportObject["early_data_header_name"] = "Sec-WebSocket-Protocol"
	case "quic":
		transportObject["type"] = "quic"
	case "grpc":
		transportObject["type"] = "grpc"
		if len(trojan.Path) > 0 {
			transportObject["service_name"] = trojan.Path
		}
	case "httpupgrade":
		transportObject["type"] = "httpupgrade"
		if len(trojan.Host) > 0 {
			transportObject["host"] = trojan.Host
		}
		if len(trojan.Path) > 0 {
			transportObject["path"] = trojan.Path
		}
	}
	return transportObject
}
