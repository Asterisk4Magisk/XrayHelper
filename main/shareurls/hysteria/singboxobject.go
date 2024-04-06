package hysteria

import (
	"XrayHelper/main/serial"
	"strconv"
	"strings"
)

// getHysteriaTlsObjectSingbox get sing-box Hysteria tls Object
func getHysteriaTlsObjectSingbox(hysteria *Hysteria) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	tlsObject.Set("enabled", true)
	tlsObject.Set("server_name", hysteria.Peer)
	insecure, _ := strconv.ParseBool(hysteria.Insecure)
	tlsObject.Set("insecure", insecure)
	var alpn serial.OrderedArray
	alpnSlice := strings.Split(hysteria.Alpn, ",")
	for _, v := range alpnSlice {
		if len(v) > 0 {
			alpn = append(alpn, v)
			tlsObject.Set("alpn", alpn)
		}
	}
	return tlsObject
}
