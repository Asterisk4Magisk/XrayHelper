package vless

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"XrayHelper/main/shareurls/addon"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

const tagVless = "vless"

type VLESS struct {
	//basic
	Remarks    string
	Id         string
	Server     string
	Port       string
	Encryption string
	Flow       string
	Network    string
	Security   string

	//addon
	Addon addon.Addon
}

func (this *VLESS) GetNodeInfo() *shareurls.NodeInfo {
	return &shareurls.NodeInfo{
		Remarks:  this.Remarks,
		Type:     "VLESS",
		Host:     this.Server,
		Port:     this.Port,
		Protocol: this.Network,
	}
}

func (this *VLESS) GetNodeInfoStr() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"VLESS"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Flow: ")+"%+v"+color.BlueString(", Network: ")+"%+v"+color.BlueString(", Id: ")+"%+v", this.Remarks, this.Server, this.Port, this.Flow, this.Network, this.Id)
}

func (this *VLESS) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", addon.GetMuxObjectXray(false))
		outboundObject.Set("protocol", "vless")
		outboundObject.Set("settings", getVLESSSettingsObjectXray(this))
		outboundObject.Set("streamSettings", addon.GetStreamSettingsObjectXray(&this.Addon, this.Network, this.Security))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "vless")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		outboundObject.Set("uuid", this.Id)
		outboundObject.Set("flow", this.Flow)
		outboundObject.Set("tls", addon.GetTlsObjectSingbox(&this.Addon, this.Security))
		outboundObject.Set("transport", addon.GetTransportObjectSingbox(&this.Addon, this.Network))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagVless).WithPathObj(*this)
	}
}
