package hysteria2

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/serial"
	"strconv"
)

// getHysteria2TlsObjectSingbox get sing-box Hysteria2 tls Object
func getHysteria2TlsObjectSingbox(hysteria2 *Hysteria2) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	tlsObject.Set("enabled", true)
	if len(hysteria2.Sni) == 0 {
		tlsObject.Set("disable_sni", true)
	} else {
		tlsObject.Set("server_name", hysteria2.Sni)
	}
	insecure, _ := strconv.ParseBool(hysteria2.Insecure)
	if builds.Config.XrayHelper.AllowInsecure || insecure {
		tlsObject.Set("insecure", true)
	} else {
		tlsObject.Set("insecure", false)
	}
	return tlsObject
}

// getHysteriaTlsObjectSingbox get sing-box Hysteria2 obfs Object
func getHysteria2ObfsObjectSingbox(hysteria2 *Hysteria2) serial.OrderedMap {
	var obfsObject serial.OrderedMap
	obfsObject.Set("type", hysteria2.Obfs)
	obfsObject.Set("password", hysteria2.ObfsPassword)
	return obfsObject
}
