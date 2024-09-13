package socks

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"XrayHelper/main/shareurls/addon"
	"fmt"
	"github.com/fatih/color"
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

func (this *Socks) GetNodeInfo() *shareurls.NodeInfo {
	return &shareurls.NodeInfo{
		Remarks:  this.Remarks,
		Type:     "Socks",
		Host:     this.Server,
		Port:     this.Port,
		Protocol: "tcp",
	}
}

func (this *Socks) GetNodeInfoStr() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Socks"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", User: ")+"%+v"+color.BlueString(", Password: ")+"%+v", this.Remarks, this.Server, this.Port, this.User, this.Password)
}

func (this *Socks) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", addon.GetMuxObjectXray(false))
		outboundObject.Set("protocol", "socks")
		outboundObject.Set("settings", getSocksSettingsObjectXray(this))
		outboundObject.Set("streamSettings", getStreamSettingsObjectXray("tcp"))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "socks")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		if len(this.User) > 0 && this.User != "null" {
			outboundObject.Set("username", this.User)
			outboundObject.Set("password", this.Password)
		}
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagSocks).WithPathObj(*this)
	}
}
