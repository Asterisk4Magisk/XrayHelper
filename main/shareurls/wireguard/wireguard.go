package wireguard

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/serial"
	"XrayHelper/main/shareurls/addon"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

const tagWireguard = "wireguard"

type Wireguard struct {
	Remarks   string
	SecretKey string
	Server    string
	Port      string
	Address   string
	Reserved  string
	PublicKey string
	Mtu       string
}

func (this *Wireguard) GetNodeInfo() *addon.NodeInfo {
	return &addon.NodeInfo{
		Remarks:  this.Remarks,
		Type:     "Wireguard",
		Host:     this.Server,
		Port:     this.Port,
		Protocol: "udp",
	}
}

func (this *Wireguard) GetNodeInfoStr() string {
	return fmt.Sprintf(color.BlueString("Remarks: ")+"%+v"+color.BlueString(", Type: ")+"Wireguard"+color.BlueString(", Server: ")+"%+v"+color.BlueString(", Port: ")+"%+v"+color.BlueString(", SecretKey: ")+"%+v", this.Remarks, this.Server, this.Port, this.SecretKey)
}

func (this *Wireguard) ToOutboundWithTag(coreType string, tag string) (*serial.OrderedMap, error) {
	switch coreType {
	case "xray":
		var outboundObject serial.OrderedMap
		outboundObject.Set("protocol", "wireguard")
		outboundObject.Set("settings", getWireguardSettingsObjectXray(this))
		outboundObject.Set("tag", tag)
		return &outboundObject, nil
	case "sing-box":
		var outboundObject serial.OrderedMap
		outboundObject.Set("type", "wireguard")
		outboundObject.Set("tag", tag)
		outboundObject.Set("server", this.Server)
		serverPort, _ := strconv.Atoi(this.Port)
		outboundObject.Set("server_port", serverPort)
		if len(this.Address) > 0 {
			var addressArr serial.OrderedArray
			for _, address := range strings.Split(this.Address, ",") {
				addressArr = append(addressArr, address)
			}
			outboundObject.Set("local_address", addressArr)
		}
		outboundObject.Set("private_key", this.SecretKey)
		outboundObject.Set("peer_public_key", this.PublicKey)
		if len(this.Reserved) > 0 {
			var reservedArr serial.OrderedArray
			for _, id := range strings.Split(this.Reserved, ",") {
				iid, _ := strconv.Atoi(id)
				reservedArr = append(reservedArr, iid)
			}
			outboundObject.Set("reserved", reservedArr)
		}
		if len(this.Mtu) > 0 {
			mtu, _ := strconv.Atoi(this.Mtu)
			outboundObject.Set("mtu", mtu)
		}
		return &outboundObject, nil
	default:
		return nil, e.New("unsupported core type " + coreType).WithPrefix(tagWireguard).WithPathObj(*this)
	}
}
