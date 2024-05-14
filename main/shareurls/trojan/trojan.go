package trojan

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"fmt"
	"github.com/fatih/color"
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
	//ws/httpupgrade/h2->host quic->security grpc->authority
	Host string
	//ws/httpupgrade/h2->path quic->key kcp->seed grpc->serviceName
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
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Trojan"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", Network: ")+"%+v"+color.BlueString(", Password: ")+"%+v", this.Remarks, this.Server, this.Port, this.Network, this.Password)
}

func (this *Trojan) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("mux", getMuxObjectXray(false))
		outboundObject.Set("protocol", "trojan")
		outboundObject.Set("settings", getTrojanSettingsObjectXray(this))
		outboundObject.Set("streamSettings", getStreamSettingsObjectXray(this))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "trojan")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		outboundObject.Set("password", this.Password)
		outboundObject.Set("tls", getTrojanTlsObjectSingbox(this))
		outboundObject.Set("transport", getTrojanTransportObjectSingbox(this))
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagTrojan).WithPathObj(*this)
	}
}
