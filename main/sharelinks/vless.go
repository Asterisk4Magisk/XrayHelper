package sharelinks

import "XrayHelper/main/errors"

type VLESS struct{}

func (this *VLESS) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *VLESS) ToOutoundJsonWithTag(tag string) string {
	// TODO
	return ""
}

func newVLESSShareLink(vlessUrl string) (ShareLink, error) {
	return nil, errors.New("TODO").WithPrefix("vless")
}
