package trojan

import (
	"strings"
)

// getTrojanTlsObjectSingbox get sing-box Trojan tls Object
func getTrojanTlsObjectSingbox(trojan *Trojan) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	if len(trojan.Security) > 0 {
		tlsObject["enabled"] = true
		tlsObject["server_name"] = trojan.Sni
		var alpn []interface{}
		alpnSlice := strings.Split(trojan.Alpn, ",")
		for _, v := range alpnSlice {
			if len(v) > 0 {
				alpn = append(alpn, v)
				tlsObject["alpn"] = alpn
			}
		}
		//utlsObject := make(map[string]interface{})
		//if len(trojan.FingerPrint) > 0 {
		//	utlsObject["enabled"] = true
		//	utlsObject["fingerprint"] = trojan.FingerPrint
		//	tlsObject["utls"] = utlsObject
		//}
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
