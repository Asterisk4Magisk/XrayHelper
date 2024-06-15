package hysteria2

import (
	"XrayHelper/main/serial"
	"strconv"
)

// getHysteria2ObfsObjectHysteria2 get Hysteria2 obfs Object
func getHysteria2ObfsObjectHysteria2(hysteria2 *Hysteria2) serial.OrderedMap {
	var obfsObject serial.OrderedMap
	if len(hysteria2.Obfs) > 0 {
		obfsObject.Set("type", hysteria2.Obfs)
		switch hysteria2.Obfs {
		case "salamander":
			var salamanderObject serial.OrderedMap
			salamanderObject.Set("password", hysteria2.ObfsPassword)
			obfsObject.Set("salamander", salamanderObject)
		}
	}
	return obfsObject
}

// getHysteria2TlsObjectHysteria2 get Hysteria2 tls Object
func getHysteria2TlsObjectHysteria2(hysteria2 *Hysteria2) serial.OrderedMap {
	var tlsObject serial.OrderedMap
	tlsObject.Set("sni", hysteria2.Sni)
	insecure, _ := strconv.ParseBool(hysteria2.Insecure)
	tlsObject.Set("insecure", insecure)
	tlsObject.Set("pinSHA256", hysteria2.PinSHA256)
	return tlsObject
}
