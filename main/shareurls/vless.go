package shareurls

import "XrayHelper/main/errors"

type VLESS struct{}

func (this *VLESS) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *VLESS) ToOutoundWithTag(tag string) interface{} {
	// TODO
	return ""
}

func newVLESSShareUrl(vlessUrl string) (ShareUrl, error) {
	return nil, errors.New("TODO").WithPrefix("vless")
}
