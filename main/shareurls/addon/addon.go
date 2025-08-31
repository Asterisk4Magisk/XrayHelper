package addon

type Addon struct {
	//addon
	//ws/httpupgrade/spilthttp/h2/h3/xhttp->host quic->security grpc->authority
	Host string
	//ws/httpupgrade/spilthttp/h2/h3/xhttp->path quic->key kcp->seed grpc->serviceName
	Path string
	//tcp/kcp/quic->headerType grpc/xhttp->mode
	Type string
	//xhttp->XHTTPObject(json)
	Extra string

	//tls
	Sni         string
	FingerPrint string
	Alpn        string
	//reality
	PublicKey     string //pbk(password)
	ShortId       string //sid
	Mldsa65Verify string //pqv
	SpiderX       string //spx
}

type NodeInfo struct {
	Remarks  string `json:"remarks"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}
