package vmess

import (
	"XrayHelper/main/errors"
	"fmt"
	"strconv"
)

type Vmess struct {
	Name        string `json:"ps"`
	Address     string `json:"add"`
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
	return fmt.Sprintf("Name: %+v, Type: Vmess, Address: %+v, Port: %+v, Network: %+v, Id: %+v", this.Name, this.Address, this.Port, this.Network, this.Id)
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
	case "v2fly":
		return nil, errors.New("v2fly TODO").WithPrefix("vmess").WithPathObj(*this)
	case "sagernet":
		return nil, errors.New("sagernet TODO").WithPrefix("vmess").WithPathObj(*this)
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("vmess").WithPathObj(*this)
	}
}
