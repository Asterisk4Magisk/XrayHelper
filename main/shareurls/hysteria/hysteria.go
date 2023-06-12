package hysteria

import (
	"XrayHelper/main/errors"
	"fmt"
	"strconv"
)

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

func (this *Hysteria) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		return nil, errors.New("xray core not support hysteria").WithPrefix("hysteria").WithPathObj(*this)
	case "v2ray":
		return nil, errors.New("v2ray core not support hysteria").WithPrefix("hysteria").WithPathObj(*this)
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "hysteria"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Host
		outboundObject["server_port"], _ = strconv.Atoi(this.Port)
		outboundObject["up_mbps"] = this.UpMBPS
		outboundObject["down_mbps"] = this.DownMBPS
		outboundObject["obfs"] = this.ObfsParam
		outboundObject["auth_str"] = this.Auth
		outboundObject["tls"] = getHysteriaTlsObjectSingbox(this)
		return outboundObject, nil
	default:
		return nil, errors.New("unsupported core type " + coreType).WithPrefix("hysteria").WithPathObj(*this)
	}
}
