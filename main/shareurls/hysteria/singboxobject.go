package hysteria

import (
	"strconv"
	"strings"
)

// getHysteriaTlsObjectSingbox get sing-box Hysteria tls Object
func getHysteriaTlsObjectSingbox(hysteria *Hysteria) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	tlsObject["enabled"] = true
	tlsObject["server_name"] = hysteria.Peer
	tlsObject["insecure"], _ = strconv.ParseBool(hysteria.Insecure)
	var alpn []interface{}
	alpnSlice := strings.Split(hysteria.Alpn, ",")
	for _, v := range alpnSlice {
		if len(v) > 0 {
			alpn = append(alpn, v)
			tlsObject["alpn"] = alpn
		}
	}
	return tlsObject
}
