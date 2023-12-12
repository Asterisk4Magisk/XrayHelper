package shadowsocksr

import (
	e "XrayHelper/main/errors"
	"fmt"
	"strconv"
)

type ShadowsocksR struct {
	Remarks    string
	Server     string
	Port       string
	Protocol   string
	ProtoParam string
	Method     string
	Obfs       string
	ObfsParam  string
	Password   string
}

func (this *ShadowsocksR) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: ShadowsocksR, Server: %+v, Port: %+v, Method: %+v, Protocol: %+v, Obfs: %+v, Password: %+v", this.Remarks, this.Server, this.Port, this.Method, this.Protocol, this.Obfs, this.Password)
}

func (this *ShadowsocksR) ToOutboundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		return nil, e.New("xray core not support shadowsocksr").WithPrefix("shadowsocksr").WithPathObj(*this)
	case "v2ray":
		return nil, e.New("v2ray core not support shadowsocksr").WithPrefix("shadowsocksr").WithPathObj(*this)
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "shadowsocksr"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Server
		outboundObject["server_port"], _ = strconv.Atoi(this.Port)
		outboundObject["method"] = this.Method
		outboundObject["password"] = this.Password
		outboundObject["obfs"] = this.Obfs
		outboundObject["obfs_param"] = this.ObfsParam
		outboundObject["protocol"] = this.Protocol
		outboundObject["protocol_param"] = this.ProtoParam
		return outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix("shadowsocksr").WithPathObj(*this)
	}
}
