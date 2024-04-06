package hysteria

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"fmt"
	"strconv"
)

const tagHysteria = "hysteria"

type Hysteria struct {
	Remarks   string
	Host      string
	Port      string
	Protocol  string
	Auth      string
	Peer      string
	Insecure  string
	UpMBPS    string
	DownMBPS  string
	Alpn      string
	Obfs      string
	ObfsParam string
}

func (this *Hysteria) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: Hysteria, Server: %+v, Port: %+v, UpMBPS: %+v, DownMBPS: %+v, ObfsParam: %+v", this.Remarks, this.Host, this.Port, this.UpMBPS, this.DownMBPS, this.ObfsParam)
}

func (this *Hysteria) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		return nil, e.New("xray core not support hysteria").WithPrefix(tagHysteria).WithPathObj(*this)
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "hysteria")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Host)
		port, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", port)
		outboundObject.Set("up_mbps", this.UpMBPS)
		outboundObject.Set("down_mbps", this.DownMBPS)
		outboundObject.Set("obfs", this.ObfsParam)
		outboundObject.Set("auth_str", this.Auth)
		outboundObject.Set("tls", getHysteriaTlsObjectSingbox(this))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagHysteria).WithPathObj(*this)
	}
}
