package shadowsocksr

import (
	"XrayHelper/main/errors"
	"fmt"
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
	return fmt.Sprintf("Name: %+v, Type: ShadowsocksR, Address: %+v, Port: %+v, Method: %+v, Protocol: %+v, Obfs: %+v, Password: %+v", this.Remarks, this.Server, this.Port, this.Method, this.Protocol, this.Obfs, this.Password)
}

func (this *ShadowsocksR) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		return nil, errors.New("xray core not support ShadowsocksR").WithPrefix("ShadowsocksR").WithPathObj(*this)
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "shadowsocksr"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Server
		outboundObject["server_port"] = this.Port
		outboundObject["method"] = this.Method
		outboundObject["password"] = this.Password
		outboundObject["obfs"] = this.Obfs
		outboundObject["obfs_param"] = this.ObfsParam
		outboundObject["protocol"] = this.Protocol
		outboundObject["protocol_param"] = this.ProtoParam
		return outboundObject, nil
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("ShadowsocksR").WithPathObj(*this)
	}
}
