package shareurls

import (
	"XrayHelper/main/errors"
	"XrayHelper/main/shareurls/shadowsocks"
	"XrayHelper/main/shareurls/vmess"
	"XrayHelper/main/utils"
	"encoding/json"
	"net/url"
	"strconv"
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
		return newShadowsocksShareUrl(strings.TrimPrefix(link, ssPrefix))
	}
	if strings.HasPrefix(link, vmessPrefix) {
		return newVmessShareUrl(strings.TrimPrefix(link, vmessPrefix))
	}
	if strings.HasPrefix(link, vlessPrefix) {
		return newVLESSShareUrl(strings.TrimPrefix(link, vlessPrefix))
	}
	if strings.HasPrefix(link, trojanPrefix) {
		return newTrojanShareUrl(strings.TrimPrefix(link, trojanPrefix))
	}
	return nil, errors.New("not a supported share link").WithPrefix("shareurls")
}

// NewShareUrl parse the url, return a ShareUrl
func newShadowsocksShareUrl(ssUrl string) (ShareUrl, error) {
	ss := new(shadowsocks.Shadowsocks)
	nodeAndName := strings.Split(ssUrl, "#")
	name, err := url.QueryUnescape(nodeAndName[1])
	if err != nil {
		return nil, errors.New("unescape shadowsocks node name failed, ", err).WithPrefix("shareurls")
	}
	ss.Name = name
	infoAndServer := strings.Split(nodeAndName[0], "@")
	addressAndPort := strings.Split(infoAndServer[1], ":")
	ss.Address = addressAndPort[0]
	port, err := strconv.Atoi(addressAndPort[1])
	if err != nil {
		return nil, errors.New("convert shadowsocks node port failed, ", err).WithPrefix("shareurls")
	}
	ss.Port = uint16(port)
	info, err := utils.DecodeBase64(infoAndServer[0])
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
	// TODO
	return nil, errors.New("trojan TODO").WithPrefix("shareurls")
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
