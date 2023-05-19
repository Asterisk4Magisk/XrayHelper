package socks

import (
	"XrayHelper/main/errors"
)

type Socks struct{}

func (this *Socks) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Socks) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	// TODO
	return nil, errors.New("TODO").WithPrefix("socks").WithPathObj(*this)
}
