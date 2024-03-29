package shadowsocks

import (
	e "XrayHelper/main/errors"
	"fmt"
	"strconv"
)

const tagShadowsocks = "shadowsocks"

type Shadowsocks struct {
	Remarks  string
	Server   string
	Port     string
	Method   string
	Password string
}

func (this *Shadowsocks) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: Shadowsocks, Server: %+v, Port: %+v, Method: %+v, Password: %+v", this.Remarks, this.Server, this.Port, this.Method, this.Password)
}

func (this *Shadowsocks) ToOutboundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "shadowsocks"
		outboundObject["settings"] = getShadowsocksSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray("tcp")
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "v2ray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "shadowsocks"
		outboundObject["settings"] = getShadowsocksSettingsObjectV2ray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectV2ray("tcp")
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "shadowsocks"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Server
		outboundObject["server_port"], _ = strconv.Atoi(this.Port)
		outboundObject["method"] = this.Method
		outboundObject["password"] = this.Password
		return outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagShadowsocks).WithPathObj(*this)
	}
}
