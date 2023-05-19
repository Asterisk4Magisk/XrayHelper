package vless

import (
	"XrayHelper/main/errors"
)

type VLESS struct{}

func (this *VLESS) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *VLESS) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	// TODO
	return nil, errors.New("TODO").WithPrefix("vless").WithPathObj(*this)
}
