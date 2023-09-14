package vmess

import (
	"XrayHelper/main/errors"
	"fmt"
	"strconv"
	"strings"
)

type String string

func (this *String) UnmarshalJSON(port []byte) error {
	*this = String(strings.ReplaceAll(string(port), "\"", ""))
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
	return fmt.Sprintf("Remarks: %+v, Type: Vmess, Server: %+v, Port: %+v, Network: %+v, Id: %+v", this.Remarks, this.Server, this.Port, this.Network, this.Id)
}

func (this *Vmess) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	if version, _ := strconv.Atoi(string(this.Version)); version < 2 {
		return nil, errors.New("unsupported vmess share link version " + this.Version).WithPrefix("vmess").WithPathObj(*this)
	}
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "vmess"
		outboundObject["settings"] = getVmessSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray(this)
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "v2ray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectV2ray(false)
		outboundObject["protocol"] = "vmess"
		outboundObject["settings"] = getVmessSettingsObjectV2ray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectV2ray(this)
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "vmess"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Server
		outboundObject["server_port"], _ = strconv.Atoi(string(this.Port))
		outboundObject["uuid"] = this.Id
		outboundObject["security"] = "auto"
		outboundObject["alter_id"], _ = strconv.Atoi(string(this.AlterId))
		outboundObject["tls"] = getVmessTlsObjectSingbox(this)
		outboundObject["transport"] = getVmessTransportObjectSingbox(this)
		return outboundObject, nil
	default:
		return nil, errors.New("unsupported core type " + coreType).WithPrefix("vmess").WithPathObj(*this)
	}
}
