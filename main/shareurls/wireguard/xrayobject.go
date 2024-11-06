package wireguard

import (
	"XrayHelper/main/serial"
	"strconv"
	"strings"
)

// getWireguardSettingsObjectXray get xray Wireguard SettingsObject
func getWireguardSettingsObjectXray(wireguard *Wireguard) serial.OrderedMap {
	var settingsObject serial.OrderedMap
	settingsObject.Set("secretKey", wireguard.SecretKey)
	if len(wireguard.Mtu) > 0 {
		mtu, _ := strconv.Atoi(wireguard.Mtu)
		settingsObject.Set("mtu", mtu)
	}
	if len(wireguard.Reserved) > 0 {
		var reservedArr serial.OrderedArray
		for _, id := range strings.Split(wireguard.Reserved, ",") {
			iid, _ := strconv.Atoi(id)
			reservedArr = append(reservedArr, iid)
		}
		settingsObject.Set("reserved", reservedArr)
	}
	if len(wireguard.Address) > 0 {
		var addressArr serial.OrderedArray
		for _, address := range strings.Split(wireguard.Address, ",") {
			addressArr = append(addressArr, address)
		}
		settingsObject.Set("address", addressArr)
	}
	var peersArr serial.OrderedArray
	var peers serial.OrderedMap
	peers.Set("endpoint", wireguard.Server+":"+wireguard.Port)
	peers.Set("publicKey", wireguard.PublicKey)
	peersArr = append(peersArr, peers)
	settingsObject.Set("peers", peersArr)
	settingsObject.Set("domainStrategy", "ForceIP")
	return settingsObject
}
