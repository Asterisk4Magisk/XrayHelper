package hysteria2

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls"
	"fmt"
	"github.com/fatih/color"
	"strconv"
)

const tagHysteria2 = "hysteria2"

type Hysteria2 struct {
	Remarks      string
	Host         string
	Port         string
	Auth         string
	Obfs         string
	ObfsPassword string
	Sni          string
	Insecure     string
	PinSHA256    string
}

func (this *Hysteria2) GetNodeInfo() *shareurls.NodeInfo {
	return &shareurls.NodeInfo{
		Remarks:  this.Remarks,
		Type:     "Hysteria2",
		Host:     this.Host,
		Port:     this.Port,
		Protocol: "udp",
	}
}

func (this *Hysteria2) GetNodeInfoStr() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Hysteria2"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Auth: ")+"%+v"+color.BlueString(", Obfs: ")+"%+v"+color.BlueString(", ObfsPassword: ")+"%+v"+color.BlueString(", PinSHA256: ")+"%+v", this.Remarks, this.Host, this.Port, this.Auth, this.Obfs, this.ObfsPassword, this.PinSHA256)
}

func (this *Hysteria2) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		return nil, e.New("xray core not support hysteria2").WithPrefix(tagHysteria2).WithPathObj(*this)
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "hysteria2")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Host)
		port, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", port)
		outboundObject.Set("obfs", getHysteria2ObfsObjectSingbox(this))
		outboundObject.Set("users", getHysteria2UsersObjectSingbox(this))
		outboundObject.Set("tls", getHysteria2TlsObjectSingbox(this))
		return &outboundObject, nil
	case "hysteria2":
		// hysteria2 will ignore tag
		var clientObject serial.OrderedMap
		clientObject.Set("server", this.Host+":"+this.Port)
		clientObject.Set("auth", this.Auth)
		clientObject.Set("obfs", getHysteria2ObfsObjectHysteria2(this))
		clientObject.Set("tls", getHysteria2TlsObjectHysteria2(this))
		return &clientObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagHysteria2).WithPathObj(*this)
	}
}
