package shareurls

import (
	"XrayHelper/main/common"
	"XrayHelper/main/errors"
	"XrayHelper/main/shareurls/shadowsocks"
	"XrayHelper/main/shareurls/socks"
	"XrayHelper/main/shareurls/trojan"
	"XrayHelper/main/shareurls/vless"
	"XrayHelper/main/shareurls/vmess"
	"encoding/json"
	"net/url"
	"strings"
)

// parseShadowsocks parse shadowsocks url
func parseShadowsocks(ssUrl string) (ShareUrl, error) {
	ss := new(shadowsocks.Shadowsocks)
	ssParse, err := url.Parse(ssUrl)
	if err != nil {
		return nil, errors.New("shadowsocks url parse err, ", err).WithPrefix("shareurls")
	}
	ss.Name = ssParse.Fragment
	ss.Address = ssParse.Hostname()
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
		ss.Address = addressAndPort[0]
		ss.Port = addressAndPort[1]
	} else {
		info, err := common.DecodeBase64(ssParse.User.Username())
		if err != nil {
			return nil, err
		}
		methodAndPassword := strings.Split(info, ":")
		ss.Method = methodAndPassword[0]
		ss.Password = methodAndPassword[1]
	}
	return ss, nil
}

// parseSocks parse socks url
func parseSocks(socksUrl string) (ShareUrl, error) {
	so := new(socks.Socks)
	soParse, err := url.Parse(socksUrl)
	if err != nil {
		return nil, errors.New("socks url parse err, ", err).WithPrefix("shareurls")
	}
	so.Name = soParse.Fragment
	so.Address = soParse.Hostname()
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
		return nil, errors.New("trojan url parse err, ", err).WithPrefix("shareurls")
	}
	tj.Name = tjParse.Fragment
	tj.Password = tjParse.User.Username()
	tj.Address = tjParse.Hostname()
	tj.Port = tjParse.Port()
	tjQuery, err := url.ParseQuery(tjParse.RawQuery)
	if err != nil {
		return nil, errors.New("trojan url parse query err, ", err).WithPrefix("shareurls")
	}
	//parse trojan network
	if types, ok := tjQuery["type"]; !ok {
		tj.Network = "tcp"
	} else if len(types) > 1 {
		return nil, errors.New("multiple trojan transport type").WithPrefix("shareurls")
	} else if tj.Network = types[0]; tj.Network == "" {
		return nil, errors.New("empty trojan transport type").WithPrefix("shareurls")
	}
	//parse trojan security
	if security, ok := tjQuery["security"]; !ok {
		tj.Security = "tls"
	} else if len(security) > 1 {
		return nil, errors.New("multiple trojan security type").WithPrefix("shareurls")
	} else if tj.Security = security[0]; tj.Security == "" {
		return nil, errors.New("empty trojan security type").WithPrefix("shareurls")
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
	case "ws":
		//parse trojan host
		if hosts, ok := tjQuery["host"]; ok && len(hosts) == 1 {
			tj.Host = hosts[0]
		}
		//parse trojan path
		if paths, ok := tjQuery["path"]; ok && len(paths) == 1 {
			tj.Path = paths[0]
		}
	case "http":
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
		//parse trojan grpc mode
		if modes, ok := tjQuery["mode"]; ok && len(modes) == 1 {
			tj.Type = modes[0]
		}
		//parse trojan grpc serviceName
		if serviceNames, ok := tjQuery["serviceName"]; ok && len(serviceNames) == 1 {
			tj.Path = serviceNames[0]
		}
	default:
		return nil, errors.New("unknown trojan transport type " + tj.Network).WithPrefix("shareurls")
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
		} else {
			tj.FingerPrint = "firefox"
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
		} else {
			tj.FingerPrint = "firefox"
		}
		//parse trojan reality PublicKey
		if publicKeys, ok := tjQuery["pbx"]; ok && len(publicKeys) == 1 {
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
		return nil, errors.New("unknown trojan security type " + tj.Security).WithPrefix("shareurls")
	}
	return tj, nil
}

// parseVLESS parse VLESS url
func parseVLESS(vlessUrl string) (ShareUrl, error) {
	vl := new(vless.VLESS)
	vlParse, err := url.Parse(vlessUrl)
	if err != nil {
		return nil, errors.New("VLESS url parse err, ", err).WithPrefix("shareurls")
	}
	vl.Name = vlParse.Fragment
	vl.Id = vlParse.User.Username()
	vl.Address = vlParse.Hostname()
	vl.Port = vlParse.Port()
	vlQuery, err := url.ParseQuery(vlParse.RawQuery)
	if err != nil {
		return nil, errors.New("VLESS url parse query err, ", err).WithPrefix("shareurls")
	}
	//parse VLESS encryption
	if encryption, ok := vlQuery["encryption"]; !ok {
		vl.Encryption = "none"
	} else if len(encryption) > 1 {
		return nil, errors.New("multiple VLESS encryption").WithPrefix("shareurls")
	} else if vl.Encryption = encryption[0]; vl.Encryption == "" {
		return nil, errors.New("empty VLESS encryption").WithPrefix("shareurls")
	}
	//parse VLESS flow
	if flows, ok := vlQuery["flow"]; !ok {
		vl.Flow = ""
	} else if len(flows) > 1 {
		return nil, errors.New("multiple VLESS flow").WithPrefix("shareurls")
	} else {
		vl.Flow = flows[0]
	}
	//parse VLESS network
	if types, ok := vlQuery["type"]; !ok {
		vl.Network = "tcp"
	} else if len(types) > 1 {
		return nil, errors.New("multiple VLESS transport type").WithPrefix("shareurls")
	} else if vl.Network = types[0]; vl.Network == "" {
		return nil, errors.New("empty VLESS transport type").WithPrefix("shareurls")
	}
	//parse VLESS security
	if security, ok := vlQuery["security"]; !ok {
		vl.Security = "tls"
	} else if len(security) > 1 {
		return nil, errors.New("multiple VLESS security type").WithPrefix("shareurls")
	} else if vl.Security = security[0]; vl.Security == "" {
		return nil, errors.New("empty VLESS security type").WithPrefix("shareurls")
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
	case "ws":
		//parse VLESS host
		if hosts, ok := vlQuery["host"]; ok && len(hosts) == 1 {
			vl.Host = hosts[0]
		}
		//parse VLESS path
		if paths, ok := vlQuery["path"]; ok && len(paths) == 1 {
			vl.Path = paths[0]
		}
	case "http":
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
		//parse VLESS grpc mode
		if modes, ok := vlQuery["mode"]; ok && len(modes) == 1 {
			vl.Type = modes[0]
		}
		//parse VLESS grpc serviceName
		if serviceNames, ok := vlQuery["serviceName"]; ok && len(serviceNames) == 1 {
			vl.Path = serviceNames[0]
		}
	default:
		return nil, errors.New("unknown VLESS transport type " + vl.Network).WithPrefix("shareurls")
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
		} else {
			vl.FingerPrint = "firefox"
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
		} else {
			vl.FingerPrint = "firefox"
		}
		//parse VLESS reality PublicKey
		if publicKeys, ok := vlQuery["pbx"]; ok && len(publicKeys) == 1 {
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
		return nil, errors.New("unknown VLESS security type " + vl.Security).WithPrefix("shareurls")
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
		return nil, errors.New("unmarshal origin json failed, ", err).WithPrefix("shareurls")
	}
	return v2, nil
}
