package vless

import (
	"strings"
)

// getVLESSTlsObjectSingbox get sing-box VLESS tls Object
func getVLESSTlsObjectSingbox(vless *VLESS) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	if len(vless.Security) > 0 {
		tlsObject["enabled"] = true
		if len(vless.Sni) > 0 {
			tlsObject["server_name"] = vless.Sni
		}
		var alpn []interface{}
		alpnSlice := strings.Split(vless.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject["alpn"] = alpn
			}
		}
		if vless.Security == "reality" {
			realityObject := make(map[string]interface{})
			realityObject["enabled"] = true
			realityObject["public_key"] = vless.PublicKey
			realityObject["short_id"] = vless.ShortId
			tlsObject["reality"] = realityObject
		}
	} else {
		tlsObject["enabled"] = false
	}
	return tlsObject
}

// getVLESSTransportObjectSingbox get sing-box VLESS transport Object
func getVLESSTransportObjectSingbox(vless *VLESS) map[string]interface{} {
	transportObject := make(map[string]interface{})
	switch vless.Network {
	case "tcp", "http", "h2":
		transportObject["type"] = "http"
		if len(vless.Host) > 0 {
			var host []interface{}
			host = append(host, vless.Host)
			transportObject["host"] = host
		}
		if len(vless.Path) > 0 {
			transportObject["path"] = vless.Path
		}
	case "ws":
		transportObject["type"] = "ws"
		if len(vless.Path) > 0 {
			transportObject["path"] = vless.Path
		}
		if len(vless.Host) > 0 {
			headersObject := make(map[string]interface{})
			headersObject["Host"] = vless.Host
			transportObject["headers"] = headersObject
		}
		transportObject["early_data_header_name"] = "Sec-WebSocket-Protocol"
	case "quic":
		transportObject["type"] = "quic"
	case "grpc":
		transportObject["type"] = "grpc"
		if len(vless.Path) > 0 {
			transportObject["service_name"] = vless.Path
		}
	case "httpupgrade":
		transportObject["type"] = "httpupgrade"
		if len(vless.Host) > 0 {
			transportObject["host"] = vless.Host
		}
		if len(vless.Path) > 0 {
			transportObject["path"] = vless.Path
		}
	}
	return transportObject
}
