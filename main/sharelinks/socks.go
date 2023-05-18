package sharelinks

import (
	"XrayHelper/main/errors"
)

type Socks struct{}

func (this *Socks) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Socks) ToOutoundJsonWithTag(tag string) string {
	// TODO
	return ""
}

func newSocksShareLink(socksUrl string) (ShareLink, error) {
	return nil, errors.New("TODO").WithPrefix("socks")
}
