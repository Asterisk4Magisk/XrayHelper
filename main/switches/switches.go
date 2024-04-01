package switches

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/switches/clash"
	"XrayHelper/main/switches/ray"
)

const tagSwitches = "switches"

// Switch implement this interface, that program can deal different core config switch
type Switch interface {
	Execute(args []string) (bool, error)
}

func NewSwitch(coreType string) (Switch, error) {
	switch coreType {
	case "xray", "sing-box":
		return new(ray.RaySwitch), nil
	case "clash.meta", "mihomo":
		return new(clash.ClashSwitch), nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagSwitches)
	}
}
