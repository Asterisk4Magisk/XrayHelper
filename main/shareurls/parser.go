package shareurls

import (
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/shareurls/hysteria"
	"XrayHelper/main/shareurls/hysteria2"
	"XrayHelper/main/shareurls/shadowsocks"
	"XrayHelper/main/shareurls/socks"
	"XrayHelper/main/shareurls/trojan"
	"XrayHelper/main/shareurls/vless"
	"XrayHelper/main/shareurls/vmess"
	"encoding/json"
	"net/url"
	"strings"
)

const tagParser = "parser"

// parseShadowsocks parse shadowsocks url
func parseShadowsocks(ssUrl string) (ShareUrl, error) {
	ss := new(shadowsocks.Shadowsocks)
	ssParse, err := url.Parse(ssUrl)
	if err != nil {
		return nil, e.New("shadowsocks url parse err, ", err).WithPrefix(tagParser)
	}
	ss.Remarks = ssParse.Fragment
	ss.Server = ssParse.Hostname()
	ss.Port = ssParse.Port()
	if ss.Port == "" {
		full, err := common.DecodeBase64(ssParse.Hostname())
		if err != nil {
			return nil, err
		}
		infoAndServer := strings.Split(full, "@")
		methodAndPassword := strings.Split(infoAndServer[0], ":")
		ss.Method = methodAndPassword[0]
		ss.Password = methodAndPassword[1]
		addressAndPort := strings.Split(infoAndServer[1], ":")
		ss.Server = addressAndPort[0]
		ss.Port = addressAndPort[1]
	} else {
		info, err := common.DecodeBase64(ssParse.User.Username())
		if err != nil {
			return nil, err
		}
		methodAndPassword := strings.Split(info, ":")
		ss.Method = methodAndPassword[0]
		ss.Password = methodAndPassword[1]
		ssQuery, err := url.ParseQuery(ssParse.RawQuery)
		if err != nil {
			return nil, e.New("shadowsocks url parse query err, ", err).WithPrefix(tagParser)
		}
		//parse shadowsocks SIP003 plugin
		if plugins, ok := ssQuery["plugin"]; ok && len(plugins) == 1 {
			plugin := strings.Split(plugins[0], ";")
			ss.Plugin = plugin[0]
			ss.PluginOpt = strings.TrimPrefix(plugins[0], plugin[0]+";")
		}
	}
	return ss, nil
}

// parseSocks parse socks url
func parseSocks(socksUrl string) (ShareUrl, error) {
	so := new(socks.Socks)
	soParse, err := url.Parse(socksUrl)
	if err != nil {
		return nil, e.New("socks url parse err, ", err).WithPrefix(tagParser)
	}
	so.Remarks = soParse.Fragment
	so.Server = soParse.Hostname()
	so.Port = soParse.Port()
	info, err := common.DecodeBase64(soParse.User.Username())
	if err != nil {
		return nil, err
	}
	userAndPassword := strings.Split(info, ":")
	so.User = userAndPassword[0]
	so.Password = userAndPassword[1]
	return so, nil
}

// parseTrojan parse trojan url
func parseTrojan(trojanUrl string) (ShareUrl, error) {
	tj := new(trojan.Trojan)
	tjParse, err := url.Parse(trojanUrl)
	if err != nil {
		return nil, e.New("trojan url parse err, ", err).WithPrefix(tagParser)
	}
	tj.Remarks = tjParse.Fragment
	tj.Password = tjParse.User.Username()
	tj.Server = tjParse.Hostname()
	tj.Port = tjParse.Port()
	tjQuery, err := url.ParseQuery(tjParse.RawQuery)
	if err != nil {
		return nil, e.New("trojan url parse query err, ", err).WithPrefix(tagParser)
	}
	//parse trojan network
	if types, ok := tjQuery["type"]; !ok {
		tj.Network = "tcp"
	} else if len(types) > 1 {
		return nil, e.New("multiple trojan transport type").WithPrefix(tagParser)
	} else if tj.Network = types[0]; tj.Network == "" {
		return nil, e.New("empty trojan transport type").WithPrefix(tagParser)
	}
	//parse trojan security
	if security, ok := tjQuery["security"]; !ok {
		tj.Security = "tls"
	} else if len(security) > 1 {
		return nil, e.New("multiple trojan security type").WithPrefix(tagParser)
	} else if tj.Security = security[0]; tj.Security == "" {
		return nil, e.New("empty trojan security type").WithPrefix(tagParser)
	}
	switch tj.Network {
	case "tcp":
		//parse trojan headerType
		if headerTypes, ok := tjQuery["headerType"]; ok && len(headerTypes) == 1 {
			if headerTypes[0] == "http" {
				if hosts, ok := tjQuery["host"]; ok && len(hosts) == 1 {
					tj.Host = hosts[0]
				}
			}
		}
	case "kcp":
		//parse trojan headerType
		if headerTypes, ok := tjQuery["headerType"]; ok && len(headerTypes) == 1 {
			tj.Type = headerTypes[0]
		}
		//parse trojan kcp seed
		if seeds, ok := tjQuery["seed"]; ok && len(seeds) == 1 {
			tj.Path = seeds[0]
		}
	case "ws", "http", "h2", "httpupgrade":
		//parse trojan host
		if hosts, ok := tjQuery["host"]; ok && len(hosts) == 1 {
			tj.Host = hosts[0]
		}
		//parse trojan path
		if paths, ok := tjQuery["path"]; ok && len(paths) == 1 {
			tj.Path = paths[0]
		}
	case "quic":
		//parse trojan headerType
		if headerTypes, ok := tjQuery["headerType"]; ok && len(headerTypes) == 1 {
			tj.Type = headerTypes[0]
		}
		//parse trojan quicSecurity
		if quicSecurity, ok := tjQuery["quicSecurity"]; ok && len(quicSecurity) == 1 {
			tj.Host = quicSecurity[0]
		}
		//parse trojan quicKey
		if quicKey, ok := tjQuery["key"]; ok && len(quicKey) == 1 {
			tj.Path = quicKey[0]
		}
	case "grpc":
		//parse trojan grpc authority
		if authority, ok := tjQuery["authority"]; ok && len(authority) == 1 {
			tj.Host = authority[0]
		}
		//parse trojan grpc mode
		if modes, ok := tjQuery["mode"]; ok && len(modes) == 1 {
			tj.Type = modes[0]
		}
		//parse trojan grpc serviceName
		if serviceNames, ok := tjQuery["serviceName"]; ok && len(serviceNames) == 1 {
			tj.Path = serviceNames[0]
		}
	default:
		return nil, e.New("unknown trojan transport type " + tj.Network).WithPrefix(tagParser)
	}
	switch tj.Security {
	case "tls":
		//parse trojan tls sni
		if sni, ok := tjQuery["sni"]; ok && len(sni) == 1 {
			tj.Sni = sni[0]
		}
		//parse trojan tls fingerprint
		if fps, ok := tjQuery["fp"]; ok && len(fps) == 1 {
			tj.FingerPrint = fps[0]
		}
		//parse trojan tls Alpn
		if alpns, ok := tjQuery["alpn"]; ok && len(alpns) == 1 {
			tj.Alpn = alpns[0]
		}
	case "reality":
		//parse trojan reality sni
		if sni, ok := tjQuery["sni"]; ok && len(sni) == 1 {
			tj.Sni = sni[0]
		}
		//parse trojan reality fingerprint
		if fps, ok := tjQuery["fp"]; ok && len(fps) == 1 {
			tj.FingerPrint = fps[0]
		}
		//parse trojan reality PublicKey
		if publicKeys, ok := tjQuery["pbk"]; ok && len(publicKeys) == 1 {
			tj.PublicKey = publicKeys[0]
		}
		//parse trojan reality ShortId
		if shortIds, ok := tjQuery["sid"]; ok && len(shortIds) == 1 {
			tj.ShortId = shortIds[0]
		}
		//parse trojan reality SpiderX
		if spiderX, ok := tjQuery["spx"]; ok && len(spiderX) == 1 {
			tj.SpiderX = spiderX[0]
		}
	default:
		return nil, e.New("unknown trojan security type " + tj.Security).WithPrefix(tagParser)
	}
	return tj, nil
}

// parseVLESS parse VLESS url
func parseVLESS(vlessUrl string) (ShareUrl, error) {
	vl := new(vless.VLESS)
	vlParse, err := url.Parse(vlessUrl)
	if err != nil {
		return nil, e.New("VLESS url parse err, ", err).WithPrefix(tagParser)
	}
	vl.Remarks = vlParse.Fragment
	vl.Id = vlParse.User.Username()
	vl.Server = vlParse.Hostname()
	vl.Port = vlParse.Port()
	vlQuery, err := url.ParseQuery(vlParse.RawQuery)
	if err != nil {
		return nil, e.New("VLESS url parse query err, ", err).WithPrefix(tagParser)
	}
	//parse VLESS encryption
	if encryption, ok := vlQuery["encryption"]; !ok {
		vl.Encryption = "none"
	} else if len(encryption) > 1 {
		return nil, e.New("multiple VLESS encryption").WithPrefix(tagParser)
	} else if vl.Encryption = encryption[0]; vl.Encryption == "" {
		return nil, e.New("empty VLESS encryption").WithPrefix(tagParser)
	}
	//parse VLESS flow
	if flows, ok := vlQuery["flow"]; ok {
		if len(flows) > 1 {
			return nil, e.New("multiple VLESS flow").WithPrefix(tagParser)
		} else {
			vl.Flow = flows[0]
		}
	}
	//parse VLESS network
	if types, ok := vlQuery["type"]; !ok {
		vl.Network = "tcp"
	} else if len(types) > 1 {
		return nil, e.New("multiple VLESS transport type").WithPrefix(tagParser)
	} else if vl.Network = types[0]; vl.Network == "" {
		return nil, e.New("empty VLESS transport type").WithPrefix(tagParser)
	}
	//parse VLESS security
	if security, ok := vlQuery["security"]; !ok {
		vl.Security = "tls"
	} else if len(security) > 1 {
		return nil, e.New("multiple VLESS security type").WithPrefix(tagParser)
	} else if vl.Security = security[0]; vl.Security == "" {
		return nil, e.New("empty VLESS security type").WithPrefix(tagParser)
	}

	switch vl.Network {
	case "tcp":
		//parse VLESS headerType
		if headerTypes, ok := vlQuery["headerType"]; ok && len(headerTypes) == 1 {
			if headerTypes[0] == "http" {
				if hosts, ok := vlQuery["host"]; ok && len(hosts) == 1 {
					vl.Host = hosts[0]
				}
			}
		}
	case "kcp":
		//parse VLESS headerType
		if headerTypes, ok := vlQuery["headerType"]; ok && len(headerTypes) == 1 {
			vl.Type = headerTypes[0]
		}
		//parse VLESS kcp seed
		if seeds, ok := vlQuery["seed"]; ok && len(seeds) == 1 {
			vl.Path = seeds[0]
		}
	case "ws", "http", "h2", "httpupgrade":
		//parse VLESS host
		if hosts, ok := vlQuery["host"]; ok && len(hosts) == 1 {
			vl.Host = hosts[0]
		}
		//parse VLESS path
		if paths, ok := vlQuery["path"]; ok && len(paths) == 1 {
			vl.Path = paths[0]
		}
	case "quic":
		//parse VLESS headerType
		if headerTypes, ok := vlQuery["headerType"]; ok && len(headerTypes) == 1 {
			vl.Type = headerTypes[0]
		}
		//parse VLESS quicSecurity
		if quicSecurity, ok := vlQuery["quicSecurity"]; ok && len(quicSecurity) == 1 {
			vl.Host = quicSecurity[0]
		}
		//parse VLESS quicKey
		if quicKey, ok := vlQuery["key"]; ok && len(quicKey) == 1 {
			vl.Path = quicKey[0]
		}
	case "grpc":
		//parse VLESS grpc authority
		if authority, ok := vlQuery["authority"]; ok && len(authority) == 1 {
			vl.Host = authority[0]
		}
		//parse VLESS grpc mode
		if modes, ok := vlQuery["mode"]; ok && len(modes) == 1 {
			vl.Type = modes[0]
		}
		//parse VLESS grpc serviceName
		if serviceNames, ok := vlQuery["serviceName"]; ok && len(serviceNames) == 1 {
			vl.Path = serviceNames[0]
		}
	default:
		return nil, e.New("unknown VLESS transport type " + vl.Network).WithPrefix(tagParser)
	}
	switch vl.Security {
	case "tls":
		//parse VLESS tls sni
		if sni, ok := vlQuery["sni"]; ok && len(sni) == 1 {
			vl.Sni = sni[0]
		}
		//parse VLESS tls fingerprint
		if fps, ok := vlQuery["fp"]; ok && len(fps) == 1 {
			vl.FingerPrint = fps[0]
		}
		//parse VLESS tls Alpn
		if alpns, ok := vlQuery["alpn"]; ok && len(alpns) == 1 {
			vl.Alpn = alpns[0]
		}
	case "reality":
		//parse VLESS reality sni
		if sni, ok := vlQuery["sni"]; ok && len(sni) == 1 {
			vl.Sni = sni[0]
		}
		//parse VLESS reality fingerprint
		if fps, ok := vlQuery["fp"]; ok && len(fps) == 1 {
			vl.FingerPrint = fps[0]
		}
		//parse VLESS reality PublicKey
		if publicKeys, ok := vlQuery["pbk"]; ok && len(publicKeys) == 1 {
			vl.PublicKey = publicKeys[0]
		}
		//parse VLESS reality ShortId
		if shortIds, ok := vlQuery["sid"]; ok && len(shortIds) == 1 {
			vl.ShortId = shortIds[0]
		}
		//parse VLESS reality SpiderX
		if spiderX, ok := vlQuery["spx"]; ok && len(spiderX) == 1 {
			vl.SpiderX = spiderX[0]
		}
	default:
		return nil, e.New("unknown VLESS security type " + vl.Security).WithPrefix(tagParser)
	}
	return vl, nil
}

// parseVmess parse Vmess url
func parseVmess(vmessUrl string) (ShareUrl, error) {
	v2 := new(vmess.Vmess)
	originJson, err := common.DecodeBase64(vmessUrl)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(originJson), v2)
	if err != nil {
		return nil, e.New("unmarshal origin json failed, ", err).WithPrefix(tagParser)
	}
	return v2, nil
}

// parseHysteria parse hysteria url
func parseHysteria(hysteriaUrl string) (ShareUrl, error) {
	ht := new(hysteria.Hysteria)
	htParse, err := url.Parse(hysteriaUrl)
	if err != nil {
		return nil, e.New("hysteria url parse err, ", err).WithPrefix(tagParser)
	}
	ht.Remarks = htParse.Fragment
	ht.Host = htParse.Hostname()
	ht.Port = htParse.Port()
	htQuery, err := url.ParseQuery(htParse.RawQuery)
	if err != nil {
		return nil, e.New("hysteria url parse query err, ", err).WithPrefix(tagParser)
	}
	//parse hysteria protocol
	if protocols, ok := htQuery["protocol"]; ok {
		if len(protocols) > 1 {
			return nil, e.New("multiple hysteria protocol").WithPrefix(tagParser)
		} else {
			ht.Protocol = protocols[0]
		}
	}
	//parse hysteria auth
	if auth, ok := htQuery["auth"]; ok {
		if len(auth) > 1 {
			return nil, e.New("multiple hysteria auth").WithPrefix(tagParser)
		} else {
			ht.Auth = auth[0]
		}
	}
	//parse hysteria peer
	if peer, ok := htQuery["peer"]; ok {
		if len(peer) > 1 {
			return nil, e.New("multiple hysteria peer").WithPrefix(tagParser)
		} else {
			ht.Peer = peer[0]
		}
	}
	//parse hysteria insecure
	if insecure, ok := htQuery["insecure"]; ok {
		if len(insecure) > 1 {
			return nil, e.New("multiple hysteria insecure").WithPrefix(tagParser)
		} else {
			ht.Insecure = insecure[0]
		}
	}
	//parse hysteria upmbps
	if upmbps, ok := htQuery["upmbps"]; ok {
		if len(upmbps) > 1 {
			return nil, e.New("multiple hysteria upmbps").WithPrefix(tagParser)
		} else {
			ht.UpMBPS = upmbps[0]
		}
	}
	//parse hysteria downmbps
	if downmbps, ok := htQuery["downmbps"]; ok {
		if len(downmbps) > 1 {
			return nil, e.New("multiple hysteria downmbps").WithPrefix(tagParser)
		} else {
			ht.DownMBPS = downmbps[0]
		}
	}
	//parse hysteria alpn
	if alpn, ok := htQuery["alpn"]; ok {
		if len(alpn) > 1 {
			return nil, e.New("multiple hysteria alpn").WithPrefix(tagParser)
		} else {
			ht.Alpn = alpn[0]
		}
	}
	//parse hysteria obfs
	if obfs, ok := htQuery["obfs"]; ok {
		if len(obfs) > 1 {
			return nil, e.New("multiple hysteria obfs").WithPrefix(tagParser)
		} else {
			ht.Obfs = obfs[0]
		}
	}
	//parse hysteria obfsParam
	if obfsParam, ok := htQuery["obfsParam"]; ok {
		if len(obfsParam) > 1 {
			return nil, e.New("multiple hysteria obfsParam").WithPrefix(tagParser)
		} else {
			ht.ObfsParam = obfsParam[0]
		}
	}
	return ht, nil
}

// parseHysteria2 parse hysteria2 url
func parseHysteria2(hysteria2Url string) (ShareUrl, error) {
	ht := new(hysteria2.Hysteria2)
	htParse, err := url.Parse(hysteria2Url)
	if err != nil {
		return nil, e.New("hysteria2 url parse err, ", err).WithPrefix(tagParser)
	}
	ht.Remarks = htParse.Fragment
	ht.Host = htParse.Hostname()
	ht.Port = htParse.Port()
	ht.Auth = htParse.User.String()
	htQuery, err := url.ParseQuery(htParse.RawQuery)
	if err != nil {
		return nil, e.New("hysteria2 url parse query err, ", err).WithPrefix(tagParser)
	}
	//parse hysteria2 obfs
	if obfs, ok := htQuery["obfs"]; ok {
		if len(obfs) > 1 {
			return nil, e.New("multiple hysteria2 obfs").WithPrefix(tagParser)
		} else {
			ht.Obfs = obfs[0]
		}
	}
	//parse hysteria2 obfs-password
	if obfsPasswords, ok := htQuery["obfs-password"]; ok {
		if len(obfsPasswords) > 1 {
			return nil, e.New("multiple hysteria2 obfs-password").WithPrefix(tagParser)
		} else {
			ht.ObfsPassword = obfsPasswords[0]
		}
	}
	//parse hysteria2 sni
	if snis, ok := htQuery["sni"]; ok {
		if len(snis) > 1 {
			return nil, e.New("multiple hysteria2 sni").WithPrefix(tagParser)
		} else {
			ht.Sni = snis[0]
		}
	}
	//parse hysteria2 insecure
	if insecure, ok := htQuery["insecure"]; ok {
		if len(insecure) > 1 {
			return nil, e.New("multiple hysteria2 insecure").WithPrefix(tagParser)
		} else {
			ht.Insecure = insecure[0]
		}
	}
	//parse hysteria2 pinSHA256
	if pinSHA256s, ok := htQuery["pinSHA256"]; ok {
		if len(pinSHA256s) > 1 {
			return nil, e.New("multiple hysteria2 pinSHA256").WithPrefix(tagParser)
		} else {
			ht.PinSHA256 = pinSHA256s[0]
		}
	}
	return ht, nil
}
