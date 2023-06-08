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
