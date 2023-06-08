package vmess

import (
	"XrayHelper/main/errors"
	"fmt"
	"strconv"
)

type Vmess struct {
	Remarks     string `json:"ps"`
	Server      string `json:"add"`
	Port        string `json:"port"`
	Id          string `json:"id"`
	AlterId     string `json:"aid"`
	Security    string `json:"scy"`
	Network     string `json:"net"`
	Type        string `json:"type"`
	Host        string `json:"host"`
	Path        string `json:"path"`
	Tls         string `json:"tls"`
	Sni         string `json:"sni"`
	FingerPrint string `json:"fp"`
	Alpn        string `json:"alpn"`
	Version     string `json:"v"`
}

func (this *Vmess) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: Vmess, Server: %+v, Port: %+v, Network: %+v, Id: %+v", this.Remarks, this.Server, this.Port, this.Network, this.Id)
}

func (this *Vmess) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	if version, _ := strconv.Atoi(this.Version); version < 2 {
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
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "vmess"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Server
		outboundObject["server_port"] = this.Port
		outboundObject["uuid"] = this.Id
		outboundObject["security"] = "auto"
		outboundObject["alter_id"], _ = strconv.Atoi(this.AlterId)
		outboundObject["tls"] = getVmessTlsObjectSingbox(this)
		outboundObject["transport"] = getVmessTransportObjectSingbox(this)
		return outboundObject, nil
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("vmess").WithPathObj(*this)
	}
}
