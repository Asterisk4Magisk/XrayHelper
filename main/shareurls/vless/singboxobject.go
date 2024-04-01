package vless

import (
	"strings"
)

// getVLESSTlsObjectSingbox get sing-box VLESS tls Object
func getVLESSTlsObjectSingbox(vless *VLESS) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	if len(vless.Security) > 0 {
		tlsObject["enabled"] = true
		tlsObject["server_name"] = vless.Sni
		var alpn []interface{}
		alpnSlice := strings.Split(vless.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject["alpn"] = alpn
			}
		}
		//utlsObject := make(map[string]interface{})
		//if len(vless.FingerPrint) > 0 {
		//	utlsObject["enabled"] = true
		//	utlsObject["fingerprint"] = vless.FingerPrint
		//	tlsObject["utls"] = utlsObject
		//}
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
		var host []interface{}
		host = append(host, vless.Host)
		transportObject["host"] = host
		transportObject["path"] = vless.Path
	case "ws":
		transportObject["type"] = "ws"
		transportObject["path"] = vless.Path
		headersObject := make(map[string]interface{})
		headersObject["Host"] = vless.Host
		transportObject["headers"] = headersObject
		transportObject["early_data_header_name"] = "Sec-WebSocket-Protocol"
	case "quic":
		transportObject["type"] = "quic"
	case "grpc":
		transportObject["type"] = "grpc"
		transportObject["service_name"] = vless.Path
	}
	return transportObject
}
