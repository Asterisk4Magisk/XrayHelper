package trojan

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"XrayHelper/main/shareurls/addon"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

const tagTrojan = "trojan"

type Trojan struct {
	//basic
	Remarks  string
	Password string
	Server   string
	Port     string
	Network  string
	Security string

	//addon
	Addon addon.Addon
}

func (this *Trojan) GetNodeInfo() *shareurls.NodeInfo {
	return &shareurls.NodeInfo{
		Remarks:  this.Remarks,
		Type:     "Trojan",
		Host:     this.Server,
		Port:     this.Port,
		Protocol: this.Network,
	}
}

func (this *Trojan) GetNodeInfoStr() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Trojan"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Network: ")+"%+v"+color.BlueString(", Password: ")+"%+v", this.Remarks, this.Server, this.Port, this.Network, this.Password)
}

func (this *Trojan) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", addon.GetMuxObjectXray(false))
		outboundObject.Set("protocol", "trojan")
		outboundObject.Set("settings", getTrojanSettingsObjectXray(this))
		outboundObject.Set("streamSettings", addon.GetStreamSettingsObjectXray(&this.Addon, this.Network, this.Security))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "trojan")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		outboundObject.Set("password", this.Password)
		outboundObject.Set("tls", addon.GetTlsObjectSingbox(&this.Addon, this.Security))
		outboundObject.Set("transport", addon.GetTransportObjectSingbox(&this.Addon, this.Network))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagTrojan).WithPathObj(*this)
	}
}
