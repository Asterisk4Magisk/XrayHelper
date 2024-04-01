package trojan

import (
	e "XrayHelper/main/errors"
	"fmt"
	"strconv"
)

const tagTrojan = "trojan"

type Trojan struct {
	//basic
	Remarks  string
	Password string
	Server   string
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
	PublicKey string //pbk
	ShortId   string //sid
	SpiderX   string //spx
}

func (this *Trojan) GetNodeInfo() string {
	return fmt.Sprintf("Remarks: %+v, Type: Trojan, Server: %+v, Port: %+v, Network: %+v, Password: %+v", this.Remarks, this.Server, this.Port, this.Network, this.Password)
}

func (this *Trojan) ToOutboundWithTag(coreType string, tag string) (interface{}, error) {
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
		outboundObject["server"] = this.Server
		outboundObject["server_port"], _ = strconv.Atoi(this.Port)
		outboundObject["password"] = this.Password
		outboundObject["tls"] = getTrojanTlsObjectSingbox(this)
		outboundObject["transport"] = getTrojanTransportObjectSingbox(this)
		return outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagTrojan).WithPathObj(*this)
	}
}
