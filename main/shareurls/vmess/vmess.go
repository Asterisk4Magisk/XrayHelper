package vmess

import (
	"XrayHelper/main/errors"
)

type Vmess struct{}

func (this *Vmess) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Vmess) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	// TODO
	return nil, errors.New("TODO").WithPrefix("vmess").WithPathObj(*this)
}
