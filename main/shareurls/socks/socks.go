package socks

import (
	"XrayHelper/main/errors"
	"fmt"
)

type Socks struct {
	Name     string
	Address  string
	Port     string
	User     string
	Password string
}

func (this *Socks) GetNodeInfo() string {
	return fmt.Sprintf("Name: %+v, Type: Socks, Address: %+v, Port: %+v", this.Name, this.Address, this.Port)
}

func (this *Socks) ToOutoundWithTag(coreType string, tag string) (interface{}, error) {
	// TODO
	return nil, errors.New("TODO").WithPrefix("socks").WithPathObj(*this)
}
