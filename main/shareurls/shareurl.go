package shareurls

import (
	"XrayHelper/main/errors"
	"XrayHelper/main/shareurls/shadowsocks"
	"XrayHelper/main/shareurls/trojan"
	"XrayHelper/main/shareurls/vmess"
	"XrayHelper/main/utils"
	"encoding/json"
	"net/url"
	"strings"
)

const (
	socksPrefix  = "socks://"
	ssPrefix     = "ss://"
	vmessPrefix  = "vmess://"
	vlessPrefix  = "vless://"
	trojanPrefix = "trojan://"
)

// ShareUrl implement this interface, that node can be converted to xray OutoundObject
type ShareUrl interface {
	GetNodeInfo() string
	ToOutoundWithTag(coreType string, tag string) (interface{}, error)
}

// NewShareUrl parse the url, return a ShareUrl
func NewShareUrl(link string) (ShareUrl, error) {
	if strings.HasPrefix(link, socksPrefix) {
		return newSocksShareUrl(strings.TrimPrefix(link, socksPrefix))
	}
	if strings.HasPrefix(link, ssPrefix) {
		return newShadowsocksShareUrl(link)
	}
	if strings.HasPrefix(link, vmessPrefix) {
		return newVmessShareUrl(strings.TrimPrefix(link, vmessPrefix))
	}
	if strings.HasPrefix(link, vlessPrefix) {
		return newVLESSShareUrl(strings.TrimPrefix(link, vlessPrefix))
	}
	if strings.HasPrefix(link, trojanPrefix) {
		return newTrojanShareUrl(link)
	}
	return nil, errors.New("not a supported share link").WithPrefix("shareurls")
}

// NewShareUrl parse the url, return a ShareUrl
func newShadowsocksShareUrl(ssUrl string) (ShareUrl, error) {
	ss := new(shadowsocks.Shadowsocks)
	ssParse, err := url.Parse(ssUrl)
	if err != nil {
		return nil, errors.New("url parse err, ", err).WithPrefix("shareurls")
	}
	ss.Name = ssParse.Fragment
	ss.Address = ssParse.Hostname()
	ss.Port = ssParse.Port()
	info, err := utils.DecodeBase64(ssParse.User.Username())
	if err != nil {
		return nil, err
	}
	methodAndPassword := strings.Split(info, ":")
	ss.Method = methodAndPassword[0]
	ss.Password = methodAndPassword[1]
	return ss, nil
}

// newSocksShareUrl parse socks url
func newSocksShareUrl(socksUrl string) (ShareUrl, error) {
	// TODO
	return nil, errors.New("socks TODO").WithPrefix("shareurls")
}

// newTrojanShareUrl parse trojan url
func newTrojanShareUrl(trojanUrl string) (ShareUrl, error) {
	tj := new(trojan.Trojan)
	tjParse, err := url.Parse(trojanUrl)
	if err != nil {
		return nil, errors.New("url parse err, ", err).WithPrefix("shareurls")
	}
	tj.Name = tjParse.Fragment
	tj.Password = tjParse.User.Username()
	tj.Address = tjParse.Hostname()
	tj.Port = tjParse.Port()
	tjQuery, err := url.ParseQuery(tjParse.RawQuery)
	if err != nil {
		return nil, errors.New("url parse query err, ", err).WithPrefix("shareurls")
	}
	//parse trojan network
	if types, ok := tjQuery["type"]; !ok {
		return nil, errors.New("cannot get trojan transport type").WithPrefix("shareurls")
	} else if len(types) > 1 {
		return nil, errors.New("multiple trojan transport type").WithPrefix("shareurls")
	} else if tj.Network = types[0]; tj.Network == "" {
		return nil, errors.New("empty trojan transport type").WithPrefix("shareurls")
	}
	//parse trojan security
	if security, ok := tjQuery["security"]; !ok {
		return nil, errors.New("cannot get trojan security type").WithPrefix("shareurls")
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

// newVLESSShareUrl parse VLESS url
func newVLESSShareUrl(vlessUrl string) (ShareUrl, error) {
	// TODO
	return nil, errors.New("vless TODO").WithPrefix("shareurls")
}

// newVmessShareUrl parse Vmess url
func newVmessShareUrl(vmessUrl string) (ShareUrl, error) {
	v2 := new(vmess.Vmess)
	originJson, err := utils.DecodeBase64(vmessUrl)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(originJson), v2)
	if err != nil {
		return nil, errors.New("unmarshal origin json failed, ", err).WithPrefix("shareurls")
	}
	return v2, nil
}
