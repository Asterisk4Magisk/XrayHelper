package addon

type Addon struct {
	//addon
	//ws/httpupgrade/spilthttp/h2->host quic->security grpc->authority
	Host string
	//ws/httpupgrade/spilthttp/h2->path quic->key kcp->seed grpc->serviceName
	Path string
	//tcp/kcp/quic->headerType grpc->mode
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
