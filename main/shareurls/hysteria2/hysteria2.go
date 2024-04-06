package hysteria2

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
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

func (this *Hysteria2) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		return nil, e.New("xray core not support hysteria2").WithPrefix(tagHysteria2).WithPathObj(*this)
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "hysteria2")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Host)
		port, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", port)
		outboundObject.Set("obfs", getHysteria2ObfsObjectSingbox(this))
		outboundObject.Set("users", getHysteria2UsersObjectSingbox(this))
		outboundObject.Set("tls", getHysteriaTlsObjectSingbox(this))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagHysteria2).WithPathObj(*this)
	}
}
