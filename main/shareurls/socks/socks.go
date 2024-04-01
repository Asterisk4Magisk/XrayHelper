package socks

import (
	e "XrayHelper/main/errors"
	"fmt"
	"strconv"
)

const tagSocks = "socks"

type Socks struct {
	Remarks  string
	Server   string
	Port     string
	User     string
	Password string
}

func (this *Socks) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: Socks, Server: %+v, Port: %+v, User: %+v, Password: %+v", this.Remarks, this.Server, this.Port, this.User, this.Password)
}

func (this *Socks) ToOutboundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "socks"
		outboundObject["settings"] = getSocksSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray("tcp")
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "socks"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Server
		outboundObject["server_port"], _ = strconv.Atoi(this.Port)
		if this.User != "null" {
			outboundObject["username"] = this.User
			outboundObject["password"] = this.Password
		}
		return outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagSocks).WithPathObj(*this)
	}
}
