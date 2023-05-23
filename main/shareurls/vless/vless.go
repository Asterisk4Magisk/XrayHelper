package vless

import (
	"XrayHelper/main/errors"
	"fmt"
)

type VLESS struct {
	//basic
	Name       string
	Id         string
	Address    string
	Port       string
	Encryption string
	Flow       string
	Network    string
	Security   string

	//addon
	//http/ws/h2->host quic->security
	Host string
	//ws/h2->path quic->key kcp->seed grpc->serviceName
	Path string
	//tcp/kcp/quic->type grpc->mode
	Type string

	//tls
	Sni         string
	FingerPrint string
	Alpn        string
	//reality
	PublicKey string //pbx
	ShortId   string //sid
	SpiderX   string //spx
}

func (this *VLESS) GetNodeInfo() string {
	return fmt.Sprintf("Name: %+v, Type: VLESS, Address: %+v, Port: %+v, Flow: %+v, Network: %+v, Id: %+v", this.Name, this.Address, this.Port, this.Flow, this.Network, this.Id)
}

func (this *VLESS) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "vless"
		outboundObject["settings"] = getVLESSSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray(this)
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "v2fly":
		return nil, errors.New("v2fly TODO").WithPrefix("VLESS").WithPathObj(*this)
	case "sagernet":
		return nil, errors.New("sagernet TODO").WithPrefix("VLESS").WithPathObj(*this)
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("VLESS").WithPathObj(*this)
	}
}