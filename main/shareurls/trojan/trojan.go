package trojan

import (
	"XrayHelper/main/errors"
)

type Trojan struct{}

func (this *Trojan) GetNodeInfo() string {
	// TODO
	return ""
}

func (this *Trojan) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	// TODO
	return nil, errors.New("TODO").WithPrefix("trojan").WithPathObj(*this)
}
