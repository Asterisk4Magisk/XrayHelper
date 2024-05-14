package shadowsocks

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

const tagShadowsocks = "shadowsocks"

type Shadowsocks struct {
	Remarks   string
	Server    string
	Port      string
	Method    string
	Password  string
	Plugin    string
	PluginOpt string
}

func (this *Shadowsocks) GetNodeInfo() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Shadowsocks"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Method: ")+"%+v"+color.BlueString(", Password: ")+"%+v", this.Remarks, this.Server, this.Port, this.Method, this.Password)
}

func (this *Shadowsocks) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", getMuxObjectXray(false))
		outboundObject.Set("protocol", "shadowsocks")
		outboundObject.Set("settings", getShadowsocksSettingsObjectXray(this))
		outboundObject.Set("streamSettings", getStreamSettingsObjectXray("tcp"))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "shadowsocks")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		outboundObject.Set("method", this.Method)
		outboundObject.Set("password", this.Password)
		if len(this.Plugin) > 0 {
			outboundObject.Set("plugin", this.Plugin)
		}
		if len(this.PluginOpt) > 0 {
			outboundObject.Set("plugin_opts", this.PluginOpt)
		}
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagShadowsocks).WithPathObj(*this)
	}
}
