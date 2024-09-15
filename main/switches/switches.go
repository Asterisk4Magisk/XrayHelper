package switches

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/switches/clash"
	"XrayHelper/main/switches/ray"
)

const tagSwitches = "switches"

// Switch implement this interface, that program can deal different core config switch
type Switch interface {
	Execute(args []string) (bool, error)
	Get(custom bool) serial.OrderedArray
	Set(custom bool, index int) error
	Choose(custom bool, index int) any
}

func NewSwitch(coreType string) (Switch, error) {
	switch coreType {
	case "xray", "sing-box", "hysteria2":
		return new(ray.RaySwitch), nil
	case "mihomo":
		return new(clash.ClashSwitch), nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagSwitches)
	}
}
