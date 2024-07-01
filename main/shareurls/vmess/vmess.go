package vmess

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls/addon"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

const tagVmess = "vmess"

type String string

func (this *String) UnmarshalJSON(str []byte) error {
	if unquote, err := strconv.Unquote(string(str)); err == nil {
		*this = String(unquote)
	} else {
		s := strings.Trim(string(str), "\"")
		s = strings.Replace(s, `\/`, `/`, -1)
		*this = String(s)
	}
	return nil
}

type Vmess struct {
	Remarks     String `json:"ps"`
	Server      String `json:"add"`
	Port        String `json:"port"`
	Id          String `json:"id"`
	AlterId     String `json:"aid"`
	Security    String `json:"scy"`
	Network     String `json:"net"`
	Type        String `json:"type"`
	Host        String `json:"host"`
	Path        String `json:"path"`
	Tls         String `json:"tls"`
	Sni         String `json:"sni"`
	FingerPrint String `json:"fp"`
	Alpn        String `json:"alpn"`
	Version     String `json:"v"`
}

func (this *Vmess) GetNodeInfo() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Vmess"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Network: ")+"%+v"+color.BlueString(", Id: ")+"%+v", this.Remarks, this.Server, this.Port, this.Network, this.Id)
}

func (this *Vmess) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	if version, _ := strconv.Atoi(string(this.Version)); version < 2 {
		return nil, e.New("unsupported vmess share link version " + this.Version).WithPrefix(tagVmess).WithPathObj(*this)
	}
	addons := &addon.Addon{Alpn: string(this.Alpn), Host: string(this.Host), Path: string(this.Path), Type: string(this.Type), Sni: string(this.Sni), FingerPrint: string(this.FingerPrint), PublicKey: "", ShortId: "", SpiderX: ""}
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", addon.GetMuxObjectXray(false))
		outboundObject.Set("protocol", "vmess")
		outboundObject.Set("settings", getVmessSettingsObjectXray(this))
		outboundObject.Set("streamSettings", addon.GetStreamSettingsObjectXray(addons, string(this.Network), string(this.Tls)))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "vmess")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(string(this.Port))
		outboundObject.Set("server_port", serverPort)
		outboundObject.Set("uuid", this.Id)
		outboundObject.Set("security", this.Security)
		alterId, _ := strconv.Atoi(string(this.AlterId))
		outboundObject.Set("alter_id", alterId)
		outboundObject.Set("tls", addon.GetTlsObjectSingbox(addons, string(this.Tls)))
		outboundObject.Set("transport", addon.GetTransportObjectSingbox(addons, string(this.Network)))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagVmess).WithPathObj(*this)
	}
}
