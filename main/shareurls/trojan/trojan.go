package trojan

import (
	"XrayHelper/main/errors"
	"fmt"
)

type Trojan struct {
	//basic
	Name     string
	Password string
	Address  string
	Port     string
	Network  string
	Security string

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

func (this *Trojan) GetNodeInfo() string {
	return fmt.Sprintf("Name: %+v, Type: Trojan, Address: %+v, Port: %+v, Network: %+v, Password: %+v", this.Name, this.Address, this.Port, this.Network, this.Password)
}

func (this *Trojan) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	switch coreType {
	case "xray":
		outboundObject := make(map[string]interface{})
		outboundObject["mux"] = getMuxObjectXray(false)
		outboundObject["protocol"] = "trojan"
		outboundObject["settings"] = getTrojanSettingsObjectXray(this)
		outboundObject["streamSettings"] = getStreamSettingsObjectXray(this)
		outboundObject["tag"] = tag
		return outboundObject, nil
	case "sing-box":
		outboundObject := make(map[string]interface{})
		outboundObject["type"] = "trojan"
		outboundObject["tag"] = tag
		outboundObject["server"] = this.Address
		outboundObject["server_port"] = this.Port
		outboundObject["password"] = this.Password
		outboundObject["tls"] = getTrojanTlsObjectSingbox(this)
		outboundObject["transport"] = getTrojanTransportObjectSingbox(this)
		return outboundObject, nil
	default:
		return nil, errors.New("not supported core type " + coreType).WithPrefix("vmess").WithPathObj(*this)
	}
}
