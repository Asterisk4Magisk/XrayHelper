package shadowsocks

import (
	"XrayHelper/main/errors"
	"fmt"
)

type Shadowsocks struct {
	Name     string
	Address  string
	Port     string
	Method   string
	Password string
}

func (this *Shadowsocks) GetNodeInfo() string {
	return fmt.Sprintf("Name: %+v, Type: Shadowsocks, Address: %+v, Port: %+v, Method: %+v, Password: %+v", this.Name, this.Address, this.Port, this.Method, this.Password)
}

func (this *Shadowsocks) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "shadowsocks"
		outboundObject["settings"] = getShadowsocksSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray("tcp")
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "singbox":
		return nil, errors.New("singbox TODO").WithPrefix("shadowsocks").WithPathObj(*this)
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("shadowsocks").WithPathObj(*this)
	}
}
