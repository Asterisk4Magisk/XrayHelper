package vmess

import (
	"XrayHelper/main/serial"
	"strings"
)

// getVmessTlsObjectSingbox get sing-box Vmess tls Object
func getVmessTlsObjectSingbox(vmess *Vmess) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	if len(vmess.Tls) > 0 {
		tlsObject.Set("enabled", true)
		if len(vmess.Sni) > 0 {
			tlsObject.Set("server_name", vmess.Sni)
		}
		var alpn serial.OrderedArray
		alpnSlice := strings.Split(string(vmess.Alpn), ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject.Set("alpn", alpn)
			}
		}
		var utlsObject serial.OrderedMap
		if len(vmess.FingerPrint) > 0 {
			utlsObject.Set("enabled", true)
			utlsObject.Set("fingerprint", vmess.FingerPrint)
			tlsObject.Set("utls", utlsObject)
		}
	} else {
		tlsObject.Set("enabled", false)
	}
	return tlsObject
}

// getVmessTransportObjectSingbox get sing-box Vmess transport Object
func getVmessTransportObjectSingbox(vmess *Vmess) serial.OrderedMap {
	var transportObject serial.OrderedMap
	switch vmess.Network {
	case "tcp", "http", "h2":
		transportObject.Set("type", "http")
		if len(vmess.Host) > 0 {
			var host serial.OrderedArray
			host = append(host, vmess.Host)
			transportObject.Set("host", host)
		}
		if len(vmess.Path) > 0 {
			transportObject.Set("path", vmess.Path)
		}
	case "ws":
		transportObject.Set("type", "ws")
		if len(vmess.Path) > 0 {
			transportObject.Set("path", vmess.Path)
		}
		if len(vmess.Host) > 0 {
			var headersObject serial.OrderedMap
			headersObject.Set("Host", vmess.Host)
			transportObject.Set("headers", headersObject)
		}
		transportObject.Set("early_data_header_name", "Sec-WebSocket-Protocol")
	case "quic":
		transportObject.Set("type", "quic")
	case "grpc":
		transportObject.Set("type", "grpc")
		if len(vmess.Path) > 0 {
			transportObject.Set("service_name", vmess.Path)
		}
	case "httpupgrade":
		transportObject.Set("type", "httpupgrade")
		if len(vmess.Host) > 0 {
			transportObject.Set("host", vmess.Host)
		}
		if len(vmess.Path) > 0 {
			transportObject.Set("path", vmess.Path)
		}
	}
	return transportObject
}
