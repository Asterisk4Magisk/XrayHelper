package sharelinks

import "XrayHelper/main/errors"

type Trojan struct{}

func (this *Trojan) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Trojan) ToOutoundWithTag(tag string) interface{} {
	// TODO
	return ""
}

func newTrojanShareLink(trojanUrl string) (ShareLink, error) {
	return nil, errors.New("TODO").WithPrefix("trojan")
}
