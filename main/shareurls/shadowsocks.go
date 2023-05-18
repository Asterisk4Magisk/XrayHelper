package shareurls

import (
	"XrayHelper/main/errors"
	"XrayHelper/main/utils"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Shadowsocks struct {
	name     string
	address  string
	port     uint16
	method   string
	password string
}

func (this *Shadowsocks) GetNodeInfo() string {
	return fmt.Sprintf("Name: %+v, Type: Shadowsocks, Address: %+v, Port: %+v, Method: %+v, Password: %+v", this.name, this.address, this.port, this.method, this.password)
}

func (this *Shadowsocks) ToOutoundWithTag(tag string) interface{} {
	outboundObject := make(map[string]interface{})
	outboundObject["mux"] = getMuxObject(false)
	outboundObject["protocol"] = "shadowsocks"
	outboundObject["settings"] = getShadowsocksSettingsObject(this)
	outboundObject["streamSettings"] = getStreamSettingsObject("tcp")
	outboundObject["tag"] = tag
	return outboundObject
}

func newShadowsocksShareUrl(ssUrl string) (ShareUrl, error) {
	ss := new(Shadowsocks)
	nodeAndName := strings.Split(ssUrl, "#")
	name, err := url.QueryUnescape(nodeAndName[1])
	if err != nil {
		return nil, errors.New("unescape node name failed, ", err).WithPrefix("shadowsocks")
	}
	ss.name = name
	infoAndServer := strings.Split(nodeAndName[0], "@")
	addressAndPort := strings.Split(infoAndServer[1], ":")
	ss.address = addressAndPort[0]
	port, err := strconv.Atoi(addressAndPort[1])
	if err != nil {
		return nil, errors.New("convert node port failed, ", err).WithPrefix("shadowsocks")
	}
	ss.port = uint16(port)
	info, err := utils.DecodeBase64(infoAndServer[0])
	if err != nil {
		return nil, err
	}
	methodAndPassword := strings.Split(info, ":")
	ss.method = methodAndPassword[0]
	ss.password = methodAndPassword[1]
	return ss, nil
}
