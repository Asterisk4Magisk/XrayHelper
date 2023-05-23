package socks

import (
	"XrayHelper/main/errors"
	"fmt"
)

type Socks struct {
	Name     string
	Address  string
	Port     string
	User     string
	Password string
}

func (this *Socks) GetNodeInfo() string {
	return fmt.Sprintf("Name: %+v, Type: Socks, Address: %+v, Port: %+v, User: %+v, Password: %+v", this.Name, this.Address, this.Port, this.User, this.Password)
}

func (this *Socks) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "socks"
		outboundObject["settings"] = getSocksSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray("tcp")
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "v2fly":
		return nil, errors.New("v2fly TODO").WithPrefix("socks").WithPathObj(*this)
	case "sagernet":
		return nil, errors.New("sagernet TODO").WithPrefix("socks").WithPathObj(*this)
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("socks").WithPathObj(*this)
	}
}
