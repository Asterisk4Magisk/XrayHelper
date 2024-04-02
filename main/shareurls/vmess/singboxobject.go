package vmess

import (
	"strings"
)

// getVmessTlsObjectSingbox get sing-box Vmess tls Object
func getVmessTlsObjectSingbox(vmess *Vmess) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	if len(vmess.Tls) > 0 {
		tlsObject["enabled"] = true
		if len(vmess.Sni) > 0 {
			tlsObject["server_name"] = vmess.Sni
		}
		var alpn []interface{}
		alpnSlice := strings.Split(string(vmess.Alpn), ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject["alpn"] = alpn
			}
		}
	} else {
		tlsObject["enabled"] = false
	}
	return tlsObject
}

// getVmessTransportObjectSingbox get sing-box Vmess transport Object
func getVmessTransportObjectSingbox(vmess *Vmess) map[string]interface{} {
	transportObject := make(map[string]interface{})
	switch vmess.Network {
	case "tcp", "http", "h2":
		transportObject["type"] = "http"
		if len(vmess.Host) > 0 {
			var host []interface{}
			host = append(host, vmess.Host)
			transportObject["host"] = host
		}
		if len(vmess.Path) > 0 {
			transportObject["path"] = vmess.Path
		}
	case "ws":
		transportObject["type"] = "ws"
		if len(vmess.Path) > 0 {
			transportObject["path"] = vmess.Path
		}
		if len(vmess.Host) > 0 {
			headersObject := make(map[string]interface{})
			headersObject["Host"] = vmess.Host
			transportObject["headers"] = headersObject
		}
		transportObject["early_data_header_name"] = "Sec-WebSocket-Protocol"
	case "quic":
		transportObject["type"] = "quic"
	case "grpc":
		transportObject["type"] = "grpc"
		if len(vmess.Path) > 0 {
			transportObject["service_name"] = vmess.Path
		}
	case "httpupgrade":
		transportObject["type"] = "httpupgrade"
		if len(vmess.Host) > 0 {
			transportObject["host"] = vmess.Host
		}
		if len(vmess.Path) > 0 {
			transportObject["path"] = vmess.Path
		}
	}
	return transportObject
}
