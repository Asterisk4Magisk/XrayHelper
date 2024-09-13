package vmessaead

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls/addon"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

const tagVmessAEAD = "vmessaead"

type VmessAEAD struct {
	//basic
	Remarks    string
	Id         string
	Server     string
	Port       string
	Encryption string
	Network    string
	Security   string

	//addon
	Addon addon.Addon
}

func (this *VmessAEAD) GetNodeInfo() *addon.NodeInfo {
	return &addon.NodeInfo{
		Remarks:  this.Remarks,
		Type:     "VMess",
		Host:     this.Server,
		Port:     this.Port,
		Protocol: this.Network,
	}
}

func (this *VmessAEAD) GetNodeInfoStr() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Vmess"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Network: ")+"%+v"+color.BlueString(", Id: ")+"%+v", this.Remarks, this.Server, this.Port, this.Network, this.Id)
}

func (this *VmessAEAD) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", addon.GetMuxObjectXray(false))
		outboundObject.Set("protocol", "vmess")
		outboundObject.Set("settings", getVmessSettingsObjectXray(this))
		outboundObject.Set("streamSettings", addon.GetStreamSettingsObjectXray(&this.Addon, this.Network, this.Security))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "vmess")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		outboundObject.Set("uuid", this.Id)
		outboundObject.Set("security", this.Encryption)
		outboundObject.Set("alter_id", 0)
		outboundObject.Set("tls", addon.GetTlsObjectSingbox(&this.Addon, this.Security))
		outboundObject.Set("transport", addon.GetTransportObjectSingbox(&this.Addon, this.Network))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagVmessAEAD).WithPathObj(*this)
	}
}
