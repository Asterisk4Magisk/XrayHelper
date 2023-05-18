package shareurls

import "XrayHelper/main/errors"

type Vmess struct{}

func (this *Vmess) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Vmess) ToOutoundWithTag(tag string) interface{} {
	// TODO
	return ""
}

func newVmessShareUrl(vmessUrl string) (ShareUrl, error) {
	return nil, errors.New("TODO").WithPrefix("vmess")
}
