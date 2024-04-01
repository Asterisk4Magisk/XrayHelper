package hysteria2

import (
	"strconv"
	"strings"
)

// getHysteriaTlsObjectSingbox get sing-box Hysteria tls Object
func getHysteriaTlsObjectSingbox(hysteria2 *Hysteria2) map[string]interface{} {
	tlsObject := make(map[string]interface{})
	tlsObject["enabled"] = true
	if len(hysteria2.Sni) == 0 {
		tlsObject["disable_sni"] = true
	} else {
		tlsObject["server_name"] = hysteria2.Sni
	}
	tlsObject["insecure"], _ = strconv.ParseBool(hysteria2.Insecure)
	return tlsObject
}

// getHysteriaTlsObjectSingbox get sing-box Hysteria2 obfs Object
func getHysteria2ObfsObjectSingbox(hysteria2 *Hysteria2) map[string]interface{} {
	obfsObject := make(map[string]interface{})
	obfsObject["type"] = hysteria2.Obfs
	obfsObject["password"] = hysteria2.ObfsPassword
	return obfsObject
}

// getHysteriaTlsObjectSingbox get sing-box Hysteria2 users Object
func getHysteria2UsersObjectSingbox(hysteria2 *Hysteria2) []interface{} {
	nameAndPassword := strings.Split(hysteria2.Auth, ":")
	var users []interface{}
	userObject := make(map[string]interface{})
	userObject["name"] = nameAndPassword[0]
	if len(nameAndPassword) == 2 {
		userObject["password"] = nameAndPassword[1]
	}
	users = append(users, userObject)
	return users
}
