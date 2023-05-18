package sharelinks

import "XrayHelper/main/errors"

type Trojan struct{}

func (this *Trojan) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Trojan) ToOutoundJsonWithTag(tag string) string {
	// TODO
	return ""
}

func newTrojanShareLink(trojanUrl string) (ShareLink, error) {
	return nil, errors.New("TODO").WithPrefix("trojan")
}
