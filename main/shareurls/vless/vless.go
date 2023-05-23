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
	// TODO
	return nil, errors.New("TODO").WithPrefix("vless").WithPathObj(*this)
}
