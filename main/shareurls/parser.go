package shareurls

import (
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/shareurls/addon"
	"XrayHelper/main/shareurls/hysteria"
	"XrayHelper/main/shareurls/hysteria2"
	"XrayHelper/main/shareurls/shadowsocks"
	"XrayHelper/main/shareurls/socks"
	"XrayHelper/main/shareurls/trojan"
	"XrayHelper/main/shareurls/vless"
	"XrayHelper/main/shareurls/vmess"
	"XrayHelper/main/shareurls/vmessaead"
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
	if len(userAndPassword) == 2 {
		so.User = userAndPassword[0]
		so.Password = userAndPassword[1]
	}
	return so, nil
}

// parseAddon parse v2ray addon
func parseAddon(aUrl string, network string, security string) (*addon.Addon, error) {
	addon := new(addon.Addon)
	parse, _ := url.Parse(aUrl)
	query, _ := url.ParseQuery(parse.RawQuery)
	switch network {
	case "tcp", "raw":
		//parse addon headerType
		if headerTypes, ok := query["headerType"]; ok && len(headerTypes) == 1 {
			if headerTypes[0] == "http" {
				if hosts, ok := query["host"]; ok && len(hosts) == 1 {
					addon.Host = hosts[0]
				}
			}
		}
	case "kcp":
		//parse addon headerType
		if headerTypes, ok := query["headerType"]; ok && len(headerTypes) == 1 {
			addon.Type = headerTypes[0]
		}
		//parse addon kcp seed
		if seeds, ok := query["seed"]; ok && len(seeds) == 1 {
			addon.Path = seeds[0]
		}
	case "ws", "http", "h2", "h3", "httpupgrade", "splithttp":
		//parse addon host
		if hosts, ok := query["host"]; ok && len(hosts) == 1 {
			addon.Host = hosts[0]
		}
		//parse addon path
		if paths, ok := query["path"]; ok && len(paths) == 1 {
			addon.Path = paths[0]
		}
	case "quic":
		//parse addon headerType
		if headerTypes, ok := query["headerType"]; ok && len(headerTypes) == 1 {
			addon.Type = headerTypes[0]
		}
		//parse addon quicSecurity
		if quicSecurity, ok := query["quicSecurity"]; ok && len(quicSecurity) == 1 {
			addon.Host = quicSecurity[0]
		}
		//parse addon quicKey
		if quicKey, ok := query["key"]; ok && len(quicKey) == 1 {
			addon.Path = quicKey[0]
		}
	case "grpc":
		//parse addon grpc authority
		if authority, ok := query["authority"]; ok && len(authority) == 1 {
			addon.Host = authority[0]
		}
		//parse addon grpc mode
		if modes, ok := query["mode"]; ok && len(modes) == 1 {
			addon.Type = modes[0]
		} else {
			addon.Type = "gun"
		}
		//parse addon grpc serviceName
		if serviceNames, ok := query["serviceName"]; ok && len(serviceNames) == 1 {
			addon.Path = serviceNames[0]
		}
	default:
		return nil, e.New("unknown v2ray addon transport type " + network).WithPrefix(tagParser)
	}
	switch security {
	case "none":
		break
	case "tls":
		//parse addon tls sni
		if sni, ok := query["sni"]; ok && len(sni) == 1 {
			addon.Sni = sni[0]
		}
		//parse addon tls fingerprint
		if fps, ok := query["fp"]; ok && len(fps) == 1 {
			addon.FingerPrint = fps[0]
		}
		//parse addon tls Alpn
		if alpns, ok := query["alpn"]; ok && len(alpns) == 1 {
			addon.Alpn = alpns[0]
		}
	case "reality":
		//parse addon reality sni
		if sni, ok := query["sni"]; ok && len(sni) == 1 {
			addon.Sni = sni[0]
		}
		//parse addon reality fingerprint
		if fps, ok := query["fp"]; ok && len(fps) == 1 {
			addon.FingerPrint = fps[0]
		} else {
			addon.FingerPrint = "chrome"
		}
		//parse addon reality PublicKey
		if publicKeys, ok := query["pbk"]; ok && len(publicKeys) == 1 {
			addon.PublicKey = publicKeys[0]
		}
		//parse addon reality ShortId
		if shortIds, ok := query["sid"]; ok && len(shortIds) == 1 {
			addon.ShortId = shortIds[0]
		}
		//parse addon reality SpiderX
		if spiderX, ok := query["spx"]; ok && len(spiderX) == 1 {
			addon.SpiderX = spiderX[0]
		}
	default:
		return nil, e.New("unknown v2ray addon security type " + security).WithPrefix(tagParser)
	}
	return addon, nil
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
	//parse addon
	if addons, err := parseAddon(trojanUrl, tj.Network, tj.Security); err != nil {
		return nil, err
	} else {
		tj.Addon = *addons
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
		vl.Security = "none"
	} else if len(security) > 1 {
		return nil, e.New("multiple VLESS security type").WithPrefix(tagParser)
	} else if vl.Security = security[0]; vl.Security == "" {
		return nil, e.New("empty VLESS security type").WithPrefix(tagParser)
	}
	//parse addon
	if addons, err := parseAddon(vlessUrl, vl.Network, vl.Security); err != nil {
		return nil, err
	} else {
		vl.Addon = *addons
	}
	return vl, nil
}

// parseVmess parse Vmess url
func parseVmess(vmessUrl string) (ShareUrl, error) {
	originJson, err := common.DecodeBase64(strings.TrimPrefix(vmessUrl, "vmess://"))
	if err != nil {
		return parseVmessAEAD(vmessUrl)
	}
	v2 := new(vmess.Vmess)
	err = json.Unmarshal([]byte(originJson), v2)
	if err != nil {
		return nil, e.New("unmarshal origin json failed, ", err).WithPrefix(tagParser)
	}
	return v2, nil
}

// parseVmessAEAD parse VmessAEAD url
func parseVmessAEAD(vmessUrl string) (ShareUrl, error) {
	vm := new(vmessaead.VmessAEAD)
	vmParse, err := url.Parse(vmessUrl)
	if err != nil {
		return nil, e.New("VmessAEAD url parse err, ", err).WithPrefix(tagParser)
	}
	vm.Remarks = vmParse.Fragment
	vm.Id = vmParse.User.Username()
	vm.Server = vmParse.Hostname()
	vm.Port = vmParse.Port()
	vmQuery, err := url.ParseQuery(vmParse.RawQuery)
	if err != nil {
		return nil, e.New("VmessAEAD url parse query err, ", err).WithPrefix(tagParser)
	}
	//parse VmessAEAD encryption
	if encryption, ok := vmQuery["encryption"]; !ok {
		vm.Encryption = "auto"
	} else if len(encryption) > 1 {
		return nil, e.New("multiple VmessAEAD encryption").WithPrefix(tagParser)
	} else if vm.Encryption = encryption[0]; vm.Encryption == "" {
		return nil, e.New("empty VmessAEAD encryption").WithPrefix(tagParser)
	}
	//parse VmessAEAD network
	if types, ok := vmQuery["type"]; !ok {
		vm.Network = "tcp"
	} else if len(types) > 1 {
		return nil, e.New("multiple VmessAEAD transport type").WithPrefix(tagParser)
	} else if vm.Network = types[0]; vm.Network == "" {
		return nil, e.New("empty VmessAEAD transport type").WithPrefix(tagParser)
	}
	//parse VmessAEAD security
	if security, ok := vmQuery["security"]; !ok {
		vm.Security = "none"
	} else if len(security) > 1 {
		return nil, e.New("multiple VmessAEAD security type").WithPrefix(tagParser)
	} else if vm.Security = security[0]; vm.Security == "" {
		return nil, e.New("empty VmessAEAD security type").WithPrefix(tagParser)
	}
	//parse addon
	if addons, err := parseAddon(vmessUrl, vm.Network, vm.Security); err != nil {
		return nil, err
	} else {
		vm.Addon = *addons
	}
	return vm, nil
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
