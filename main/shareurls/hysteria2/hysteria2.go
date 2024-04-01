package hysteria2

import (
	e "XrayHelper/main/errors"
	"fmt"
	"strconv"
)

const tagHysteria2 = "hysteria2"

type Hysteria2 struct {
	Remarks      string
	Host         string
	Port         string
	Auth         string
	Obfs         string
	ObfsPassword string
	Sni          string
	Insecure     string
	PinSHA256    string
}

func (this *Hysteria2) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: Hysteria2, Server: %+v, Port: %+v, Auth: %+v, Obfs: %+v, ObfsPassword: %+v, PinSHA256: %+v", this.Remarks, this.Host, this.Port, this.Auth, this.Obfs, this.ObfsPassword, this.PinSHA256)
}

func (this *Hysteria2) ToOutboundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		return nil, e.New("xray core not support hysteria2").WithPrefix(tagHysteria2).WithPathObj(*this)
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "hysteria2"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Host
		outboundObject["server_port"], _ = strconv.Atoi(this.Port)
		outboundObject["obfs"] = getHysteria2ObfsObjectSingbox(this)
		outboundObject["users"] = getHysteria2UsersObjectSingbox(this)
		outboundObject["tls"] = getHysteriaTlsObjectSingbox(this)
		return outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagHysteria2).WithPathObj(*this)
	}
}
